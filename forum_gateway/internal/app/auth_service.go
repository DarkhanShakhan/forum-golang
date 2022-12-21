package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/url"
)

var oauthStateString = "pseudo-random"

func SignInHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getSignIn(w, r)
	case http.MethodPost:
		postSignIn(w, r)
	default:
		w.WriteHeader(http.StatusBadRequest)
	}
}

func getSignIn(w http.ResponseWriter, r *http.Request) {
	templ, err := template.ParseFiles("web/sign_in.html")
	if err != nil {
		fmt.Println(err)
	}
	templ.Execute(w, nil)
}

func postSignIn(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if r.Context().Value("authorised") == true {
		w.WriteHeader(400)
		return
	}
	credentials := Credentials{Email: r.FormValue("email"), Password: r.FormValue("password")}
	requestBody, err := json.Marshal(credentials)
	if err != nil {
		fmt.Println(1)
		fmt.Println(err)
		return
	}
	requestUrl := "http://localhost:8081/sign_in"
	req, err := http.NewRequest(http.MethodPost, requestUrl, bytes.NewReader(requestBody))
	if err != nil {
		fmt.Println(err)
	}
	client := http.Client{}
	response, err := client.Do(req)
	if err != nil {
		fmt.Println(2)
		fmt.Println(err)
		return
	}
	fmt.Println(response.StatusCode)
	session := Session{}
	err = json.NewDecoder(response.Body).Decode(&session)
	// FIXME: empty session check
	if err != nil {
		fmt.Println(3)
		fmt.Println(err)
		return
	}
	if session.Token != "" {
		cookie := http.Cookie{Name: "token", Value: session.Token}
		http.SetCookie(w, &cookie)
		http.Redirect(w, r, "/posts", 304)
		return
	}
	w.Write([]byte("invalid"))
}

type Session struct {
	Token      string `json:"token"`
	UserId     int64  `json:"user_id"`
	ExpiryDate string `json:"expiry_date"`
}

type Credentials struct {
	Id       int64  `json:"id,omitempty"`
	Name     string `json:"name,omitempty"`
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getSignUp(w, r)
	case http.MethodPost:
		postSignUp(w, r)
	default:
		w.WriteHeader(http.StatusBadRequest)
	}
}

func getSignUp(w http.ResponseWriter, r *http.Request) {
	templ, err := template.ParseFiles("web/sign_up.html")
	if err != nil {
		fmt.Println(err)
	}
	templ.Execute(w, nil)
}

func postSignUp(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	// FIXME: check password and username and email
	credentials := Credentials{Name: r.FormValue("name"), Email: r.FormValue("email"), Password: r.FormValue("password")}
	requestBody, err := json.Marshal(credentials)
	if err != nil {
		fmt.Println(1)
		fmt.Println(err)
		return
	}
	requestUrl := "http://localhost:8081/sign_up"
	req, err := http.NewRequest(http.MethodPost, requestUrl, bytes.NewReader(requestBody))
	if err != nil {
		fmt.Println(err)
	}
	client := http.Client{}
	response, err := client.Do(req)
	if err != nil {
		fmt.Println(2)
		fmt.Println(err)
		return
	}
	fmt.Println(response.StatusCode)
}

func SignOutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// if cookie, err := r.Cookie("token"); err == nil {
	// 	token := cookie.Value
	// 	Session{Token: to}
	// }
}

func SignInGoogleHandler(w http.ResponseWriter, r *http.Request) {
	url := AuthCodeURL(oauthStateString)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func AuthCodeURL(state string) string {
	var buf bytes.Buffer
	buf.WriteString("https://accounts.google.com/o/oauth2/auth")
	v := url.Values{"response_type": {"code"}, "client_id": {""}}
	v.Set("redirect_uri", "http://localhost:8082/google_callback")
	v.Set("scope", "https://www.googleapis.com/auth/userinfo.email https://www.googleapis.com/auth/userinfo.profile")
	v.Set("state", state)
	buf.WriteByte('?')
	buf.WriteString(v.Encode())
	return buf.String()
}

func GoogleCallbackHandler(w http.ResponseWriter, r *http.Request) {
	content, err := getUserInfo(r.FormValue("state"), r.FormValue("code"))
	if err != nil {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	// send email to forum_auth
	// FIXME: create authGoogle endpoint
	fmt.Fprintf(w, "Content: %s\n", content)
}

func getUserInfo(state string, code string) ([]byte, error) {
	fmt.Println(state)
	fmt.Println(code)
	if state != oauthStateString {
		return nil, fmt.Errorf("invalid oauth state")
	}
	token, err := exchange(code)
	if err != nil {
		return nil, fmt.Errorf("code exchange failed: %s", err.Error())
	}
	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed getting user info: %s", err.Error())
	}
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed reading response body: %s", err.Error())
	}
	return contents, nil
}

func exchange(code string) (Token, error) {
	var buf bytes.Buffer
	buf.WriteString("https://oauth2.googleapis.com/token?")
	v := url.Values{"grant_type": {"authorization_code"}, "code": {code}}
	v.Set("redirect_uri", "http://localhost:8082/google_callback")
	v.Set("client_id", "")
	v.Set("client_secret", "")
	buf.WriteString(v.Encode())
	url := buf.String()
	req, _ := http.NewRequest("POST", url, nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return Token{}, err
	}
	defer resp.Body.Close()
	var token Token
	bytes, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(bytes, &token)
	return token, nil
}

type Token struct {
	AccessToken string `json:"access_token"`
}
