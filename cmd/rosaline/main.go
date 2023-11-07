package main

import (
	"fmt"
	"os"
	"rosaline/cmd/rosaline/interfaces"

	"github.com/fred1268/go-clap/clap"
)

type config struct {
	Mode string `clap:"--mode,-m"`
}

func main() {
	cfg := &config{Mode: "cli"}

	if _, err := clap.Parse(os.Args[1:], cfg); err != nil {
		fmt.Println(err)
		return
	}

	if cfg.Mode != "cli" && cfg.Mode != "uci" {
		fmt.Println("invalid mode:", cfg.Mode)
		return
	}

	var engineInterface interfaces.EngineInterface
	if cfg.Mode == "cli" {
		engineInterface = interfaces.NewCliProtocolHandler()
	} else {
		engineInterface = interfaces.NewUciProtocolHandler()
	}

	engineInterface.Loop()
}
