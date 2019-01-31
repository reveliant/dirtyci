package main

import (
	"flag"
	"github.com/reveliant/dirtyci/server"
)

func main() {
	var mode = "release"
	var help = flag.Bool("h", false, "Show this help message")
	var debug = flag.Bool("d", false, "Enable debug mode")
	var filename = flag.String("c", "config.toml", "Configuration file path")
	flag.Parse()

	if *help {
		flag.PrintDefaults()
		return
	}

	if *debug {
		mode = "debug"
	}
	server.SetMode(mode)

	var router = server.NewRouter()
	router.LoadConfig(*filename)
	router.LoadPlugins()
	router.Home(server.Redirect("https://github.com/reveliant/dirty-ci"))
	router.Run()
}
