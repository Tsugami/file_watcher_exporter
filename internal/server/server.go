package server

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	_ "github.com/joho/godotenv/autoload"
)

type Server struct {
	config *Config
	server *http.Server
}

func NewServer(config *Config) *Server {
	srv := &Server{
		config: config,
	}

	srv.server = &http.Server{
		Addr:         fmt.Sprintf(":%d", config.Port),
		Handler:      srv.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return srv
}

func (s *Server) ListenAndServe() error {
	monitor := NewMonitor(s.config)
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()

		fmt.Printf("Iniciando servidor HTTP em :%d\n", s.config.Port)
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Erro ao iniciar o servidor HTTP: %v\n", err)
		}
	}()

	RegisterMetrics()

	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			monitor.CheckFilesRecursively()
			time.Sleep(15 * time.Second)
		}
	}()

	wg.Wait()

	return nil
}
