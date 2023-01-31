package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"forum_gateway/internal/entity"
	"io/ioutil"
	"net/http"
	"net/url"
)

var oauthStateString = "pseudo-random"

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
	h.APIResponse(w, http.StatusOK, entity.Response{}, "templates/registration.html")
}

func (h *Handler) postSignUp(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	credentials := entity.GetCredentials(r)
	confirm_password := r.FormValue("confirm_password")
	ok, message := credentials.ValidateSignUp(confirm_password)
	if !ok {
		h.errLog.Println(message)
		h.APIResponse(w, http.StatusBadRequest, entity.Response{ErrorMessage: message}, "templates/registration.html")
		return
	}
	ctx, cancel := getTimeout(r.Context())
	defer cancel()
	errChan := make(chan error)
	var err error
	go h.auUcase.SignUp(ctx, credentials, errChan)
	select {
	case err = <-errChan:
		if err != nil {
			h.errLog.Println(err)
			switch err {
			case entity.ErrEmailExists:
				h.APIResponse(w, http.StatusBadRequest, entity.Response{ErrorMessage: "user with a given email already exists"}, "templates/registration.html")
				return
			case entity.ErrRequestTimeout:
				h.APIResponse(w, http.StatusRequestTimeout, entity.Response{ErrorMessage: "Request Timeout"}, "web/error.go")
			default:
				h.APIResponse(w, http.StatusInternalServerError, entity.Response{ErrorMessage: "Internal Server Error"}, "web/errog.go")
			}
		}
	case <-ctx.Done():
		err = ctx.Err()
		h.errLog.Println(err)
		h.APIResponse(w, http.StatusRequestTimeout, entity.Response{ErrorMessage: "Request Timeout"}, "web/error.html")
		return
	}
	http.Redirect(w, r, "/sign_in", http.StatusFound)
}

// SIGN IN
func (h *Handler) SignInHandler(w http.ResponseWriter, r *http.Request) {
	if r.Context().Value("authorised") == true {
		h.APIResponse(w, http.StatusForbidden, entity.Response{ErrorMessage: "Forbidden"}, "web/error.html")
		return
	}
	switch r.Method {
	case http.MethodGet:
		h.getSignIn(w, r)
	case http.MethodPost:
		h.postSignIn(w, r)
	default:
		h.APIResponse(w, http.StatusBadRequest, entity.Response{ErrorMessage: "Bad request"}, "web/error.html")
	}
}

func (h *Handler) getSignIn(w http.ResponseWriter, r *http.Request) {
	h.APIResponse(w, http.StatusOK, entity.Response{}, "templates/login.html")
}

func (h *Handler) postSignIn(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	credentials := entity.GetCredentials(r)
	ok, message := credentials.ValidateSignIn()
	if !ok {
		h.errLog.Println(message)
		h.APIResponse(w, http.StatusBadRequest, entity.Response{ErrorMessage: message}, "templates/login.html")
		return
	}
	// FIXME: validate credentials
	ctx, cancel := getTimeout(r.Context())
	defer cancel()
	sessionChan := make(chan entity.SessionResult)
	var sessionRes entity.SessionResult

	go h.auUcase.SignIn(ctx, credentials, sessionChan)

	select {
	case <-ctx.Done():
		err := ctx.Err()
		h.errLog.Println(err)
		h.APIResponse(w, http.StatusRequestTimeout, entity.Response{ErrorMessage: "Request Timeout"}, "web/error.html")
		return
	case sessionRes = <-sessionChan:
		err := sessionRes.Err
		if err != nil {
			h.errLog.Println(err)
			switch err {
			//FIXME: error for unauthorised user?? 401
			case entity.ErrNotFound:
				h.APIResponse(w, http.StatusBadRequest, entity.Response{ErrorMessage: "User with a given email doesn't exist"}, "templates/login.html")
			case entity.ErrInvalidPassword:
				h.APIResponse(w, http.StatusUnauthorized, entity.Response{ErrorMessage: "Invalid password"}, "templates/login.html")
			case entity.ErrRequestTimeout:
				h.APIResponse(w, http.StatusRequestTimeout, entity.Response{ErrorMessage: "Request Timeout"}, "web/error.html")
			default:
				h.APIResponse(w, http.StatusInternalServerError, entity.Response{ErrorMessage: "Internal Server Error"}, "web/error.html")
			}
			return
		}
	}
	if sessionRes.Session.Token != "" {
		cookie := http.Cookie{
			Name:    "token",
			Expires: sessionRes.Session.ExpiryTime,
			Value:   sessionRes.Session.Token,
			Path:    "/",
		}
		http.SetCookie(w, &cookie)
		http.Redirect(w, r, "/posts", 303)
		return
	}
	h.APIResponse(w, http.StatusInternalServerError, entity.Response{ErrorMessage: "Internal Server Error"}, "web/error.html")
}

func (h *Handler) SignOutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Context().Value("authorised") == false {
		h.APIResponse(w, http.StatusForbidden, entity.Response{ErrorMessage: "Forbidden"}, "templates/errors.html")
		return
	}
	if r.Method != http.MethodPost {
		h.APIResponse(w, http.StatusMethodNotAllowed, entity.Response{ErrorMessage: "Invalid Method"}, "templates/errors.html")
		return
	}

	cookie, err := r.Cookie("token")
	if err != nil {
		h.APIResponse(w, http.StatusForbidden, entity.Response{ErrorMessage: "Forbidden"}, "templates/errors.html")
		return
	}
	token := cookie.Value
	session := entity.Session{Token: token}
	ctx, cancel := getTimeout(r.Context())
	defer cancel()
	errChan := make(chan error)
	go h.auUcase.SignOut(ctx, session, errChan)
	select {
	case <-ctx.Done():
		err = ctx.Err()
		h.errLog.Println(err)
		h.APIResponse(w, http.StatusRequestTimeout, entity.Response{ErrorMessage: "Request Timeout"}, "templates/errors.html")
		return
	case err = <-errChan:
		switch err {
		case entity.ErrRequestTimeout:
			h.APIResponse(w, http.StatusRequestTimeout, entity.Response{ErrorMessage: "Request Timeout"}, "templates/errors.go")
		case nil:
			http.Redirect(w, r, "/posts", 303) // FIXME: set cookie to empty session
		default:
			h.APIResponse(w, http.StatusInternalServerError, entity.Response{ErrorMessage: "Internal Server Error"}, "templates/errors.go")
		}
	}
}

func (h *Handler) SignInGoogleHandler(w http.ResponseWriter, r *http.Request) {
	if r.Context().Value("authorised") == true {
		h.APIResponse(w, http.StatusForbidden, entity.Response{ErrorMessage: "Forbidden"}, "templates/errors.html")
		return
	}
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

func (h *Handler) GoogleCallbackHandler(w http.ResponseWriter, r *http.Request) {
	content, err := getUserInfo(r.FormValue("state"), r.FormValue("code"), "https://www.googleapis.com/oauth2/v2/userinfo?access_token=")
	if err != nil {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	creds := entity.Credentials{}
	json.Unmarshal(content, &creds)
	ctx, cancel := getTimeout(r.Context())
	defer cancel()
	sessionChan := make(chan entity.SessionResult)
	var sessionRes entity.SessionResult

	go h.auUcase.OAuth(ctx, creds, sessionChan)

	select {
	case <-ctx.Done():
		err := ctx.Err()
		h.errLog.Println(err)
		h.APIResponse(w, http.StatusRequestTimeout, entity.Response{ErrorMessage: "Request Timeout"}, "web/error.html")
		return
	case sessionRes = <-sessionChan:
		err := sessionRes.Err
		if err != nil {
			h.errLog.Println(err)
			switch err {
			case entity.ErrNotFound:
				h.APIResponse(w, http.StatusBadRequest, entity.Response{ErrorMessage: "User with a given email doesn't exist"}, "templates/login.html")
			case entity.ErrInvalidPassword:
				h.APIResponse(w, http.StatusBadRequest, entity.Response{ErrorMessage: "Invalid password"}, "templates/login.html")
			case entity.ErrRequestTimeout:
				h.APIResponse(w, http.StatusRequestTimeout, entity.Response{ErrorMessage: "Request Timeout"}, "web/error.html")
			default:
				h.APIResponse(w, http.StatusInternalServerError, entity.Response{ErrorMessage: "Internal Server Error"}, "web/error.html")
			}
			return
		}
	}
	if sessionRes.Session.Token != "" {
		cookie := http.Cookie{
			Name:    "token",
			Expires: sessionRes.Session.ExpiryTime,
			Value:   sessionRes.Session.Token,
			Path:    "/",
		}
		http.SetCookie(w, &cookie)
		http.Redirect(w, r, "/posts", 303)
		return
	}
	h.APIResponse(w, http.StatusInternalServerError, entity.Response{ErrorMessage: "Internal Server Error"}, "web/error.html")
}

func getUserInfo(state string, code string, url string) ([]byte, error) {
	if state != oauthStateString {
		return nil, fmt.Errorf("invalid oauth state")
	}
	token, err := exchange(code)
	if err != nil {
		return nil, fmt.Errorf("code exchange failed: %s", err.Error())
	}
	response, err := http.Get(url + token.AccessToken)
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

func (h *Handler) SignInGithubHandler(w http.ResponseWriter, r *http.Request) {
	if r.Context().Value("authorised") == true {
		h.APIResponse(w, http.StatusForbidden, entity.Response{ErrorMessage: "Forbidden"}, "templates/errors.html")
		return
	}
	url := AuthCodeURLGit(oauthStateString)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func AuthCodeURLGit(state string) string {
	var buf bytes.Buffer
	buf.WriteString("https://github.com/login/oauth/authorize")
	v := url.Values{"client_id": {""}}
	v.Set("redirect_uri", "http://localhost:8082/github_callback")
	v.Set("scope", "repo, user") //FIXME: problem with email
	v.Set("state", state)
	buf.WriteByte('?')
	buf.WriteString(v.Encode())
	return buf.String()
}
func exchangeGit(code string) (Token, error) {
	var buf bytes.Buffer
	buf.WriteString("https://github.com/login/oauth/access_token?")
	v := url.Values{"code": {code}}
	v.Set("redirect_uri", "http://localhost:8082/github_callback")
	v.Set("client_id", "")
	v.Set("client_secret", "")
	buf.WriteString(v.Encode())
	url := buf.String()
	req, _ := http.NewRequest("POST", url, nil)
	req.Header.Set("Accept", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return Token{}, err
	}
	defer resp.Body.Close()
	var token Token
	bytes, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(bytes))
	json.Unmarshal(bytes, &token)
	return token, nil
}

func (h *Handler) GithubCallbackHandler(w http.ResponseWriter, r *http.Request) {
	content, err := getUserInfoGit(r.FormValue("state"), r.FormValue("code"), "https://api.github.com/user")
	if err != nil {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	creds := entity.Credentials{}
	//FIXME: verify credentials
	json.Unmarshal(content, &creds)
	ctx, cancel := getTimeout(r.Context())
	defer cancel()
	sessionChan := make(chan entity.SessionResult)
	var sessionRes entity.SessionResult

	go h.auUcase.OAuth(ctx, creds, sessionChan)

	select {
	case <-ctx.Done():
		err := ctx.Err()
		h.errLog.Println(err)
		h.APIResponse(w, http.StatusRequestTimeout, entity.Response{ErrorMessage: "Request Timeout"}, "web/error.html")
		return
	case sessionRes = <-sessionChan:
		err := sessionRes.Err
		if err != nil {
			h.errLog.Println(err)
			switch err {
			case entity.ErrNotFound:
				h.APIResponse(w, http.StatusBadRequest, entity.Response{ErrorMessage: "User with a given email doesn't exist"}, "templates/login.html")
			case entity.ErrInvalidPassword:
				h.APIResponse(w, http.StatusBadRequest, entity.Response{ErrorMessage: "Invalid password"}, "templates/login.html")
			case entity.ErrRequestTimeout:
				h.APIResponse(w, http.StatusRequestTimeout, entity.Response{ErrorMessage: "Request Timeout"}, "web/error.html")
			default:
				h.APIResponse(w, http.StatusInternalServerError, entity.Response{ErrorMessage: "Internal Server Error"}, "web/error.html")
			}
			return
		}
	}
	if sessionRes.Session.Token != "" {
		cookie := http.Cookie{
			Name:    "token",
			Expires: sessionRes.Session.ExpiryTime,
			Value:   sessionRes.Session.Token,
			Path:    "/",
		}
		http.SetCookie(w, &cookie)
		http.Redirect(w, r, "/posts", 303)
		return
	}
	h.APIResponse(w, http.StatusInternalServerError, entity.Response{ErrorMessage: "Internal Server Error"}, "web/error.html")
}

func getUserInfoGit(state string, code string, url string) ([]byte, error) {
	if state != oauthStateString {
		return nil, fmt.Errorf("invalid oauth state")
	}
	token, err := exchangeGit(code)
	if err != nil {
		return nil, fmt.Errorf("code exchange failed: %s", err.Error())
	}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", fmt.Sprintf("bearer %s", token.AccessToken))
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed getting user info: %s", err.Error())
	}
	defer response.Body.Close()
	fmt.Println(response.StatusCode)
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed reading response body: %s", err.Error())
	}
	fmt.Println(string(contents))
	return contents, nil
}

type Token struct {
	AccessToken string `json:"access_token"`
}
