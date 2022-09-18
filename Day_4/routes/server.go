package routes

import (
	"api-mvc/config"
	"api-mvc/controller"
	"api-mvc/db/postgres"
	rds "api-mvc/db/redis"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func Start() {
	db, err := postgres.NewClient()
	if err != nil {
		log.Panicf("db: %v", err)
	}
	redis, err := rds.NewClient()
	if err != nil {
		log.Panicf("redis: %v", err)
	}

	controller := controller.NewController(redis.Conn(), db.Conn())

	e := Routes(&controller)

	idleConnsClosed := make(chan struct{})
	go func() {
		defer close(idleConnsClosed)

		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		signal.Notify(sigint, syscall.SIGTERM)

		<-sigint

		err := e.Shutdown(context.Background())
		if err != nil {
			log.Println(err.Error())
		}
	}()

	err = e.Start(fmt.Sprintf(":%d", config.Cfg().APPPort))
	if err != nil && err != http.ErrServerClosed {
		log.Println(err)
	}

	<-idleConnsClosed

	log.Println("stopped server gracefully")
}
