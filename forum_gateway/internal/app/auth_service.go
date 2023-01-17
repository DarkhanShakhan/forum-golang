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
	h.APIResponse(w, http.StatusOK, entity.Response{}, "web/sign_up.html")
}

func (h *Handler) postSignUp(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	credentials := entity.GetCredentials(r)
	confirm_password := r.FormValue("confirm_password")
	ok, message := credentials.ValidateSignUp(confirm_password)
	if !ok {
		h.errLog.Println(message)
		h.APIResponse(w, http.StatusBadRequest, entity.Response{ErrorMessage: message}, "web/sign_up.html")
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
				h.APIResponse(w, http.StatusBadRequest, entity.Response{ErrorMessage: "user with a given email already exists"}, "web/sign_up.html")
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
	h.APIResponse(w, http.StatusOK, entity.Response{}, "web/sign_in.html")
}

func (h *Handler) postSignIn(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	credentials := entity.GetCredentials(r)
	ok, message := credentials.ValidateSignIn()
	if !ok {
		h.errLog.Println(message)
		h.APIResponse(w, http.StatusBadRequest, entity.Response{ErrorMessage: message}, "web/sign_in.html")
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
			case entity.ErrNotFound:
				h.APIResponse(w, http.StatusBadRequest, entity.Response{ErrorMessage: "User with a given email doesn't exist"}, "web/sign_in.html")
			case entity.ErrInvalidPassword:
				h.APIResponse(w, http.StatusBadRequest, entity.Response{ErrorMessage: "Invalid password"}, "web/sign_in.html")
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
		http.Redirect(w, r, "/posts", 302)
		return
	}
	h.APIResponse(w, http.StatusInternalServerError, entity.Response{ErrorMessage: "Internal Server Error"}, "web/error.html")
}

func (h *Handler) SignOutHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Context().Value("authorised"))
	if r.Context().Value("authorised") == false {
		h.APIResponse(w, http.StatusForbidden, entity.Response{ErrorMessage: "Forbidden"}, "web/error.html")
		return
	}
	if r.Method != http.MethodDelete {
		h.APIResponse(w, http.StatusBadRequest, entity.Response{ErrorMessage: "Bad request"}, "web/error.html")
		return
	}

	cookie, err := r.Cookie("token")
	if err != nil {
		h.APIResponse(w, http.StatusForbidden, entity.Response{ErrorMessage: "Forbidden"}, "web/error.html")
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
		h.APIResponse(w, http.StatusRequestTimeout, entity.Response{ErrorMessage: "Request Timeout"}, "web/error.html")
		return
	case err = <-errChan:
		switch err {
		case entity.ErrRequestTimeout:
			h.APIResponse(w, http.StatusRequestTimeout, entity.Response{ErrorMessage: "Request Timeout"}, "web/error.go")
		case nil:
			http.Redirect(w, r, "/posts", 302) // FIXME: set cookie to empty session
		default:
			h.APIResponse(w, http.StatusInternalServerError, entity.Response{ErrorMessage: "Internal Server Error"}, "web/error.go")
		}
	}
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
