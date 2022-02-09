package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"

	"channels-instagram-dm/api"
	"channels-instagram-dm/mq"
	"channels-instagram-dm/repository"
	"channels-instagram-dm/service"
	"channels-instagram-dm/sync"
	"github.com/gorilla/mux"
)

type Config struct {
	AppPort   string
	ProfPort  string
	DBHost    string
	DBName    string
	SlotsURI  string
	MQHost    string
	MQCluster string
	MQClient  string
}

func main() {
	cfg := Config{
		AppPort:   os.Getenv("APPLICATION_PORT"),
		ProfPort:  os.Getenv("PROF_PORT"),
		DBHost:    os.Getenv("DB_HOST"),
		DBName:    os.Getenv("DB_NAME"),
		SlotsURI:  os.Getenv("SLOTS_URI"),
		MQHost:    os.Getenv("MQ_HOST"),
		MQCluster: os.Getenv("MQ_CLUSTER"),
		MQClient:  os.Getenv("MQ_CLIENT"),
	}

	if cfg.AppPort == "" {
		log.Fatal("Environment variable 'APPLICATION_PORT' should not be empty")
	}

	if cfg.DBHost == "" {
		log.Fatal("Environment variable 'DB_HOST' should not be empty")
	}

	if cfg.DBName == "" {
		log.Fatal("Environment variable 'DB_NAME' should not be empty")
	}

	if cfg.SlotsURI == "" {
		log.Fatal("Environment variable 'SLOTS_URI' should not be empty")
	}

	if cfg.MQHost == "" {
		log.Fatal("Environment variable 'MQ_HOST' should not be empty")
	}

	if cfg.MQCluster == "" {
		log.Fatal("Environment variable 'MQ_CLUSTER' should not be empty")
	}

	if cfg.MQClient == "" {
		log.Fatal("Environment variable 'MQ_CLIENT' should not be empty")
	}

	mainContext, mainCancel := context.WithCancel(context.Background())

	logger := NewLogger(os.Stdout, "")

	syncer := Syncer()

	eventBus := EventBus()

	repositoryFactory, err := repository.Factory(
		mainContext,
		logger.Copy("REPOSITORY"),
		cfg.DBHost,
		cfg.DBName,
	)
	if err != nil {
		panic(err)
	}

	serviceFactory, err := service.Factory(
		mainContext,
		logger.Copy("IG_SERVICE"),
		cfg.SlotsURI,
	)
	if err != nil {
		panic(err)
	}

	mqFactory, err := mq.Factory(
		mainContext,
		logger.Copy("MQ"),
		cfg.MQHost,
		cfg.MQCluster,
		cfg.MQClient,
	)
	if err != nil {
		panic(err)
	}

	runtimeContext := api.RuntimeContext(
		mainContext,
		repositoryFactory,
		serviceFactory,
		logger,
		mqFactory,
		eventBus,
		syncer,
	)

	router := mux.NewRouter()
	api.InitRoutes(
		runtimeContext.WithLogger(
			runtimeContext.Logger().Copy("API"),
		),
		router,
	)

	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			c := context.WithValue(req.Context(), "session", api.NewSession(req))
			next.ServeHTTP(w, req.WithContext(c))
		})
	})

	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.AppPort),
		Handler: router,
	}

	go func() {
		if err := httpServer.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	sync.Run(
		runtimeContext.WithLogger(
			runtimeContext.Logger().Copy("SYNC"),
		),
	)

	profiling(cfg.ProfPort)

	logger.Info("Server was started", nil)

	// Shutting down
	c := make(chan os.Signal, 0)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	mainCancel()

	select {
	case <-runtimeContext.Syncer().Done():
		logger.Info("Server stopped as done", nil)
	case <-time.After(30 * time.Second):
		logger.Info("Server stopped by timeout", nil)
	}

}

func profiling(profPort string) {
	if profPort == "" {
		return
	}

	go func() {
		r := http.NewServeMux()

		r.HandleFunc("/debug/pprof/", pprof.Index)
		r.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
		r.HandleFunc("/debug/pprof/profile", pprof.Profile)
		r.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
		r.HandleFunc("/debug/pprof/trace", pprof.Trace)

		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", profPort), r))
	}()
}
