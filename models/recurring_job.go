package models

import (
	"encoding/json"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// RecurringJob 重复任务模型
// 用于存储重复任务的配置信息，以便在系统重启后恢复任务
type RecurringJob struct {
	gorm.Model
	Name           string     `gorm:"uniqueIndex;size:255" json:"name"` // 任务名称，唯一
	JobKey         string     `gorm:"size:255" json:"job_key"`          // 任务键，用于在gocron中识别任务
	FunctionName   string     `gorm:"size:255" json:"function_name"`    // 函数名称
	CronExpression string     `gorm:"size:100" json:"cron_expression"`  // Cron表达式
	Args           string     `gorm:"type:text" json:"args"`            // 参数，JSON格式
	LastRunAt      *time.Time `json:"last_run_at"`                      // 上次执行时间
	NextRunAt      *time.Time `json:"next_run_at"`                      // 下次执行时间
	Times          int        `json:"times"`                            // 执行次数限制(0表示无限)
	TimesRun       int        `json:"times_run"`                        // 已执行次数
	Status         string     `gorm:"size:50" json:"status"`            // 状态(active,paused,completed)
	ErrorCount     int        `json:"error_count"`                      // 错误次数
	LastError      string     `gorm:"type:text" json:"last_error"`      // 最后一次错误
}

// SetArgs 设置任务参数
func (r *RecurringJob) SetArgs(args interface{}) error {
	if args == nil {
		r.Args = ""
		return nil
	}

	b, err := json.Marshal(args)
	if err != nil {
		return err
	}

	r.Args = string(b)
	return nil
}

// GetArgs 获取任务参数
func (r *RecurringJob) GetArgs(dest interface{}) error {
	if r.Args == "" {
		return nil
	}

	return json.Unmarshal([]byte(r.Args), dest)
}

// RecurringJobExecution 重复任务执行记录
// 用于记录每次任务执行的情况
type RecurringJobExecution struct {
	gorm.Model
	RecurringJobID uint       `json:"recurring_job_id"`        // 关联的重复任务ID
	StartedAt      time.Time  `json:"started_at"`              // 开始执行时间
	FinishedAt     *time.Time `json:"finished_at"`             // 结束执行时间
	Success        bool       `json:"success"`                 // 是否成功
	Error          string     `gorm:"type:text" json:"error"`  // 错误信息
	Output         string     `gorm:"type:text" json:"output"` // 输出信息
	Duration       int64      `json:"duration"`                // 执行持续时间(毫秒)
}

// Info 记录一条信息级别的输出，自动添加时间戳
func (e *RecurringJobExecution) Info(format string, args ...interface{}) {
	e.logWithLevel("INFO", format, args...)
}

// Warning 记录一条警告级别的输出，自动添加时间戳
func (e *RecurringJobExecution) Warning(format string, args ...interface{}) {
	e.logWithLevel("WARN", format, args...)
}

// Error 记录一条错误级别的输出，自动添加时间戳
func (e *RecurringJobExecution) LogError(format string, args ...interface{}) {
	e.logWithLevel("ERROR", format, args...)
}

// Debug 记录一条调试级别的输出，自动添加时间戳
func (e *RecurringJobExecution) Debug(format string, args ...interface{}) {
	e.logWithLevel("DEBUG", format, args...)
}

// logWithLevel 内部方法，添加时间戳和日志级别
func (e *RecurringJobExecution) logWithLevel(level, format string, args ...interface{}) {
	timestamp := time.Now().Format("2006-01-02 15:04:05.000")
	message := fmt.Sprintf(format, args...)
	// 使用清晰的格式，确保级别标记更加突出
	logLine := fmt.Sprintf("[%s] [%s] %s", timestamp, level, message)

	if e.Output != "" {
		e.Output += "\n"
	}
	e.Output += logLine
}
