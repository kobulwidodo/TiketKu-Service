package config

import (
	"go-clean/src/lib/nsq"
	"go-clean/src/lib/redis"
	"go-clean/src/lib/sql"
	"time"
)

type Application struct {
	Meta    ApplicationMeta
	Gin     GinConfig
	SQL     sql.Config
	Redis   redis.Config
	Nsq     nsq.Config
	Workers WorkersConfig
}

type ApplicationMeta struct {
	Title       string
	Description string
	Host        string
	BasePath    string
	Version     string
}

type GinConfig struct {
	Port            string
	Mode            string
	Timeout         time.Duration
	ShutdownTimeout time.Duration
	CORS            CORSConfig
	Meta            ApplicationMeta
}

type WorkersConfig struct {
	BookingWorker WorkerConfig
}

type WorkerConfig struct {
	Topic           string
	Channel         string
	LookupdAddress  string
	ShutdownTimeout int
}

type CORSConfig struct {
	Mode string
}

func Init() Application {
	return Application{}
}
