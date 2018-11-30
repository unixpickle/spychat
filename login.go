package main

import (
	"net/http"
	"net/url"
)

// HandleLogin serves the login page.
func (s *Server) HandleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		user := r.FormValue("username")
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

// Authenticated checks if the request is authenticated.
func (s *Server) Authenticated(r *http.Request) bool {
	sess, _ := s.CookieStore.Get(r, "sessid")
	val, _ := sess.Values["authenticated"].(bool)
	return val
}
