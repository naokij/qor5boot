package recurring

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/qor5/admin/v3/presets"
	"github.com/qor5/web/v3"
	v "github.com/qor5/x/v3/ui/vuetify"
	vx "github.com/qor5/x/v3/ui/vuetifyx"
	h "github.com/theplant/htmlgo"
	"gorm.io/gorm"

	"github.com/naokij/qor5boot/models"
	"github.com/robfig/cron/v3"
)

func init() {
	// 初始化随机数种子（Go 1.20+ 已经不需要显式设置种子）
}

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
	modelBuilder.MenuIcon("mdi-clipboard-clock")

	manager := &RecurringJobManager{
		taskManager:  taskManager,
		pb:           pb,
		modelBuilder: modelBuilder,
	}

	// 注册一些示例函数
	manager.registerSampleFunctions()

	// 注册管理界面
	manager.registerAdminUI()

	// 注册执行记录管理界面
	manager.registerExecutionUI()

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

		// 记录开始信息
		execution.Info("任务开始执行")

		// 记录主要信息
		execution.Info("执行日志任务: %s", message)

		// 模拟一些处理步骤
		time.Sleep(100 * time.Millisecond)
		execution.Debug("执行步骤1: 准备数据")

		time.Sleep(200 * time.Millisecond)
		execution.Debug("执行步骤2: 处理数据")

		// 示例警告信息
		if len(message) > 100 {
			execution.Warning("消息内容过长: %d 字符", len(message))
		}

		// 模拟随机错误（10%概率）
		if rand.Intn(10) == 0 {
			execution.LogError("随机错误发生")
			return fmt.Errorf("随机错误")
		}

		// 记录完成信息
		execution.Info("任务执行完成")

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

// 注册执行记录管理界面
func (m *RecurringJobManager) registerExecutionUI() {
	// 创建执行记录模型构建器
	executionBuilder := m.pb.Model(&models.RecurringJobExecution{})
	executionBuilder.Label("RecurringJobLogs")
	executionBuilder.MenuIcon("mdi-history")

	// 配置列表视图
	executionBuilder.Listing("ID", "RecurringJobID", "StartedAt", "Duration", "Success", "Error")

	// 添加过滤功能
	executionBuilder.Listing().FilterDataFunc(func(ctx *web.EventContext) vx.FilterData {
		// 查询所有任务用于过滤器
		var jobs []models.RecurringJob
		if err := m.taskManager.db.Order("name").Find(&jobs).Error; err != nil {
			log.Printf("获取任务列表失败: %v", err)
		}

		// 构建任务选项
		jobOptions := []*vx.SelectItem{}
		for _, job := range jobs {
			jobOptions = append(jobOptions, &vx.SelectItem{
				Text:  job.Name,
				Value: fmt.Sprintf("%d", job.ID),
			})
		}

		return []*vx.FilterItem{
			{
				Key:          "recurring_job_id",
				Label:        "任务名称",
				ItemType:     vx.ItemTypeSelect,
				Options:      jobOptions,
				SQLCondition: `recurring_job_id %s ?`,
			},
			{
				Key:      "success",
				Label:    "执行结果",
				ItemType: vx.ItemTypeSelect,
				Options: []*vx.SelectItem{
					{Text: "成功", Value: "true"},
					{Text: "失败", Value: "false"},
				},
				SQLCondition: `success %s ?`,
			},
		}
	})

	// 添加过滤标签页
	executionBuilder.Listing().FilterTabsFunc(func(ctx *web.EventContext) []*presets.FilterTab {
		// 获取所有执行记录数量（可用于界面显示）
		var allCount int64
		m.taskManager.db.Model(&models.RecurringJobExecution{}).Count(&allCount)

		// 基础标签页
		tabs := []*presets.FilterTab{
			{
				Label: "全部记录",
				ID:    "all",
				Query: url.Values{"all": []string{"1"}},
			},
			{
				Label: "成功记录",
				ID:    "success",
				Query: url.Values{"success": []string{"true"}},
			},
			{
				Label: "失败记录",
				ID:    "error",
				Query: url.Values{"success": []string{"false"}},
			},
		}

		return tabs
	})

	// 格式化持续时间显示
	executionBuilder.Listing().Field("Duration").ComponentFunc(func(obj interface{}, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		execution := obj.(*models.RecurringJobExecution)
		duration := time.Duration(execution.Duration) * time.Millisecond
		if duration < time.Second {
			return h.Td(h.Text(fmt.Sprintf("%dms", execution.Duration)))
		} else if duration < time.Minute {
			return h.Td(h.Text(fmt.Sprintf("%.2fs", float64(execution.Duration)/1000)))
		} else {
			minutes := duration / time.Minute
			seconds := (duration % time.Minute) / time.Second
			return h.Td(h.Text(fmt.Sprintf("%d分%d秒", minutes, seconds)))
		}
	})

	// 格式化成功/失败状态显示
	executionBuilder.Listing().Field("Success").ComponentFunc(func(obj interface{}, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		execution := obj.(*models.RecurringJobExecution)
		var (
			text  string
			color string
		)

		if execution.Success {
			text = "成功"
			color = "success"
		} else {
			text = "失败"
			color = "error"
		}

		return h.Td(v.VChip(h.Text(text)).Color(color))
	})

	// 关联任务名称显示
	executionBuilder.Listing().Field("RecurringJobID").ComponentFunc(func(obj interface{}, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		execution := obj.(*models.RecurringJobExecution)

		// 查询关联的任务名称
		var job models.RecurringJob
		if err := m.taskManager.db.First(&job, execution.RecurringJobID).Error; err != nil {
			return h.Td(h.Text(fmt.Sprintf("#%d", execution.RecurringJobID)))
		}

		// 使用过滤器链接而不是直接链接到任务详情
		filterLink := fmt.Sprintf("/recurring-job-executions?f_recurring_job_id=%d", job.ID)
		return h.Td(h.A(h.Text(job.Name)).Attr("href", filterLink))
	})

	// 配置详情视图
	executionBuilder.Detailing("RecurringJobID", "StartedAt", "FinishedAt", "Duration", "Success", "Error", "Output")

	// 格式化输出字段的显示
	executionBuilder.Detailing().Field("Output").ComponentFunc(func(obj interface{}, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		execution := obj.(*models.RecurringJobExecution)

		if execution.Output == "" {
			return h.Div(h.Text("无输出内容")).Class("grey--text")
		}

		// 添加调试信息
		debugInfo := h.Div().Class("mb-3").Children(
			h.Div().Class("text-caption").Text("输出内容长度: " + strconv.Itoa(len(execution.Output))),
		)

		// 处理多行文本，按日志级别添加颜色
		lines := strings.Split(execution.Output, "\n")
		formattedLines := []h.HTMLComponent{}

		for _, line := range lines {
			// 设置不同日志级别的颜色样式
			var colorClass string
			var colorStyle string
			var prefix string

			if strings.Contains(line, "[INFO]") {
				colorClass = "blue--text text--darken-3"
				colorStyle = "color: #0D47A1 !important;" // 深蓝色
				prefix = "ℹ️ "
			} else if strings.Contains(line, "[WARN]") {
				colorClass = "amber--text text--darken-4"
				colorStyle = "color: #FF6F00 !important;" // 深橙色
				prefix = "⚠️ "
			} else if strings.Contains(line, "[ERROR]") {
				colorClass = "red--text text--darken-4"
				colorStyle = "color: #B71C1C !important;" // 深红色
				prefix = "❌ "
			} else if strings.Contains(line, "[DEBUG]") {
				colorClass = "grey--text text--darken-2"
				colorStyle = "color: #424242 !important;" // 深灰色
				prefix = "🔍 "
			} else {
				colorStyle = "color: #000000;"
			}

			// 创建带颜色的日志行
			logLine := h.Div().
				Text(prefix+line).
				Attr("style", colorStyle).
				Class(colorClass + " log-line py-1")

			formattedLines = append(formattedLines, logLine)
		}

		// 添加调试信息：显示第一行日志的完整内容
		if len(lines) > 0 {
			debugInfo.AppendChildren(
				h.Div().Class("text-caption mt-1").Text("第一行日志: " + lines[0]),
			)
		}

		// 用卡片容器包装所有日志行
		return v.VCard(
			v.VCardTitle(h.Text("执行输出日志")).Class("subtitle-1 py-2"),
			debugInfo,
			v.VDivider(),
			v.VCardText(
				formattedLines...,
			).Class("pa-2"),
		).Elevation(1).Class("log-container overflow-auto").Attr("style", "max-height: 500px; font-family: monospace;")
	})
}

// 注册管理界面
func (m *RecurringJobManager) registerAdminUI() {
	// 配置列表视图
	m.modelBuilder.Listing("ID", "Name", "FunctionName", "CronExpression", "Runs", "Status", "LastRunAt", "NextRunAt", "ErrorCount", "Actions")

	// 添加状态过滤功能
	m.modelBuilder.Listing().FilterDataFunc(func(ctx *web.EventContext) vx.FilterData {
		return []*vx.FilterItem{
			{
				Key:      "status",
				Label:    "状态",
				ItemType: vx.ItemTypeSelect,
				Options: []*vx.SelectItem{
					{Text: "活跃", Value: "active"},
					{Text: "已暂停", Value: "paused"},
					{Text: "已完成", Value: "completed"},
					{Text: "错误", Value: "error"},
				},
				SQLCondition: `status %s ?`,
			},
		}
	})

	// 添加过滤标签页
	m.modelBuilder.Listing().FilterTabsFunc(func(ctx *web.EventContext) []*presets.FilterTab {
		return []*presets.FilterTab{
			{
				Label: "活跃任务",
				ID:    "active",
				Query: url.Values{"status": []string{"active"}},
			},
			{
				Label: "全部任务",
				ID:    "all",
				Query: url.Values{"all": []string{"1"}},
			},
			{
				Label: "已暂停",
				ID:    "paused",
				Query: url.Values{"status": []string{"paused"}},
			},
			{
				Label: "已完成",
				ID:    "completed",
				Query: url.Values{"status": []string{"completed"}},
			},
			{
				Label: "错误任务",
				ID:    "error",
				Query: url.Values{"status": []string{"error"}},
			},
		}
	})

	// 配置编辑视图
	m.modelBuilder.Editing("Name", "FunctionName", "CronExpression", "Times", "Args")

	// 为CronExpression添加组件
	m.modelBuilder.Editing().Field("CronExpression").ComponentFunc(func(obj interface{}, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		job, ok := obj.(*models.RecurringJob)
		if !ok {
			return nil
		}

		// 准备Cron表达式示例
		examples := []struct {
			Expr string
			Desc string
		}{
			{"0 0 * * *", "每天午夜执行"},
			{"0 0 * * 1", "每周一午夜执行"},
			{"0 8 * * 1-5", "每个工作日8点执行"},
			{"0 0,12 * * *", "每天0点和12点执行"},
			{"0 */4 * * *", "每4小时执行一次"},
			{"*/10 * * * *", "每10分钟执行一次"},
			{"0 0 1 * *", "每月1号午夜执行"},
		}

		exampleItems := []h.HTMLComponent{}
		for _, e := range examples {
			exampleItems = append(exampleItems,
				h.Div(
					h.Strong(e.Expr),
					h.Text(" - "+e.Desc),
				).Class("mb-1"),
			)
		}

		return h.Div(
			v.VTextField().
				Label("Cron表达式").
				Hint("例如: 0 0 * * * (每天午夜执行)").
				Attr(web.VField("CronExpression", job.CronExpression)...),
			h.Div(
				h.Div().Text("常用Cron表达式示例:").Class("text-subtitle-2 mt-3"),
				h.Div(exampleItems...),
				h.Div(
					h.A().Text("Cron表达式在线测试工具").
						Href("https://crontab.guru/").
						Target("_blank"),
				).Class("mt-2"),
				h.Div(
					h.Text("Cron表达式格式（标准5字段）："),
					h.Code("分 时 日 月 周"),
				).Class("mt-2"),
			).Class("text-caption mt-2"),
		)
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

	// 使用原地更新方法修改SaveFunc
	m.modelBuilder.Editing().SaveFunc(func(obj interface{}, id string, ctx *web.EventContext) (err error) {
		job, ok := obj.(*models.RecurringJob)
		if !ok {
			return fmt.Errorf("无效的任务对象")
		}

		// 添加详细日志以调试表单数据传递
		log.Printf("==== 表单提交调试信息 ====")
		log.Printf("表单所有数据: %+v", ctx.R.Form)
		log.Printf("job初始状态: ID=%d, cronExpr=%s",
			job.ID, job.CronExpression)

		// 校验参数
		if job.Name == "" || job.FunctionName == "" {
			return fmt.Errorf("名称和函数名是必填项")
		}

		if job.CronExpression == "" {
			return fmt.Errorf("Cron表达式不能为空")
		}

		// 验证Cron表达式格式
		parser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
		if _, err := parser.Parse(job.CronExpression); err != nil {
			return fmt.Errorf("Cron表达式格式无效: %v", err)
		}

		// 添加任务之前最终检查
		log.Printf("任务最终状态: name=%s, cronExpr=%s",
			job.Name, job.CronExpression)

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

			_, err := m.taskManager.AddJob(
				job.Name,
				job.FunctionName,
				args,
				job.Times,
				job.CronExpression,
			)
			if err != nil {
				return err
			}
		} else {
			// 更新现有任务（使用原地更新逻辑）
			var jobID uint64
			if jobID, err = strconv.ParseUint(id, 10, 32); err != nil {
				return fmt.Errorf("任务ID格式错误: %w", err)
			}

			// 准备参数
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

			// 记录日志
			log.Printf("原地更新任务 ID=%s, 新名称=%s", id, job.Name)

			// 调用UpdateJob进行原地更新，保留原有状态和统计信息
			_, err := m.taskManager.UpdateJob(
				uint(jobID),
				job.Name,
				job.FunctionName,
				args,
				job.Times,
				job.CronExpression,
				true, // 总是保留原有状态和统计信息
			)
			if err != nil {
				return fmt.Errorf("更新任务失败: %w", err)
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

// 在RecurringJob详情页添加最近执行记录
func (m *RecurringJobManager) registerExtraUI() {
	// TODO: 在未来版本实现
}
