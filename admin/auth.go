package admin

import (
	_ "embed"
	"net/http"

	"github.com/markbates/goth/providers/dingtalk"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/google"
	"github.com/markbates/goth/providers/microsoftonline"
	"github.com/markbates/goth/providers/wecom"
	"github.com/naokij/qor5boot/models"
	"github.com/qor5/admin/v3/activity"
	plogin "github.com/qor5/admin/v3/login"
	"github.com/qor5/admin/v3/presets"
	"github.com/qor5/admin/v3/role"
	"github.com/qor5/web/v3"
	"github.com/qor5/x/v3/i18n"
	"github.com/qor5/x/v3/login"
	. "github.com/theplant/htmlgo"
	"gorm.io/gorm"
)

var (
	loginSecret                = getEnvWithDefault("LOGIN_SECRET", "")
	loginGoogleKey             = getEnvWithDefault("LOGIN_GOOGLE_KEY", "")
	loginGoogleSecret          = getEnvWithDefault("LOGIN_GOOGLE_SECRET", "")
	loginMicrosoftOnlineKey    = getEnvWithDefault("LOGIN_MICROSOFTONLINE_KEY", "")
	loginMicrosoftOnlineSecret = getEnvWithDefault("LOGIN_MICROSOFTONLINE_SECRET", "")
	loginGithubKey             = getEnvWithDefault("LOGIN_GITHUB_KEY", "")
	loginGithubSecret          = getEnvWithDefault("LOGIN_GITHUB_SECRET", "")
	loginWecomCorpid           = getEnvWithDefault("LOGIN_WECOM_CORPID", "")
	loginWecomSecret           = getEnvWithDefault("LOGIN_WECOM_SECRET", "")
	loginWecomAgentid          = getEnvWithDefault("LOGIN_WECOM_AGENTID", "")
	loginDingtalkAppid         = getEnvWithDefault("LOGIN_DINGTALK_APPID", "")
	loginDingtalkSecret        = getEnvWithDefault("LOGIN_DINGTALK_SECRET", "")
	loginDingtalkCorpId        = getEnvWithDefault("LOGIN_DINGTALK_CORPID", "")
	baseURL                    = getEnvWithDefault("BASE_URL", "")
	recaptchaSiteKey           = getEnvWithDefault("RECAPTCHA_SITE_KEY", "")
	recaptchaSecret            = getEnvWithDefault("RECAPTCHA_SECRET_KEY", "")
	loginInitialUserEmail      = getEnvWithDefault("LOGIN_INITIAL_USER_EMAIL", "")
	loginInitialUserPassword   = getEnvWithDefault("LOGIN_INITIAL_USER_PASSWORD", "")
)

func getCurrentUser(r *http.Request) (u *models.User) {
	u, ok := login.GetCurrentUser(r).(*models.User)
	if !ok {
		return nil
	}

	return u
}

func initLoginSessionBuilder(db *gorm.DB, pb *presets.Builder, ab *activity.Builder) *plogin.SessionBuilder {
	loginBuilder := plogin.New(pb).
		DB(db).
		UserModel(&models.User{}).
		Secret(loginSecret).
		OAuthProviders(func() []*login.Provider {
			var providers []*login.Provider

			if loginGoogleKey != "" && loginGoogleSecret != "" {
				providers = append(providers, &login.Provider{
					Goth: google.New(loginGoogleKey, loginGoogleSecret, baseURL+"/auth/callback?provider="+models.OAuthProviderGoogle),
					Key:  models.OAuthProviderGoogle,
					Text: "LoginProviderGoogleText",
					Logo: RawHTML(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 48 48" width="24px" height="24px"><path fill="#fbc02d" d="M43.611,20.083H42V20H24v8h11.303c-1.649,4.657-6.08,8-11.303,8c-6.627,0-12-5.373-12-12	s5.373-12,12-12c3.059,0,5.842,1.154,7.961,3.039l5.657-5.657C34.046,6.053,29.268,4,24,4C12.955,4,4,12.955,4,24s8.955,20,20,20	s20-8.955,20-20C44,22.659,43.862,21.35,43.611,20.083z"></path><path fill="#e53935" d="M6.306,14.691l6.571,4.819C14.655,15.108,18.961,12,24,12c3.059,0,5.842,1.154,7.961,3.039	l5.657-5.657C34.046,6.053,29.268,4,24,4C16.318,4,9.656,8.337,6.306,14.691z"></path><path fill="#4caf50" d="M24,44c5.166,0,9.86-1.977,13.409-5.192l-6.19-5.238C29.211,35.091,26.715,36,24,36	c-5.202,0-9.619-3.317-11.283-7.946l-6.522,5.025C9.505,39.556,16.227,44,24,44z"></path><path fill="#1565c0" d="M43.611,20.083L43.595,20L42,20H24v8h11.303c-0.792,2.237-2.231,4.166-4.087,5.571	c0.001-0.001,0.002-0.001,0.003-0.002l6.19,5.238C36.971,39.205,44,34,44,24C44,22.659,43.862,21.35,43.611,20.083z"></path></svg>`),
				})
			}

			if loginMicrosoftOnlineKey != "" && loginMicrosoftOnlineSecret != "" {
				providers = append(providers, &login.Provider{
					Goth: microsoftonline.New(loginMicrosoftOnlineKey, loginMicrosoftOnlineSecret, baseURL+"/auth/callback"),
					Key:  models.OAuthProviderMicrosoftOnline,
					Text: "LoginProviderMicrosoftText",
					Logo: RawHTML(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 48 48" width="24px" height="24px"><path fill="#f35325" d="M2 2h20v20H2z"/><path fill="#81bc06" d="M24 2h20v20H24z"/><path fill="#05a6f0" d="M2 24h20v20H2z"/><path fill="#ffba08" d="M24 24h20v20H24z"/></svg>`),
				})
			}

			if loginGithubKey != "" && loginGithubSecret != "" {
				providers = append(providers, &login.Provider{
					Goth: github.New(loginGithubKey, loginGithubSecret, baseURL+"/auth/callback?provider="+models.OAuthProviderGithub),
					Key:  models.OAuthProviderGithub,
					Text: "LoginProviderGithubText",
					Logo: RawHTML(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 96 96" width="24px" height="24px"><path fill-rule="evenodd" clip-rule="evenodd" d="M48.854 0C21.839 0 0 22 0 49.217c0 21.756 13.993 40.172 33.405 46.69 2.427.49 3.316-1.059 3.316-2.362 0-1.141-.08-5.052-.08-9.127-13.59 2.934-16.42-5.867-16.42-5.867-2.184-5.704-5.42-7.17-5.42-7.17-4.448-3.015.324-3.015.324-3.015 4.934.326 7.523 5.052 7.523 5.052 4.367 7.496 11.404 5.378 14.235 4.074.404-3.178 1.699-5.378 3.074-6.6-10.839-1.141-22.243-5.378-22.243-24.283 0-5.378 1.94-9.778 5.014-13.2-.485-1.222-2.184-6.275.486-13.038 0 0 4.125-1.304 13.426 5.052a46.97 46.97 0 0 1 12.214-1.63c4.125 0 8.33.571 12.213 1.63 9.302-6.356 13.427-5.052 13.427-5.052 2.67 6.763.97 11.816.485 13.038 3.155 3.422 5.015 7.822 5.015 13.2 0 18.905-11.404 23.06-22.324 24.283 1.78 1.548 3.316 4.481 3.316 9.126 0 6.6-.08 11.897-.08 13.526 0 1.304.89 2.853 3.316 2.364 19.412-6.52 33.405-24.935 33.405-46.691C97.707 22 75.788 0 48.854 0z" fill="#24292f"/></svg>`),
				})
			}

			if loginWecomCorpid != "" && loginWecomSecret != "" && loginWecomAgentid != "" {
				providers = append(providers, &login.Provider{
					Goth: wecom.New(loginWecomCorpid, loginWecomSecret, loginWecomAgentid, baseURL+"/auth/callback?provider="+models.OAuthProviderWecom),
					Key:  models.OAuthProviderWecom,
					Text: "LoginProviderWecomText",
					Logo: RawHTML(`<svg width="24" height="24" viewBox="0 0 202 169" xmlns="http://www.w3.org/2000/svg"><g fill="none" fill-rule="evenodd"><path d="M134.19 137.564a2.774 2.774 0 00.583.538 35.667 35.667 0 0110.599 19.059 11.555 11.555 0 00.388 2.35 11.555 11.555 0 002.986 5.144c4.532 4.53 11.879 4.53 16.41 0 4.532-4.532 4.532-11.88 0-16.411a11.55 11.55 0 00-6.988-3.333 35.667 35.667 0 01-19.857-11.054h-.002a2.775 2.775 0 00-4.12 3.707" fill="#FB6500"/><path d="M170.88 146.48a2.774 2.774 0 00.538-.583 35.667 35.667 0 0119.059-10.599 11.555 11.555 0 002.35-.389 11.555 11.555 0 005.143-2.986c4.531-4.532 4.531-11.879 0-16.41-4.532-4.532-11.879-4.532-16.41 0a11.55 11.55 0 00-3.334 6.988 35.667 35.667 0 01-11.054 19.857l.001.002a2.775 2.775 0 003.706 4.12" fill="#0082EF"/><path d="M179.795 109.79a2.774 2.774 0 00-.583-.538 35.667 35.667 0 01-10.599-19.059 11.555 11.555 0 00-.388-2.351 11.555 11.555 0 00-2.986-5.143c-4.532-4.531-11.88-4.531-16.411 0-4.531 4.532-4.531 11.879 0 16.41a11.55 11.55 0 006.989 3.334 35.667 35.667 0 0119.857 11.054l.002-.001a2.775 2.775 0 004.119-3.706" fill="#2DBC00"/><path d="M143.105 100.874a2.774 2.774 0 00-.538.583 35.667 35.667 0 01-19.059 10.599 11.555 11.555 0 00-2.35.388 11.555 11.555 0 00-5.144 2.986c-4.53 4.532-4.53 11.88 0 16.411 4.532 4.531 11.88 4.531 16.411 0a11.55 11.55 0 003.333-6.989 35.667 35.667 0 0111.054-19.857v-.002a2.775 2.775 0 00-3.707-4.119" fill="#FC0"/><path d="M160.36 44.175c-3.228-6.632-7.565-12.788-12.89-18.298C134.014 11.96 115.176 2.997 94.426.637A97.692 97.692 0 0083.45 0c-3.39 0-6.924.197-10.506.587-20.843 2.266-39.786 11.183-53.341 25.11-5.35 5.496-9.71 11.638-12.96 18.252C2.236 52.925 0 62.453 0 72.27c0 12.636 3.845 25.096 11.12 36.034 4.115 6.186 11.15 14.181 17.172 19.08l-3.336 13.795-.96 3.885c-.17.382-.3.782-.392 1.2-.057.253-.075.52-.1.785-.019.201-.06.395-.06.602a6.422 6.422 0 006.422 6.422c1.16 0 2.233-.33 3.172-.868l.089-.049c.137-.08.277-.153.408-.243l23.82-11.956c5.128 1.471 10.215 2.412 15.578 3.002 3.486.384 7.025.578 10.518.578 3.563 0 7.255-.215 10.975-.638a93.368 93.368 0 0021.03-4.915 11.564 11.564 0 01-2.118-.945c-4.085-2.337-6.201-6.783-5.781-11.183a78.463 78.463 0 01-14.91 3.188 81.649 81.649 0 01-9.196.536c-2.922 0-5.884-.163-8.802-.484-.61-.067-1.216-.16-1.823-.24-4.008-.535-7.971-1.352-11.8-2.474a8.057 8.057 0 00-2.453-.375c-1.32 0-2.6.35-3.9 1.023-.168.087-.334.16-.504.258l-15.283 9.007-.017.01c-.316.184-.495.253-.66.253-.537 0-.973-.454-.973-1.012l.566-2.283.65-2.48 2.26-8.608c.146-.528.296-1.162.296-1.874a6.16 6.16 0 00-2.509-4.968 64.836 64.836 0 01-6.124-5.143c-3.317-3.147-6.23-6.582-8.677-10.262-5.762-8.662-8.808-18.49-8.808-28.426 0-7.71 1.764-15.208 5.24-22.284 2.597-5.286 6.095-10.208 10.398-14.63C41.67 24.14 57.337 16.8 74.645 14.917a81.755 81.755 0 018.805-.494c2.972 0 6.066.18 9.196.536 17.227 1.96 32.805 9.336 43.865 20.774 4.282 4.431 7.762 9.364 10.34 14.663 3.42 7.031 5.155 14.468 5.155 22.105 0 .794-.05 1.585-.089 2.377 4.47-2.735 10.375-2.198 14.243 1.67.196.195.36.41.538.615.128-1.627.2-3.259.2-4.895 0-9.724-2.2-19.176-6.54-28.093" fill="#0079DE"/></g></svg>`),
				})
			}

			if loginDingtalkAppid != "" && loginDingtalkSecret != "" {
				providers = append(providers, &login.Provider{
					Goth: dingtalk.New(loginDingtalkAppid, loginDingtalkSecret, baseURL+"/auth/callback?provider="+models.OAuthProviderDingtalk, loginDingtalkCorpId, "openid", "corpid"),
					Key:  models.OAuthProviderDingtalk,
					Text: "LoginProviderDingtalkText",
					Logo: RawHTML(`<svg width="56" height="24" xmlns="http://www.w3.org/2000/svg"><path d="M14.474 11.239C10.651 8.287 6.341 4.378 1.616.185c-.373-.33-.703-.2-.87.239-1.064 2.812-.031 5.31 1.636 6.75 1.675 1.444 4.164 2.775 5.69 3.481.059.028.007.115-.053.088C5.217 9.51 3.269 8.618.625 6.58c-.284-.22-.572-.134-.607.294-.216 2.688 1.504 4.8 3.424 5.513 1.189.442 2.488.686 3.693.832.064.007.05.1-.014.1-1.553.005-3.43-.367-5.054-.991-.343-.132-.462.142-.41.359.278 1.135 1.683 2.872 3.92 3.237.417.068.863.077 1.26.065.095-.003.117.05.087.122l-1.586 2.714c-.052.086-.02.156.088.156h2.005c.092 0 .15.06.102.138-.048.079-2.677 4.439-2.801 4.647-.109.181.02.32.223.171l8.846-6.476c.11-.081.083-.19-.071-.19h-1.798c-.118 0-.143-.079-.063-.159.08-.08 2.04-2.033 2.736-2.764.726-.763 1.097-2.161-.13-3.11m28.633-6.29h-2.031L39.74 7.545c-.071.147.019.303.175.303h6.061l.213-1.733h-3.68l.599-1.167zm-3.884 6.665h1.604l-.25 2.043h-1.605l-.213 1.733h1.606l-.524 4.285c-.115.935.838 1.417 1.778 1.25 1.289-.23 1.79-.395 2.86-.887l.231-1.883c-.607.292-1.847.65-2.499.801-.209.05-.395-.073-.37-.282l.404-3.284h2.805l.213-1.733h-2.805l.25-2.043h2.805l.213-1.732h-6.29l-.214 1.732zm16.209-5.498h-7.678l-.213 1.733h4.138l-1.298 10.576a.35.35 0 01-.339.301h-3.278l-.231 1.884h3.956c.832 0 1.59-.675 1.691-1.507l1.382-11.254h1.657l.213-1.733zM24.164 4.948h-2.03l-1.338 2.597c-.07.147.02.303.176.303h6.061l.213-1.733h-3.68l.598-1.167zm-3.884 6.665h1.604l-.25 2.043H20.03l-.213 1.733h1.606l-.524 4.285c-.115.935.838 1.417 1.778 1.25 1.289-.23 1.791-.395 2.86-.887l.232-1.883c-.608.292-1.848.65-2.5.801-.209.05-.395-.073-.37-.282l.405-3.284h2.804l.213-1.733h-2.806l.251-2.043h2.805l.213-1.732h-6.291l-.213 1.732h.001zM28.6 7.848h4.138l-1.3 10.576a.35.35 0 01-.338.301h-3.278l-.23 1.884h3.955c.832 0 1.59-.675 1.692-1.507L34.62 7.848h1.657l.213-1.733h-7.68l-.212 1.733z" fill="#007FFF" fill-rule="evenodd"/></svg>`),
				})
			}

			return providers
		}()...).
		HomeURLFunc(func(r *http.Request, user interface{}) string {
			return "/"
		}).
		MaxRetryCount(5).
		// TODO online  to set  true
		Recaptcha(false, login.RecaptchaConfig{
			SiteKey:   recaptchaSiteKey,
			SecretKey: recaptchaSecret,
		}).
		WrapBeforeSetPassword(func(in login.HookFunc) login.HookFunc {
			return func(r *http.Request, user interface{}, extraVals ...interface{}) error {
				if err := in(r, user, extraVals...); err != nil {
					return err
				}
				password := extraVals[0].(string)
				if len(password) < 12 {
					return &login.NoticeError{
						Level:   login.NoticeLevel_Error,
						Message: "Password cannot be less than 12 characters",
					}
				}
				return nil
			}
		}).
		OAuthIdentifier(models.OAuthProviderWecom, login.OAuthIdentifierUserID).
		OAuthIdentifier(models.OAuthProviderGoogle, login.OAuthIdentifierEmail).
		OAuthIdentifier(models.OAuthProviderMicrosoftOnline, login.OAuthIdentifierEmail).
		OAuthIdentifier(models.OAuthProviderGithub, login.OAuthIdentifierEmail).
		OAuthIdentifier(models.OAuthProviderDingtalk, login.OAuthIdentifierNickName).
		AfterOAuthComplete(func(r *http.Request, user interface{}, extraVals ...interface{}) error {
			return nil
		}).
		WrapAfterOAuthComplete(func(in login.HookFunc) login.HookFunc {
			return func(r *http.Request, user interface{}, extraVals ...interface{}) error {
				return nil
			}
		}).TOTP(false).MaxRetryCount(0)

	loginBuilder.LoginPageFunc(plogin.NewAdvancedLoginPage(func(ctx *web.EventContext, config *plogin.AdvancedLoginPageConfig) (*plogin.AdvancedLoginPageConfig, error) {
		// 从自定义消息中获取登录标题
		adminMsg := i18n.MustGetModuleMessages(ctx.R, I18nAdminKey, Messages_zh_CN).(*Messages)
		config.TitleLabel = adminMsg.LoginTitleLabel
		config.PageTitle = adminMsg.LoginTitleLabel

		config.BrandLogo = RawHTML(string(LoginLogoSVG))

		// 只使用CSS隐藏并替换原始SVG
		ctx.Injector.HeadHTML(`
		<style>
			/* 完全隐藏原始的SVG */
			.v-row .hidden-md-and-down .position-absolute svg {
				display: none !important;
			}
			/* 在父容器添加我们的SVG作为背景 */
			.v-row .hidden-md-and-down .position-absolute {
				background-image: url('/admin/assets/login_logo.svg');
				background-repeat: no-repeat;
				background-size: contain;
				width: 61px;
				height: 27px;
			}
		</style>
		`)
		config.OAuthProviderDisplay = func(provider *login.Provider) plogin.OAuthProviderDisplay {
			return plogin.OAuthProviderDisplay{
				Logo: provider.Logo,
				Text: i18n.T(ctx.R, I18nAdminKey, provider.Text),
			}
		}
		return config, nil
	})(loginBuilder.ViewHelper(), pb))

	genInitialUser(db)

	return plogin.NewSessionBuilder(loginBuilder, db).
		Activity(ab.RegisterModel(&models.User{})).
		IsPublicUser(func(u interface{}) bool {
			return false
		}).
		TablePrefix("cms_").
		AutoMigrate()
}

func genInitialUser(db *gorm.DB) {
	email := loginInitialUserEmail
	password := loginInitialUserPassword
	if email == "" || password == "" {
		return
	}

	var count int64
	if err := db.Model(&models.User{}).Where("account = ?", email).Count(&count).Error; err != nil {
		panic(err)
	}

	if count > 0 {
		return
	}
	if err := initDefaultRoles(db); err != nil {
		panic(err)
	}

	user := &models.User{
		Name:   email,
		Status: models.StatusActive,
		LDAPUserPass: models.LDAPUserPass{
			Account:  email,
			Password: password,
		},
	}
	user.EncryptPassword()
	if err := db.Create(user).Error; err != nil {
		panic(err)
	}
	if err := grantUserRole(db, user.ID, models.RoleAdmin); err != nil {
		panic(err)
	}
}

func grantUserRole(db *gorm.DB, userID uint, roleName string) error {
	var roleID int
	if err := db.Table("roles").Where("name = ?", roleName).Pluck("id", &roleID).Error; err != nil {
		panic(err)
	}
	return db.Table("user_role_join").Create(
		&map[string]interface{}{
			"user_id": userID,
			"role_id": roleID,
		}).Error
}

func initDefaultRoles(db *gorm.DB) error {
	var cnt int64
	if err := db.Model(&role.Role{}).Count(&cnt).Error; err != nil {
		return err
	}

	if cnt == 0 {
		var roles []*role.Role
		for _, r := range models.DefaultRoles {
			roles = append(roles, &role.Role{
				Name: r,
			})
		}

		if err := db.Create(roles).Error; err != nil {
			return err
		}
	}

	return nil
}
