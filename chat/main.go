package main

import (
  "log"
  "net/http"
  "text/template"
  "path/filepath"
  "sync"
  "flag"
  // "github.com/pigcafe/chat/trace"
  // "os"
)

type templateHandler struct {
  once sync.Once
  filename string
  templ *template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  t.once.Do(func() {
    t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
  })
  t.templ.Execute(w, r)
}

func main() {
  var addr = flag.String("addr", ":8080", "Application address")
  flag.Parse() // Interpretate flag.
  r := newRoom()
  // r.tracer = trace.New(os.Stdout)
  http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"}))
  http.Handle("/login", &templateHandler{filename: "login.html"})
  http.Handle("/room", r)
  // Starting chat room.
  go r.run()
  // Starting web server.
  log.Println("Starting web server. Port:", *addr)
  if err := http.ListenAndServe(*addr, nil); err != nil {
    log.Fatal("ListenAndServe:", err)
  }
}
