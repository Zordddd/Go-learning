package handler

import "net/http"

func TaskHandler(w http.ResponseWriter, r *http.Request) {
	Method := r.Method
	switch Method {
	case http.MethodGet:
		GetTasksHandler(w, r)
	case http.MethodPost:
		CreateTaskHandler(w, r)
	}
}
