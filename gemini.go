package main

import (
	"context"
	"log"
	textTemplate "text/template"
	"time"

	"git.sr.ht/~adnano/go-gemini"
	"git.sr.ht/~adnano/go-gemini/certificate"
)

func CreateGeminiServer() *gemini.Server {
	certificates := &certificate.Store{}
	certificates.Register("dun.ai")
	if err := certificates.Load("/var/lib/gemini/certs"); err != nil {
		log.Fatal(err)
	}

	gServer := &gemini.Server{
		Handler:        gemini.LoggingMiddleware(CreateGeminiMux()),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		GetCertificate: certificates.Get,
	}
	return gServer
}

func CreateGeminiMux() *gemini.Mux {
	geminiMux := &gemini.Mux{}
	gIndex := textTemplate.Must(textTemplate.ParseFS(templates, "templates/index.gmi"))
	geminiMux.HandleFunc("/", func(c context.Context, rw gemini.ResponseWriter, r *gemini.Request) {
		cv, err := ReadCV()
		if err != nil {
			return
		}
		var projects []Project
		db.Order("\"order\"").Find(&projects)
		gIndex.Execute(rw, map[string]interface{}{
			"cv":       cv,
			"projects": projects,
		})
		// fmt.Fprintln(rw, "Welcome!")
		// fmt.Fprintln(rw)
		// fmt.Fprintln(rw, "=> /projects My projects")
		// fmt.Fprintln(rw, "=> /stuff Stuff")
	})

	return geminiMux
}
