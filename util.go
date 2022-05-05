package main

import (
	"html/template"
	"math/rand"
	"net/http"
	"strings"
)

func HeaderHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/styles/") {
			rw.Header().Set("cache-control", "max-age=86400, public")
		}
		rw.Header()["Date"] = nil // Don't send the date !
		rw.Header().Set("x-xss-Protection", "1; mode=block")
		rw.Header().Set("content-security-policy", "default-src 'self'; img-src 'self'; object-src 'none'; script-src 'self'; style-src 'self'; frame-ancestors 'self'; base-uri 'self'; form-action 'self';")
		rw.Header().Set("x-frame-options", "DENY")
		rw.Header().Set("x-content-type-options", "nosniff")
		next.ServeHTTP(rw, r)
	})
}

func templateHelper(w http.ResponseWriter, data interface{}, path string) {
	t, err := template.ParseFS(index, path)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(200)
	err = t.Execute(w, data)
	if err != nil {
		w.WriteHeader(500)
		return
	}
}

func isAuthorised(user string, pass string) bool {
	for _, usr := range users {
		if usr.Username == user && usr.Password == pass {
			return true
		}
	}
	return false
}

func generateCsrf() string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyz_./ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	s := make([]rune, 40)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}
