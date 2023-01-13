package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"forum_gateway/internal/entity"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"unicode"
)

// SIGN UP
func (h *Handler) SignUpHandler(w http.ResponseWriter, r *http.Request) {
	if r.Context().Value("authorised") == true {
		h.APIResponse(w, http.StatusForbidden, entity.Response{ErrorMessage: "Forbidden"}, "web/error.html")
		return
	}
	switch r.Method {
	case http.MethodGet:
		h.getSignUp(w, r)
	case http.MethodPost:
		h.postSignUp(w, r)
	default:
		h.APIResponse(w, http.StatusBadRequest, entity.Response{ErrorMessage: "Bad request"}, "web/error.html")
	}
}

func (h *Handler) getSignUp(w http.ResponseWriter, r *http.Request) {
	h.APIResponse(w, http.StatusOK, entity.Response{}, "web/sign_up.html")
}

func (h *Handler) postSignUp(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	name := r.FormValue("name")
	email := r.FormValue("email")
	password := r.FormValue("password")
	confirm_password := r.FormValue("confirm_password")
	ok, message := checkCredentials(name, email, password, confirm_password)
	if !ok {
		h.errLog.Println(message)
		h.APIResponse(w, http.StatusBadRequest, entity.Response{ErrorMessage: message}, "web/sign_up.html")
	}
	credentials := entity.Credentials{Name: name, Email: email, Password: password}
	requestBody, err := json.Marshal(credentials)
	if err != nil {
		h.errLog.Println(err)
		h.APIResponse(w, http.StatusBadRequest, entity.Response{ErrorMessage: "bad request"}, "web/error.html")
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

func checkCredentials(name, email, password, confirm_password string) (bool, string) {
	if !validName(name) {
		return false, "Invalid name format\nName should be at least 5 symbols long and shouldn't contain empty space"
	} else if !validEmail(email) {
		return false, "Invalid email format"
	} else if !validPassword(password) {
		return false, "Invalid password format\nPassword should contain at least one number, one uppercase letter, one lowercase letter, one symbol or punctuation and at least 8 symbols"
	} else if password != confirm_password {
		return false, "Passwords don't match"
	}
	return true, ""
}

func validName(name string) bool {
	if len(name) < 5 {
		return false
	}
	for _, ch := range name {
		if ch != ' ' {
			return false
		}
	}
	return true
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
	credentials := entity.Credentials{Email: r.FormValue("email"), Password: r.FormValue("password")}
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
	session := entity.Session{}
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
