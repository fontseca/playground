package main

import (
  "context"
  "fmt"
  "log"
  "net"
  "net/http"
  "playground"
  "slices"
  "time"
)

func main() {
  mux := http.NewServeMux()

  mux.HandleFunc("GET /playground/engine.js", func(w http.ResponseWriter, r *http.Request) { http.ServeFile(w, r, "engine.js") })
  mux.HandleFunc("GET /playground/stylesheet.css", func(w http.ResponseWriter, r *http.Request) { http.ServeFile(w, r, "stylesheet.css") })

  playgroundCtx, playgroundCtxCanceler := context.WithCancel(context.Background())
  defer playgroundCtxCanceler()
  mux.HandleFunc("POST /playground.request", func(w http.ResponseWriter, r *http.Request) {
    playground.Scanner(playgroundCtx, w, r)
  })

  mux.HandleFunc("GET /", playground.Renderer)
  mux.HandleFunc("POST /", playground.Renderer)

  server := http.Server{
    Handler:           mux,
    IdleTimeout:       1 * time.Minute,
    ReadTimeout:       5 * time.Second,
    WriteTimeout:      5 * time.Second,
    MaxHeaderBytes:    1024,
    ReadHeaderTimeout: 0,
  }

  listener, err := net.Listen("tcp", ":0")
  if nil != err {
    log.Fatalf("net.Listen(...) failed: %v", err)
  }

  defer listener.Close()

  var (
    addrs []net.Addr
    ip    string
  )

  addrs, _ = net.InterfaceAddrs()
  for addr := range slices.Values(addrs) {
    ipnet, ok := addr.(*net.IPNet)
    if ok && !ipnet.IP.IsLoopback() && nil != ipnet.IP.To4() {
      ip = ipnet.IP.String()
    }
  }

  fmt.Printf("running fontseca.dev/playground server at %v:%v\n", ip, listener.Addr().(*net.TCPAddr).Port)
  log.Fatalf("server.Serve(...) failed: %v", server.Serve(listener))
}
