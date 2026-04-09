package app

import (
	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
)

func setupStaticAndTemplates(r *gin.Engine) {
	r.HTMLRender = createRender()
	r.Static("/static", "static")
}

// загружаем шаблоны страниц
func createRender() multitemplate.Renderer {
	r := multitemplate.NewRenderer()
	// загрузка цельных страниц
	r.AddFromFiles("index", "templates/base.html", "templates/index.html")
	r.AddFromFiles("register", "templates/base.html", "templates/auth/register.html")
	r.AddFromFiles("profile", "templates/base.html", "templates/profile.html")

	// загрузка частей страниц
	r.AddFromFiles("verify", "templates/partials/verify-code-form.html")
	r.AddFromFiles("auth-error", "templates/partials/auth-error.html")
	r.AddFromFiles("generate", "templates/partials/route-result.html")

	// админка
	r.AddFromFiles("admin-tables", "templates/base.html", "templates/admin/tables.html")
	r.AddFromFiles("admin-view", "templates/base.html", "templates/admin/view.html")
	r.AddFromFiles("admin-edit", "templates/base.html", "templates/admin/edit.html")
	r.AddFromFiles("admin-create", "templates/base.html", "templates/admin/create.html")
	r.AddFromFiles("admin-place-create", "templates/base.html", "templates/admin/place-create.html")

	return r
}
