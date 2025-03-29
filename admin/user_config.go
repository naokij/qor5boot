package admin

import (
	"fmt"
	"net/url"
	"slices"
	"strconv"
	"strings"

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

		// 处理已删除用户的查询
		// 首先检查 filter_deleted 参数（从过滤器来的）
		showDeleted := false
		if ctx.R.FormValue("f_deleted") == "1" {
			showDeleted = true
		}
		// 然后检查 deleted 参数（从标签页来的）
		if ctx.R.FormValue("deleted") == "1" {
			showDeleted = true
		}

		// 无论是否查询已删除用户，都预加载 Roles 关联
		qdb = qdb.Preload("Roles")

		if showDeleted {
			qdb = qdb.Unscoped().Where("deleted_at IS NOT NULL")
		} else {
			qdb = qdb.Where("deleted_at IS NULL")
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
		msgr := i18n.MustGetModuleMessages(ctx.R, I18nAdminKey, Messages_zh_CN).(*Messages)

		return []*vx.FilterItem{
			{
				Key:          "name",
				Label:        msgr.Name,
				ItemType:     vx.ItemTypeString,
				SQLCondition: `name %s ?`,
			},
			{
				Key:          "account",
				Label:        msgr.Account,
				ItemType:     vx.ItemTypeString,
				SQLCondition: `account %s ?`,
			},
			{
				Key:      "status",
				Label:    msgr.Status,
				ItemType: vx.ItemTypeSelect,
				Options: []*vx.SelectItem{
					{Text: msgr.UserStatusActive, Value: models.StatusActive},
					{Text: msgr.UserStatusInactive, Value: models.StatusInactive},
				},
				SQLCondition: `status %s ?`,
			},
			{
				Key:      "deleted",
				Label:    msgr.UserDeletedFilter,
				ItemType: vx.ItemTypeSelect,
				Options: []*vx.SelectItem{
					{Text: msgr.UserDeletedYes, Value: "1"},
					{Text: msgr.UserDeletedNo, Value: "0"},
				},
				// 这个条件在 SearchFunc 中单独处理，这里不需要设置 SQLCondition
			},
		}
	})

	cl.FilterTabsFunc(func(ctx *web.EventContext) []*presets.FilterTab {
		msgr := i18n.MustGetModuleMessages(ctx.R, I18nAdminKey, Messages_zh_CN).(*Messages)

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
			{
				Label: msgr.DeletedUsersTab,
				ID:    "deleted",
				Query: url.Values{"deleted": []string{"1"}},
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
		var text string
		var color string
		msgr := i18n.MustGetModuleMessages(ctx.R, I18nAdminKey, Messages_zh_CN).(*Messages)

		// 不区分大小写比较
		statusLower := strings.ToLower(u.Status)
		switch statusLower {
		case strings.ToLower(models.StatusActive):
			text = msgr.UserStatusActive
			color = "success"
		case strings.ToLower(models.StatusInactive):
			text = msgr.UserStatusInactive
			color = "error"
		default:
			if u.Status == "" {
				text = msgr.UserStatusUnset
				color = "grey"
			} else {
				text = u.Status
				color = "grey"
			}
		}

		return h.Td(v.VChip(h.Text(text)).Color(color))
	})

	// 为每行添加操作按钮
	cl.Field("ID").ComponentFunc(func(obj interface{}, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		u := obj.(*models.User)
		id := fmt.Sprintf("%d", u.ID)
		msgr := i18n.MustGetModuleMessages(ctx.R, I18nAdminKey, Messages_zh_CN).(*Messages)

		// 是否在查看已删除用户
		isDeletedView := ctx.R.FormValue("f_deleted") == "1"

		// 如果不是已删除用户视图，只显示ID
		if !isDeletedView {
			return h.Td(h.Text(id))
		}

		// 创建标准ID列
		idCell := h.Td(h.Text(id))

		// 如果是初始管理员用户，不显示恢复/删除按钮
		if u.GetAccountName() == loginInitialUserEmail {
			return idCell
		}

		var buttons []h.HTMLComponent

		// 在已删除视图添加恢复按钮
		buttons = append(buttons, v.VBtn(msgr.RestoreUserBtn).
			Size("small").
			Color("success").
			Attr("@click", web.Plaid().
				EventFunc("restore_user").
				Query("id", id).
				Go()).
			Class("mr-2"))

		// 在已删除视图添加永久删除按钮
		buttons = append(buttons, v.VBtn(msgr.PermanentDeleteBtn).
			Size("small").
			Color("error").
			Attr("@click", web.Plaid().
				EventFunc("permanent_delete_user").
				Query("id", id).
				Go()))

		return h.Td(
			h.Div(
				h.Text(id),
				h.Div(buttons...).Class("mt-2"),
			),
		)
	})

	// 注册恢复用户事件
	user.RegisterEventFunc("restore_user", func(ctx *web.EventContext) (r web.EventResponse, err error) {
		msgr := i18n.MustGetModuleMessages(ctx.R, I18nAdminKey, Messages_zh_CN).(*Messages)
		id := ctx.R.FormValue("id")
		if id == "" {
			ctx.Flash = msgr.NoUserID
			r.Reload = true
			return
		}

		// 使用Unscoped恢复被软删除的用户
		if err = db.Unscoped().Model(&models.User{}).Where("id = ?", id).Update("deleted_at", nil).Error; err != nil {
			ctx.Flash = msgr.RestoreUserFailed + err.Error()
		} else {
			ctx.Flash = msgr.RestoreUserSuccess
		}

		r.Reload = true
		return
	})

	// 注册永久删除用户事件
	user.RegisterEventFunc("permanent_delete_user", func(ctx *web.EventContext) (r web.EventResponse, err error) {
		msgr := i18n.MustGetModuleMessages(ctx.R, I18nAdminKey, Messages_zh_CN).(*Messages)
		id := ctx.R.FormValue("id")
		if id == "" {
			ctx.Flash = msgr.NoUserID
			r.Reload = true
			return
		}

		// 永久删除用户
		if err = db.Unscoped().Delete(&models.User{}, "id = ?", id).Error; err != nil {
			ctx.Flash = msgr.PermanentDeleteFailed + err.Error()
		} else {
			ctx.Flash = msgr.PermanentDeleteSuccess
		}

		r.Reload = true
		return
	})

	// 注册软删除用户事件
	user.RegisterEventFunc("delete_user", func(ctx *web.EventContext) (r web.EventResponse, err error) {
		msgr := i18n.MustGetModuleMessages(ctx.R, I18nAdminKey, Messages_zh_CN).(*Messages)
		id := ctx.R.FormValue("id")
		if id == "" {
			ctx.Flash = msgr.NoUserID
			r.Reload = true
			return
		}

		uid, err1 := strconv.Atoi(id)
		if err1 != nil {
			ctx.Flash = msgr.InvalidUserObject
			r.Reload = true
			return
		}

		// 检查是否为初始管理员用户，不允许删除
		var user models.User
		if err = db.First(&user, uid).Error; err != nil {
			ctx.Flash = msgr.UserNotFound + ": " + err.Error()
			r.Reload = true
			return
		}

		if user.GetAccountName() == loginInitialUserEmail {
			ctx.Flash = msgr.CantDeleteAdminUser
			r.Reload = true
			return
		}

		// 软删除用户
		if err = db.Delete(&models.User{}, uid).Error; err != nil {
			ctx.Flash = msgr.DeleteUserFailed + err.Error()
		} else {
			ctx.Flash = msgr.DeleteUserSuccess
		}

		r.Reload = true
		return
	})

	ed := user.Editing("Type", "Name", "OAuthProvider", "OAuthIdentifier", "Account", "Password", "Status", "Roles", "Company")

	// 使用FetchFunc来阻止编辑已删除用户
	ed.FetchFunc(func(obj interface{}, id string, ctx *web.EventContext) (interface{}, error) {
		msgr := i18n.MustGetModuleMessages(ctx.R, I18nAdminKey, Messages_zh_CN).(*Messages)

		if id == "" {
			// 创建新用户时设置默认值
			u, ok := obj.(*models.User)
			if ok {
				// 设置默认状态为活跃
				if u.Status == "" {
					u.Status = models.StatusActive
				}
			}
			return obj, nil
		}

		// 使用Unscoped查询包括已删除的用户
		u := &models.User{}
		if err := db.Unscoped().First(u, id).Error; err != nil {
			return nil, fmt.Errorf(msgr.UserNotFound+": %v", err)
		}

		return u, nil
	})

	// 使用SaveFunc处理保存逻辑
	ed.SaveFunc(func(obj interface{}, id string, ctx *web.EventContext) (err error) {
		msgr := i18n.MustGetModuleMessages(ctx.R, I18nAdminKey, Messages_zh_CN).(*Messages)

		// 执行原始的表单处理逻辑
		u, ok := obj.(*models.User)
		if !ok {
			return fmt.Errorf(msgr.InvalidUserObject)
		}
		if u.DeletedAt.Valid {
			return fmt.Errorf(msgr.CantEditDeletedUser)
		}

		// 确保状态字段有值
		if u.Status == "" {
			u.Status = models.StatusActive
		}

		// 保存到数据库
		if id == "" {
			// 创建新用户
			if err := db.Create(u).Error; err != nil {
				return err
			}
		} else {
			// 更新现有用户
			if err := db.Save(u).Error; err != nil {
				return err
			}
		}

		return nil
	})

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

	ed.Field("Password").ComponentFunc(func(obj interface{}, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		//密码控件 - 回到最简单的实现方式
		msgr := i18n.MustGetModuleMessages(ctx.R, I18nAdminKey, Messages_zh_CN).(*Messages)
		return v.VTextField().
			Attr(web.VField(field.Name, "")...).
			Label(field.Label).
			Type("password").
			Placeholder(msgr.PasswordPlaceholder).
			Hint(msgr.PasswordMinLengthHint).
			ErrorMessages(field.Errors...)
	}).SetterFunc(func(obj interface{}, field *presets.FieldContext, ctx *web.EventContext) (err error) {
		u := obj.(*models.User)
		if v := ctx.R.FormValue(field.Name); v != "" {
			// 验证密码长度至少为12位
			if len(v) < 12 {
				msgr := i18n.MustGetModuleMessages(ctx.R, I18nAdminKey, Messages_zh_CN).(*Messages)
				return fmt.Errorf(msgr.PasswordMinLengthError)
			}
			u.Password = v
			u.EncryptPassword()
		}
		return nil
	})

	ed.Field("Status").ComponentFunc(func(obj interface{}, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		u := obj.(*models.User)
		// 如果是新用户且状态为空，则默认设置为活跃
		if u.ID == 0 && u.Status == "" {
			u.Status = models.StatusActive
		}

		return v.VSelect().Attr(web.VField(field.Name, field.Value(obj))...).
			Label(field.Label).
			Items([]string{models.StatusActive, models.StatusInactive})
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

		// 允许修改 admin@admin.com 用户，但不允许删除其 Admin 角色
		if u.GetAccountName() == loginInitialUserEmail {
			// 检查是否移除了 Admin 角色
			rids := ctx.R.Form[field.Name]
			hasAdminRole := false

			// 查找所有角色
			var roles []role.Role
			if err = db.Find(&roles).Error; err != nil {
				return err
			}

			// 找到 Admin 角色的 ID
			var adminRoleID uint
			for _, r := range roles {
				if r.Name == models.RoleAdmin {
					adminRoleID = r.ID
					break
				}
			}

			// 检查提交的角色中是否包含 Admin 角色
			for _, id := range rids {
				uid, err1 := strconv.Atoi(id)
				if err1 != nil {
					continue
				}
				if uint(uid) == adminRoleID {
					hasAdminRole = true
					break
				}
			}

			// 如果移除了 Admin 角色，则拒绝请求
			if !hasAdminRole {
				return perm.PermissionDenied
			}
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

}
