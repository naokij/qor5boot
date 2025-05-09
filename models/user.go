package models

import (
	"time"

	"github.com/qor5/admin/v3/role"
	"github.com/qor5/x/v3/login"
	"gorm.io/gorm"
)

const (
	RoleAdmin   = "Admin"
	RoleManager = "Manager"
	RoleEditor  = "Editor"
	RoleViewer  = "Viewer"

	OAuthProviderGoogle          = "google"
	OAuthProviderMicrosoftOnline = "microsoftonline"
	OAuthProviderGithub          = "github"
	OAuthProviderWecom           = "wecom"

	StatusActive   = "active"
	StatusInactive = "inactive"
)

var DefaultRoles = []string{
	RoleAdmin,
	RoleManager,
	RoleEditor,
	RoleViewer,
}

var OAuthProviders = []string{
	OAuthProviderGoogle,
	OAuthProviderMicrosoftOnline,
	OAuthProviderGithub,
	OAuthProviderWecom,
}

type User struct {
	gorm.Model

	Name             string
	Status           string
	Company          string
	RegistrationDate time.Time
	Roles            []role.Role `gorm:"many2many:user_role_join;"`
	UpdatedAt        time.Time
	CreatedAt        time.Time

	// Username is email
	LDAPUserPass
	login.OAuthInfo
	login.SessionSecure
}

func (u User) GetName() string {
	return u.Name
}

func (u User) GetID() uint {
	return u.ID
}

func (u User) GetRoles() (rs []string) {
	for _, r := range u.Roles {
		rs = append(rs, r.Name)
	}
	if len(rs) == 0 {
		rs = []string{RoleViewer}
	}
	return
}

func (u User) IsOAuthUser() bool {
	return u.OAuthProvider != ""
}
