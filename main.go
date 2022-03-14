package main

import (
	"context"
	"fmt"
	"html/template"
	textTemplate "text/template"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"git.sr.ht/~adnano/go-gemini"
	"git.sr.ht/~adnano/go-gemini/certificate"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type TemplateRenderer struct {
	templates *template.Template
}

func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	// cv, err := ReadCV()
	// if err != nil {
	//     panic(err)
	// }
	// fmt.Println("!!!!!!", data)
	// if viewContext, isMap := data.(map[string]interface{}); isMap {
	//     viewContext["foo"] = "bar"
	// }
	// fmt.Println("!!!!!!", data)
	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
	if err := LoadConfig(); err != nil {
		panic(err)
	}
	if err := RepoInit(); err != nil {
		panic(err)
	}
	e := echo.New()
	e.Debug = true
	e.Static("/static", "static")
	renderer := &TemplateRenderer{template.Must(template.ParseGlob("templates/*.html"))}
	e.Renderer = renderer
	e.Use(middleware.Recover())
	// e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
	//     StackSize: 1 << 15, // 32 KB
	//     LogLevel:  log.ERROR,
	// }))
	// e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
	//     return func(c echo.Context) error {
	//         return echo.NewHTTPError(http.StatusOK, "Error")
	//     }
	// })
	e.GET("/", func(c echo.Context) error {
		cv, err := ReadCV()
		if err != nil {
			return err
		}
		var projects []Project
		db.Order("\"order\"").Find(&projects)
		return c.Render(http.StatusOK, "index.html", map[string]interface{}{
			"cv":       cv,
			"projects": projects,
		})
		// return c.JSON(http.StatusOK, struct {
		//     Success bool `json:"success"`
		// }{true})
	})

	// Gemini server
	certificates := &certificate.Store{}
	certificates.Register("dun.ai")
	if err := certificates.Load("/var/lib/gemini/certs"); err != nil {
		log.Fatal(err)
	}

	files, err := os.ReadDir("/var/lib/gemini/certs")
	if err != nil {
		log.Fatal(err)
	}
	for k, v := range files {
		fmt.Println(k, v)
	}
	fmt.Println(certificates.Get("dun.ai"))

	mux := &gemini.Mux{}
	gIndex := textTemplate.Must(textTemplate.ParseFiles("templates/index.gmi"))
	mux.HandleFunc("/", func(c context.Context, rw gemini.ResponseWriter, r *gemini.Request) {
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
	gServer := &gemini.Server{
		Handler: gemini.LoggingMiddleware(mux),
		ReadTimeout: 10 * time.Second,
		WriteTimeout: 10 * time.Second,
		GetCertificate: certificates.Get,
	}

	// e.GET("/cv", func(c echo.Context) error {
	//     cv, err := ReadCV()
	//     if err != nil {
	//         return err
	//     }
	//     pdf, err := cv.GeneratePDF()
	//     if err != nil {
	//         return err
	//     }
	//     return c.Blob(http.StatusOK, "application/pdf", pdf.GetBytesPdf())
	// })

	// Background task

	go func() {
		UpdateStars()
		ticker := time.NewTicker(1 * time.Hour)
		for {
			select {
			case <-ticker.C:
				UpdateStars()
			}
		}
	}()

	// Run everything

	errch := make(chan error, 1)

	go func() {
		errch <- e.Start(":8888")
	}()

	go func() {
		errch <- gServer.ListenAndServe(context.Background())
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	// e.Server.Addr = ":8888"
	select {
	case err := <-errch:
		log.Fatal(err)
	case <-c:
		log.Println("Terminating...")
		ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
		defer cancel()
		if err := gServer.Shutdown(ctx); err != nil {
			log.Fatal(err)
		}
		if err := e.Shutdown(ctx); err != nil {
			log.Fatal(err)
		}
		break
	}
	log.Println("Finished")
}
