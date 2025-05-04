package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"voice-based-tracker/pkg/audio"
	"voice-based-tracker/pkg/fileops"

	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"
)

const (
	fileName                   = "audio-text"
	defaultHTTPServerAddress   = ":8080"
	timeoutForGracefulShutdown = 30 * time.Second
)

func main() {
	handler := NewVoiceHandler()
	router := NewRouter(handler)

	RunHTTPServer(context.Background(), router)
}

func RunHTTPServer(
	ctx context.Context,
	router *chi.Mux,
) error {
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	server := &http.Server{
		Addr:    defaultHTTPServerAddress,
		Handler: router,
	}

	log.Printf("starting http server, address: %s\n", server.Addr)

	cCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	var srvErr error

	go func() {
		srvErr = server.ListenAndServe()

		defer cancel()

		log.Println("http server stopped, address")
	}()

	<-cCtx.Done()

	log.Println("terminal signal received or unexpected error happened")

	if srvErr != nil {
		if errors.Is(srvErr, http.ErrServerClosed) {
			return nil
		}

		return errors.Wrap(srvErr, "http server errored")
	}

	shutdownCtx, shutdownCancel := context.WithTimeout(
		context.Background(),
		timeoutForGracefulShutdown,
	)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		return errors.Wrap(err, "http server shutdown failed")
	}

	log.Println("http server shutdown done")

	return nil
}

type VoiceHandler struct {
	stopCh    chan bool
	samplesCh chan int

	aw                   *os.File
	ac                   audio.AudioCapturer
	missingSamplesFiller func(int)
}

func NewVoiceHandler() *VoiceHandler {
	stopCh := make(chan bool)
	samplesCh := make(chan int)

	aw, missingSamplesFiller := fileops.NewAudioWriter(fileName)
	ac := audio.NewAudioCapture()

	return &VoiceHandler{
		stopCh:               stopCh,
		samplesCh:            samplesCh,
		aw:                   aw,
		ac:                   ac,
		missingSamplesFiller: missingSamplesFiller,
	}
}

func (vh *VoiceHandler) StartRecording(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	_ = ctx

	err := vh.ac.StartRecording(vh.aw, vh.stopCh, vh.samplesCh)
	if err != nil {
		log.Printf("failed to start recording: %v\n", err)
	}
}

func (vh *VoiceHandler) StopRecording(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	_ = ctx

	err := vh.ac.StopRecording(vh.stopCh, vh.samplesCh, vh.missingSamplesFiller)
	if err != nil {
		log.Printf("failed to stop recording: %v\n", err)
	}
}

func NewRouter(vh *VoiceHandler) *chi.Mux {
	router := chi.NewRouter()

	router.Group(func(r chi.Router) {
		r.Route("/recording", func(r chi.Router) {
			r.Get("/start", vh.StartRecording)
			r.Get("/stop", vh.StopRecording)
		})
	})

	return router
}
