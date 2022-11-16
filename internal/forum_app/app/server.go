package app

import (
	"encoding/json"
	"fmt"
	"forum/internal/forum_app/entity"
	"log"
	"net/http"
	"strconv"
)

type Server struct {
	ucase UserUsecase
	pcase PostUsecase
	ccase CommentUsecase
}

func NewServer(ucase UserUsecase, pcase PostUsecase, ccase CommentUsecase) *Server {
	return &Server{ucase, pcase, ccase}
}

func (s *Server) UsersAllHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		log.Println("UsersAllHandler: method not allowed (" + r.Method + ")")
		w.WriteHeader(405)
		return
	}
	users, err := s.ucase.FetchAll()
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		log.Println("UsersAllHandler" + err.Error())
		w.WriteHeader(500)
		return
	}
	response, err := json.Marshal(users)
	if err != nil {
		log.Println("UsersAllHandler" + err.Error())
		w.WriteHeader(500)
		return
	}
	w.Write(response)
}

func (s *Server) UserByIdHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		log.Println("UsersByIdHandler: method not allowed (" + r.Method + ")")
		w.WriteHeader(405)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := r.ParseForm(); err != nil {
		log.Println("UserByIdHandler: " + err.Error())
		w.WriteHeader(500)
		return
	}
	id, err := strconv.Atoi(r.Form.Get("id"))
	if err != nil {
		log.Println("UserByIdHandler: " + err.Error())
		w.WriteHeader(400)
		return
	}
	user, err := s.ucase.FetchById(id)
	response, err := json.Marshal(user)
	if err != nil {
		log.Println("UserByIdHandler: " + err.Error())
		w.WriteHeader(500)
		return
	}
	w.Write(response)
}

func (s *Server) PostByIdHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		log.Println("PostByIdHandler: method not allowed (" + r.Method + ")")
		w.WriteHeader(405)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := r.ParseForm(); err != nil {
		log.Println("PostByIdHandler: " + err.Error())
		w.WriteHeader(500)
		return
	}
	id, err := strconv.Atoi(r.Form.Get("id"))
	if err != nil {
		log.Println("PostByIdHandler: " + err.Error())
		w.WriteHeader(400)
		return
	}
	post, err := s.pcase.FetchById(id)
	response, err := json.Marshal(post)
	if err != nil {
		log.Println("PostByIdHandler: " + err.Error())
		w.WriteHeader(500)
		return
	}
	w.Write(response)
}

func (s *Server) PostsAllHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		log.Println("PostsAllHandler: method not allowed (" + r.Method + ")")
		w.WriteHeader(405)
		return
	}
	posts, err := s.pcase.FetchAll()
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		log.Println("PostsAllHandler" + err.Error())
		w.WriteHeader(500)
		return
	}
	response, err := json.Marshal(posts)
	if err != nil {
		log.Println("PostsAllHandler" + err.Error())
		w.WriteHeader(500)
		return
	}
	w.Write(response)
}

func (s *Server) CategoryPostsHandler(w http.ResponseWriter, r *http.Request) {}

func (s *Server) StorePostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		log.Println("StorePostHandler: method not allowed (" + r.Method + ")")
		w.WriteHeader(405)
		return
	}
	var post entity.Post
	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		log.Println("StorePostHandler: bad request")
		w.WriteHeader(400)
		return
	}
	id, err := s.pcase.Store(post)
	if err != nil {
		log.Println("StorePostHandler: " + err.Error())
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	w.Write([]byte(fmt.Sprintf("{\"id\":%d}", id)))
}

func (s *Server) StorePostReactionHandler(w http.ResponseWriter, r *http.Request) {}

func (s *Server) UpdatePostReactionHandler(w http.ResponseWriter, r *http.Request) {}

func (s *Server) DeletePostReactionHandler(w http.ResponseWriter, r *http.Request) {}

func (s *Server) CommentByIdHandler(w http.ResponseWriter, r *http.Request) {}

func (s *Server) StoreCommentHandler(w http.ResponseWriter, r *http.Request) {}

func (s *Server) StoreCommentReactionHandler(w http.ResponseWriter, r *http.Request) {}

func (s *Server) UpdateCommentReactionHandler(w http.ResponseWriter, r *http.Request) {}

func (s *Server) DeleteCommentReaction(w http.ResponseWriter, r *http.Request) {}
