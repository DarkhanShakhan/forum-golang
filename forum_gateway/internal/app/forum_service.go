package app

import (
	"errors"
	"fmt"
	"forum_gateway/internal/entity"
	"net/http"
	"strconv"
	"strings"
)

func (h *Handler) PostsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.APIResponse(w, http.StatusMethodNotAllowed, entity.Response{ErrorMessage: "Invalid method"}, "web/error.html")
		return
	}
	ctx, cancel := getTimeout(r.Context())
	defer cancel()
	response := entity.Response{}
	responseChan := make(chan entity.Response)
	go h.forumUcase.FetchPosts(ctx, responseChan)
	select {
	case <-ctx.Done():
		err := ctx.Err()
		h.errLog.Println(err)
		h.APIResponse(w, http.StatusRequestTimeout, entity.Response{ErrorMessage: "Request Timeout"}, "web/error.html")
		return
	case response = <-responseChan:
		err := response.Err
		switch err {
		case entity.ErrInternalServer:
			h.APIResponse(w, http.StatusInternalServerError, entity.Response{ErrorMessage: "Internal Server Error"}, "web/error.html")
		case entity.ErrRequestTimeout:
			h.APIResponse(w, http.StatusRequestTimeout, entity.Response{ErrorMessage: "Request Timeout"}, "web/error.html")
		case nil:
			var auth interface{} = r.Context().Value("authorised")
			response.AuthStatus, _ = auth.(bool)
			h.APIResponse(w, http.StatusOK, response, "web/posts.html")
		}
	}
}

func (h *Handler) PostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.APIResponse(w, http.StatusMethodNotAllowed, entity.Response{ErrorMessage: "Invalid method"}, "web/error.html")
		return
	}
	post_id, err := getID(r.URL.String(), "posts")
	if err != nil {
		h.errLog.Println(err)
		h.APIResponse(w, http.StatusNotFound, entity.Response{ErrorMessage: err.Error()}, "web/error.html")
		return
	}
	ctx, cancel := getTimeout(r.Context())
	defer cancel()
	response := entity.Response{}
	responseChan := make(chan entity.Response)
	go h.forumUcase.FetchPost(ctx, post_id, responseChan)
	select {
	case <-ctx.Done():
		err := ctx.Err()
		h.errLog.Println(err)
		h.APIResponse(w, http.StatusRequestTimeout, entity.Response{ErrorMessage: "Request Timeout"}, "web/error.html")
		return
	case response = <-responseChan:
		err := response.Err
		switch err {
		case entity.ErrInternalServer:
			h.APIResponse(w, http.StatusInternalServerError, entity.Response{ErrorMessage: "Internal Server Error"}, "web/error.html")
		case entity.ErrRequestTimeout:
			h.APIResponse(w, http.StatusRequestTimeout, entity.Response{ErrorMessage: "Request Timeout"}, "web/error.html")
		case entity.ErrNotFound:
			h.APIResponse(w, http.StatusNotFound, entity.Response{ErrorMessage: "Not Found"}, "web/error.html")
		case nil:
			var auth interface{} = r.Context().Value("authorised")
			response.AuthStatus, _ = auth.(bool)
			h.APIResponse(w, http.StatusOK, response, "web/post.html")
		}
	}
}

func (h *Handler) CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Context().Value("authorised") == false {
		h.APIResponse(w, http.StatusForbidden, entity.Response{ErrorMessage: "Forbidden"}, "web/error.html")
		return
	}
	switch r.Method {
	case http.MethodGet:
		h.getCreatePost(w, r)
	case http.MethodPost:
		h.postCreatePost(w, r)
	default:
		h.APIResponse(w, http.StatusBadRequest, entity.Response{ErrorMessage: "Bad request"}, "web/error.html")
	}
}

func (h *Handler) getCreatePost(w http.ResponseWriter, r *http.Request) {
	h.APIResponse(w, http.StatusOK, entity.Response{}, "web/post_create.html")
}

func (h *Handler) postCreatePost(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	post, err := entity.GetPost(r)
	if err != nil {
		h.APIResponse(w, http.StatusOK, entity.Response{ErrorMessage: err.Error()}, "web/post_create.html")
		return
	}
	var id interface{} = r.Context().Value("user_id")
	post.User.Id = id.(int64)
	ctx, cancel := getTimeout(r.Context())
	defer cancel()
	resChan := make(chan entity.Result)
	var res entity.Result
	go h.forumUcase.StorePost(ctx, post, resChan)
	select {
	case <-ctx.Done():
		err := ctx.Err()
		h.errLog.Println(err)
		h.APIResponse(w, http.StatusRequestTimeout, entity.Response{ErrorMessage: "Request Timeout"}, "web/error.html")
		return
	case res = <-resChan:
		switch res.Err {
		case entity.ErrBadRequest:
			h.APIResponse(w, http.StatusBadRequest, entity.Response{ErrorMessage: "Bad Request"}, "web/error.html")
		case nil:
			http.Redirect(w, r, fmt.Sprintf("/posts/%d", res.Id), 302)
		case entity.ErrRequestTimeout:
			h.APIResponse(w, http.StatusRequestTimeout, entity.Response{ErrorMessage: "Request Timeout"}, "web/error.html")
		default:
			h.APIResponse(w, http.StatusInternalServerError, entity.Response{ErrorMessage: "Internal Server Error"}, "web/error.html")

		}
	}
}

func (h *Handler) CreateCommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Context().Value("authorised") == false {
		h.APIResponse(w, http.StatusForbidden, entity.Response{ErrorMessage: "Forbidden"}, "web/error.html")
		return
	}
	if r.Method != http.MethodPost {
		h.APIResponse(w, http.StatusMethodNotAllowed, entity.Response{ErrorMessage: "Invalid method"}, "web/error.html")
		return
	}
	r.ParseForm()
	commentRes := entity.GetComment(r)
	if commentRes.Err != nil {
		h.errLog.Println(commentRes.Err)
		if commentRes.Err == entity.ErrEmptyComment {
			h.APIResponse(w, http.StatusBadRequest, entity.Response{ErrorMessage: "Empty comment"}, "web/error.html")
		} else {
			h.APIResponse(w, http.StatusBadRequest, entity.Response{ErrorMessage: "Bad Request"}, "web/error.html")
		}
		return
	}

	ctx, cancel := getTimeout(r.Context())
	defer cancel()
	resChan := make(chan entity.Result)
	var res entity.Result
	go h.forumUcase.StoreComment(ctx, commentRes.Comment, resChan)
	select {
	case <-ctx.Done():
		err := ctx.Err()
		h.errLog.Println(err)
		h.APIResponse(w, http.StatusRequestTimeout, entity.Response{ErrorMessage: "Request Timeout"}, "web/error.html")
		return
	case res = <-resChan:
		switch res.Err {
		case entity.ErrBadRequest:
			h.APIResponse(w, http.StatusBadRequest, entity.Response{ErrorMessage: "Bad Request"}, "web/error.html")
		case nil:
			http.Redirect(w, r, fmt.Sprintf("/posts/%d", commentRes.Comment.Post.Id), 302)
		case entity.ErrRequestTimeout:
			h.APIResponse(w, http.StatusRequestTimeout, entity.Response{ErrorMessage: "Request Timeout"}, "web/error.html")
		default:
			h.APIResponse(w, http.StatusInternalServerError, entity.Response{ErrorMessage: "Internal Server Error"}, "web/error.html")
		}
	}
}

func getID(urlPath, endpoint string) (int, error) {
	path := strings.Split(urlPath, "/")
	if path[len(path)-2] != endpoint {
		return 0, errors.New("Not Found")
	}
	id, err := strconv.Atoi(path[len(path)-1])
	if err != nil {
		return 0, errors.New("Not Found")
	}
	return id, nil
}

// func UserHandler(w http.ResponseWriter, r *http.Request) {
// 	if r.Method != http.MethodGet {
// 		w.WriteHeader(http.StatusBadRequest)
// 		return
// 	}

// 	requestUrl := fmt.Sprintf("http://localhost:8080/user?id=%s", getID(r.URL.Path))
// 	req, err := http.NewRequest("GET", requestUrl, nil)
// 	if err != nil {
// 		w.WriteHeader(http.StatusInternalServerError)
// 		return
// 	}
// 	client := http.Client{}
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		w.WriteHeader(http.StatusInternalServerError)
// 		return
// 	}
// 	defer resp.Body.Close()
// 	var target interface{}
// 	err = json.NewDecoder(resp.Body).Decode(&target)
// 	if err != nil {
// 		w.WriteHeader(http.StatusInternalServerError)
// 		return
// 	}
// 	templ, err := template.ParseFiles("web/user.html")
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	templ.Execute(w, target)
// }

// func CategoryHandler(w http.ResponseWriter, r *http.Request) {
// 	if r.Method != http.MethodGet {
// 		w.WriteHeader(http.StatusBadRequest)
// 		return
// 	}

// 	requestUrl := fmt.Sprintf("http://localhost:8080/category?id=%s", getID(r.URL.Path))
// 	req, err := http.NewRequest("GET", requestUrl, nil)
// 	if err != nil {
// 		w.WriteHeader(http.StatusInternalServerError)
// 		return
// 	}
// 	client := http.Client{}
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		w.WriteHeader(http.StatusInternalServerError)
// 		return
// 	}
// 	defer resp.Body.Close()
// 	var target interface{}
// 	err = json.NewDecoder(resp.Body).Decode(&target)
// 	if err != nil {
// 		w.WriteHeader(http.StatusInternalServerError)
// 		return
// 	}
// 	templ, err := template.ParseFiles("web/category.html")
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	templ.Execute(w, target)
// }

// func PostCreateHandler(w http.ResponseWriter, r *http.Request) {
// 	switch r.Method {
// 	case http.MethodGet:
// 		getPostCreate(w, r)
// 	case http.MethodPost:
// 		postPostCreate(w, r)
// 	default:
// 		w.WriteHeader(http.StatusBadRequest) //?

// 	}
// }

// func getPostCreate(w http.ResponseWriter, r *http.Request) {
// 	templ, err := template.ParseFiles("web/post_create.html")
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	templ.Execute(w, nil)
// }

// func postPostCreate(w http.ResponseWriter, r *http.Request) {
// 	r.ParseForm()
// 	res := parsePost(r.Form)
// 	fmt.Println(res)
// 	req, err := http.NewRequest(http.MethodPost, "http://localhost:8080/post/save", strings.NewReader(res))
// 	fmt.Println(err)
// 	client := http.Client{}
// 	resp, _ := client.Do(req)
// 	// defer resp.Body.Close()
// 	var target post
// 	err = json.NewDecoder(resp.Body).Decode(&target)
// 	if err != nil {
// 		w.WriteHeader(http.StatusInternalServerError)
// 		return
// 	}

// 	redirectUrl := fmt.Sprintf("/post/%d", target.Id)
// 	http.Redirect(w, r, redirectUrl, 302)
// }

// func parsePost(form url.Values) string {
// 	post := post{}
// 	post.Title = form["title"][0]
// 	for _, a := range form["category"] {
// 		id, _ := strconv.Atoi(a)
// 		post.Category = append(post.Category, category{Id: id})
// 	}
// 	post.User.Id = 1
// 	res, _ := json.Marshal(post)
// 	return string(res)
// }

// type post struct {
// 	Id       int        `json:"id,omitempty"`
// 	User     user       `json:"user,omitempty"`
// 	Category []category `json:"categories,omitempty"`
// 	Title    string     `json:"title,omitempty"`
// }

// type user struct {
// 	Id int `json:"id,omitempty"`
// }

// type category struct {
// 	Id int `json:"id,omitempty"`
// }

// func CommentCreateHandler(w http.ResponseWriter, r *http.Request) {}
