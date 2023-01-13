package app

import (
	"encoding/json"
	"forum_app/internal/entity"
	"net/http"
)

func (h *Handler) APIResponse(w http.ResponseWriter, code int, response entity.Response) {
	if code == 204 {
		w.WriteHeader(204)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		h.errLog.Println(err)
		w.WriteHeader(500)
		w.Write([]byte(`{"error":"Internal Server Error"}`))
		return
	}
	w.WriteHeader(code)
	w.Write(jsonResponse)
}
