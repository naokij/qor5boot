package models

import (
	"fmt"
	"log"
	"time"

	"github.com/qor5/x/v3/login"
	"gorm.io/gorm"
)

/*
LDAPUserPass 说明:

基本结构与功能:
LDAPUserPass 是系统认证机制的核心实现，支持混合认证模式（LDAP+本地）。该结构体嵌入到User结构中，
提供完整的用户认证、密码管理和TOTP双因素认证功能。

认证流程:
1. 用户登录时，优先尝试通过LDAP验证（如果LDAP已启用）
2. 若LDAP验证失败或发生错误，自动回退到本地密码验证
3. 支持账户锁定策略，防止暴力破解

LDAP集成:
- 通过admin包中的authenticateWithLDAP函数进行实际的LDAP认证
- 支持动态启用/禁用LDAP功能
- 认证成功后，在本地数据库中保存用户记录

安全特性:
- 密码存储采用哈希加密，不存储明文
- 支持账户锁定机制，多次登录失败后自动锁定
- 支持密码重置令牌，包含创建和过期机制
- 集成TOTP双因素认证

实现的接口方法:
1. 用户查找与账户管理:
   - FindUser: 通过账户名在数据库中查找用户
   - GetAccountName: 获取账户名

2. 认证与密码管理:
   - IsPasswordCorrect: 验证用户密码（LDAP优先，本地备选）
   - EncryptPassword: 加密用户密码
   - SetPassword: 设置新密码

3. 账户锁定机制:
   - GetLoginRetryCount: 获取登录失败次数
   - GetLocked: 检查账户是否锁定
   - LockUser: 锁定用户账户
   - UnlockUser: 解锁用户账户
   - IncreaseRetryCount: 增加登录失败计数

4. 密码重置:
   - GenerateResetPasswordToken: 生成密码重置令牌
   - GenerateResetPasswordTokenExpiration: 设置令牌过期时间
   - ConsumeResetPasswordToken: 使用密码重置令牌
   - GetResetPasswordToken: 获取密码重置令牌信息

5. TOTP双因素认证:
   - GetTOTPSecret: 获取TOTP密钥
   - GetIsTOTPSetup: 检查是否设置了TOTP
   - SetTOTPSecret: 设置TOTP密钥
   - SetIsTOTPSetup: 设置TOTP状态
   - SetLastUsedTOTPCode: 记录最后使用的TOTP代码
   - GetLastUsedTOTPCode: 获取最后使用的TOTP代码

LDAP配置:
系统通过以下变量控制LDAP功能:
- ldapEnabled: 是否启用LDAP
- ldapServer: LDAP服务器地址
- authenticateWithLDAP: LDAP认证函数

这些配置通过SetLDAPConfig函数从admin包中导入
*/

// LDAPUserPass 嵌入到User结构体，用于支持LDAP认证
type LDAPUserPass struct {
	Account  string `gorm:"index:,unique,where:account!='' and deleted_at is null"`
	Password string `gorm:"size:60"`
	// UnixNano string
	PassUpdatedAt               string
	LoginRetryCount             int
	Locked                      bool
	LockedAt                    *time.Time
	ResetPasswordToken          string `gorm:"index:,unique,where:reset_password_token!=''"`
	ResetPasswordTokenCreatedAt *time.Time
	ResetPasswordTokenExpiredAt *time.Time
	TOTPSecret                  string
	IsTOTPSetup                 bool
	LastUsedTOTPCode            string
	LastTOTPCodeUsedAt          *time.Time
}

// FindUser 查找用户
func (up *LDAPUserPass) FindUser(db *gorm.DB, model interface{}, account string) (user interface{}, err error) {
	err = db.Where("account = ?", account).
		First(model).
		Error
	if err != nil {
		return nil, err
	}
	return model, nil
}

// GetAccountName 获取账号名称
func (up *LDAPUserPass) GetAccountName() string {
	return up.Account
}

// GetLoginRetryCount 获取登录重试次数
func (up *LDAPUserPass) GetLoginRetryCount() int {
	return up.LoginRetryCount
}

// GetLocked 检查用户是否被锁定
func (up *LDAPUserPass) GetLocked() bool {
	if !up.Locked {
		return false
	}
	return up.Locked && up.LockedAt != nil && time.Now().Sub(*up.LockedAt) <= time.Hour
}

// GetTOTPSecret 获取TOTP密钥
func (up *LDAPUserPass) GetTOTPSecret() string {
	return up.TOTPSecret
}

// GetIsTOTPSetup 检查是否设置了TOTP
func (up *LDAPUserPass) GetIsTOTPSetup() bool {
	return up.IsTOTPSetup
}

// EncryptPassword 加密密码
func (up *LDAPUserPass) EncryptPassword() {
	// 调用原始UserPass的加密方法
	userPass := login.UserPass{
		Password: up.Password,
	}
	userPass.EncryptPassword()
	up.Password = userPass.Password
	up.PassUpdatedAt = fmt.Sprint(time.Now().UnixNano())
}

// IsPasswordCorrect 验证密码是否正确
func (up *LDAPUserPass) IsPasswordCorrect(password string) bool {
	// 如果LDAP未启用，直接使用本地认证
	if !ldapEnabled || ldapServer == "" {
		log.Printf("LDAP未启用，使用本地密码验证用户: %s", up.Account)
		userPass := login.UserPass{
			Password: up.Password,
		}
		return userPass.IsPasswordCorrect(password)
	}

	log.Printf("开始验证用户 %s 的密码，尝试LDAP认证", up.Account)

	// 尝试LDAP认证
	authenticated, err := authenticateWithLDAP(up.Account, password)
	if err != nil {
		log.Printf("LDAP认证过程出错: %v，回退到本地认证", err)
		// LDAP认证出错，回退到本地认证
		userPass := login.UserPass{
			Password: up.Password,
		}
		isCorrect := userPass.IsPasswordCorrect(password)
		log.Printf("本地认证结果: %v", isCorrect)
		return isCorrect
	}

	if authenticated {
		// LDAP认证成功
		log.Printf("用户 %s LDAP认证成功", up.Account)
		return true
	}

	// LDAP认证失败，回退到本地认证
	log.Printf("用户 %s LDAP认证失败，回退到本地认证", up.Account)
	userPass := login.UserPass{
		Password: up.Password,
	}
	isCorrect := userPass.IsPasswordCorrect(password)
	log.Printf("本地认证结果: %v", isCorrect)
	return isCorrect
}

// GetPasswordUpdatedAt 获取密码更新时间
func (up *LDAPUserPass) GetPasswordUpdatedAt() string {
	return up.PassUpdatedAt
}

// LockUser 锁定用户
func (up *LDAPUserPass) LockUser(db *gorm.DB, model interface{}) error {
	lockedAt := time.Now()
	if err := db.Model(model).Where("account = ?", up.Account).Updates(map[string]interface{}{
		"locked":    true,
		"locked_at": &lockedAt,
	}).Error; err != nil {
		return err
	}

	up.Locked = true
	up.LockedAt = &lockedAt

	return nil
}

// UnlockUser 解锁用户
func (up *LDAPUserPass) UnlockUser(db *gorm.DB, model interface{}) error {
	if err := db.Model(model).Where("account = ?", up.Account).Updates(map[string]interface{}{
		"locked":            false,
		"login_retry_count": 0,
		"locked_at":         nil,
	}).Error; err != nil {
		return err
	}

	up.Locked = false
	up.LoginRetryCount = 0
	up.LockedAt = nil

	return nil
}

// IncreaseRetryCount 增加重试计数
func (up *LDAPUserPass) IncreaseRetryCount(db *gorm.DB, model interface{}) error {
	if err := db.Model(model).Where("account = ?", up.Account).Updates(map[string]interface{}{
		"login_retry_count": gorm.Expr("coalesce(login_retry_count,0) + 1"),
	}).Error; err != nil {
		return err
	}
	up.LoginRetryCount++

	return nil
}

// GenerateResetPasswordTokenExpiration 生成重置密码令牌过期时间
func (up *LDAPUserPass) GenerateResetPasswordTokenExpiration(db *gorm.DB) (createdAt time.Time, expiredAt time.Time) {
	createdAt = db.NowFunc()
	return createdAt, createdAt.Add(10 * time.Minute)
}

// GenerateResetPasswordToken 生成重置密码令牌
func (up *LDAPUserPass) GenerateResetPasswordToken(db *gorm.DB, model interface{}) (token string, err error) {
	// 调用原始UserPass的方法
	userPass := login.UserPass{
		Account: up.Account,
	}
	token, err = userPass.GenerateResetPasswordToken(db, model)
	if err != nil {
		return "", err
	}
	// 同步状态
	up.ResetPasswordToken = userPass.ResetPasswordToken
	up.ResetPasswordTokenCreatedAt = userPass.ResetPasswordTokenCreatedAt
	up.ResetPasswordTokenExpiredAt = userPass.ResetPasswordTokenExpiredAt
	return token, nil
}

// ConsumeResetPasswordToken 消费重置密码令牌
func (up *LDAPUserPass) ConsumeResetPasswordToken(db *gorm.DB, model interface{}) error {
	err := db.Model(model).
		Where("account = ?", up.Account).
		Updates(map[string]interface{}{
			"reset_password_token_expired_at": time.Now(),
		}).
		Error
	if err != nil {
		return err
	}
	return nil
}

// GetResetPasswordToken 获取重置密码令牌
func (up *LDAPUserPass) GetResetPasswordToken() (token string, createdAt *time.Time, expired bool) {
	if up.ResetPasswordTokenExpiredAt != nil && time.Now().Sub(*up.ResetPasswordTokenExpiredAt) > 0 {
		return "", nil, true
	}
	return up.ResetPasswordToken, up.ResetPasswordTokenCreatedAt, false
}

// SetPassword 设置密码
func (up *LDAPUserPass) SetPassword(db *gorm.DB, model interface{}, password string) error {
	up.Password = password
	up.EncryptPassword()
	err := db.Model(model).
		Where("account = ?", up.Account).
		Updates(map[string]interface{}{
			"password":        up.Password,
			"pass_updated_at": up.PassUpdatedAt,
		}).
		Error
	if err != nil {
		return err
	}
	return nil
}

// SetTOTPSecret 设置TOTP密钥
func (up *LDAPUserPass) SetTOTPSecret(db *gorm.DB, model interface{}, key string) error {
	if err := db.Model(model).Where("account = ?", up.Account).Updates(map[string]interface{}{
		"totp_secret": key,
	}).Error; err != nil {
		return err
	}

	up.TOTPSecret = key

	return nil
}

// SetIsTOTPSetup 设置TOTP是否已设置
func (up *LDAPUserPass) SetIsTOTPSetup(db *gorm.DB, model interface{}, v bool) error {
	if err := db.Model(model).Where("account = ?", up.Account).Updates(map[string]interface{}{
		"is_totp_setup": v,
	}).Error; err != nil {
		return err
	}

	up.IsTOTPSetup = v

	return nil
}

// SetLastUsedTOTPCode 设置最后使用的TOTP代码
func (up *LDAPUserPass) SetLastUsedTOTPCode(db *gorm.DB, model interface{}, passcode string) error {
	now := time.Now()
	if err := db.Model(model).Where("account = ?", up.Account).Updates(map[string]interface{}{
		"last_used_totp_code":    passcode,
		"last_totp_code_used_at": &now,
	}).Error; err != nil {
		return err
	}

	up.LastUsedTOTPCode = passcode
	up.LastTOTPCodeUsedAt = &now

	return nil
}

// GetLastUsedTOTPCode 获取最后使用的TOTP代码
func (up *LDAPUserPass) GetLastUsedTOTPCode() (code string, usedAt *time.Time) {
	return up.LastUsedTOTPCode, up.LastTOTPCodeUsedAt
}

// 注意：这些LDAP相关变量需要在admin包中设置，这里只是声明用于编译
var (
	ldapEnabled          bool
	ldapServer           string
	authenticateWithLDAP func(email, password string) (bool, error)
)

// SetLDAPConfig 从外部设置LDAP配置
func SetLDAPConfig(enabled bool, server string, authFunc func(email, password string) (bool, error)) {
	ldapEnabled = enabled
	ldapServer = server
	authenticateWithLDAP = authFunc
}
