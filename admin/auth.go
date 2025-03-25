package admin

import (
	_ "embed"
	"errors"
	"net/http"

	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/google"
	"github.com/markbates/goth/providers/microsoftonline"
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
		WrapAfterOAuthComplete(func(in login.HookFunc) login.HookFunc {
			return func(r *http.Request, user interface{}, extraVals ...interface{}) error {
				if err := in(r, user, extraVals...); err != nil {
					return err
				}

				u := user.(goth.User)
				if u.Email == "" {
					return nil
				}

				// 查找是否有匹配的用户账号
				var existingUser models.User
				err := db.Where("account = ?", u.Email).First(&existingUser).Error

				if errors.Is(err, gorm.ErrRecordNotFound) {
					// 没有找到匹配用户，拒绝登录
					return &login.NoticeError{
						Level:   login.NoticeLevel_Error,
						Message: "您的账号未在系统中注册。请联系管理员授权后再使用OAuth登录。",
					}
				} else if err != nil {
					// 数据库查询错误
					return err
				}

				// 找到匹配用户，更新OAuth信息
				existingUser.OAuthProvider = u.Provider
				existingUser.OAuthUserID = u.UserID
				existingUser.OAuthIdentifier = u.Email
				existingUser.OAuthAvatar = u.AvatarURL

				if err := db.Save(&existingUser).Error; err != nil {
					return err
				}

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
