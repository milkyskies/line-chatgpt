package http

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/milkyskies/line-chatgpt/internal/linebot"
)

type Handler struct {
	Router *mux.Router
	Server *http.Server
}

func NewHandler() *Handler {
	h := &Handler{
		Router: mux.NewRouter(),
	}

	h.mapRoutes()
	h.Router.Use(JSONMiddleware)

	h.Server = &http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: h.Router,
	}

	return h
}

func (h *Handler) mapRoutes() {
	h.Router.HandleFunc("/alive", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "I am alive")
	})

	h.Router.HandleFunc("/line", func(w http.ResponseWriter, r *http.Request) {
		bot, err := linebot.NewBot()
		if err != nil {
			fmt.Println("Error creating LINE bot")
			return
		}

		bot.HandleRequest(r)
	})
}

func (h *Handler) Serve() error {
	go func() {
		if err := h.Server.ListenAndServe(); err != nil {
			log.Println(err.Error())
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	h.Server.Shutdown(ctx)

	log.Println("Shut down gracefully")
	return nil
}
