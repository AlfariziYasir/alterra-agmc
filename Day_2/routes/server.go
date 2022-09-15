package routes

import (
	"api-mvc/config"
	"api-mvc/controller"
	"api-mvc/db"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func Start() {
	db, _ := db.NewClient()

	controller := controller.Controller{DB: db.Conn()}

	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", config.Cfg().APPPort),
		Handler: Routes(&controller),
	}

	idleConnsClosed := make(chan struct{})
	go func() {
		defer close(idleConnsClosed)

		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		signal.Notify(sigint, syscall.SIGTERM)

		<-sigint

		err := httpServer.Shutdown(context.Background())
		if err != nil {
			log.Println(err.Error())
		}
	}()

	log.Printf("\nstarting server on %s\n", httpServer.Addr)
	err := httpServer.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Println(err)
	}

	<-idleConnsClosed

	log.Println("stopped server gracefully")
}
