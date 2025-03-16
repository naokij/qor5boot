package admin

import (
	"fmt"
	"net/url"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/qor5/admin/v3/activity"
	plogin "github.com/qor5/admin/v3/login"
	"github.com/qor5/admin/v3/presets"
	"github.com/qor5/admin/v3/presets/gorm2op"
	"github.com/qor5/admin/v3/role"
	"github.com/qor5/web/v3"
	"github.com/qor5/x/v3/i18n"
	"github.com/qor5/x/v3/perm"
	v "github.com/qor5/x/v3/ui/vuetify"
	vx "github.com/qor5/x/v3/ui/vuetifyx"
	"github.com/sunfmin/reflectutils"
	h "github.com/theplant/htmlgo"
	"gorm.io/gorm"

	"github.com/naokij/qor5boot/models"
)

func configUser(b *presets.Builder, ab *activity.Builder, db *gorm.DB, loginSessionBuilder *plogin.SessionBuilder) {
	user := b.Model(&models.User{})
	defer func() { ab.RegisterModel(user) }()

	user.Listing().SearchFunc(func(ctx *web.EventContext, params *presets.SearchParams) (result *presets.SearchResult, err error) {
		u := getCurrentUser(ctx.R)
		qdb := db

		if currentRoles := u.GetRoles(); !slices.Contains(currentRoles, models.RoleAdmin) {
			qdb = db.Joins("inner join user_role_join urj on users.id = urj.user_id inner join roles r on r.id = urj.role_id").
				Group("users.id").
				Having("COUNT(CASE WHEN r.name in (?) THEN 1 END) = 0", []string{models.RoleAdmin, models.RoleManager})
		}

		for i, condition := range params.SQLConditions {
			if condition.Query == "(id::text) IN (?)" {
				params.SQLConditions[i].Query = "(users.id::text) IN (?)"
			}
		}

		return gorm2op.DataOperator(qdb).Search(ctx, params)
	})

	cl := user.Listing("ID", "Name", "Account", "Status", "Roles").
		SearchColumns("name", "account").
		PerPage(10)

	cl.FilterDataFunc(func(ctx *web.EventContext) vx.FilterData {
		return []*vx.FilterItem{
			{
				Key:          "name",
				Label:        "Name",
				ItemType:     vx.ItemTypeString,
				SQLCondition: `name %s ?`,
			},
			{
				Key:          "account",
				Label:        "Account",
				ItemType:     vx.ItemTypeString,
				SQLCondition: `account %s ?`,
			},
			{
				Key:      "status",
				Label:    "Status",
				ItemType: vx.ItemTypeSelect,
				Options: []*vx.SelectItem{
					{Text: "Active", Value: models.StatusActive},
					{Text: "Inactive", Value: models.StatusInactive},
				},
				SQLCondition: `status %s ?`,
			},
		}
	})

	cl.FilterTabsFunc(func(ctx *web.EventContext) []*presets.FilterTab {
		msgr := i18n.MustGetModuleMessages(ctx.R, I18nExampleKey, Messages_zh_CN).(*Messages)

		return []*presets.FilterTab{
			{
				Label: msgr.FilterTabsAll,
				ID:    "all",
				Query: url.Values{"all": []string{"1"}},
			},
			{
				Label: msgr.FilterTabsActive,
				ID:    "active",
				Query: url.Values{"status": []string{models.StatusActive}},
			},
		}
	})

	cl.Field("Roles").ComponentFunc(func(obj interface{}, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		u := obj.(*models.User)
		var roles []string
		for _, r := range u.Roles {
			roles = append(roles, r.Name)
		}
		return h.Td(h.Text(strings.Join(roles, ", ")))
	})

	cl.Field("Status").ComponentFunc(func(obj interface{}, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		u := obj.(*models.User)
		var color string
		switch u.Status {
		case models.StatusActive:
			color = "green"
		case models.StatusInactive:
			color = "red"
		}
		return h.Td(v.VChip(h.Text(u.Status)).Color(color))
	})

	ed := user.Editing("Type", "Name", "OAuthProvider", "OAuthIdentifier", "Account", "Status", "Roles", "Company")

	ed.Field("Type").ComponentFunc(func(obj interface{}, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		u := obj.(*models.User)
		if u.ID == 0 {
			return nil
		}

		var accountType string
		if u.IsOAuthUser() {
			accountType = "OAuth Account"
		} else {
			accountType = "Main Account"
		}

		return h.Div(
			v.VRow(
				v.VCol(
					h.Text(accountType),
				).Class("text-left deep-orange--text"),
			),
		).Class("mb-2")
	})

	ed.Field("OAuthProvider").Label("OAuth Provider").ComponentFunc(func(obj interface{}, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		u := obj.(*models.User)
		if !u.IsOAuthUser() && u.ID != 0 {
			return nil
		} else {
			return v.VSelect().Attr(web.VField(field.Name, field.Value(obj))...).
				Label(field.Label).
				Items(models.OAuthProviders)
		}
	})

	ed.Field("OAuthIdentifier").Label("OAuth Identifier").ComponentFunc(func(obj interface{}, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		u := obj.(*models.User)
		if !u.IsOAuthUser() {
			return nil
		} else {
			return v.VTextField().Attr(web.VField(field.Name, field.Value(obj))...).Label(field.Label).ErrorMessages(field.Errors...).Disabled(true)
		}
	})

	ed.Field("Status").ComponentFunc(func(obj interface{}, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		return v.VSelect().Attr(web.VField(field.Name, field.Value(obj))...).
			Label(field.Label).
			Items([]string{"active", "inactive"})
	})

	ed.Field("Roles").ComponentFunc(func(obj interface{}, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		selectedItems := []v.DefaultOptionItem{}
		var values []string
		u, ok := obj.(*models.User)
		if ok {
			var roles []role.Role
			db.Model(u).Association("Roles").Find(&roles)
			for _, r := range roles {
				values = append(values, fmt.Sprint(r.ID))
				selectedItems = append(selectedItems, v.DefaultOptionItem{
					Text:  r.Name,
					Value: fmt.Sprint(r.ID),
				})
			}
		}

		var roles []role.Role
		db.Find(&roles)
		allRoleItems := []v.DefaultOptionItem{}
		for _, r := range roles {
			allRoleItems = append(allRoleItems, v.DefaultOptionItem{
				Text:  r.Name,
				Value: fmt.Sprint(r.ID),
			})
		}

		return v.VAutocomplete().Label(field.Label).Chips(true).
			Items(allRoleItems).ItemTitle("text").ItemValue("value").
			Multiple(true).Attr(web.VField(field.Name, values)...).
			ErrorMessages(field.Errors...).
			Disabled(field.Disabled)
	}).SetterFunc(func(obj interface{}, field *presets.FieldContext, ctx *web.EventContext) (err error) {
		u, ok := obj.(*models.User)
		if !ok {
			return
		}
		if u.GetAccountName() == loginInitialUserEmail {
			return perm.PermissionDenied
		}
		rids := ctx.R.Form[field.Name]
		var roles []role.Role
		for _, id := range rids {
			uid, err1 := strconv.Atoi(id)
			if err1 != nil {
				continue
			}
			roles = append(roles, role.Role{
				Model: gorm.Model{ID: uint(uid)},
			})
		}

		if u.ID == 0 {
			err = reflectutils.Set(obj, field.Name, roles)
		} else {
			err = db.Model(u).Association(field.Name).Replace(roles)
		}
		if err != nil {
			return
		}
		return
	})

	dp := user.Detailing("ID", "Name", "Account", "Status", "Roles", "Company", "CreatedAt", "UpdatedAt").
		Drawer(true)

	dp.Field("Status").ComponentFunc(func(obj interface{}, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		u := obj.(*models.User)
		var color string
		switch u.Status {
		case models.StatusActive:
			color = "green"
		case models.StatusInactive:
			color = "red"
		}
		return v.VChip(h.Text(u.Status)).Color(color)
	})

	dp.Field("Roles").ComponentFunc(func(obj interface{}, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		u := obj.(*models.User)
		var roles []string
		for _, r := range u.Roles {
			roles = append(roles, r.Name)
		}
		return h.Text(strings.Join(roles, ", "))
	})

	dp.Field("CreatedAt").ComponentFunc(func(obj interface{}, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		u := obj.(*models.User)
		return h.Text(u.CreatedAt.Format(time.RFC3339))
	})

	dp.Field("UpdatedAt").ComponentFunc(func(obj interface{}, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		u := obj.(*models.User)
		return h.Text(u.UpdatedAt.Format(time.RFC3339))
	})

}
