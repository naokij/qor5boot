package admin

import (
	"context"
	"embed"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"slices"
	"strconv"
	"syscall"
	"time"

	"github.com/iancoleman/strcase"
	"github.com/pkg/errors"
	"github.com/qor5/web/v3"
	"github.com/qor5/x/v3/i18n"
	"github.com/qor5/x/v3/login"
	"github.com/qor5/x/v3/perm"
	v "github.com/qor5/x/v3/ui/vuetify"
	h "github.com/theplant/htmlgo"
	"golang.org/x/text/language"
	"golang.org/x/text/language/display"
	"gorm.io/gorm"

	"github.com/naokij/qor5boot/admin/recurring"
	"github.com/naokij/qor5boot/models"
	"github.com/qor5/admin/v3/activity"
	plogin "github.com/qor5/admin/v3/login"
	"github.com/qor5/admin/v3/presets"
	"github.com/qor5/admin/v3/presets/gorm2op"
	"github.com/qor5/admin/v3/role"
	"github.com/qor5/admin/v3/tiptap"
	"github.com/qor5/admin/v3/utils"
	"github.com/qor5/admin/v3/worker"
)

//go:embed assets
var assets embed.FS

// 全局的重复任务管理器
var recurringJobManager *recurring.RecurringJobManager

// StopRecurringJobManager 停止重复任务管理器
func StopRecurringJobManager() {
	if recurringJobManager != nil {
		log.Println("正在关闭全局重复任务管理器...")
		recurringJobManager.Stop()
		log.Println("全局重复任务管理器已关闭")
	}
}

type Config struct {
	pb                  *presets.Builder
	loginSessionBuilder *plogin.SessionBuilder
	db                  *gorm.DB
}

func (c *Config) GetPresetsBuilder() *presets.Builder {
	return c.pb
}

func (c *Config) GetLoginSessionBuilder() *plogin.SessionBuilder {
	return c.loginSessionBuilder
}

var (
	dbReset                   = getEnvWithDefault("DB_RESET", "")
	resetAndImportInitialData = getEnvWithDefaultBool("RESET_AND_IMPORT_INITIAL_DATA", false)
)

func getEnvWithDefaultBool(key string, defaultValue bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	boolValue, err := strconv.ParseBool(value)
	if err != nil {
		return defaultValue
	}
	return boolValue
}

func NewConfig(db *gorm.DB, enableWork bool) Config {
	// 初始化LDAP配置
	initLDAP()

	if err := db.AutoMigrate(
		&models.User{},
		&role.Role{},
		&perm.DefaultDBPolicy{},
	); err != nil {
		panic(err)
	}

	// @snippet_begin(ActivityExample)
	ab := activity.New(db, func(ctx context.Context) (*activity.User, error) {
		u := ctx.Value(login.UserKey).(*models.User)
		return &activity.User{
			ID:     fmt.Sprint(u.ID),
			Name:   u.Name,
			Avatar: "",
		}, nil
	}).
		WrapLogModelInstall(func(in presets.ModelInstallFunc) presets.ModelInstallFunc {
			return func(pb *presets.Builder, mb *presets.ModelBuilder) (err error) {
				err = in(pb, mb)
				if err != nil {
					return
				}
				mb.Listing().WrapSearchFunc(func(in presets.SearchFunc) presets.SearchFunc {
					return func(ctx *web.EventContext, params *presets.SearchParams) (result *presets.SearchResult, err error) {
						u := getCurrentUser(ctx.R)
						if rs := u.GetRoles(); !slices.Contains(rs, models.RoleAdmin) {
							params.SQLConditions = append(params.SQLConditions, &presets.SQLCondition{
								Query: "user_id = ?",
								Args:  []interface{}{fmt.Sprint(u.ID)},
							})
						}
						return in(ctx, params)
					}
				})
				return
			}
		}).
		TablePrefix("cms_").
		AutoMigrate()

	// ab.Model(l).SkipDelete().SkipCreate()
	// @snippet_end

	b := presets.New().DataOperator(gorm2op.DataOperator(db)).RightDrawerWidth("700")
	defer b.Build()

	b.ExtraAsset("/tiptap.css", "text/css", tiptap.ThemeGithubCSSComponentsPack())

	initPermission(b, db)

	b.GetI18n().
		SupportLanguages(language.SimplifiedChinese, language.English).
		RegisterForModule(language.SimplifiedChinese, presets.ModelsI18nModuleKey, Messages_zh_CN_ModelsI18nModuleKey).
		RegisterForModule(language.SimplifiedChinese, I18nAdminKey, Messages_zh_CN).
		RegisterForModule(language.English, I18nAdminKey, Messages_en_US).
		GetSupportLanguagesFromRequestFunc(func(r *http.Request) []language.Tag {
			// // Example:
			// user := getCurrentUser(r)
			// var supportedLanguages []language.Tag
			// for _, role := range user.GetRoles() {
			//	switch role {
			//	case "English Group":
			//		supportedLanguages = append(supportedLanguages, language.English)
			//	case "Chinese Group":
			//		supportedLanguages = append(supportedLanguages, language.SimplifiedChinese)
			//	}
			// }
			// return supportedLanguages
			return b.GetI18n().GetSupportLanguages()
		})

	// 添加Dashboard模型，放在最前面
	type Dashboard struct{}
	dashboardB := b.Model(&Dashboard{}).Label("Dashboard")
	dashboardB.Listing().PageFunc(func(ctx *web.EventContext) (r web.PageResponse, err error) {
		msgr := i18n.MustGetModuleMessages(ctx.R, I18nAdminKey, Messages_zh_CN).(*Messages)
		r.PageTitle = msgr.Dashboard
		r.Body = DashboardBody()
		return
	})
	// 添加图标
	dashboardB.MenuIcon("mdi-view-dashboard")

	utils.Install(b)

	configMenuOrder(b)

	roleBuilder := role.New(db).
		Resources([]*v.DefaultOptionItem{
			{Text: "All", Value: "*"},
			{Text: "Settings", Value: "*:settings:*"},
			{Text: "Customers", Value: "*:customers:*"},
			{Text: "ActivityLogs", Value: "*:activity_logs:*"},
			{Text: "Workers", Value: "*:workers:*"},
		}).
		AfterInstall(func(pb *presets.Builder, mb *presets.ModelBuilder) error {
			mb.Listing().SearchFunc(func(ctx *web.EventContext, params *presets.SearchParams) (result *presets.SearchResult, err error) {
				u := getCurrentUser(ctx.R)
				qdb := db
				// If the current user doesn't has 'admin' role, do not allow them to view admin and manager roles
				// We didn't do this on permission because of we are not supporting the permission on listing page
				if currentRoles := u.GetRoles(); !slices.Contains(currentRoles, models.RoleAdmin) {
					qdb = db.Where("name NOT IN (?)", []string{models.RoleAdmin, models.RoleManager})
				}
				return gorm2op.DataOperator(qdb).Search(ctx, params)
			})
			return nil
		})
	if enableWork {
		w := worker.New(db)
		defer w.Listen()
		addJobs(w)
		b.Use(w.Activity(ab))

		// 添加重复任务支持
		recurringJobManager = recurring.NewRecurringJobManager(db, b)
		if err := recurringJobManager.Start(); err != nil {
			log.Printf("启动重复任务管理器失败: %v", err)
		}

		// 程序退出时关闭管理器，但不立即关闭
		// 这个方法不会阻塞，只会在程序收到信号时执行
		go func() {
			// 创建一个信号通道
			c := make(chan os.Signal, 1)
			signal.Notify(c, os.Interrupt, syscall.SIGTERM)

			// 等待信号
			<-c
			log.Println("正在停止重复任务管理器...")
			recurringJobManager.Stop()
		}()
	}

	loginSessionBuilder := initLoginSessionBuilder(db, b, ab)

	configBrand(b)

	profileBuilder := configProfile(db, ab, loginSessionBuilder)

	configUser(b, ab, db, loginSessionBuilder)
	b.Use(
		ab,
		roleBuilder,
		loginSessionBuilder,
		profileBuilder,
	)

	if resetAndImportInitialData {
		tbs := GetNonIgnoredTableNames(db)
		EmptyDB(db, tbs)
		InitDB(db, tbs)
	}

	return Config{
		pb:                  b,
		loginSessionBuilder: loginSessionBuilder,
		db:                  db,
	}
}

func configMenuOrder(b *presets.Builder) {
	// 按照期望顺序配置菜单项
	b.MenuOrder(
		"profile",
		"Dashboard", // 将Dashboard放在第二位置
		b.MenuGroup("User Management").SubItems(
			"User",
			"Role",
		).Icon("mdi-account-multiple"),
		b.MenuGroup("TaskManagement").SubItems(
			"Worker",
			"RecurringJobs",
			"RecurringJobExecutions",
		).Icon("mdi-clock-outline"),
		"ActivityLogs",
	)
}

func configBrand(b *presets.Builder) {
	b.SwitchLanguageFunc(func(ctx *web.EventContext) h.HTMLComponent {
		supportLanguages := b.GetI18n().GetSupportLanguagesFromRequest(ctx.R)

		if len(b.GetI18n().GetSupportLanguages()) <= 1 || len(supportLanguages) == 0 {
			return nil
		}

		queryName := b.GetI18n().GetQueryName()
		msgr := presets.MustGetMessages(ctx.R)

		if len(supportLanguages) == 1 {
			return h.Template().Children(
				h.Div(
					v.VList(
						v.VListItem(
							web.Slot(
								v.VIcon("mdi-translate").Size(v.SizeSmall).Class("mr-4 ml-1"),
							).Name("prepend"),
							v.VListItemTitle(
								h.Div(h.Text(fmt.Sprintf("%s%s %s", msgr.Language, msgr.Colon, display.Self.Name(supportLanguages[0])))).Role("button"),
							),
						).Class("pa-0").Density(v.DensityCompact),
					).Class("pa-0 ma-n4 mt-n6"),
				).Attr("@click", web.Plaid().MergeQuery(true).Query(queryName, supportLanguages[0].String()).Go()),
			)
		}

		// 使用语言简写作为显示文本，使用translate图标
		currentLanguage := "EN"
		lang := ctx.R.FormValue(queryName)
		if lang == "" {
			lang = b.GetI18n().GetCurrentLangFromCookie(ctx.R)
		}
		switch lang {
		case language.SimplifiedChinese.String():
			currentLanguage = "中文"
		case language.English.String():
			currentLanguage = "EN"
		case language.Japanese.String():
			currentLanguage = "JP"
		}

		var languages []h.HTMLComponent
		for _, tag := range supportLanguages {
			var langText string
			switch tag.String() {
			case language.SimplifiedChinese.String():
				langText = "中文"
			case language.English.String():
				langText = "English"
			case language.Japanese.String():
				langText = "日本語"
			default:
				langText = display.Self.Name(tag)
			}

			languages = append(languages,
				h.Div(
					v.VListItem(
						v.VListItemTitle(
							h.Div(h.Text(langText)),
						),
					).Attr("@click", web.Plaid().MergeQuery(true).Query(queryName, tag.String()).Go()),
				),
			)
		}

		return v.VMenu().Children(
			h.Template().Attr("v-slot:activator", "{isActive, props}").Children(
				h.Div(
					v.VBtn("").Children(
						v.VIcon("mdi-translate"),
						h.Span(currentLanguage).Class("ml-2"),
						v.VIcon("mdi-menu-down").Class("ml-1"),
					).Attr("variant", "text").Class("i18n-switcher-btn"),
				).Attr("v-bind", "props").Style("display: inline-block;"),
			),
			v.VList(
				languages...,
			).Density(v.DensityCompact),
		)
	})

	b.BrandFunc(func(ctx *web.EventContext) h.HTMLComponent {
		msgr := i18n.MustGetModuleMessages(ctx.R, I18nAdminKey, Messages_zh_CN).(*Messages)
		logo := "/admin/assets/logo.svg"

		now := time.Now()
		nextEvenHour := time.Date(now.Year(), now.Month(), now.Day(), now.Hour()+1+(now.Hour()%2), 0, 0, 0, now.Location())
		diff := int(nextEvenHour.Sub(now).Seconds())
		hours := diff / 3600
		minutes := (diff % 3600) / 60
		seconds := diff % 60
		countdown := fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)

		return h.Div(
			v.VRow(
				v.VCol(h.A(h.Img(logo).Attr("width", "80")).Href("/")),
				v.VCol(h.H1(msgr.Demo)).Class("pt-4"),
			),
			// ).Density(DensityCompact),
			h.If(dbReset != "",
				h.Div(
					h.Span(msgr.DBResetTipLabel),
					v.VIcon("schedule").Size(v.SizeXSmall),
					// .Left(true),
					h.Span(countdown).Id("countdown"),
				).Class("pt-1 pb-2"),
				v.VDivider(),
				h.Script("function updateCountdown(){const now=new Date();const nextEvenHour=new Date(now);nextEvenHour.setHours(nextEvenHour.getHours()+(nextEvenHour.getHours()%2===0?2:1),0,0,0);const timeLeft=nextEvenHour-now;const hours=Math.floor(timeLeft/(60*60*1000));const minutes=Math.floor((timeLeft%(60*60*1000))/(60*1000));const seconds=Math.floor((timeLeft%(60*1000))/1000);const countdownElem=document.getElementById(\"countdown\");countdownElem.innerText=`${hours.toString().padStart(2,\"0\")}:${minutes.toString().padStart(2,\"0\")}:${seconds.toString().padStart(2,\"0\")}`}updateCountdown();setInterval(updateCountdown,1000);"),
			),
		).Class("mb-n4 mt-n2")
	}).HomePageFunc(func(ctx *web.EventContext) (r web.PageResponse, err error) {
		msgr := i18n.MustGetModuleMessages(ctx.R, I18nAdminKey, Messages_zh_CN).(*Messages)
		r.PageTitle = msgr.Dashboard
		r.Body = DashboardBody()
		return
	}).NotFoundPageLayoutConfig(&presets.LayoutConfig{
		NotificationCenterInvisible: true,
	})
}

func configProfile(db *gorm.DB, ab *activity.Builder, lsb *plogin.SessionBuilder) *plogin.ProfileBuilder {
	return plogin.NewProfileBuilder(
		func(ctx context.Context) (*plogin.Profile, error) {
			evCtx := web.MustGetEventContext(ctx)
			u := getCurrentUser(evCtx.R)
			if u == nil {
				return nil, perm.PermissionDenied
			}
			notifiCounts, err := ab.GetNotesCounts(ctx, "", nil)
			if err != nil {
				return nil, err
			}
			user := &plogin.Profile{
				ID:   fmt.Sprint(u.ID),
				Name: u.Name,
				// Avatar: "",
				Roles:  u.GetRoles(),
				Status: strcase.ToCamel(u.Status),
				Fields: []*plogin.ProfileField{
					{Name: "Email", Value: u.Account},
					{Name: "Company", Value: u.Company},
				},
				NotifCounts: notifiCounts,
			}
			if u.OAuthAvatar != "" {
				user.Avatar = u.OAuthAvatar
			}
			return user, nil
		},
		func(ctx context.Context, newName string) error {
			evCtx := web.MustGetEventContext(ctx)
			u := getCurrentUser(evCtx.R)
			if u == nil {
				return perm.PermissionDenied
			}
			u.Name = newName
			if err := db.Save(u).Error; err != nil {
				return errors.Wrap(err, "failed to update user name")
			}
			return nil
		},
	).SessionBuilder(lsb).CustomizeButtons(func(ctx context.Context, buttons ...h.HTMLComponent) ([]h.HTMLComponent, error) {
		// 添加修改密码按钮
		msgr := i18n.MustGetModuleMessages(web.MustGetEventContext(ctx).R, I18nAdminKey, Messages_zh_CN).(*Messages)

		changePasswordBtn := v.VBtn(msgr.ChangePassword).
			Variant(v.VariantTonal).
			Color(v.ColorPrimary).
			OnClick(plogin.OpenChangePasswordDialogEvent)

		// 将修改密码按钮插入到原有按钮之前
		newButtons := append([]h.HTMLComponent{changePasswordBtn}, buttons...)

		return newButtons, nil
	})
}
