package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {

	if len(os.Args) != 2 {
		panic("Please pass port as argument to goserve")
	}

	r := chi.NewRouter()
	r.Use(middleware.RealIP)
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)
	r.Use(middleware.NoCache)
	FileServer(r)

	go func() {
		time.Sleep(5 * time.Millisecond)
		fmt.Println("Running GoServe port: " + os.Args[1])
	}()

	log.Fatal(http.ListenAndServe(":"+os.Args[1], r))
}

func FileServer(router *chi.Mux) {
	root := "./"
	fs := http.FileServer(http.Dir(root))

	router.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		if _, err := os.Stat(root + r.RequestURI); os.IsNotExist(err) {
			http.StripPrefix(r.RequestURI, fs).ServeHTTP(w, r)
		} else {
			fs.ServeHTTP(w, r)
		}
	})
}
