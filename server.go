package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

func LogHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// fmt.Println("Request hit", r.URL.Path)
		h.ServeHTTP(w, r)
	})
}

func ServeStatic() {
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	fmt.Println("Serving static files")
	err := http.ListenAndServe("127.0.0.1:3000", LogHandler(http.DefaultServeMux))

	if err != nil {
		log.Fatal(err)
	}
}

func handleFileRead(filename string, contentType string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		file, err := os.ReadFile("./static/" + filename)
		if err != nil {
			w.WriteHeader(500)
			fmt.Println(err.Error())
			return
		}
		w.Header().Set("Content-Type", contentType)
		w.Write(file)
	})
}

func handleDataRead(serv *DripServ) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		serv.Lock()
		defer serv.Unlock()

		str, err := json.Marshal(serv)
		if err != nil {
			fmt.Println(err.Error())
			w.WriteHeader(500)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(str))
	})
}
func handleImages(serv *DripServ) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		serv.Lock()
		defer serv.Unlock()

		pathChunks := strings.Split(r.URL.Path, "/")
		id := strings.TrimSuffix(pathChunks[len(pathChunks)-1], ".png")

		drip, found := serv.dripsMap[id]
		if !found || !drip.hasImage() {
			w.WriteHeader(404)
			return
		}

		w.Header().Set("Content-Type", "image/png")
		w.Write(drip.image)

	})
}

func handleStatic(serv *DripServ) http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/", handleFileRead("index.html", "text/html"))
	mux.Handle("/index.html", handleFileRead("index.html", "text/html"))
	mux.Handle("/index.js", handleFileRead("index.js", "text/javascript"))
	mux.Handle("/leaflet.canvas-markers.js", handleFileRead("leaflet.canvas-markers.js", "text/javascript"))
	mux.Handle("/style.css", handleFileRead("style.css", "text/css"))
	mux.Handle("/images/", handleImages(serv))
	mux.Handle("/data.json", handleDataRead(serv))

	return mux
}

func ServeData(addr string, port int, serv *DripServ) {
	fullAddr := fmt.Sprintf("%v:%v", addr, port)

	fmt.Printf("Serving data at %v\n", fullAddr)
	err := http.ListenAndServe(fullAddr, LogHandler(handleStatic(serv)))

	if err != nil {
		log.Fatal(err)
	}
}
