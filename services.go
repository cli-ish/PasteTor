package main

import (
	"log"
	"net/http"
	"pastetor/logic"
	"strings"
)

func routeIndexGet(w http.ResponseWriter) {
	templateHelper(w, nil, "templates/index.gohtml")
}

func routeFaqGet(w http.ResponseWriter) {
	templateHelper(w, nil, "templates/faq.gohtml")
}

func routeIndexPost(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		w.Header().Set("Location", "/")
		w.WriteHeader(301)
		return
	}
	data := r.Form.Get("paste_content")
	id, err := logic.AddNote(data, 0)
	if err != nil {
		w.Header().Set("Location", "/")
		w.WriteHeader(301)
		return
	}
	log.Println("New Note created: " + id)
	w.Header().Set("Location", "/p/"+id)
	w.WriteHeader(301)
}

func routeNoteGet(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/p/")
	note, err := logic.GetNote(id, 0)
	if err != nil {
		w.WriteHeader(404)
		return
	}
	data := struct {
		Note logic.Note
	}{
		Note: note,
	}
	templateHelper(w, data, "templates/paste.gohtml")
}

func routeNoteRawGet(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/raw/")
	note, err := logic.GetNote(id, 0)
	if err != nil {
		w.WriteHeader(404)
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	_, _ = w.Write([]byte(note.Data))
}

func routeNoteReportGet(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/report/")
	data := struct {
		Message string
		Id      string
	}{Id: id}
	note, err := logic.GetNote(id, 0)
	if err != nil {
		w.WriteHeader(404)
	}
	if note.State != 0 {
		// Don't show that this could be allowed.
		data.Message = "Already reported"
		templateHelper(w, data, "templates/report.gohtml")
		return
	}
	err = logic.ReportNote(note)
	if err != nil {
		w.WriteHeader(500)
		_, _ = w.Write([]byte("Error"))
		return
	}
	data.Message = "Reported!"
	templateHelper(w, data, "templates/report.gohtml")
}

func routeAdminGet(w http.ResponseWriter, r *http.Request) {
	username, password, ok := r.BasicAuth()
	if !ok || !isAuthorised(username, password) {
		w.Header().Add("WWW-Authenticate", `Basic realm="Give username and password"`)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	notes, err := logic.GetReportedNotes(0)
	if err != nil {
		return
	}
	csrf := generateCsrf()
	c := http.Cookie{
		Name:     "csrf",
		Value:    csrf,
		Path:     "/management",
		HttpOnly: true,
	}
	http.SetCookie(w, &c)
	data := struct {
		Notes []string
		Csrf  string
	}{notes, csrf}
	templateHelper(w, data, "templates/report_list.gohtml")
}

func routeAdminAllowedGet(w http.ResponseWriter, r *http.Request) {
	username, password, ok := r.BasicAuth()
	if !ok || !isAuthorised(username, password) {
		w.Header().Add("WWW-Authenticate", `Basic realm="Give username and password"`)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	notes, err := logic.GetAllowedNotes(0)
	if err != nil {
		return
	}
	csrf := generateCsrf()
	c := http.Cookie{
		Name:     "csrf",
		Value:    csrf,
		Path:     "/management",
		HttpOnly: true,
	}
	http.SetCookie(w, &c)
	data := struct {
		Notes []string
		Csrf  string
	}{notes, csrf}
	templateHelper(w, data, "templates/allowed_list.gohtml")
}

func routeAdminAllowNoteGet(w http.ResponseWriter, r *http.Request) {
	username, password, ok := r.BasicAuth()
	if !ok || !isAuthorised(username, password) {
		w.Header().Add("WWW-Authenticate", `Basic realm="Give username and password"`)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	csrf := r.URL.Query().Get("csrf")
	csrfc, err := r.Cookie("csrf")
	if err != nil || len(csrfc.Value) != 40 || csrf != csrfc.Value {
		w.WriteHeader(500)
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/management/allow/")
	note, err := logic.GetNote(id, 0)
	if err != nil {
		w.WriteHeader(404)
	}
	err = logic.AllowNote(note, 0)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	w.Header().Add("Location", "/management/")
	w.WriteHeader(301)
}

func routeAdminDeleteNoteGet(w http.ResponseWriter, r *http.Request) {
	username, password, ok := r.BasicAuth()
	if !ok || !isAuthorised(username, password) {
		w.Header().Add("WWW-Authenticate", `Basic realm="Give username and password"`)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	csrf := r.URL.Query().Get("csrf")
	csrfc, err := r.Cookie("csrf")
	if err != nil || len(csrfc.Value) != 40 || csrf != csrfc.Value {
		w.WriteHeader(500)
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/management/delete/")
	note, err := logic.GetNote(id, 0)
	if err != nil {
		w.WriteHeader(404)
		return
	}
	err = logic.DeleteNote(note, 0)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	w.Header().Add("Location", "/management/")
	w.WriteHeader(301)
}

func routeAdminUnallowNoteGet(w http.ResponseWriter, r *http.Request) {
	username, password, ok := r.BasicAuth()
	if !ok || !isAuthorised(username, password) {
		w.Header().Add("WWW-Authenticate", `Basic realm="Give username and password"`)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	csrf := r.URL.Query().Get("csrf")
	csrfc, err := r.Cookie("csrf")
	if err != nil || len(csrfc.Value) != 40 || csrf != csrfc.Value {
		w.WriteHeader(500)
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/management/unallow/")
	note, err := logic.GetNote(id, 0)
	if err != nil {
		w.WriteHeader(404)
	}
	err = logic.UnallowNote(note, 0)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	w.Header().Add("Location", "/management/allowed/")
	w.WriteHeader(301)
}
