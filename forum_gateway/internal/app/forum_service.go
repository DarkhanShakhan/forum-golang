package app

import (
	"encoding/json"
	"fmt"
	"forum_gateway/internal/entity"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"text/template"
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
		h.APIResponse(w, http.StatusMethodNotAllowed, entity.Response{ErrorMessage: "Invalid method"}, "/web/error.html")
		return
	}

	requestUrl := fmt.Sprintf("http://localhost:8080/post?id=%s", getID(r.URL.Path))
	req, err := http.NewRequest("GET", requestUrl, nil)
	if err != nil {
		h.APIResponse(w, http.StatusInternalServerError, entity.Response{ErrorMessage: "Internal Server Error"}, "/web/error.html")
		return
	}
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		h.APIResponse(w, http.StatusInternalServerError, entity.Response{ErrorMessage: "Internal Server Error"}, "web/error.html")
		return
	}
	defer resp.Body.Close()
	var response entity.Response
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		h.APIResponse(w, http.StatusInternalServerError, entity.Response{ErrorMessage: "Internal Server Error"}, "web/error.html")
		return
	}
	h.APIResponse(w, http.StatusOK, response, "web/post.html")
}

func UsersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// FIXME: check authorised
	url := "http://localhost:8080/users"
	client := http.Client{}
	response, err := client.Get(url)
	if err != nil {
		fmt.Println("here1")
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer response.Body.Close()
	var target interface{}
	err = json.NewDecoder(response.Body).Decode(&target)
	if err != nil {
		fmt.Println("here2")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	templ, err := template.New("posts").Parse(`<html>{{range.}}Name: <a href="http://localhost:8082/user/{{.id}}">{{.name}}</a><br> Email: {{.email}} <br><br>{{end}}</html>`)
	if err != nil {
		fmt.Println(err)
	}
	templ.ExecuteTemplate(w, "posts", target)
}

func getID(urlPath string) string {
	path := strings.Split(urlPath, "/")
	id := path[len(path)-1]
	return id
}

func UserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	requestUrl := fmt.Sprintf("http://localhost:8080/user?id=%s", getID(r.URL.Path))
	req, err := http.NewRequest("GET", requestUrl, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	var target interface{}
	err = json.NewDecoder(resp.Body).Decode(&target)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	templ, err := template.ParseFiles("web/user.html")
	if err != nil {
		fmt.Println(err)
	}
	templ.Execute(w, target)
}

func CategoryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	requestUrl := fmt.Sprintf("http://localhost:8080/category?id=%s", getID(r.URL.Path))
	req, err := http.NewRequest("GET", requestUrl, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	var target interface{}
	err = json.NewDecoder(resp.Body).Decode(&target)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	templ, err := template.ParseFiles("web/category.html")
	if err != nil {
		fmt.Println(err)
	}
	templ.Execute(w, target)
}

func PostCreateHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getPostCreate(w, r)
	case http.MethodPost:
		postPostCreate(w, r)
	default:
		w.WriteHeader(http.StatusBadRequest) //?

	}
}

func getPostCreate(w http.ResponseWriter, r *http.Request) {
	templ, err := template.ParseFiles("web/post_create.html")
	if err != nil {
		fmt.Println(err)
	}
	templ.Execute(w, nil)
}

func postPostCreate(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	res := parsePost(r.Form)
	fmt.Println(res)
	req, err := http.NewRequest(http.MethodPost, "http://localhost:8080/post/save", strings.NewReader(res))
	fmt.Println(err)
	client := http.Client{}
	resp, _ := client.Do(req)
	// defer resp.Body.Close()
	var target post
	err = json.NewDecoder(resp.Body).Decode(&target)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	redirectUrl := fmt.Sprintf("/post/%d", target.Id)
	http.Redirect(w, r, redirectUrl, 302)
}

func parsePost(form url.Values) string {
	post := post{}
	post.Title = form["title"][0]
	for _, a := range form["category"] {
		id, _ := strconv.Atoi(a)
		post.Category = append(post.Category, category{Id: id})
	}
	post.User.Id = 1
	res, _ := json.Marshal(post)
	return string(res)
}

type post struct {
	Id       int        `json:"id,omitempty"`
	User     user       `json:"user,omitempty"`
	Category []category `json:"categories,omitempty"`
	Title    string     `json:"title,omitempty"`
}

type user struct {
	Id int `json:"id,omitempty"`
}

type category struct {
	Id int `json:"id,omitempty"`
}

func CommentCreateHandler(w http.ResponseWriter, r *http.Request) {}
