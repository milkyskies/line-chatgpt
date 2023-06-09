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
)

type Handler struct {
	Router             *mux.Router
	Server             *http.Server
	LineWebhookHandler http.Handler
}

func NewHandler(lineWebhookHandler http.Handler) (*Handler, error) {
	h := &Handler{
		Router:             mux.NewRouter(),
		LineWebhookHandler: lineWebhookHandler,
	}

	h.mapRoutes()
	// h.Router.Use(JSONMiddleware)

	h.Server = &http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: h.Router,
	}

	return h, nil
}

func (h *Handler) mapRoutes() {
	h.Router.HandleFunc("/alive", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "I am alive")
	})

	h.Router.Handle("/line", h.LineWebhookHandler)

	h.Router.HandleFunc("/audio_replies/{audio_file:.+}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		audioFile := vars["audio_file"]
		filePath := fmt.Sprintf("content/whisper/audio_replies/%s.m4a", audioFile)

		w.Header().Set("Content-Type", "audio/mpeg")
		http.ServeFile(w, r, filePath)
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
	if err := h.Server.Shutdown(ctx); err != nil {
		return err
	}

	log.Println("Shut down gracefully")
	return nil
}
