package main

import (
	"bytes"
	"os"

	"github.com/common-nighthawk/go-figure"
	"github.com/dimiro1/banner"
	"github.com/minebarteksa/clean-show/config"
	"github.com/minebarteksa/clean-show/infrastructure/database"
	"github.com/minebarteksa/clean-show/logger"
	"github.com/minebarteksa/clean-show/registry"
)

func init() {
	config.LoadEnvConfig()
}

func main() {
	figure.NewColorFigure("Go Clean Show", "", "green", true).Print()
	banner.Init(os.Stdout, true, true, bytes.NewBufferString(bannerContent))

	logger.InitProduction()
	if config.Env.Debug {
		logger.InitDebug()
		logger.Log.Infow("Running in Debug")
	}

	db := database.NewSqlDB()
	logger.Log.Infow("Connected to the SQL Database")

	r := registry.NewRegistry(db)
	r.Start()
}

const bannerContent = `
GoVersion: {{ .GoVersion }}
GOOS: {{ .GOOS }}
GOARCH: {{ .GOARCH }}
NumCPU: {{ .NumCPU }}
GOPATH: {{ .GOPATH }}
GOROOT: {{ .GOROOT }}
Compiler: {{ .Compiler }}
ENV: {{ .Env "GOPATH" }}
Now: {{ .Now "Monday, 2 Jan 2006" }}
`
