package main

import (
	"context"
	"go-clean/src/business/domain"
	"go-clean/src/business/usecase"
	"go-clean/src/handler/rest"
	"go-clean/src/handler/worker"
	"go-clean/src/lib/auth"
	"go-clean/src/lib/configreader"
	"go-clean/src/lib/log"
	"go-clean/src/lib/midtrans"
	"go-clean/src/lib/nsq"
	tracer "go-clean/src/lib/otel"
	"go-clean/src/lib/redis"
	"go-clean/src/lib/sql"
	"go-clean/src/utils/config"

	_ "go-clean/docs/swagger"

	"github.com/spf13/cobra"
)

// @contact.name   Rakhmad Giffari Nurfadhilah
// @contact.url    https://fadhilmail.tech/
// @contact.email  rakhmadgiffari14@gmail.com

// @securitydefinitions.apikey BearerAuth
// @in header
// @name Authorization

const (
	configFile string = "./etc/cfg/config.json"
)

func main() {
	cfg := config.Init()
	configReader := configreader.Init(configreader.Options{
		ConfigFile: configFile,
	})
	configReader.ReadConfig(&cfg)

	log := log.Init(log.Config{
		Level: "debug",
	})

	redis := redis.Init(cfg.Redis)

	oteltracer := tracer.Init(cfg.OtelTracer)
	defer oteltracer.Shutdown(context.Background())

	nsq := nsq.Init(cfg.Nsq)

	midtrans := midtrans.Init(cfg.Midtrans)

	db := sql.Init(cfg.SQL)

	d := domain.Init(db, redis, midtrans, log, oteltracer)

	auth := auth.Init()

	rootCmd := &cobra.Command{Use: "app"}

	restCmd := &cobra.Command{
		Use:   "rest",
		Short: "Run the REST API Server",
		Run: func(cmd *cobra.Command, args []string) {
			uc := usecase.Init(auth, d, nsq, log, oteltracer)
			r := rest.Init(cfg.Gin, uc, auth, log, oteltracer)
			r.Run()
		},
	}

	bookingWorker := &cobra.Command{
		Use:   "booking-worker",
		Short: "Run the Booking Worker",
		Run: func(cmd *cobra.Command, args []string) {
			uc := usecase.Init(auth, d, nil, log, oteltracer)
			w := worker.Init(cfg.Workers.BookingWorker, uc, log)
			w.Run()
		},
	}

	rootCmd.AddCommand(restCmd, bookingWorker)

	if err := rootCmd.Execute(); err != nil {
		log.Error(context.Background(), err)
	}
}
