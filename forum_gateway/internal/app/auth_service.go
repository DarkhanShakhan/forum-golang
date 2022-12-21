package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"unicode"
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

// SIGN UP
func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	if r.Context().Value("authorised") == true {
		w.WriteHeader(403) // FIXME: not sure
		w.Write([]byte("you are already authorised. you can sign out"))
		return
	}
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
		w.WriteHeader(500)
	}
	templ.Execute(w, nil)
}

func postSignUp(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	name := r.FormValue("name")
	email := r.FormValue("email")
	password := r.FormValue("password")
	rep_password := r.FormValue("rep_password")
	if len(name) == 0 || !validEmail(email) || !validPassword(password) || password != rep_password {
		// FIXME: send detailed error message
		log.Println("invalid request: credentials verification")
		templ, err := template.ParseFiles("web/sign_up.html")
		if err != nil {
			w.WriteHeader(500)
			return
		}
		w.WriteHeader(400)
		templ.Execute(w, "invalid request: credentials verification")
		return
	}
	credentials := Credentials{Name: name, Email: email, Password: password}
	requestBody, err := json.Marshal(credentials)
	if err != nil {
		log.Println(err)
		templ, err := template.ParseFiles("web/sign_up.html")
		if err != nil {
			w.WriteHeader(500)
			return
		}
		w.WriteHeader(400)
		templ.Execute(w, "invalid request")
		return
	}
	requestUrl := "http://localhost:8081/sign_up"
	req, err := http.NewRequest(http.MethodPost, requestUrl, bytes.NewReader(requestBody))
	if err != nil {
		log.Println(err)
		w.WriteHeader(500)

		return
	}
	client := http.Client{}
	response, err := client.Do(req)
	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
		templ, err := template.ParseFiles("web/sign_up.html")
		if err != nil {
			return
		}
		templ.Execute(w, "internal server error")
		return
	}
	fmt.Println(response.Body) // FIXME: check for errors
	http.Redirect(w, r, "/sign_in", http.StatusFound)
}

func validEmail(email string) bool {
	return regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`).MatchString(email)
}

func validPassword(pass string) bool {
	var (
		upp, low, num, sym bool
		tot                uint8
	)
	for _, char := range pass {
		switch {
		case unicode.IsUpper(char):
			upp = true
			tot++
		case unicode.IsLower(char):
			low = true
			tot++
		case unicode.IsNumber(char):
			num = true
			tot++
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			sym = true
			tot++
		default:
			return false
		}
	}
	if !upp || !low || !num || !sym || tot < 8 {
		return false
	}
	return true
}

type response struct {
	Err     string      `json:"error,omitempty"`
	Content interface{} `json:"content,omitempty"`
}

func SignOutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// if cookie, err := r.Cookie("token"); err == nil {
	// 	token := cookie.Value
	// 	session := Session{Token: token}
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
