package main

import (
	"bytes"
	"os"

	"github.com/common-nighthawk/go-figure"
	"github.com/dimiro1/banner"
	"github.com/jmoiron/sqlx"
	"github.com/minebarteksa/clean-show/config"
	"github.com/minebarteksa/clean-show/infrastructure/database"
	"github.com/minebarteksa/clean-show/infrastructure/security"
	. "github.com/minebarteksa/clean-show/logger"
	"github.com/minebarteksa/clean-show/registry"
)

func init() {
	config.LoadEnvConfig()
}

func main() {
	figure.NewColorFigure("Go Clean Show", "", "green", true).Print()
	banner.Init(os.Stdout, true, true, bytes.NewBufferString(bannerContent))

	InitProduction()
	if config.Env.Debug {
		InitDebug()
		Log.Infow("Running in Debug")
	}

	db, err := sqlx.Connect(config.Env.DBDriver, config.Env.DBSource)
	if err != nil {
		Log.Fatalw("failed to connect to the database", "err", err)
	}

	err = db.Ping()
	if err != nil {
		Log.Fatalw("failed to ping databse", "err", err)
	}
	sql := database.NewSqlDB(db)
	Log.Infow("Connected to the SQL Database")

	hasher := security.NewArgon2idHasher()

	r := registry.NewRegistry(sql, hasher)
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
