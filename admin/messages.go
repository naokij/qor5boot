package admin

import (
	"github.com/qor5/x/v3/i18n"
)

const I18nAdminKey i18n.ModuleKey = "I18nAdminKey"

type Messages struct {
	// Sidebar
	SidebarTitle string

	// Common
	IndexPage         string
	SearchPlaceholder string
	PleaseSelect      string
	Cancel            string
	OK                string
	Add               string
	Submit            string
	Edit              string
	EditAlt           string
	Show              string
	ExtraExport       string

	// 通用过滤器标签
	FilterTabsAll            string
	FilterTabsActive         string
	FilterTabsHasUnreadNotes string

	// Demo相关
	Demo              string
	DemoTips          string
	DemoUsernameLabel string
	DemoPasswordLabel string
	DBResetTipLabel   string

	// 登录相关
	LoginTitleLabel                string
	LoginProviderGoogleText        string
	LoginProviderMicrosoftText     string
	LoginProviderGithubText        string
	OAuthCompleteInfoTitle         string
	OAuthCompleteInfoPositionLabel string
	OAuthCompleteInfoAgreeLabel    string
	OAuthCompleteInfoBackLabel     string
	PasswordMinLengthHint          string
	PasswordMinLengthError         string
	PasswordPlaceholder            string

	// 用户信息
	Name           string
	Email          string
	Company        string
	Status         string
	ChangePassword string
	LoginSessions  string

	// 菜单项
	Dashboard      string
	TaskManagement string

	// Activity Admin
	Users                       string
	Resource                    string
	Account                     string
	Session                     string
	ActivityDate                string
	Action                      string
	RevokeToken                 string
	RevokeAllTheUsersSession    string
	ConfirmRevokeUsersSession   string
	RevokeAllTokens             string
	AllSessionsOfUserRevoked    string
	SuccessfullyRevokedSession  string
	FailedToRevokeSession       string
	FailedToRevokeSessionNoUser string
	ConfirmRevokeAllSession     string
	SuccessfullyRevokedAllToken string
	FailedToRevokeAllToken      string
	RevocationFailed            string
	UsernameLoggedInTime        string
	LoginSuccess                string
	LoginSuccessWithAuthID      string
	LoginFailed                 string
	LoginFailedWithAuthID       string
	LoginFailedWithMessage      string
	SessionExtendSuccess        string
	SessionExtendFailed         string
	LastLogin                   string
	RecentActivities            string
	UserActivities              string

	// User Admin
	Role            string
	Roles           string
	RoleName        string
	EnName          string
	EnPermissions   string
	AkPermissions   string
	EmailSignedName string
	DisplayName     string
	SiteName        string
	UserName        string
	UserNameEmail   string
	Descriptions    string
	CreatePassword  string
	UserPhotoLabel  string
	RoleResources   string
	SignupAt        string
	ConfirmedAt     string
	LastPasswordAt  string
	LastSignin      string
	LastFailAt      string
	LastLocation    string
	ResetPassword   string
	ConfirmPassword string
	FailedAttempt   string
	OnLeaveFrom     string
	OnLeaveTo       string
	UserStatus      string
	LeaveRequested  string
	OnLeave         string
	Blocked         string
	ProfilePage     string
	Permissions     string
	UserStatus1     string
	UserStatus2     string
	UserStatus3     string
	UserStatus4     string
	UserStatus5     string

	// SEO Admin
	SEOSettings string
	SEO         string
	Default     string
	TagsFor     string
	Variables   string

	// Product Admin
	Products string

	// Tables Tab
	TabAllItems   string
	TabActiveOnly string
	Filter        string
	FilterBtn     string
	ClearFilter   string
}

var Messages_en_US = &Messages{
	// Sidebar
	SidebarTitle: "QOR5Boot",

	// Common
	IndexPage:         "Index",
	SearchPlaceholder: "Search and press enter",
	PleaseSelect:      "Please Select",
	Cancel:            "Cancel",
	OK:                "OK",
	Add:               "Add",
	Submit:            "Submit",
	Edit:              "Edit",
	EditAlt:           "Edit",
	Show:              "Show",
	ExtraExport:       "Export",

	// 通用过滤器标签
	FilterTabsAll:            "All",
	FilterTabsActive:         "Active",
	FilterTabsHasUnreadNotes: "Has Unread Notes",

	// Demo相关
	Demo:              "Demo",
	DemoTips:          "Demo Tips",
	DemoUsernameLabel: "Username",
	DemoPasswordLabel: "Password",
	DBResetTipLabel:   "DB Reset Tip",

	// 登录相关
	LoginTitleLabel:                "Qor5boot Management System",
	LoginProviderGoogleText:        "Login with Google",
	LoginProviderMicrosoftText:     "Login with Microsoft",
	LoginProviderGithubText:        "Login with GitHub",
	OAuthCompleteInfoTitle:         "OAuth Complete Info",
	OAuthCompleteInfoPositionLabel: "Position",
	OAuthCompleteInfoAgreeLabel:    "Agree",
	OAuthCompleteInfoBackLabel:     "Back",
	PasswordMinLengthHint:          "Password Min Length Hint",
	PasswordMinLengthError:         "Password Min Length Error",
	PasswordPlaceholder:            "Password Placeholder",

	// 用户信息
	Name:           "Name",
	Email:          "Email",
	Company:        "Company",
	Status:         "Status",
	ChangePassword: "Change Password",
	LoginSessions:  "Login Sessions",

	// 菜单项
	Dashboard:      "Dashboard",
	TaskManagement: "Task Management",

	// Activity Admin
	Users:                       "Users",
	Resource:                    "Resource",
	Account:                     "Account",
	Session:                     "Session",
	ActivityDate:                "Activity Date",
	Action:                      "Action",
	RevokeToken:                 "Revoke Token",
	RevokeAllTheUsersSession:    "Revoke All The User's Session",
	ConfirmRevokeUsersSession:   "Are you sure to Revoke the User's Session?",
	RevokeAllTokens:             "Revoke All Tokens",
	AllSessionsOfUserRevoked:    "All sessions of user %s revoked",
	SuccessfullyRevokedSession:  "Successfully Revoked Session",
	FailedToRevokeSession:       "Failed to Revoke Session",
	FailedToRevokeSessionNoUser: "Failed to Revoke Session: No User Selected",
	ConfirmRevokeAllSession:     "Are you sure to Revoke All The Session?",
	SuccessfullyRevokedAllToken: "Successfully Revoked All Token",
	FailedToRevokeAllToken:      "Failed to Revoke All Token",
	RevocationFailed:            "Revocation failed",
	UsernameLoggedInTime:        "%s logged in %s",
	LoginSuccess:                "Login Success",
	LoginSuccessWithAuthID:      "Login Success: %s",
	LoginFailed:                 "Login Failed",
	LoginFailedWithAuthID:       "Login Failed: %s",
	LoginFailedWithMessage:      "Login Failed: %s",
	SessionExtendSuccess:        "Session Extended Success",
	SessionExtendFailed:         "Session Extended Failed",
	LastLogin:                   "Last Login",
	RecentActivities:            "Recent Activities",
	UserActivities:              "User Activities",

	// User Admin
	Role:            "Role",
	Roles:           "Roles",
	RoleName:        "Role Name",
	EnName:          "English Name",
	EnPermissions:   "English Permissions",
	AkPermissions:   "AK Permissions",
	EmailSignedName: "Email Signed Name",
	DisplayName:     "Display Name",
	SiteName:        "Site Name",
	UserName:        "User Name",
	UserNameEmail:   "Email",
	Descriptions:    "Descriptions",
	CreatePassword:  "Create Password",
	UserPhotoLabel:  "Photo",
	RoleResources:   "Role Resources",
	SignupAt:        "Signup At",
	ConfirmedAt:     "Confirm At",
	LastPasswordAt:  "Last Password Change At",
	LastSignin:      "Last Login At",
	LastFailAt:      "Last Failed At",
	LastLocation:    "Last Location",
	ResetPassword:   "Reset Password",
	ConfirmPassword: "Confirm Password",
	FailedAttempt:   "Failed attempt",
	OnLeaveFrom:     "On leave from",
	OnLeaveTo:       "On leave to",
	UserStatus:      "Status",
	LeaveRequested:  "Leave Requested",
	OnLeave:         "On Leave",
	Blocked:         "Blocked",
	ProfilePage:     "Profile",
	Permissions:     "Permissions",
	UserStatus1:     "Active",
	UserStatus2:     "Leave Requested",
	UserStatus3:     "On Leave",
	UserStatus4:     "Suspended",
	UserStatus5:     "Blocked",

	// SEO Admin
	SEOSettings: "SEO Settings",
	SEO:         "SEO",
	Default:     "Default",
	TagsFor:     "Tags For",
	Variables:   "Variables",

	// Product Admin
	Products: "Products",

	// Tables Tab
	TabAllItems:   "All Items",
	TabActiveOnly: "Active Only",
	Filter:        "Filter",
	FilterBtn:     "Filter",
	ClearFilter:   "Clear Filter",
}

var Messages_zh_CN = &Messages{
	// Sidebar
	SidebarTitle: "QOR5Boot",

	// Common
	IndexPage:         "首页",
	SearchPlaceholder: "搜索并回车",
	PleaseSelect:      "请选择",
	Cancel:            "取消",
	OK:                "确定",
	Add:               "添加",
	Submit:            "保存",
	Edit:              "编辑",
	EditAlt:           "编辑",
	Show:              "查看",
	ExtraExport:       "导出",

	// 通用过滤器标签
	FilterTabsAll:            "全部",
	FilterTabsActive:         "活跃",
	FilterTabsHasUnreadNotes: "有未读笔记",

	// Demo相关
	Demo:              "示例",
	DemoTips:          "示例提示",
	DemoUsernameLabel: "用户名",
	DemoPasswordLabel: "密码",
	DBResetTipLabel:   "DB重置提示",

	// 登录相关
	LoginTitleLabel:                "Qor5boot 管理系统",
	LoginProviderGoogleText:        "使用Google登录",
	LoginProviderMicrosoftText:     "使用Microsoft登录",
	LoginProviderGithubText:        "使用GitHub登录",
	OAuthCompleteInfoTitle:         "OAuth完成信息",
	OAuthCompleteInfoPositionLabel: "位置",
	OAuthCompleteInfoAgreeLabel:    "同意",
	OAuthCompleteInfoBackLabel:     "返回",
	PasswordMinLengthHint:          "密码最小长度提示",
	PasswordMinLengthError:         "密码最小长度错误",
	PasswordPlaceholder:            "密码占位符",

	// 用户信息
	Name:           "名称",
	Email:          "电子邮件",
	Company:        "公司",
	Status:         "状态",
	ChangePassword: "更改密码",
	LoginSessions:  "登录会话",

	// 菜单项
	Dashboard:      "仪表盘",
	TaskManagement: "任务管理",

	// Activity Admin
	Users:                       "用户",
	Resource:                    "资源",
	Account:                     "账户",
	Session:                     "会话",
	ActivityDate:                "活动日期",
	Action:                      "操作",
	RevokeToken:                 "撤销令牌",
	RevokeAllTheUsersSession:    "撤销所有用户的会话",
	ConfirmRevokeUsersSession:   "确定要撤销用户的会话吗？",
	RevokeAllTokens:             "撤销所有令牌",
	AllSessionsOfUserRevoked:    "用户 %s 的所有会话已被撤销",
	SuccessfullyRevokedSession:  "成功撤销会话",
	FailedToRevokeSession:       "撤销会话失败",
	FailedToRevokeSessionNoUser: "撤销会话失败：未选择用户",
	ConfirmRevokeAllSession:     "确定要撤销所有会话吗？",
	SuccessfullyRevokedAllToken: "成功撤销所有令牌",
	FailedToRevokeAllToken:      "撤销所有令牌失败",
	RevocationFailed:            "撤销失败",
	UsernameLoggedInTime:        "%s 在 %s 登录",
	LoginSuccess:                "登录成功",
	LoginSuccessWithAuthID:      "登录成功：%s",
	LoginFailed:                 "登录失败",
	LoginFailedWithAuthID:       "登录失败：%s",
	LoginFailedWithMessage:      "登录失败：%s",
	SessionExtendSuccess:        "会话延长成功",
	SessionExtendFailed:         "会话延长失败",
	LastLogin:                   "上次登录",
	RecentActivities:            "最近活动",
	UserActivities:              "用户活动",

	// User Admin
	Role:            "角色",
	Roles:           "角色",
	RoleName:        "角色名称",
	EnName:          "英文名称",
	EnPermissions:   "英文权限",
	AkPermissions:   "AK权限",
	EmailSignedName: "电子签名名称",
	DisplayName:     "显示名称",
	SiteName:        "站点名称",
	UserName:        "用户名",
	UserNameEmail:   "电子邮件",
	Descriptions:    "描述",
	CreatePassword:  "创建密码",
	UserPhotoLabel:  "照片",
	RoleResources:   "角色资源",
	SignupAt:        "注册时间",
	ConfirmedAt:     "确认时间",
	LastPasswordAt:  "上次密码更改时间",
	LastSignin:      "上次登录时间",
	LastFailAt:      "上次失败时间",
	LastLocation:    "上次位置",
	ResetPassword:   "重置密码",
	ConfirmPassword: "确认密码",
	FailedAttempt:   "失败尝试",
	OnLeaveFrom:     "离开时间",
	OnLeaveTo:       "离开时间",
	UserStatus:      "状态",
	LeaveRequested:  "请求离开",
	OnLeave:         "离开",
	Blocked:         "已阻止",
	ProfilePage:     "个人资料",
	Permissions:     "权限",
	UserStatus1:     "活跃",
	UserStatus2:     "请求离开",
	UserStatus3:     "离开",
	UserStatus4:     "暂停",
	UserStatus5:     "已阻止",

	// SEO Admin
	SEOSettings: "SEO设置",
	SEO:         "SEO",
	Default:     "默认",
	TagsFor:     "为",
	Variables:   "变量",

	// Product Admin
	Products: "产品",

	// Tables Tab
	TabAllItems:   "所有项目",
	TabActiveOnly: "仅活跃",
	Filter:        "过滤",
	FilterBtn:     "过滤",
	ClearFilter:   "清除过滤",
}

type Messages_ModelsI18nModuleKey struct {
	QOR5Example string
	Roles       string
	Users       string
	Dashboard   string

	Posts          string
	PostsID        string
	PostsTitle     string
	PostsHeroImage string
	PostsBody      string
	Example        string
	Settings       string
	Post           string
	PostsBodyImage string

	SeoPost             string
	SeoVariableTitle    string
	SeoVariableSiteName string

	PageBuilder              string
	Pages                    string
	SharedContainers         string
	DemoContainers           string
	Templates                string
	PageCategories           string
	ECManagement             string
	ECDashboard              string
	Orders                   string
	InputDemos               string
	Products                 string
	NestedFieldDemos         string
	SiteManagement           string
	SEO                      string
	UserManagement           string
	Profile                  string
	FeaturedModelsManagement string
	Customers                string
	ListModels               string
	MicrositeModels          string
	Workers                  string
	RecurringJobs            string
	RecurringJobLogs         string
	TaskManagement           string

	// 重复任务相关字段
	RecurringJobsName           string
	RecurringJobsFunctionName   string
	RecurringJobsCronExpression string
	RecurringJobsTimes          string
	RecurringJobsArgs           string
	RecurringJobsStatus         string
	RecurringJobsLastRunAt      string
	RecurringJobsNextRunAt      string
	RecurringJobsErrorCount     string
	RecurringJobsActions        string
	RecurringJobsRuns           string

	// 重复任务状态值
	RecurringJobsStatusActive    string
	RecurringJobsStatusPaused    string
	RecurringJobsStatusCompleted string
	RecurringJobsStatusError     string

	// 重复任务编辑表单
	RecurringJobsEditFunctionName   string
	RecurringJobsEditCronExpression string
	RecurringJobsEditTimes          string
	RecurringJobsEditArgs           string

	// 重复任务过滤标签
	RecurringJobsTabAll       string
	RecurringJobsTabActive    string
	RecurringJobsTabPaused    string
	RecurringJobsTabCompleted string
	RecurringJobsTabError     string

	// 重复任务日志相关
	RecurringJobLogsID         string
	RecurringJobLogsJobID      string
	RecurringJobLogsStartedAt  string
	RecurringJobLogsFinishedAt string
	RecurringJobLogsDuration   string
	RecurringJobLogsSuccess    string
	RecurringJobLogsError      string
	RecurringJobLogsOutput     string

	// 重复任务日志过滤标签
	RecurringJobLogsTabAll     string
	RecurringJobLogsTabSuccess string
	RecurringJobLogsTabFailed  string

	// 操作按钮
	RecurringJobsPause  string
	RecurringJobsResume string
	RecurringJobsRun    string
	RecurringJobsDelete string

	PagesID         string
	PagesTitle      string
	PagesSlug       string
	PagesLocale     string
	PagesNotes      string
	PagesDraftCount string
	PagesPath       string
	PagesOnline     string
	PagesVersion    string
	PagesVersions   string
	PagesStartAt    string
	PagesEndAt      string
	PagesOption     string
	PagesLive       string

	Page                   string
	PagesStatus            string
	PagesSchedule          string
	PagesCategoryID        string
	PagesTemplateSelection string
	PagesEditContainer     string

	WebHeader       string
	WebHeadersColor string
	Header          string
	Navigation      string
	Content         string

	WebFooter            string
	WebFootersEnglishUrl string
	Footer               string

	VideoBanner                       string
	VideoBannersAddTopSpace           string
	VideoBannersAddBottomSpace        string
	VideoBannersAnchorID              string
	VideoBannersVideo                 string
	VideoBannersBackgroundVideo       string
	VideoBannersMobileBackgroundVideo string
	VideoBannersVideoCover            string
	VideoBannersMobileVideoCover      string
	VideoBannersHeading               string
	VideoBannersPopupText             string
	VideoBannersText                  string
	VideoBannersLinkText              string
	VideoBannersLink                  string

	Heading                   string
	HeadingsAddTopSpace       string
	HeadingsAddBottomSpace    string
	HeadingsAnchorID          string
	HeadingsHeading           string
	HeadingsFontColor         string
	HeadingsBackgroundColor   string
	HeadingsLink              string
	HeadingsLinkText          string
	HeadingsLinkDisplayOption string
	HeadingsText              string

	BrandGrid                string
	BrandGridsAddTopSpace    string
	BrandGridsAddBottomSpace string
	BrandGridsAnchorID       string
	BrandGridsBrands         string

	ListContent                   string
	ListContentsAddTopSpace       string
	ListContentsAddBottomSpace    string
	ListContentsAnchorID          string
	ListContentsBackgroundColor   string
	ListContentsItems             string
	ListContentsLink              string
	ListContentsLinkText          string
	ListContentsLinkDisplayOption string

	ImageContainer                           string
	ImageContainersAddTopSpace               string
	ImageContainersAddBottomSpace            string
	ImageContainersAnchorID                  string
	ImageContainersBackgroundColor           string
	ImageContainersTransitionBackgroundColor string
	ImageContainersImage                     string
	Image                                    string

	InNumber                string
	InNumbersAddTopSpace    string
	InNumbersAddBottomSpace string
	InNumbersAnchorID       string
	InNumbersHeading        string
	InNumbersItems          string
	InNumbers               string

	ContactForm                    string
	ContactFormsAddTopSpace        string
	ContactFormsAddBottomSpace     string
	ContactFormsAnchorID           string
	ContactFormsHeading            string
	ContactFormsText               string
	ContactFormsSendButtonText     string
	ContactFormsFormButtonText     string
	ContactFormsMessagePlaceholder string
	ContactFormsNamePlaceholder    string
	ContactFormsEmailPlaceholder   string
	ContactFormsThankyouMessage    string
	ContactFormsActionUrl          string
	ContactFormsPrivacyPolicy      string

	ActivityActionLogIn         string
	ActivityActionExtendSession string

	PagesPage string

	User                  string
	Role                  string
	LoginSession          string
	Dictionary            string
	RecurringJob          string
	RecurringJobExecution string
	Worker                string
	WorkerJob             string
}

var Messages_en_US_ModelsI18nModuleKey = &Messages_ModelsI18nModuleKey{
	Posts:               "Posts",
	PostsID:             "ID",
	PostsTitle:          "Title",
	PostsHeroImage:      "Hero Image",
	PostsBody:           "Body",
	Example:             "QOR5 Demo",
	Settings:            "SEO Settings",
	Post:                "Post",
	PostsBodyImage:      "Body Image",
	SeoPost:             "Post",
	SeoVariableTitle:    "Title",
	SeoVariableSiteName: "Site Name",

	QOR5Example: "QOR5 Example",
	Roles:       "Roles",
	Users:       "Users",
	Dashboard:   "Dashboard",

	PageBuilder:              "Page Builder Menu",
	Pages:                    "Pages",
	SharedContainers:         "Shared Containers",
	DemoContainers:           "Demo Containers",
	Templates:                "Templates",
	PageCategories:           "Page Categories",
	ECManagement:             "E-Commerce Management",
	ECDashboard:              "E-Commerce Dashboard",
	Orders:                   "Orders",
	InputDemos:               "Input Demos",
	Products:                 "Products",
	NestedFieldDemos:         "Nested Field Demos",
	SiteManagement:           "Site Management",
	SEO:                      "SEO",
	Profile:                  "Profile",
	UserManagement:           "User Management",
	FeaturedModelsManagement: "Featured Models Management",
	Customers:                "Customers",
	ListModels:               "List Models",
	MicrositeModels:          "Microsite Models",
	Workers:                  "Workers",
	RecurringJobs:            "Recurring Tasks",
	RecurringJobLogs:         "Recurring Task Logs",
	TaskManagement:           "Task Management",

	// 重复任务相关字段
	RecurringJobsName:           "Task Name",
	RecurringJobsFunctionName:   "Function Name",
	RecurringJobsCronExpression: "Cron Expression",
	RecurringJobsTimes:          "Run Limit",
	RecurringJobsArgs:           "Arguments",
	RecurringJobsStatus:         "Status",
	RecurringJobsLastRunAt:      "Last Run At",
	RecurringJobsNextRunAt:      "Next Run At",
	RecurringJobsErrorCount:     "Error Count",
	RecurringJobsActions:        "Actions",
	RecurringJobsRuns:           "Runs",

	// 重复任务状态值
	RecurringJobsStatusActive:    "Active",
	RecurringJobsStatusPaused:    "Paused",
	RecurringJobsStatusCompleted: "Completed",
	RecurringJobsStatusError:     "Error",

	// 重复任务编辑表单
	RecurringJobsEditFunctionName:   "Function Name",
	RecurringJobsEditCronExpression: "Cron Expression",
	RecurringJobsEditTimes:          "Run Limit",
	RecurringJobsEditArgs:           "Arguments",

	// 重复任务过滤标签
	RecurringJobsTabAll:       "All Tasks",
	RecurringJobsTabActive:    "Active Tasks",
	RecurringJobsTabPaused:    "Paused Tasks",
	RecurringJobsTabCompleted: "Completed Tasks",
	RecurringJobsTabError:     "Error Tasks",

	// 重复任务日志相关
	RecurringJobLogsID:         "ID",
	RecurringJobLogsJobID:      "Task Name",
	RecurringJobLogsStartedAt:  "Started At",
	RecurringJobLogsFinishedAt: "Finished At",
	RecurringJobLogsDuration:   "Duration",
	RecurringJobLogsSuccess:    "Status",
	RecurringJobLogsError:      "Error",
	RecurringJobLogsOutput:     "Output",

	// 重复任务日志过滤标签
	RecurringJobLogsTabAll:     "All Records",
	RecurringJobLogsTabSuccess: "Success Records",
	RecurringJobLogsTabFailed:  "Failed Records",

	// 操作按钮
	RecurringJobsPause:  "Pause",
	RecurringJobsResume: "Resume",
	RecurringJobsRun:    "Run Now",
	RecurringJobsDelete: "Delete",

	PagesID:         "ID",
	PagesTitle:      "Title",
	PagesSlug:       "Slug",
	PagesLocale:     "Locale",
	PagesNotes:      "Notes",
	PagesDraftCount: "Draft Count",
	PagesPath:       "Path",
	PagesOnline:     "Online",
	PagesVersion:    "Version",
	PagesVersions:   "Versions",
	PagesStartAt:    "Start At",
	PagesEndAt:      "End At",
	PagesOption:     "Option",
	PagesLive:       "Live",

	Page:                   "Page",
	PagesStatus:            "Status",
	PagesSchedule:          "PagesSchedule",
	PagesCategoryID:        "PagesCategoryID",
	PagesTemplateSelection: "PagesTemplateSelection",
	PagesEditContainer:     "PagesEditContainer",

	WebHeader:       "WebHeader",
	WebHeadersColor: "WebHeadersColor",
	Header:          "Header",
	Navigation:      "Navigation",
	Content:         "Content",

	WebFooter:            "WebFooter",
	WebFootersEnglishUrl: "WebFootersEnglishUrl",
	Footer:               "Footer",

	VideoBanner:                       "VideoBanner",
	VideoBannersAddTopSpace:           "VideoBannersAddTopSpace",
	VideoBannersAddBottomSpace:        "VideoBannersAddBottomSpace",
	VideoBannersAnchorID:              "VideoBannersAnchorID",
	VideoBannersVideo:                 "VideoBannersVideo",
	VideoBannersBackgroundVideo:       "VideoBannersBackgroundVideo",
	VideoBannersMobileBackgroundVideo: "VideoBannersMobileBackgroundVideo",
	VideoBannersVideoCover:            "VideoBannersVideoCover",
	VideoBannersMobileVideoCover:      "VideoBannersMobileVideoCover",
	VideoBannersHeading:               "VideoBannersHeading",
	VideoBannersPopupText:             "VideoBannersPopupText",
	VideoBannersText:                  "VideoBannersText",
	VideoBannersLinkText:              "VideoBannersLinkText",
	VideoBannersLink:                  "VideoBannersLink",

	Heading:                   "Heading",
	HeadingsAddTopSpace:       "HeadingsAddTopSpace",
	HeadingsAddBottomSpace:    "HeadingsAddBottomSpace",
	HeadingsAnchorID:          "HeadingsAnchorID",
	HeadingsHeading:           "HeadingsHeading",
	HeadingsFontColor:         "HeadingsFontColor",
	HeadingsBackgroundColor:   "HeadingsBackgroundColor",
	HeadingsLink:              "HeadingsLink",
	HeadingsLinkText:          "HeadingsLinkText",
	HeadingsLinkDisplayOption: "HeadingsLinkDisplayOption",
	HeadingsText:              "HeadingsText",

	BrandGrid:                "BrandGrid",
	BrandGridsAddTopSpace:    "BrandGridsAddTopSpace",
	BrandGridsAddBottomSpace: "BrandGridsAddBottomSpace",
	BrandGridsAnchorID:       "BrandGridsAnchorID",
	BrandGridsBrands:         "BrandGridsBrands",

	ListContent:                   "ListContent",
	ListContentsAddTopSpace:       "ListContentsAddTopSpace",
	ListContentsAddBottomSpace:    "ListContentsAddBottomSpace",
	ListContentsAnchorID:          "ListContentsAnchorID",
	ListContentsBackgroundColor:   "ListContentsBackgroundColor",
	ListContentsItems:             "ListContentsItems",
	ListContentsLink:              "ListContentsLink",
	ListContentsLinkText:          "ListContentsLinkText",
	ListContentsLinkDisplayOption: "ListContentsLinkDisplayOption",

	ImageContainer:                           "ImageContainer",
	ImageContainersAddTopSpace:               "ImageContainersAddTopSpace",
	ImageContainersAddBottomSpace:            "ImageContainersAddBottomSpace",
	ImageContainersAnchorID:                  "ImageContainersAnchorID",
	ImageContainersBackgroundColor:           "ImageContainersBackgroundColor",
	ImageContainersTransitionBackgroundColor: "ImageContainersTransitionBackgroundColor",
	ImageContainersImage:                     "ImageContainersImage",
	Image:                                    "Image",

	InNumber:                "InNumber",
	InNumbersAddTopSpace:    "InNumbersAddTopSpace",
	InNumbersAddBottomSpace: "InNumbersAddBottomSpace",
	InNumbersAnchorID:       "InNumbersAnchorID",
	InNumbersHeading:        "InNumbersHeading",
	InNumbersItems:          "InNumbersItems",
	InNumbers:               "InNumbers",

	ContactForm:                    "ContactForm",
	ContactFormsAddTopSpace:        "ContactFormsAddTopSpace",
	ContactFormsAddBottomSpace:     "ContactFormsAddBottomSpace",
	ContactFormsAnchorID:           "ContactFormsAnchorID",
	ContactFormsHeading:            "ContactFormsHeading",
	ContactFormsText:               "ContactFormsText",
	ContactFormsSendButtonText:     "ContactFormsSendButtonText",
	ContactFormsFormButtonText:     "ContactFormsFormButtonText",
	ContactFormsMessagePlaceholder: "ContactFormsMessagePlaceholder",
	ContactFormsNamePlaceholder:    "ContactFormsNamePlaceholder",
	ContactFormsEmailPlaceholder:   "ContactFormsEmailPlaceholder",
	ContactFormsThankyouMessage:    "ContactFormsThankyouMessage",
	ContactFormsActionUrl:          "ContactFormsActionUrl",
	ContactFormsPrivacyPolicy:      "ContactFormsPrivacyPolicy",

	ActivityActionLogIn:         "登录",
	ActivityActionExtendSession: "延长会话",

	PagesPage: "Page",

	User:                  "User",
	Role:                  "Role",
	LoginSession:          "Login Session",
	Dictionary:            "Dictionary",
	RecurringJob:          "Recurring Job",
	RecurringJobExecution: "Recurring Job Execution",
	Worker:                "Worker",
	WorkerJob:             "Worker Job",
}

var Messages_zh_CN_ModelsI18nModuleKey = &Messages_ModelsI18nModuleKey{
	Posts:          "帖子 示例",
	PostsID:        "ID",
	PostsTitle:     "标题",
	PostsHeroImage: "主图",
	PostsBody:      "内容",
	Example:        "QOR5演示",
	Settings:       "SEO 设置",
	Post:           "帖子",
	PostsBodyImage: "内容图片",

	SeoPost:             "帖子",
	SeoVariableTitle:    "标题",
	SeoVariableSiteName: "站点名称",

	QOR5Example: "QOR5 示例",
	Roles:       "权限管理",
	Users:       "用户管理",
	Dashboard:   "仪表盘",

	PageBuilder:              "页面管理菜单",
	Pages:                    "页面管理",
	SharedContainers:         "公用组件",
	DemoContainers:           "示例组件",
	Templates:                "模板页面",
	PageCategories:           "目录管理",
	ECManagement:             "电子商务管理",
	ECDashboard:              "电子商务仪表盘",
	Orders:                   "订单管理",
	InputDemos:               "表单 示例",
	Products:                 "产品管理",
	NestedFieldDemos:         "嵌套表单 示例",
	SiteManagement:           "站点管理菜单",
	SEO:                      "SEO 管理",
	UserManagement:           "用户管理菜单",
	Profile:                  "个人页面",
	FeaturedModelsManagement: "特色模块管理菜单",
	Customers:                "Customers 示例",
	ListModels:               "发布带排序及分页模块 示例",
	MicrositeModels:          "Microsite 示例",
	Workers:                  "后台工作管理",
	RecurringJobs:            "重复任务",
	RecurringJobLogs:         "重复任务日志",
	TaskManagement:           "任务管理",

	PagesID:         "ID",
	PagesTitle:      "标题",
	PagesSlug:       "Slug",
	PagesLocale:     "地区",
	PagesNotes:      "备注",
	PagesDraftCount: "草稿数",
	PagesPath:       "路径",
	PagesOnline:     "在线",
	PagesVersion:    "版本",
	PagesVersions:   "版本",
	PagesStartAt:    "开始时间",
	PagesEndAt:      "结束时间",
	PagesOption:     "选项",
	PagesLive:       "发布状态",

	Page:                   "Page",
	PagesStatus:            "状态",
	PagesSchedule:          "PagesSchedule",
	PagesCategoryID:        "PagesCategoryID",
	PagesTemplateSelection: "PagesTemplateSelection",
	PagesEditContainer:     "PagesEditContainer",

	WebHeader:       "WebHeader",
	WebHeadersColor: "WebHeadersColor",
	Header:          "Header",
	Navigation:      "Navigation",
	Content:         "Content",

	WebFooter:            "WebFooter",
	WebFootersEnglishUrl: "WebFootersEnglishUrl",
	Footer:               "Footer",

	VideoBanner:                       "VideoBanner",
	VideoBannersAddTopSpace:           "VideoBannersAddTopSpace",
	VideoBannersAddBottomSpace:        "VideoBannersAddBottomSpace",
	VideoBannersAnchorID:              "VideoBannersAnchorID",
	VideoBannersVideo:                 "VideoBannersVideo",
	VideoBannersBackgroundVideo:       "VideoBannersBackgroundVideo",
	VideoBannersMobileBackgroundVideo: "VideoBannersMobileBackgroundVideo",
	VideoBannersVideoCover:            "VideoBannersVideoCover",
	VideoBannersMobileVideoCover:      "VideoBannersMobileVideoCover",
	VideoBannersHeading:               "VideoBannersHeading",
	VideoBannersPopupText:             "VideoBannersPopupText",
	VideoBannersText:                  "VideoBannersText",
	VideoBannersLinkText:              "VideoBannersLinkText",
	VideoBannersLink:                  "VideoBannersLink",

	Heading:                   "Heading",
	HeadingsAddTopSpace:       "HeadingsAddTopSpace",
	HeadingsAddBottomSpace:    "HeadingsAddBottomSpace",
	HeadingsAnchorID:          "HeadingsAnchorID",
	HeadingsHeading:           "HeadingsHeading",
	HeadingsFontColor:         "HeadingsFontColor",
	HeadingsBackgroundColor:   "HeadingsBackgroundColor",
	HeadingsLink:              "HeadingsLink",
	HeadingsLinkText:          "HeadingsLinkText",
	HeadingsLinkDisplayOption: "HeadingsLinkDisplayOption",
	HeadingsText:              "HeadingsText",

	BrandGrid:                "BrandGrid",
	BrandGridsAddTopSpace:    "BrandGridsAddTopSpace",
	BrandGridsAddBottomSpace: "BrandGridsAddBottomSpace",
	BrandGridsAnchorID:       "BrandGridsAnchorID",
	BrandGridsBrands:         "BrandGridsBrands",

	ListContent:                   "ListContent",
	ListContentsAddTopSpace:       "ListContentsAddTopSpace",
	ListContentsAddBottomSpace:    "ListContentsAddBottomSpace",
	ListContentsAnchorID:          "ListContentsAnchorID",
	ListContentsBackgroundColor:   "ListContentsBackgroundColor",
	ListContentsItems:             "ListContentsItems",
	ListContentsLink:              "ListContentsLink",
	ListContentsLinkText:          "ListContentsLinkText",
	ListContentsLinkDisplayOption: "ListContentsLinkDisplayOption",

	ImageContainer:                           "ImageContainer",
	ImageContainersAddTopSpace:               "ImageContainersAddTopSpace",
	ImageContainersAddBottomSpace:            "ImageContainersAddBottomSpace",
	ImageContainersAnchorID:                  "ImageContainersAnchorID",
	ImageContainersBackgroundColor:           "ImageContainersBackgroundColor",
	ImageContainersTransitionBackgroundColor: "ImageContainersTransitionBackgroundColor",
	ImageContainersImage:                     "ImageContainersImage",
	Image:                                    "Image",

	InNumber:                "InNumber",
	InNumbersAddTopSpace:    "InNumbersAddTopSpace",
	InNumbersAddBottomSpace: "InNumbersAddBottomSpace",
	InNumbersAnchorID:       "InNumbersAnchorID",
	InNumbersHeading:        "InNumbersHeading",
	InNumbersItems:          "InNumbersItems",
	InNumbers:               "InNumbers",

	ContactForm:                    "ContactForm",
	ContactFormsAddTopSpace:        "ContactFormsAddTopSpace",
	ContactFormsAddBottomSpace:     "ContactFormsAddBottomSpace",
	ContactFormsAnchorID:           "ContactFormsAnchorID",
	ContactFormsHeading:            "ContactFormsHeading",
	ContactFormsText:               "ContactFormsText",
	ContactFormsSendButtonText:     "ContactFormsSendButtonText",
	ContactFormsFormButtonText:     "ContactFormsFormButtonText",
	ContactFormsMessagePlaceholder: "ContactFormsMessagePlaceholder",
	ContactFormsNamePlaceholder:    "ContactFormsNamePlaceholder",
	ContactFormsEmailPlaceholder:   "ContactFormsEmailPlaceholder",
	ContactFormsThankyouMessage:    "ContactFormsThankyouMessage",
	ContactFormsActionUrl:          "ContactFormsActionUrl",
	ContactFormsPrivacyPolicy:      "ContactFormsPrivacyPolicy",

	ActivityActionLogIn:         "登录",
	ActivityActionExtendSession: "延长会话",

	PagesPage: "Page",

	// 重复任务相关字段
	RecurringJobsName:           "任务名称",
	RecurringJobsFunctionName:   "函数名称",
	RecurringJobsCronExpression: "Cron表达式",
	RecurringJobsTimes:          "执行次数限制",
	RecurringJobsArgs:           "参数",
	RecurringJobsStatus:         "状态",
	RecurringJobsLastRunAt:      "上次执行时间",
	RecurringJobsNextRunAt:      "下次执行时间",
	RecurringJobsErrorCount:     "错误次数",
	RecurringJobsActions:        "操作",
	RecurringJobsRuns:           "执行次数",

	// 重复任务状态值
	RecurringJobsStatusActive:    "活跃",
	RecurringJobsStatusPaused:    "已暂停",
	RecurringJobsStatusCompleted: "已完成",
	RecurringJobsStatusError:     "错误",

	// 重复任务编辑表单
	RecurringJobsEditFunctionName:   "函数名称",
	RecurringJobsEditCronExpression: "Cron表达式",
	RecurringJobsEditTimes:          "执行次数限制",
	RecurringJobsEditArgs:           "参数",

	// 重复任务过滤标签
	RecurringJobsTabAll:       "全部任务",
	RecurringJobsTabActive:    "活跃任务",
	RecurringJobsTabPaused:    "已暂停",
	RecurringJobsTabCompleted: "已完成",
	RecurringJobsTabError:     "错误任务",

	// 重复任务日志相关
	RecurringJobLogsID:         "ID",
	RecurringJobLogsJobID:      "任务名称",
	RecurringJobLogsStartedAt:  "开始时间",
	RecurringJobLogsFinishedAt: "结束时间",
	RecurringJobLogsDuration:   "持续时间",
	RecurringJobLogsSuccess:    "状态",
	RecurringJobLogsError:      "错误",
	RecurringJobLogsOutput:     "输出",

	// 重复任务日志过滤标签
	RecurringJobLogsTabAll:     "全部记录",
	RecurringJobLogsTabSuccess: "成功记录",
	RecurringJobLogsTabFailed:  "失败记录",

	// 操作按钮
	RecurringJobsPause:  "暂停",
	RecurringJobsResume: "恢复",
	RecurringJobsRun:    "立即执行",
	RecurringJobsDelete: "删除",

	User:                  "用户",
	Role:                  "角色",
	LoginSession:          "登录会话",
	Dictionary:            "字典",
	RecurringJob:          "重复任务",
	RecurringJobExecution: "重复任务执行",
	Worker:                "后台工作",
	WorkerJob:             "后台工作任务",
}
