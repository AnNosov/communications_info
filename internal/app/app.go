package app

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/AnNosov/communications_info/config"
	v1 "github.com/AnNosov/communications_info/internal/controller/http/v1"
	"github.com/AnNosov/communications_info/internal/usecase"
	"github.com/AnNosov/communications_info/pkg/http/server"

	"github.com/gorilla/mux"
)

func Run(cfg *config.Config) {

	cAction := usecase.New(cfg)

	handler := mux.NewRouter()
	v1.NewRouter(handler, *cAction)
	httpServer := server.New(handler, server.Port(cfg.HttpServer.Port))

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM, syscall.SIGINT) // у меня работает только с SIGINT (ctrl + c)

	select {
	case s := <-interrupt:
		log.Println("app - Run - signal: " + s.String())
	case err := <-httpServer.Notify():
		log.Println("app - Run - httpServer.Notify: ", err.Error())
	}

	err := httpServer.Shutdown()
	if err != nil {
		log.Println("app - Run - httpServer.Shutdown: ", err.Error())
	}
}
