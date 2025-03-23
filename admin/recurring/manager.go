// Package recurring 提供了一个完整的重复任务管理系统，支持定时任务的创建、调度、执行和监控。
// 该系统基于 gocron 实现，提供了以下主要功能：
// 1. 支持多种时间间隔的任务调度（秒、分钟、小时、天、周）
// 2. 支持任务执行次数限制
// 3. 支持任务的暂停、恢复和立即执行
// 4. 提供任务执行历史记录和错误追踪
// 5. 支持并发安全的任务管理
package recurring

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/go-co-op/gocron"
	"gorm.io/gorm"

	"github.com/naokij/qor5boot/models"
)

// 错误定义
var (
	// ErrJobNotFound 表示找不到指定的任务
	ErrJobNotFound = errors.New("找不到指定任务")
	// ErrJobAlreadyExists 表示任务名称已存在
	ErrJobAlreadyExists = errors.New("任务已存在")
	// ErrInvalidFunction 表示函数未注册
	ErrInvalidFunction = errors.New("无效的函数")
	// ErrDuplicateName 表示任务名称已存在
	ErrDuplicateName = errors.New("任务名称已存在")
)

// JobFunc 定义了任务函数的类型签名
// 参数说明：
// - ctx: 任务执行的上下文，包含超时控制
// - args: 任务参数，JSON格式的字节数组
// - execution: 任务执行记录，用于记录执行状态和结果
type JobFunc func(ctx context.Context, args []byte, execution *models.RecurringJobExecution) error

// TaskManager 是任务管理器的核心结构体，负责所有任务的调度和执行
// 字段说明：
// - db: 数据库连接，用于持久化任务和执行记录
// - scheduler: gocron调度器，负责任务的定时执行
// - functions: 已注册的任务函数映射表
// - jobs: 当前运行的任务映射表
// - jobModels: 任务模型映射表
// - mu: 并发锁，用于保证并发安全
// - isRunning: 管理器运行状态标志
// - defaultLogger: 默认日志记录器
type TaskManager struct {
	db            *gorm.DB
	scheduler     *gocron.Scheduler
	functions     map[string]JobFunc
	jobs          map[string]*gocron.Job
	jobModels     map[string]*models.RecurringJob
	mu            sync.RWMutex
	isRunning     bool
	defaultLogger *log.Logger
}

// NewTaskManager 创建一个新的任务管理器实例
// 参数：
// - db: 数据库连接
// 返回：
// - *TaskManager: 新创建的任务管理器实例
func NewTaskManager(db *gorm.DB) *TaskManager {
	// 初始化gocron调度器
	scheduler := gocron.NewScheduler(time.Local)

	// 创建任务管理器
	manager := &TaskManager{
		db:            db,
		scheduler:     scheduler,
		functions:     make(map[string]JobFunc),
		jobs:          make(map[string]*gocron.Job),
		jobModels:     make(map[string]*models.RecurringJob),
		defaultLogger: log.Default(),
	}

	// 自动迁移数据库模型
	err := db.AutoMigrate(&models.RecurringJob{}, &models.RecurringJobExecution{})
	if err != nil {
		log.Printf("迁移重复任务模型失败: %v", err)
	}

	return manager
}

// RegisterFunction 注册一个新的任务函数
// 参数：
// - name: 函数名称，用于在任务中引用
// - fn: 任务函数实现
func (m *TaskManager) RegisterFunction(name string, fn JobFunc) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.functions[name] = fn
}

// Start 启动任务管理器，加载所有活动的任务并开始调度
// 返回：
// - error: 启动过程中的错误信息
func (m *TaskManager) Start() error {
	if m.isRunning {
		return nil
	}

	// 从数据库加载所有活动的任务
	var jobs []models.RecurringJob
	err := m.db.Where("status = ?", "active").Find(&jobs).Error
	if err != nil {
		return fmt.Errorf("加载任务失败: %w", err)
	}

	// 恢复所有任务
	for _, job := range jobs {
		_, err := m.scheduleJob(&job)
		if err != nil {
			log.Printf("恢复任务 %s 失败: %v", job.Name, err)
		}
	}

	// 启动调度器
	m.scheduler.StartAsync()
	m.isRunning = true
	log.Println("重复任务管理器已启动")
	return nil
}

// Stop 停止任务管理器，清理所有资源并等待任务完成
func (m *TaskManager) Stop() {
	if !m.isRunning {
		return
	}

	log.Println("正在关闭重复任务管理器...")

	// 获取锁以确保其他操作暂停
	m.mu.Lock()
	defer m.mu.Unlock()

	// 停止调度器并等待所有任务完成
	m.scheduler.Stop()

	// 清理资源
	for key := range m.jobs {
		delete(m.jobs, key)
	}

	for key := range m.jobModels {
		delete(m.jobModels, key)
	}

	m.isRunning = false
	log.Println("重复任务管理器已完全停止")
}

// AddJob 添加一个新的重复任务
// 参数：
// - name: 任务名称
// - functionName: 要执行的函数名称
// - args: 执行函数的参数
// - times: 执行次数限制(0表示无限)
// - cronExpression: Cron表达式
// 返回：
// - *models.RecurringJob: 创建的任务对象
// - error: 创建过程中的错误信息
func (m *TaskManager) AddJob(name, functionName string, args interface{}, times int, cronExpression string) (*models.RecurringJob, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 检查名称是否已存在
	var count int64
	if err := m.db.Unscoped().Model(&models.RecurringJob{}).Where("name = ? AND deleted_at IS NULL", name).Count(&count).Error; err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, ErrDuplicateName
	}

	// 检查函数是否已注册
	if _, ok := m.functions[functionName]; !ok {
		return nil, ErrInvalidFunction
	}

	// 创建任务对象
	job := models.RecurringJob{
		Name:           name,
		JobKey:         fmt.Sprintf("%s_%d", name, time.Now().UnixNano()),
		FunctionName:   functionName,
		CronExpression: cronExpression,
		Times:          times,
		Status:         "active",
	}

	// 设置参数
	if err := job.SetArgs(args); err != nil {
		return nil, err
	}

	// 保存到数据库
	if err := m.db.Create(&job).Error; err != nil {
		return nil, err
	}

	// 如果任务管理器已启动，立即调度任务
	if m.isRunning {
		scheduledJob, err := m.scheduleJob(&job)
		if err != nil {
			// 调度失败，更新状态
			m.db.Model(&job).Updates(map[string]interface{}{
				"status":     "error",
				"last_error": err.Error(),
			})
			return &job, err
		}

		// 获取下次执行时间
		nextRun := scheduledJob.NextRun()
		m.db.Model(&job).Update("next_run_at", nextRun)
	}

	return &job, nil
}

// RemoveJob 移除指定的任务
// 参数：
// - name: 要移除的任务名称
// 返回：
// - error: 移除过程中的错误信息
func (m *TaskManager) RemoveJob(name string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 查找任务
	var job models.RecurringJob
	err := m.db.Where("name = ?", name).First(&job).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrJobNotFound
		}
		return err
	}

	// 从调度器中移除任务
	if scheduledJob, exists := m.jobs[job.JobKey]; exists {
		m.scheduler.RemoveByReference(scheduledJob)
		delete(m.jobs, job.JobKey)
		delete(m.jobModels, job.JobKey)
	}

	// 从数据库中真正物理删除（不是软删除）
	return m.db.Unscoped().Delete(&job).Error
}

// PauseJob 暂停指定的任务
// 参数：
// - name: 要暂停的任务名称
// 返回：
// - error: 暂停过程中的错误信息
func (m *TaskManager) PauseJob(name string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 查找任务
	var job models.RecurringJob
	err := m.db.Where("name = ?", name).First(&job).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrJobNotFound
		}
		return err
	}

	// 从调度器中移除任务
	if scheduledJob, exists := m.jobs[job.JobKey]; exists {
		m.scheduler.RemoveByReference(scheduledJob)
		delete(m.jobs, job.JobKey)
	}

	// 更新状态为暂停
	return m.db.Model(&job).Updates(map[string]interface{}{
		"status": "paused",
	}).Error
}

// ResumeJob 恢复已暂停的任务
// 参数：
// - name: 要恢复的任务名称
// 返回：
// - error: 恢复过程中的错误信息
func (m *TaskManager) ResumeJob(name string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 查找任务
	var job models.RecurringJob
	err := m.db.Where("name = ?", name).First(&job).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrJobNotFound
		}
		return err
	}

	// 如果任务不是暂停状态，返回错误
	if job.Status != "paused" {
		return errors.New("只能恢复处于暂停状态的任务")
	}

	// 更新状态为活动
	err = m.db.Model(&job).Updates(map[string]interface{}{
		"status": "active",
	}).Error
	if err != nil {
		return err
	}

	// 重新调度任务
	_, err = m.scheduleJob(&job)
	return err
}

// GetJob 获取指定任务的详细信息
// 参数：
// - name: 任务名称
// 返回：
// - *models.RecurringJob: 任务对象
// - error: 获取过程中的错误信息
func (m *TaskManager) GetJob(name string) (*models.RecurringJob, error) {
	var job models.RecurringJob
	err := m.db.Where("name = ?", name).First(&job).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrJobNotFound
		}
		return nil, err
	}
	return &job, nil
}

// ListJobs 列出所有任务
// 返回：
// - []models.RecurringJob: 任务列表
// - error: 获取过程中的错误信息
func (m *TaskManager) ListJobs() ([]models.RecurringJob, error) {
	var jobs []models.RecurringJob
	err := m.db.Find(&jobs).Error
	return jobs, err
}

// RunJobNow 立即执行一次指定的任务
// 参数：
// - name: 要执行的任务名称
// 返回：
// - error: 执行过程中的错误信息
func (m *TaskManager) RunJobNow(name string) error {
	var job models.RecurringJob
	err := m.db.Where("name = ?", name).First(&job).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrJobNotFound
		}
		return err
	}

	// 执行任务
	go m.executeJob(&job)
	return nil
}

// scheduleJob 内部方法，用于调度任务
// 参数：
// - job: 要调度的任务对象
// 返回：
// - *gocron.Job: 调度后的任务对象
// - error: 调度过程中的错误信息
func (m *TaskManager) scheduleJob(job *models.RecurringJob) (*gocron.Job, error) {
	// 检查函数是否已注册
	_, ok := m.functions[job.FunctionName]
	if !ok {
		return nil, ErrInvalidFunction
	}

	// 创建执行函数
	execFn := func() {
		m.executeJob(job)
	}

	// 始终使用Cron表达式调度
	if job.CronExpression == "" {
		return nil, fmt.Errorf("Cron表达式不能为空")
	}

	// 使用Cron表达式调度
	// 注意：使用标准5字段Cron表达式（分 时 日 月 周）
	scheduledJob, err := m.scheduler.Cron(job.CronExpression).Do(execFn)
	if err != nil {
		return nil, fmt.Errorf("无效的Cron表达式: %w", err)
	}

	// 设置执行次数限制 - 考虑已执行的次数
	if job.Times > 0 {
		remainingRuns := job.Times - job.TimesRun
		if remainingRuns <= 0 {
			// 任务已完成，将状态更新为"completed"
			m.db.Model(&models.RecurringJob{}).Where("id = ?", job.ID).Update("status", "completed")
			return scheduledJob, nil // 返回但不实际调度
		}
		scheduledJob.LimitRunsTo(remainingRuns)
	}

	// 存储任务引用
	m.jobs[job.JobKey] = scheduledJob
	m.jobModels[job.JobKey] = job

	return scheduledJob, nil
}

// executeJob 内部方法，用于执行任务
// 参数：
// - job: 要执行的任务对象
func (m *TaskManager) executeJob(job *models.RecurringJob) {
	// 锁定任务执行，防止并发问题
	m.mu.Lock()
	// 获取最新任务数据，避免使用过期数据
	var updatedJob models.RecurringJob
	if err := m.db.First(&updatedJob, job.ID).Error; err != nil {
		log.Printf("任务执行前获取最新数据失败: %v", err)
		m.mu.Unlock()
		return
	}

	// 检查任务是否已完成或暂停
	if updatedJob.Status != "active" {
		log.Printf("任务 %s 状态为 %s，跳过执行", updatedJob.Name, updatedJob.Status)
		m.mu.Unlock()
		return
	}

	// 检查是否已达到执行次数限制
	if updatedJob.Times > 0 && updatedJob.TimesRun >= updatedJob.Times {
		log.Printf("任务 %s 已达到执行次数限制 (%d/%d)，标记为完成",
			updatedJob.Name, updatedJob.TimesRun, updatedJob.Times)
		m.db.Model(&models.RecurringJob{}).Where("id = ?", updatedJob.ID).Update("status", "completed")

		// 从调度器中移除任务
		if scheduledJob, exists := m.jobs[updatedJob.JobKey]; exists {
			m.scheduler.RemoveByReference(scheduledJob)
			delete(m.jobs, updatedJob.JobKey)
			delete(m.jobModels, updatedJob.JobKey)
		}
		m.mu.Unlock()
		return
	}
	m.mu.Unlock()

	// 创建执行记录
	execution := &models.RecurringJobExecution{
		RecurringJobID: job.ID,
		StartedAt:      time.Now(),
	}
	m.db.Create(execution)

	// 获取任务函数
	fn, ok := m.functions[job.FunctionName]
	if !ok {
		finishExecution(m.db, execution, false, ErrInvalidFunction.Error(), "")
		return
	}

	// 创建上下文
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()

	// 执行任务函数
	err := fn(ctx, []byte(job.Args), execution)

	// 更新执行记录
	finishTime := time.Now()
	duration := finishTime.Sub(execution.StartedAt).Milliseconds()
	success := err == nil

	errorMsg := ""
	if err != nil {
		errorMsg = err.Error()
	}

	// 更新执行记录
	execution.FinishedAt = &finishTime
	execution.Duration = duration
	execution.Success = success
	execution.Error = errorMsg
	m.db.Save(execution)

	// 这里我们使用事务来确保原子更新
	tx := m.db.Begin()

	// 再次获取最新的任务数据
	if err := tx.First(&updatedJob, job.ID).Error; err != nil {
		log.Printf("获取最新任务数据失败: %v", err)
		tx.Rollback()
		return
	}

	// 更新任务记录
	newTimesRun := updatedJob.TimesRun + 1
	updates := map[string]interface{}{
		"last_run_at": execution.StartedAt,
		"times_run":   newTimesRun,
	}

	if !success {
		updates["error_count"] = updatedJob.ErrorCount + 1
		updates["last_error"] = errorMsg
	}

	// 如果有gocron的任务引用，获取下次执行时间
	if scheduledJob, exists := m.jobs[job.JobKey]; exists {
		nextRun := scheduledJob.NextRun()
		updates["next_run_at"] = nextRun
	}

	// 直接更新数据库
	if err := tx.Model(&models.RecurringJob{}).Where("id = ?", job.ID).Updates(updates).Error; err != nil {
		log.Printf("更新任务状态失败: %v", err)
		tx.Rollback()
		return
	}

	// 检查是否达到执行次数限制
	if updatedJob.Times > 0 && newTimesRun >= updatedJob.Times {
		if err := tx.Model(&models.RecurringJob{}).Where("id = ?", job.ID).Update("status", "completed").Error; err != nil {
			log.Printf("更新任务状态为completed失败: %v", err)
			tx.Rollback()
			return
		}

		// 提交事务
		tx.Commit()

		// 从调度器中移除任务 (在事务外执行，避免锁定)
		m.mu.Lock()
		if scheduledJob, exists := m.jobs[job.JobKey]; exists {
			m.scheduler.RemoveByReference(scheduledJob)
			delete(m.jobs, job.JobKey)
			delete(m.jobModels, job.JobKey)
		}
		m.mu.Unlock()
	} else {
		// 提交事务
		tx.Commit()
	}
}

// finishExecution 内部方法，用于完成执行记录
// 参数：
// - db: 数据库连接
// - execution: 执行记录对象
// - success: 是否执行成功
// - errorMsg: 错误信息
// - output: 输出信息
func finishExecution(db *gorm.DB, execution *models.RecurringJobExecution, success bool, errorMsg, output string) {
	now := time.Now()
	execution.FinishedAt = &now
	execution.Duration = now.Sub(execution.StartedAt).Milliseconds()
	execution.Success = success
	execution.Error = errorMsg
	execution.Output = output
	db.Save(execution)
}

// UpdateJob 更新现有任务的配置并重新调度，保留原有的统计信息和状态
func (m *TaskManager) UpdateJob(jobID uint, name, functionName string, args interface{}, times int, cronExpression string, keepStatus bool) (*models.RecurringJob, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 查找现有任务
	var job models.RecurringJob
	if err := m.db.First(&job, jobID).Error; err != nil {
		return nil, err
	}

	// 如果任务名称改变了，需要检查新名称是否可用
	if job.Name != name {
		var count int64
		if err := m.db.Unscoped().Model(&models.RecurringJob{}).Where("name = ? AND deleted_at IS NULL AND id != ?", name, jobID).Count(&count).Error; err != nil {
			return nil, err
		}
		if count > 0 {
			return nil, ErrDuplicateName
		}
	}

	// 保存原有的统计信息和状态
	originalStatus := job.Status
	originalTimesRun := job.TimesRun
	originalErrorCount := job.ErrorCount
	originalLastError := job.LastError
	originalLastRunAt := job.LastRunAt

	// 从调度器中移除当前任务（如果存在）
	if scheduledJob, exists := m.jobs[job.JobKey]; exists {
		m.scheduler.RemoveByReference(scheduledJob)
		delete(m.jobs, job.JobKey)
		delete(m.jobModels, job.JobKey)
	}

	// 更新任务配置
	job.Name = name
	job.FunctionName = functionName
	job.CronExpression = cronExpression
	job.Times = times

	// 如果不保持状态，且原状态是completed，则重置为active
	if !keepStatus && originalStatus == "completed" {
		job.Status = "active"
	}

	// 设置参数
	if err := job.SetArgs(args); err != nil {
		return nil, err
	}

	// 恢复原有的统计信息
	job.TimesRun = originalTimesRun
	job.ErrorCount = originalErrorCount
	job.LastError = originalLastError
	job.LastRunAt = originalLastRunAt

	// 保存更新后的任务
	if err := m.db.Save(&job).Error; err != nil {
		return nil, err
	}

	// 如果任务状态是active，重新调度
	if job.Status == "active" && m.isRunning {
		scheduledJob, err := m.scheduleJob(&job)
		if err != nil {
			// 调度失败，更新状态
			m.db.Model(&job).Updates(map[string]interface{}{
				"status":     "error",
				"last_error": err.Error(),
			})
			return &job, err
		}

		// 获取下次执行时间
		nextRun := scheduledJob.NextRun()
		m.db.Model(&job).Update("next_run_at", nextRun)
	}

	return &job, nil
}
