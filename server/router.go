package server

import (
	"os"
	"log"
	"net/http"
	"plugin"
	"github.com/gin-gonic/gin"
)

type Router struct {
	*gin.Engine
	config *Config
}

func NewRouter() *Router {
	var router = new(Router)
	router.Engine = gin.New()
	router.Engine.Use(gin.ErrorLogger())
	router.Engine.Use(gin.Logger())
	router.Engine.Use(router.configHandler())
	router.Engine.Static("/static", "static")
	return router
}

func (router *Router) LoadConfig(filename string) {
	router.config = NewConfig()

	var homeDir, _ = os.LookupEnv("HOME")
	router.config.SetDefaults(Repository{
		PublicKeyPath: homeDir + "/.ssh/id_rsa.pub",
		PrivateKeyPath: homeDir + "/.ssh/id_rsa",
		RemoteName: "origin",
		RemoteBranch: "master",
		LocalBranch: "master",
	})

	var err = router.config.Load(filename)
	if err != nil {
		log.Fatalf("Cannot load configuration file '%s'\n", filename)
	}
}

func (router *Router) LoadPlugins() {
	for path, url := range router.config.Plugins {
		p, err := plugin.Open(router.config.PluginsDir + "/" + path + ".so")
		if err != nil {
			log.Printf("Cannot load plugin '%s'\n", path)
			continue
		}

		handler, err := p.Lookup("Handler")
		if err != nil {
			log.Printf("Cannot lookup Handler in plugin '%s'\n", path)
			continue
		}

		router.AddHook(url, router.attachPlugin(handler.(func(*gin.Context) *string)))
	}
}

func (router *Router) attachPlugin(handler func(*gin.Context) *string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var url = handler(ctx)
		if url != nil {
			ctx.Status(router.PullRepo(ctx, *url))
		}
	}
}

func (router *Router) configHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Set("config", router.config)
	}
}

func (router *Router) Home(handler gin.HandlerFunc) {
	router.Engine.GET("/", handler)
	router.Engine.HEAD("/", handler)
}

func (router *Router) AddHook(url string, handler gin.HandlerFunc) {
	router.Engine.GET(url, handler)
	router.Engine.POST(url, handler)
}


func (router *Router) PullRepo(ctx *gin.Context, url string) int {
	if router.config == nil {
		// No configuration available
		return http.StatusInternalServerError
	}

	var repo = router.config.FindRepo(url)
	if repo == nil {
		// No match for repository
		return http.StatusNotFound
	}

	if ctx.Request.Method != "POST" {
		// Nothing to do, OK
		return http.StatusOK
	}

	// Pull can be long: use goroutine and respond Accepted
	go repo.Pull()
	return http.StatusAccepted
}
