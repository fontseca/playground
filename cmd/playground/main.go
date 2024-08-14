package main

import (
  "fontseca.dev/playground"
  "log/slog"
  "net/http"
  "time"
)

func main() {
  mux := http.NewServeMux()

  mux.HandleFunc("GET /playground/engine.js", func(w http.ResponseWriter, r *http.Request) { http.ServeFile(w, r, "engine.js") })
  mux.HandleFunc("GET /playground/stylesheet.css", func(w http.ResponseWriter, r *http.Request) { http.ServeFile(w, r, "stylesheet.css") })

  mux.HandleFunc("POST /playground.request", playground.Scanner)
  mux.HandleFunc("GET /", playground.Renderer)

  server := http.Server{
    Addr:              "0.0.0.0:52368",
    Handler:           mux,
    IdleTimeout:       1 * time.Minute,
    ReadTimeout:       5 * time.Second,
    WriteTimeout:      5 * time.Second,
    MaxHeaderBytes:    1024,
    ReadHeaderTimeout: 0,
  }

  slog.Info("running fontseca.dev's playground", slog.String("addr", server.Addr))

  slog.Error(server.ListenAndServe().Error())
}
