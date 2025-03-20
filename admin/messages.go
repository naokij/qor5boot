package admin

import (
	"github.com/qor5/x/v3/i18n"
)

const I18nAdminKey i18n.ModuleKey = "I18nAdminKey"

type Messages struct {
	FilterTabsAll                  string
	FilterTabsHasUnreadNotes       string
	FilterTabsActive               string
	DemoTips                       string
	DemoUsernameLabel              string
	DemoPasswordLabel              string
	LoginProviderGoogleText        string
	LoginProviderMicrosoftText     string
	LoginProviderGithubText        string
	OAuthCompleteInfoTitle         string
	OAuthCompleteInfoPositionLabel string
	OAuthCompleteInfoAgreeLabel    string
	OAuthCompleteInfoBackLabel     string
	Demo                           string
	DBResetTipLabel                string
	Name                           string
	Email                          string
	Company                        string
	Role                           string
	Status                         string
	ChangePassword                 string
	LoginSessions                  string
	PasswordMinLengthHint          string
	PasswordMinLengthError         string
	PasswordPlaceholder            string
	Dashboard                      string
}

var Messages_en_US = &Messages{
	FilterTabsAll:                  "All",
	FilterTabsHasUnreadNotes:       "Has Unread Notes",
	FilterTabsActive:               "Active",
	DemoTips:                       "Please note that the database would be reset every even hour.",
	DemoUsernameLabel:              "Demo Username: ",
	DemoPasswordLabel:              "Demo Password: ",
	LoginProviderGoogleText:        "Login with Google",
	LoginProviderMicrosoftText:     "Login with Microsoft",
	LoginProviderGithubText:        "Login with Github",
	OAuthCompleteInfoTitle:         "Complete your information",
	OAuthCompleteInfoPositionLabel: "Position(Optional)",
	OAuthCompleteInfoAgreeLabel:    "Subscribe to QOR5 newsletter(Optional)",
	OAuthCompleteInfoBackLabel:     "Back to login",
	Demo:                           "DEMO",
	DBResetTipLabel:                "Database reset countdown",
	Name:                           "Name",
	Email:                          "Email",
	Company:                        "Company",
	Role:                           "Role",
	Status:                         "Status",
	ChangePassword:                 "Change Password",
	LoginSessions:                  "Login Sessions",
	PasswordMinLengthHint:          "Password must be at least 12 characters long",
	PasswordMinLengthError:         "Password must be at least 12 characters long",
	PasswordPlaceholder:            "Enter password",
	Dashboard:                      "Dashboard",
}

var Messages_zh_CN = &Messages{
	FilterTabsAll:                  "全部",
	FilterTabsHasUnreadNotes:       "未读备注",
	FilterTabsActive:               "有效",
	DemoTips:                       "请注意，数据库将每隔偶数小时重置一次。",
	DemoUsernameLabel:              "演示账户：",
	DemoPasswordLabel:              "演示密码：",
	LoginProviderGoogleText:        "使用Google登录",
	LoginProviderMicrosoftText:     "使用Microsoft登录",
	LoginProviderGithubText:        "使用Github登录",
	OAuthCompleteInfoTitle:         "请填写您的信息",
	OAuthCompleteInfoPositionLabel: "职位（可选）",
	OAuthCompleteInfoAgreeLabel:    "订阅QOR5新闻（可选）",
	OAuthCompleteInfoBackLabel:     "返回登录",
	Demo:                           "演示",
	DBResetTipLabel:                "数据库重置倒计时",
	Name:                           "姓名",
	Email:                          "邮箱",
	Company:                        "公司",
	Role:                           "角色",
	Status:                         "状态",
	ChangePassword:                 "修改密码",
	LoginSessions:                  "登录会话",
	PasswordMinLengthHint:          "密码长度至少需要12位",
	PasswordMinLengthError:         "密码长度至少需要12位",
	PasswordPlaceholder:            "输入密码",
	Dashboard:                      "仪表盘",
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
	Workers:                  "后台工作进程管理",

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
}
