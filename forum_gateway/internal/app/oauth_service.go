package app

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"forum_gateway/internal/entity"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

func (h *Handler) SignInOAuthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Context().Value("authorised") == true {
		h.APIResponse(w, http.StatusForbidden, entity.Response{ErrorMessage: "Forbidden"}, "templates/errors.html")
		return
	}
	method := getOAuthMethod(r.URL.String())
	if method == invalid {
		h.errLog.Println(errors.New("invalid OAuth method"))
		h.APIResponse(w, http.StatusNotFound, entity.Response{ErrorMessage: "Not Found"}, "templates/errors.html")
		return
	}
	url := h.oauths[method].AuthUrl()
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func getOAuthMethod(url string) method {
	temp := strings.Split(strings.Split(url, "?")[0], "/")

	switch temp[len(temp)-1] {
	case "google":
		return google
	case "github":
		return github
	default:
		return invalid
	}
}

func (h *Handler) CallBackHandler(w http.ResponseWriter, r *http.Request) {
	method := getOAuthMethod(r.URL.String())
	if method == invalid {
		h.errLog.Println("invalid OAuth method")
		h.APIResponse(w, http.StatusNotFound, entity.Response{ErrorMessage: "Not Found"}, "templates/errors.html")
		return
	}
	if r.FormValue("state") != h.oauths[method].State() {
		h.errLog.Println("invalid OAuth state")
		http.Redirect(w, r, "/sign-up", http.StatusTemporaryRedirect)
		return
	}
	tokenRes := h.oauths[method].Token(r.FormValue("code"))
	if tokenRes.Error != nil {
		h.errLog.Println(tokenRes.Error)
		http.Redirect(w, r, "/sign-up", http.StatusTemporaryRedirect)
		return
	}
	credsRes := h.oauths[method].Credentials(tokenRes.Token)
	if credsRes.Err != nil {
		h.errLog.Println(credsRes.Err)
		http.Redirect(w, r, "/sign-up", http.StatusTemporaryRedirect)
		return
	}
	ctx, cancel := getTimeout(r.Context())
	defer cancel()
	sessionChan := make(chan entity.SessionResult)
	var sessionRes entity.SessionResult

	go h.auUcase.OAuth(ctx, credsRes.Credentials, sessionChan)

	select {
	case <-ctx.Done():
		err := ctx.Err()
		h.errLog.Println(err)
		h.APIResponse(w, http.StatusRequestTimeout, entity.Response{ErrorMessage: "Request Timeout"}, "templates/errors.html")
		return
	case sessionRes = <-sessionChan:
		err := sessionRes.Err
		if err != nil {
			h.errLog.Println(err)
			switch err {
			case entity.ErrRequestTimeout:
				h.APIResponse(w, http.StatusRequestTimeout, entity.Response{ErrorMessage: "Request Timeout"}, "templates/errors.html")
			default:
				h.APIResponse(w, http.StatusInternalServerError, entity.Response{ErrorMessage: "Internal Server Error"}, "templates/errors.html")
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
		http.Redirect(w, r, "/posts", http.StatusSeeOther)
		return
	}
	h.errLog.Println("internal server error")
	h.APIResponse(w, http.StatusInternalServerError, entity.Response{ErrorMessage: "Internal Server Error"}, "templates/errors.html")
}

type Token struct {
	AccessToken string `json:"access_token"`
}

type TokenResult struct {
	Token Token
	Error error
}

//oauth
type OAuth interface {
	Token(code string) TokenResult
	Credentials(token Token) entity.CredentialsResult
	AuthUrl() string
	State() string
}

func NewOAuth(clientId, clientSecret, state string, m method) (OAuth, error) {
	switch m {
	case github:
		return NewGitHub(clientId, clientSecret, state), nil
	case google:
		return NewGoogle(clientId, clientSecret, state), nil
	default:
		return nil, errors.New("invalid method")
	}
}

type method string

const (
	github  method = "GitHub"
	google         = "Google"
	invalid        = "Invalid"
)

//github oauth
type GitHub struct {
	clientId     string
	clientSecret string
	state        string
	authUrl      string
}

func NewGitHub(clientId, clientSecret, state string) *GitHub {
	gh := GitHub{
		clientId:     clientId,
		clientSecret: clientSecret,
		state:        state,
	}
	gh.setAuthUrl()
	return &gh
}

func (gh *GitHub) AuthUrl() string {
	return gh.authUrl
}

func (gh *GitHub) setAuthUrl() {
	var buf bytes.Buffer
	buf.WriteString("https://github.com/login/oauth/authorize")
	v := url.Values{"client_id": {gh.clientId}}
	v.Set("redirect_uri", "http://localhost:8082/callback/github")
	v.Set("scope", "user")
	v.Set("state", gh.state)
	buf.WriteByte('?')
	buf.WriteString(v.Encode())
	gh.authUrl = buf.String()
}

func (gh *GitHub) Token(code string) TokenResult {
	var buf bytes.Buffer
	buf.WriteString("https://github.com/login/oauth/access_token?")
	v := url.Values{"code": {code}}
	v.Set("redirect_uri", "http://localhost:8082/callback/github")
	v.Set("client_id", gh.clientId)
	v.Set("client_secret", gh.clientSecret)
	buf.WriteString(v.Encode())
	url := buf.String()
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return TokenResult{Error: err}
	}
	req.Header.Set("Accept", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return TokenResult{Error: err}
	}
	defer resp.Body.Close()
	var token Token
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return TokenResult{Error: err}
	}
	err = json.Unmarshal(bytes, &token)
	if err != nil {
		return TokenResult{Error: err}
	}
	return TokenResult{Token: token}
}

func (gh *GitHub) Credentials(token Token) entity.CredentialsResult {
	info, err := gh.getScopeInfo(token, "https://api.github.com/user")
	if err != nil {
		return entity.CredentialsResult{Err: err}
	}
	creds := entity.Credentials{}
	err = json.Unmarshal(info, &creds)
	if err != nil {
		return entity.CredentialsResult{Err: err}
	}
	if creds.Email == "" {
		info, err = gh.getScopeInfo(token, "https://api.github.com/user/emails") //gets email if not visible
		tempCreds := []entity.Credentials{}
		err = json.Unmarshal(info, &tempCreds)
		if err != nil {
			return entity.CredentialsResult{Err: err}
		}
		creds.Email = tempCreds[0].Email
	}
	creds.Name = creds.Login
	return entity.CredentialsResult{Credentials: creds}
}

func (gh *GitHub) State() string {
	return gh.state
}

func (gh *GitHub) getScopeInfo(token Token, url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("bearer %s", token.AccessToken))
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

//google oauth
type Google struct {
	clientId     string
	clientSecret string
	state        string
	authUrl      string
}

func NewGoogle(clientId, clientSecret, state string) *Google {
	g := Google{
		clientId:     clientId,
		clientSecret: clientSecret,
		state:        state,
	}
	g.setAuthUrl()
	return &g
}
func (g *Google) AuthUrl() string {
	return g.authUrl
}

func (g *Google) setAuthUrl() {
	var buf bytes.Buffer
	buf.WriteString("https://accounts.google.com/o/oauth2/auth")
	v := url.Values{"response_type": {"code"}, "client_id": {g.clientId}}
	v.Set("redirect_uri", "http://localhost:8082/callback/google")
	v.Set("scope", "https://www.googleapis.com/auth/userinfo.email https://www.googleapis.com/auth/userinfo.profile")
	v.Set("state", g.state)
	buf.WriteByte('?')
	buf.WriteString(v.Encode())
	g.authUrl = buf.String()
}

func (gh *Google) Token(code string) TokenResult {
	var buf bytes.Buffer
	buf.WriteString("https://oauth2.googleapis.com/token?")
	v := url.Values{"grant_type": {"authorization_code"}, "code": {code}}
	v.Set("redirect_uri", "http://localhost:8082/callback/google")
	v.Set("client_id", gh.clientId)
	v.Set("client_secret", gh.clientSecret)
	buf.WriteString(v.Encode())
	url := buf.String()
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return TokenResult{Error: err}
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return TokenResult{Error: err}
	}
	defer resp.Body.Close()
	var token Token
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return TokenResult{Error: err}
	}
	err = json.Unmarshal(bytes, &token)
	if err != nil {
		return TokenResult{Error: err}
	}
	return TokenResult{Token: token}
}

func (gh *Google) Credentials(token Token) entity.CredentialsResult {
	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		return entity.CredentialsResult{Err: err}
	}
	defer response.Body.Close()
	content, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return entity.CredentialsResult{Err: err}
	}
	creds := entity.Credentials{}
	err = json.Unmarshal(content, &creds)
	if err != nil {
		return entity.CredentialsResult{Err: err}
	}
	return entity.CredentialsResult{Credentials: creds}
}
func (g *Google) State() string {
	return g.state
}
