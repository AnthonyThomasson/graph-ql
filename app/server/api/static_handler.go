package api

import (
	"bytes"
	"encoding/json"
	"html/template"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
)

type manifest struct {
	MainCSS struct {
		File string `json:"file"`
		Src  string `json:"src"`
	} `json:"src/main.css"`
	MainTS struct {
		File string `json:"file"`
		Src  string `json:"src"`
	} `json:"src/main.ts"`
}

func setStaticHandlers(r *chi.Mux) {
	staticFileServer := http.FileServer(http.Dir("static"))

	r.Get("/", serveIndex)
	r.Get("/index.html", serveIndex)
	r.Handle("/*", http.StripPrefix("/", staticFileServer))
}

func serveIndex(w http.ResponseWriter, r *http.Request) {
	indexHTML, err := template.ParseFiles("static/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonData, err := os.ReadFile("static/manifest.json")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var manifest manifest
	err = json.Unmarshal(jsonData, &manifest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var html bytes.Buffer
	indexHTML.Execute(&html, struct {
		CSS  string
		JS   string
		Mode string
	}{
		CSS:  manifest.MainCSS.File,
		JS:   manifest.MainTS.File,
		Mode: os.Getenv("MODE"),
	})
	w.Write(html.Bytes())
}
