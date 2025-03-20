package recurring

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/qor5/admin/v3/presets"
	"github.com/qor5/web/v3"
	v "github.com/qor5/x/v3/ui/vuetify"
	h "github.com/theplant/htmlgo"
	"gorm.io/gorm"

	"github.com/naokij/qor5boot/models"
)

// RecurringJobManager 处理重复任务的管理器
type RecurringJobManager struct {
	taskManager  *TaskManager
	pb           *presets.Builder
	modelBuilder *presets.ModelBuilder
}

// NewRecurringJobManager 创建重复任务管理器
func NewRecurringJobManager(db *gorm.DB, pb *presets.Builder) *RecurringJobManager {
	taskManager := NewTaskManager(db)

	// 创建模型并设置标签，只在这里注册一次
	modelBuilder := pb.Model(&models.RecurringJob{})
	modelBuilder.Label("RecurringJobs")

	manager := &RecurringJobManager{
		taskManager:  taskManager,
		pb:           pb,
		modelBuilder: modelBuilder,
	}

	// 注册一些示例函数
	manager.registerSampleFunctions()

	// 注册管理界面
	manager.registerAdminUI()

	return manager
}

// Start 启动管理器
func (m *RecurringJobManager) Start() error {
	return m.taskManager.Start()
}

// Stop 停止管理器
func (m *RecurringJobManager) Stop() {
	m.taskManager.Stop()
}

// 注册示例函数
func (m *RecurringJobManager) registerSampleFunctions() {
	// 日志函数 - 简单地记录一条消息
	m.taskManager.RegisterFunction("log", func(ctx context.Context, args []byte, execution *models.RecurringJobExecution) error {
		var message string
		if len(args) > 0 {
			if err := json.Unmarshal(args, &message); err != nil {
				message = string(args)
			}
		} else {
			message = "执行定时日志任务"
		}

		log.Printf("[重复任务日志] %s", message)
		return nil
	})

	// 测试函数 - 可以随机成功或失败
	m.taskManager.RegisterFunction("test", func(ctx context.Context, args []byte, execution *models.RecurringJobExecution) error {
		// 这里可以做一些测试工作
		log.Printf("[重复任务测试] 执行测试任务")

		// 等待一些时间模拟工作
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(2 * time.Second):
		}

		// 这里可以添加一些随机逻辑测试错误处理等
		return nil
	})

	// 失败函数 - 总是返回错误
	m.taskManager.RegisterFunction("fail", func(ctx context.Context, args []byte, execution *models.RecurringJobExecution) error {
		log.Printf("[重复任务失败] 执行失败任务")
		return fmt.Errorf("这个任务总是失败")
	})
}

// 注册管理界面
func (m *RecurringJobManager) registerAdminUI() {
	// 配置列表视图 - 移除自定义的Actions列
	m.modelBuilder.Listing("ID", "Name", "FunctionName", "Interval", "Runs", "Status", "LastRunAt", "NextRunAt", "ErrorCount", "Actions")

	// 配置编辑视图
	m.modelBuilder.Editing("Name", "FunctionName", "Interval", "Unit", "Times", "Args")

	// 为Interval字段创建显示组件(合并Interval和Unit)
	m.modelBuilder.Listing().Field("Interval").ComponentFunc(func(obj interface{}, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		job, ok := obj.(*models.RecurringJob)
		if !ok {
			return nil
		}

		// 根据单位类型生成对应的中文显示
		unitText := ""
		switch job.Unit {
		case "second", "seconds":
			unitText = "秒"
		case "minute", "minutes":
			unitText = "分钟"
		case "hour", "hours":
			unitText = "小时"
		case "day", "days":
			unitText = "天"
		case "week", "weeks":
			unitText = "周"
		default:
			unitText = job.Unit
		}

		// 生成间隔显示文本
		return h.Td(h.Text(fmt.Sprintf("每 %d %s", job.Interval, unitText)))
	})

	// 为Runs字段创建显示组件(合并Times和TimesRun)
	m.modelBuilder.Listing().Field("Runs").ComponentFunc(func(obj interface{}, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		job, ok := obj.(*models.RecurringJob)
		if !ok {
			return nil
		}

		if job.Times > 0 {
			// 有限次数的任务
			return h.Td(h.Text(fmt.Sprintf("%d / %d", job.TimesRun, job.Times)))
		} else {
			// 无限次数的任务
			return h.Td(h.Text(fmt.Sprintf("%d / ∞", job.TimesRun)))
		}
	})

	// 为FunctionName字段创建选择器
	m.modelBuilder.Editing().Field("FunctionName").ComponentFunc(func(obj interface{}, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		job, ok := obj.(*models.RecurringJob)
		if !ok {
			return nil
		}

		// 获取可用函数列表
		options := []v.DefaultOptionItem{
			{Text: "日志函数", Value: "log"},
			{Text: "测试函数", Value: "test"},
			{Text: "失败函数", Value: "fail"},
		}

		return v.VSelect().
			Label("函数名").
			Items(options).
			ItemTitle("text").
			ItemValue("value").
			Attr(web.VField("FunctionName", job.FunctionName)...)
	})

	// 为Unit字段创建选择器
	m.modelBuilder.Editing().Field("Unit").ComponentFunc(func(obj interface{}, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		job, ok := obj.(*models.RecurringJob)
		if !ok {
			return nil
		}

		options := []v.DefaultOptionItem{
			{Text: "秒", Value: "second"},
			{Text: "分钟", Value: "minute"},
			{Text: "小时", Value: "hour"},
			{Text: "天", Value: "day"},
			{Text: "周", Value: "week"},
		}

		return v.VSelect().
			Label("时间单位").
			Items(options).
			ItemTitle("text").
			ItemValue("value").
			Attr(web.VField("Unit", job.Unit)...)
	})

	// 为Actions字段创建操作按钮
	m.modelBuilder.Listing().Field("Actions").ComponentFunc(func(obj interface{}, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		job, ok := obj.(*models.RecurringJob)
		if !ok {
			return nil
		}

		var buttons []h.HTMLComponent

		// 立即执行按钮
		if job.Status != "completed" {
			buttons = append(buttons, v.VBtn("").
				Icon(true).
				Color("primary").
				Size("small").
				Children(
					v.VIcon("mdi-play"),
				).
				Attr("@click", web.Plaid().
					EventFunc("presets_ExecuteJob").
					Query("id", fmt.Sprintf("%d", job.ID)).
					Go()).
				Attr("title", "立即执行").
				Class("mr-2"))
		}

		// 暂停/恢复按钮
		switch job.Status {
		case "active":
			buttons = append(buttons, v.VBtn("").
				Icon(true).
				Color("warning").
				Size("small").
				Children(
					v.VIcon("mdi-pause"),
				).
				Attr("@click", web.Plaid().
					EventFunc("presets_PauseJob").
					Query("id", fmt.Sprintf("%d", job.ID)).
					Go()).
				Attr("title", "暂停").
				Class("mr-2"))
		case "paused":
			buttons = append(buttons, v.VBtn("").
				Icon(true).
				Color("success").
				Size("small").
				Children(
					v.VIcon("mdi-play-pause"),
				).
				Attr("@click", web.Plaid().
					EventFunc("presets_ResumeJob").
					Query("id", fmt.Sprintf("%d", job.ID)).
					Go()).
				Attr("title", "恢复").
				Class("mr-2"))
		}

		// 删除按钮已移除，将由ListingBuilder自动处理

		return h.Td(h.Div(buttons...).Class("d-flex justify-center"))
	})

	// 添加事件处理
	m.modelBuilder.RegisterEventFunc("listing", func(ctx *web.EventContext) (r web.EventResponse, err error) {
		return
	})

	// 注册立即执行事件
	m.modelBuilder.RegisterEventFunc("presets_ExecuteJob", func(ctx *web.EventContext) (r web.EventResponse, err error) {
		id := ctx.R.URL.Query().Get("id")
		if id == "" {
			ctx.Flash = "未找到任务ID"
			r.Reload = true
			return
		}

		var job models.RecurringJob
		if err := m.taskManager.db.First(&job, id).Error; err == nil {
			if err := m.taskManager.RunJobNow(job.Name); err != nil {
				ctx.Flash = err.Error()
			} else {
				ctx.Flash = fmt.Sprintf("任务 %s 已加入执行队列", job.Name)
			}
		} else {
			ctx.Flash = "找不到指定任务"
		}
		r.Reload = true
		return
	})

	// 注册暂停事件
	m.modelBuilder.RegisterEventFunc("presets_PauseJob", func(ctx *web.EventContext) (r web.EventResponse, err error) {
		id := ctx.R.URL.Query().Get("id")
		if id == "" {
			ctx.Flash = "未找到任务ID"
			r.Reload = true
			return
		}

		var job models.RecurringJob
		if err := m.taskManager.db.First(&job, id).Error; err == nil {
			if err := m.taskManager.PauseJob(job.Name); err != nil {
				ctx.Flash = err.Error()
			} else {
				ctx.Flash = fmt.Sprintf("任务 %s 已暂停", job.Name)
			}
		} else {
			ctx.Flash = "找不到指定任务"
		}
		r.Reload = true
		return
	})

	// 注册恢复事件
	m.modelBuilder.RegisterEventFunc("presets_ResumeJob", func(ctx *web.EventContext) (r web.EventResponse, err error) {
		id := ctx.R.URL.Query().Get("id")
		if id == "" {
			ctx.Flash = "未找到任务ID"
			r.Reload = true
			return
		}

		var job models.RecurringJob
		if err := m.taskManager.db.First(&job, id).Error; err == nil {
			if err := m.taskManager.ResumeJob(job.Name); err != nil {
				ctx.Flash = err.Error()
			} else {
				ctx.Flash = fmt.Sprintf("任务 %s 已恢复", job.Name)
			}
		} else {
			ctx.Flash = "找不到指定任务"
		}
		r.Reload = true
		return
	})

	m.modelBuilder.RegisterEventFunc("presets_DoDelete", func(ctx *web.EventContext) (r web.EventResponse, err error) {
		id := ctx.R.URL.Query().Get("id")
		if id == "" {
			ctx.Flash = "未找到任务ID"
			r.Reload = true
			return
		}

		var jobID uint64
		if jobID, err = strconv.ParseUint(id, 10, 32); err != nil {
			ctx.Flash = "任务ID格式错误"
			r.Reload = true
			return
		}

		var fullJob models.RecurringJob
		if err = m.taskManager.db.First(&fullJob, uint(jobID)).Error; err != nil {
			ctx.Flash = "找不到指定任务"
			r.Reload = true
			return
		}

		if err = m.taskManager.RemoveJob(fullJob.Name); err != nil {
			ctx.Flash = "删除任务失败：" + err.Error()
		} else {
			ctx.Flash = "任务已成功删除"
		}

		r.Reload = true
		return
	})

	// 添加新任务处理
	m.modelBuilder.Editing().SaveFunc(func(obj interface{}, id string, ctx *web.EventContext) (err error) {
		job, ok := obj.(*models.RecurringJob)
		if !ok {
			return fmt.Errorf("无效的任务对象")
		}

		// 校验参数
		if job.Name == "" || job.FunctionName == "" || job.Interval <= 0 {
			return fmt.Errorf("名称、函数名和间隔都是必填项")
		}

		// 添加任务
		if id == "" {
			// 创建新任务
			var args interface{}
			if job.Args != "" {
				// 尝试解析JSON
				if json.Valid([]byte(job.Args)) {
					var jsonArgs interface{}
					if err := json.Unmarshal([]byte(job.Args), &jsonArgs); err == nil {
						args = jsonArgs
					} else {
						args = job.Args
					}
				} else {
					args = job.Args
				}
			}

			_, err := m.taskManager.AddJob(job.Name, job.FunctionName, job.Interval, job.Unit, args, job.Times)
			if err != nil {
				return err
			}
		} else {
			// 更新现有任务
			// 首先移除旧任务
			existingJob, err := m.taskManager.GetJob(job.Name)
			if err != nil {
				return err
			}

			err = m.taskManager.RemoveJob(job.Name)
			if err != nil {
				return err
			}

			// 然后添加新任务
			var args interface{}
			if job.Args != "" {
				if json.Valid([]byte(job.Args)) {
					var jsonArgs interface{}
					if err := json.Unmarshal([]byte(job.Args), &jsonArgs); err == nil {
						args = jsonArgs
					} else {
						args = job.Args
					}
				} else {
					args = job.Args
				}
			}

			_, err = m.taskManager.AddJob(job.Name, job.FunctionName, job.Interval, job.Unit, args, job.Times)
			if err != nil {
				// 如果添加新任务失败，尝试恢复旧任务
				if existingJob.Status == "active" {
					m.taskManager.AddJob(
						existingJob.Name,
						existingJob.FunctionName,
						existingJob.Interval,
						existingJob.Unit,
						existingJob.Args,
						existingJob.Times,
					)
				}
				return err
			}
		}

		return nil
	})

	// 添加删除处理
	m.modelBuilder.Editing().DeleteFunc(func(obj interface{}, id string, ctx *web.EventContext) (err error) {
		job, ok := obj.(*models.RecurringJob)
		if !ok {
			return fmt.Errorf("无效的任务对象")
		}

		// 获取完整的任务信息
		var fullJob models.RecurringJob
		if err := m.taskManager.db.First(&fullJob, job.ID).Error; err != nil {
			return fmt.Errorf("找不到指定任务: %v", err)
		}

		return m.taskManager.RemoveJob(fullJob.Name)
	})
}
