package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func startHttpServer(idx int, addr string, message string) *http.Server {
	log.Printf("Starting server %d on %s\n", idx, addr)
	engine := gin.Default()
	engine.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, message)
	})
	httpServer := http.Server{
		Addr:    addr,
		Handler: engine,
	}
	go func() {
		if err := httpServer.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()
	return &httpServer
}

func waitForShutdown(servers []*http.Server, duration time.Duration) {
	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
	<-done
	log.Println("Shutdown initiated")
	timedContext, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()
	wg := sync.WaitGroup{}
	wg.Add(len(servers))
	for i, s := range servers {
		go func(idx int, server *http.Server) {
			log.Printf("Shutdown server %d\n", idx)
			server.Shutdown(timedContext)
			wg.Done()
		}(i, s)
	}
	wg.Wait()
}

func main() {
	numberOfServers := 10
	log.Printf("Starting servers %d...\n", numberOfServers)
	servers := make([]*http.Server, numberOfServers)
	for i := 0; i < numberOfServers; i++ {
		servers[i] = startHttpServer(i, fmt.Sprintf(":808%d", i), fmt.Sprintf("Hello from server %d", i))
	}
	waitForShutdown(servers, time.Minute)
	log.Println("Application shutdown")
}
