package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gorilla/context"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/unixpickle/essentials"
)

type ServerFlags struct {
	Port        int
	AssetDir    string
	TemplateDir string
}

func main() {
	var flags ServerFlags
	flag.IntVar(&flags.Port, "port", 8080, "server port number")
	flag.StringVar(&flags.AssetDir, "assets", "assets", "asset directory")
	flag.StringVar(&flags.TemplateDir, "templates", "templates", "template directory")
	flag.Parse()

	server := &Server{
		Flags:        &flags,
		SessionTable: NewSessionTable(),
		CookieStore: sessions.NewCookieStore(securecookie.GenerateRandomKey(16),
			securecookie.GenerateRandomKey(16)),
	}
	http.Handle("/assets/", context.ClearHandler(server.AssetHandler()))

	handlers := map[string]http.HandlerFunc{
		"/":      server.HandleRoot,
		"/login": server.HandleLogin,
	}
	for path, handler := range handlers {
		http.Handle(path, context.ClearHandler(handler))
	}

	log.Println("Attempting to listen on port", flags.Port, "...")
	err := http.ListenAndServe(":"+strconv.Itoa(flags.Port), nil)
	if err != nil {
		essentials.Die(err)
	}
}

type Server struct {
	Flags *ServerFlags

	SessionTable *SessionTable
	CookieStore  *sessions.CookieStore
}

func (s *Server) AssetHandler() http.Handler {
	return http.StripPrefix("/assets/", http.FileServer(http.Dir(s.Flags.AssetDir)))
}

func (s *Server) HandleRoot(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		s.serveTemplate(w, "404", map[string]string{"path": r.URL.Path})
		return
	}
	if s.Authenticated(r) {
		s.serveTemplate(w, "index", nil)
	} else {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}

func (s *Server) serveTemplate(w http.ResponseWriter, name string, data interface{}) {
	path := filepath.Join(s.Flags.TemplateDir, name+".html")
	temp, err := template.New(name + ".html").ParseFiles(path)
	if err != nil {
		log.Println("error in template load:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := temp.Execute(w, data); err != nil {
		log.Println("error in template execution:", err)
	}
}
