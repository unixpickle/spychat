package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"
	"net/url"
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
		"/":      server.Root,
		"/login": server.Login,
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

func (s *Server) Root(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		s.serveTemplate(w, "404", map[string]string{"path": r.URL.Path})
		return
	}
	if s.authenticated(r) {
		s.serveTemplate(w, "index", nil)
	} else {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}

func (s *Server) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		user := r.FormValue("user")
		password := r.FormValue("password")
		sess, err := NewSession(user, password)
		if err != nil {
			http.Redirect(w, r, "/login?error="+url.QueryEscape(err.Error()),
				http.StatusSeeOther)
			return
		}
		id := s.SessionTable.Add(sess)
		rawSess, _ := s.CookieStore.Get(r, "sessid")
		rawSess.Values["authenticated"] = true
		rawSess.Values["id"] = id
		rawSess.Save(r, w)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	s.serveTemplate(w, "login", map[string]string{"error": r.FormValue("error")})
}

func (s *Server) serveTemplate(w http.ResponseWriter, name string, data interface{}) {
	path := filepath.Join(s.Flags.TemplateDir, name+".html")
	temp, err := template.New(name).ParseFiles(path)
	if err != nil {
		log.Println("error in template load:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := temp.Execute(w, data); err != nil {
		log.Println("error in template execution:", err)
	}
}

func (s *Server) authenticated(r *http.Request) bool {
	sess, _ := s.CookieStore.Get(r, "sessid")
	val, _ := sess.Values["authenticated"].(bool)
	return val
}
