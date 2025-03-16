package integration_test

import (
	"net/http"
	"testing"

	. "github.com/qor5/web/v3/multipartestutils"
	"github.com/theplant/gofixtures"
	"gorm.io/gorm"

	"github.com/naokij/qor5boot/admin"
	"github.com/naokij/qor5boot/models"
	"github.com/qor5/admin/v3/role"
)

var TestDB *gorm.DB

var profileData = gofixtures.Data(gofixtures.Sql(`
INSERT INTO public.users (id, created_at, updated_at, deleted_at, name, company, status, registration_date, account, password, encrypted_password, o_auth_provider, o_auth_identifier) VALUES (1, '2024-06-18 03:24:28.001791 +00:00', '2024-06-19 07:07:18.502134 +00:00', null, 'qor@theplant.jp', 'Test Company', 'active', '2024-06-18', 'qor@theplant.jp', '', '$2a$10$XKsTcchE1r1X5MyTD0k1keyUwub23DXsjSIQW73MtXfoiqrqbXAnu', '', '');
INSERT INTO public.roles (id, created_at, updated_at, deleted_at, name) VALUES (1, '2024-08-23 08:43:32.969461 +00:00', '2024-09-12 06:25:17.533058 +00:00', null, 'Admin');
INSERT INTO public.user_role_join (user_id, role_id) VALUES (1, 1);
`, []string{`user_role_join`, `roles`, "users"}))

func TestProfile(t *testing.T) {
	h := admin.TestHandler(TestDB, &models.User{
		Model: gorm.Model{ID: 1},
		Roles: []role.Role{
			{
				Name: models.RoleAdmin,
			},
		},
	})
	dbr, _ := TestDB.DB()

	cases := []TestCase{
		{
			Name:  "View profile",
			Debug: true,
			ReqFunc: func() *http.Request {
				profileData.TruncatePut(dbr)
				req := NewMultipartBuilder().
					PageURL("/profile").
					BuildEventFuncRequest()
				return req
			},
			ExpectPageBodyContainsInOrder: []string{`qor@theplant.jp`, `Test Company`},
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			RunCase(t, c, h)
		})
	}
}
