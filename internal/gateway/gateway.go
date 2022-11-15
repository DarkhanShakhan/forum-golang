package gateway

import (
	"encoding/json"
	"net/http"
	"strconv"
)

type APIGateway struct {
	ucase UserUsecase
	pcase PostUsecase
	ccase CommentUsecase
}

func NewAPIGateway(ucase UserUsecase, pcase PostUsecase, ccase CommentUsecase) *APIGateway {
	return &APIGateway{ucase, pcase, ccase}
}

func (a *APIGateway) MainHandler(w http.ResponseWriter, r *http.Request) {
	posts, _ := a.pcase.FetchAll()
	blob, _ := json.Marshal(posts)
	w.Header().Set("Content-Type", "application/json")
	w.Write(blob)
}

func (a *APIGateway) UserHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	id, _ := strconv.Atoi(path[len(path)-1:])
	user, _ := a.ucase.FetchById(id)
	blob, _ := json.Marshal(user)
	w.Header().Set("Content-Type", "application/json")
	w.Write(blob)
}
