package app

import (
	"forum_gateway/internal/entity"
	"html/template"
	"net/http"
)

func (h *Handler) APIResponse(w http.ResponseWriter, code int, response entity.Response, filename string) {
	templ, err := template.ParseFiles(filename)
	if err != nil {
		h.errLog.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
		return
	}
	w.WriteHeader(code)
	templ.Execute(w, response)
}
