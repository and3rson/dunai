package main

import (
	"context"
	"github.com/and3rson/dunai/dunai"
	"log"
	"os"
	"os/signal"
	"time"
)

func main() {
	if err := dunai.LoadConfig(); err != nil {
		panic(err)
	}
	// if err := RepoInit(); err != nil {
	// 	panic(err)
	// }

	// CV
	cv, err := dunai.ReadCV()
	if err != nil {
		panic(err)
	}

	// HTTP server
	httpServer := dunai.CreateHTTPServer(cv)

	// Gemini server
	gServer := dunai.CreateGeminiServer(cv)

	// Background task
	go func() {
		dunai.UpdateStars(cv)
		ticker := time.NewTicker(1 * time.Hour)
		for {
			select {
			case <-ticker.C:
				dunai.UpdateStars(cv)
			}
		}
	}()

	// Run everything
	errch := make(chan error, 1)
	go func() {
		errch <- httpServer.ListenAndServe()
	}()
	go func() {
		errch <- gServer.ListenAndServe(context.Background())
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	select {
	case err := <-errch:
		log.Fatal(err)
	case <-c:
		log.Println("Terminating...")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := gServer.Shutdown(ctx); err != nil {
			log.Fatal(err)
		}
		if err := httpServer.Shutdown(ctx); err != nil {
			log.Fatal(err)
		}
		break
	}
	log.Println("Finished")
}
