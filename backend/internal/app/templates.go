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

	// загрузка частей страниц
	r.AddFromFiles("verify", "templates/partials/verify-code-form.html")
	r.AddFromFiles("auth-error", "templates/partials/auth-error.html")
	r.AddFromFiles("generate", "templates/partials/route-result.html")

	return r
}
