package recurring

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/naokij/qor5boot/models"
)

// 错误定义
var (
	ErrJobNotFound      = errors.New("找不到指定任务")
	ErrJobAlreadyExists = errors.New("任务已存在")
	ErrInvalidInterval  = errors.New("无效的时间间隔")
	ErrInvalidUnit      = errors.New("无效的时间单位")
	ErrInvalidFunction  = errors.New("无效的函数")
)

// 注册的函数类型
type JobFunc func(ctx context.Context, args []byte, execution *models.RecurringJobExecution) error

// TaskManager 任务管理器
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

// NewTaskManager 创建一个新的任务管理器
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

// RegisterFunction 注册一个可以被调度的函数
func (m *TaskManager) RegisterFunction(name string, fn JobFunc) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.functions[name] = fn
}

// Start 启动任务管理器
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

// Stop 停止任务管理器
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
func (m *TaskManager) AddJob(name, functionName string, interval int, unit string, args interface{}, times int) (*models.RecurringJob, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 检查函数是否已注册
	if _, ok := m.functions[functionName]; !ok {
		return nil, ErrInvalidFunction
	}

	// 检查任务名是否已存在
	var count int64
	m.db.Model(&models.RecurringJob{}).Where("name = ?", name).Count(&count)
	if count > 0 {
		return nil, ErrJobAlreadyExists
	}

	// 验证时间间隔
	if interval <= 0 {
		return nil, ErrInvalidInterval
	}

	// 验证时间单位
	validUnits := map[string]bool{
		"second": true, "minute": true, "hour": true, "day": true, "week": true,
		"seconds": true, "minutes": true, "hours": true, "days": true, "weeks": true,
	}
	if !validUnits[unit] {
		return nil, ErrInvalidUnit
	}

	// 创建任务记录
	job := &models.RecurringJob{
		Name:         name,
		JobKey:       uuid.New().String(), // 生成唯一的任务键
		FunctionName: functionName,
		Interval:     interval,
		Unit:         unit,
		Times:        times,
		Status:       "active",
	}

	// 设置参数
	if err := job.SetArgs(args); err != nil {
		return nil, err
	}

	// 保存到数据库
	if err := m.db.Create(job).Error; err != nil {
		return nil, err
	}

	// 调度任务
	scheduledJob, err := m.scheduleJob(job)
	if err != nil {
		// 如果调度失败，更新任务状态并返回错误
		m.db.Model(job).Updates(map[string]interface{}{
			"status":     "error",
			"last_error": err.Error(),
		})
		return nil, err
	}

	// 记录下次执行时间
	nextRun := scheduledJob.NextRun()
	m.db.Model(job).Update("next_run_at", nextRun)

	return job, nil
}

// RemoveJob 移除指定任务
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

	// 从数据库中软删除
	return m.db.Delete(&job).Error
}

// PauseJob 暂停任务
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

// ResumeJob 恢复任务
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

// GetJob 获取任务信息
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
func (m *TaskManager) ListJobs() ([]models.RecurringJob, error) {
	var jobs []models.RecurringJob
	err := m.db.Find(&jobs).Error
	return jobs, err
}

// RunJobNow 立即执行一次任务
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

// 调度任务
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

	// 创建调度器
	scheduler := m.scheduler.Every(job.Interval)

	// 设置时间单位
	switch job.Unit {
	case "second", "seconds":
		scheduler = scheduler.Second()
	case "minute", "minutes":
		scheduler = scheduler.Minute()
	case "hour", "hours":
		scheduler = scheduler.Hour()
	case "day", "days":
		scheduler = scheduler.Day()
	case "week", "weeks":
		scheduler = scheduler.Week()
	default:
		return nil, ErrInvalidUnit
	}

	// 调度任务
	scheduledJob, err := scheduler.Do(execFn)
	if err != nil {
		return nil, err
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

// 执行任务
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

// 完成执行记录
func finishExecution(db *gorm.DB, execution *models.RecurringJobExecution, success bool, errorMsg, output string) {
	now := time.Now()
	execution.FinishedAt = &now
	execution.Duration = now.Sub(execution.StartedAt).Milliseconds()
	execution.Success = success
	execution.Error = errorMsg
	execution.Output = output
	db.Save(execution)
}
