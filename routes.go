package main

import (
	"net/http"
)

func routes(r *http.ServeMux) {
	r.Handle("/styles/", http.FileServer(http.FS(styles)))
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			routeIndexGet(w)
		case "POST":
			routeIndexPost(w, r)
		default:
			w.WriteHeader(405)
		}
	})

	r.HandleFunc("/p/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			routeNoteGet(w, r)
		default:
			w.WriteHeader(405)
		}
	})

	r.HandleFunc("/report/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			routeNoteReportGet(w, r)
		default:
			w.WriteHeader(405)
		}
	})

	r.HandleFunc("/raw/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			routeNoteRawGet(w, r)
		default:
			w.WriteHeader(405)
		}
	})

	r.HandleFunc("/faq", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			routeFaqGet(w)
		default:
			w.WriteHeader(405)
		}
	})

	r.HandleFunc("/management/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			routeAdminGet(w, r)
		default:
			w.WriteHeader(405)
		}
	})

	r.HandleFunc("/management/allowed/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			routeAdminAllowedGet(w, r)
		default:
			w.WriteHeader(405)
		}
	})

	r.HandleFunc("/management/allow/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			routeAdminAllowNoteGet(w, r)
		default:
			w.WriteHeader(405)
		}
	})

	r.HandleFunc("/management/unallow/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			routeAdminUnallowNoteGet(w, r)
		default:
			w.WriteHeader(405)
		}
	})

	r.HandleFunc("/management/delete/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			routeAdminDeleteNoteGet(w, r)
		default:
			w.WriteHeader(405)
		}
	})
}
