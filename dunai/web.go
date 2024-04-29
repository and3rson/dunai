package dunai

import (
	"fmt"
	htmlTemplate "html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/russross/blackfriday/v2"
)

var IndexTemplate = MustLoadTemplate("index")

func MustLoadTemplate(name string) *htmlTemplate.Template {
	data, err := templates.ReadFile(fmt.Sprintf("templates/%s.html", name))
	if err != nil {
		panic(err)
	}
	template, err := htmlTemplate.New("index").Funcs(htmlTemplate.FuncMap{
		"markdown": func(s string) htmlTemplate.HTML {
			result := htmlTemplate.HTML(blackfriday.Run([]byte(s)))
			// Strip the <p> and </p> tags
			return result[3 : len(result)-5]
		},
	}).Parse(string(data))
	if err != nil {
		log.Fatal(err)
	}
	return template
}

func CreateHTTPServer(cv *CV) *http.Server {
	return &http.Server{
		Addr:    "0.0.0.0:8888",
		Handler: CreateWebRouter(cv),
	}
}

func CreateWebRouter(cv *CV) *mux.Router {
	webMux := mux.NewRouter()
	webMux.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		// var projects []Project
		// db.Order("\"order\"").Find(&projects)
		rw.Header().Add("Content-Type", "text/html; charset=utf-8")
		rw.WriteHeader(200)
		// rw.Write([]byte("Test\n"))
		if err := IndexTemplate.Execute(rw, map[string]interface{}{
			"cv": cv,
		}); err != nil {
			log.Printf("render template: %v", err)
		}
	})
	webMux.PathPrefix("/static").Handler(http.FileServer(http.FS(static)))
	return webMux
}
