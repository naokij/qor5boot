package admin

import (
	_ "embed"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/qor5/x/v3/login"
	"github.com/qor5/x/v3/sitemap"
	"gorm.io/gorm"

	"github.com/naokij/qor5boot/models"
	"github.com/qor5/admin/v3/role"
)

//go:embed assets/favicon.ico
var favicon []byte

func TestHandlerComplex(db *gorm.DB, u *models.User, enableWork bool) (http.Handler, Config) {
	mux := http.NewServeMux()
	c := NewConfig(db, enableWork)
	if u == nil {
		u = &models.User{
			Model: gorm.Model{ID: 888},
			Roles: []role.Role{
				{
					Name: models.RoleAdmin,
				},
			},
		}
	}
	m := login.MockCurrentUser(u)
	mux.Handle("/", m(c.pb))
	return mux, c
}

func TestHandler(db *gorm.DB, u *models.User) http.Handler {
	mux, _ := TestHandlerComplex(db, u, false)
	return mux
}

func TestHandlerWorker(db *gorm.DB, u *models.User) http.Handler {
	mux, _ := TestHandlerComplex(db, u, true)
	return mux
}

func Router(db *gorm.DB) http.Handler {
	c := NewConfig(db, true)

	mux := http.NewServeMux()
	c.loginSessionBuilder.Mount(mux)
	//	mux.Handle("/frontstyle.css", c.pb.GetWebBuilder().PacksHandler("text/css", web.ComponentsPack(`
	// :host {
	//	all: initial;
	//	display: block;
	// div {
	//	background-color:orange;
	// }
	// `)))

	mux.Handle("/", c.pb)
	mux.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		w.Write(favicon)
		return
	})

	// example of sitemap and robot
	sitemap.SiteMap("product").RegisterRawString("https://dev.qor5.com/admin", "/product").MountTo(mux)
	robot := sitemap.Robots()
	robot.Agent(sitemap.AlexaAgent).Allow("/product1", "/product2").Disallow("/admin")
	robot.Agent(sitemap.GoogleAgent).Disallow("/admin")
	robot.MountTo(mux)

	cr := chi.NewRouter()
	cr.Use(
		c.loginSessionBuilder.Middleware(),
		withRoles(db),
		securityMiddleware(),
	)
	cr.Mount("/", mux)
	return cr
}
