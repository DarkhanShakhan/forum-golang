package app

import (
	"forum_gateway/internal/entity"
	"html/template"
	"net/http"
)

func APIResponse(w http.ResponseWriter, code int, response entity.Response, filename string) {
	templ, err := template.ParseFiles(filename)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(code)
	templ.Execute(w, response)
}
