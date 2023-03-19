package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type EnvConfig struct {
	Debug    bool
	Port     uint16
	DBDriver string
	DBSource string
}

var Env *EnvConfig

func LoadEnvConfig() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	port, err := strconv.ParseUint(os.Getenv("PORT"), 10, 16)
	if err != nil {
		panic(err)
	}

	Env = &EnvConfig{
		Debug:    os.Getenv("DEBUG") == "true",
		Port:     uint16(port),
		DBDriver: os.Getenv("DB_DRIVER"),
		DBSource: os.Getenv("DB_SOURCE"),
	}
}
