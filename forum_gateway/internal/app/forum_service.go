package app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"text/template"
)

func PostsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	//FIXME: check authorised
	url := "http://localhost:8080/posts"
	client := http.Client{}
	response, err := client.Get(url)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer response.Body.Close()
	var target interface{}
	err = json.NewDecoder(response.Body).Decode(&target)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	templ, err := template.New("posts").Parse(`<html>{{range.}} Title: {{.title}}, Categories: {{.categories}}<br> {{end}}</html>`)
	if err != nil {
		fmt.Println(err)
	}
	templ.ExecuteTemplate(w, "posts", target)
}

func PostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	requestUrl := fmt.Sprintf("http://localhost:8080/post?id=%s", getID(r.URL.Path))
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
	templ, err := template.New("post").Parse(`<html>Title: {{.title}}<br>Categories: {{range.categories}} {{.title}} {{end}}<br>Content: {{.content}}<br> User: <a href="http://localhost:8082/user/{{.user.id}}">{{.user.name}}</a> </html>`)
	if err != nil {
		fmt.Println(err)
	}
	templ.ExecuteTemplate(w, "post", target)
}

func UsersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		fmt.Println("here0")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	//FIXME: check authorised
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
	templ, err := template.New("post").Parse(`<html>Name: {{.name}}<br>Posts: {{.posts}}<br>Likes: {{.total_post_likes}}<br>Dislikes: {{.total_post_dislikes}}</html>`)
	if err != nil {
		fmt.Println(err)
	}
	templ.ExecuteTemplate(w, "post", target)
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
