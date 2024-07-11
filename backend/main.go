package main

import (
	"embed"
	"fmt"
	"github.com/diggerhq/digger/backend/bootstrap"
	"github.com/diggerhq/digger/backend/ci_backends"
	"github.com/diggerhq/digger/backend/config"
	"github.com/diggerhq/digger/backend/controllers"
        "github.com/diggerhq/digger/backend/middleware"
        "github.com/diggerhq/digger/backend/models"
        "github.com/diggerhq/digger/backend/utils"

)

//go:embed templates
var templates embed.FS

func main() {
	ghController := controllers.DiggerController{
		CiBackendProvider:    ci_backends.DefaultBackendProvider{},
		GithubClientProvider: utils.DiggerGithubRealClientProvider{},
	}
	r := bootstrap.Bootstrap(templates, ghController)
	cfg := config.New()
	web := controllers.WebController{Config: cfg}

	projectsGroup := r.Group("/projects")
	projectsGroup.Use(middleware.GetWebMiddleware())
	projectsGroup.GET("/", web.ProjectsPage)
	projectsGroup.GET("/:projectid/details", web.ProjectDetailsPage)
	projectsGroup.POST("/:projectid/details", web.ProjectDetailsUpdatePage)

	runsGroup := r.Group("/runs")
	runsGroup.Use(middleware.GetWebMiddleware())
	runsGroup.GET("/", web.RunsPage)
	runsGroup.GET("/:runid/details", web.RunDetailsPage)

	reposGroup := r.Group("/repos")
	reposGroup.Use(middleware.GetWebMiddleware())
	reposGroup.GET("/", web.ReposPage)

	repoGroup := r.Group("/repo")
	repoGroup.Use(middleware.GetWebMiddleware())
	repoGroup.GET("/", web.ReposPage)
	repoGroup.GET("/:repoid/", web.UpdateRepoPage)
	repoGroup.POST("/:repoid/", web.UpdateRepoPage)

	policiesGroup := r.Group("/policies")
	policiesGroup.Use(middleware.GetWebMiddleware())
	policiesGroup.GET("/", web.PoliciesPage)
	policiesGroup.GET("/add", web.AddPolicyPage)
	policiesGroup.POST("/add", web.AddPolicyPage)
	policiesGroup.GET("/:policyid/details", web.PolicyDetailsPage)
	policiesGroup.POST("/:policyid/details", web.PolicyDetailsUpdatePage)

	authorized := r.Group("/")
	authorized.Use(middleware.GetApiMiddleware(), middleware.AccessLevel(models.AccessPolicyType, models.AdminPolicyType))

	admin := r.Group("/")
	admin.Use(middleware.GetApiMiddleware(), middleware.AccessLevel(models.AdminPolicyType))
	r.GET("/", controllers.Home)
	port := config.GetPort()
	r.Run(fmt.Sprintf(":%d", port))
}
