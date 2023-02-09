package app

import (
	"forum_gateway/internal/entity"
	"net/http"
)

var oauthStateString = "pseudo-random"

// SIGN UP
func (h *Handler) SignUpHandler(w http.ResponseWriter, r *http.Request) {
	if r.Context().Value("authorised") == true {
		h.APIResponse(w, http.StatusForbidden, entity.Response{ErrorMessage: "Forbidden"}, "templates/errors.html")
		return
	}
	switch r.Method {
	case http.MethodGet:
		h.getSignUp(w, r)
	case http.MethodPost:
		h.postSignUp(w, r)
	default:
		h.APIResponse(w, http.StatusMethodNotAllowed, entity.Response{ErrorMessage: "Invalid Method"}, "templates/errors.html")
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
				h.APIResponse(w, http.StatusRequestTimeout, entity.Response{ErrorMessage: "Request Timeout"}, "templates/errors.html")
			default:
				h.APIResponse(w, http.StatusInternalServerError, entity.Response{ErrorMessage: "Internal Server Error"}, "templates/error.html")
			}
		}
	case <-ctx.Done():
		err = ctx.Err()
		h.errLog.Println(err)
		h.APIResponse(w, http.StatusRequestTimeout, entity.Response{ErrorMessage: "Request Timeout"}, "templates/errors.html")
		return
	}
	http.Redirect(w, r, "/sign-in", http.StatusSeeOther)
}

// SIGN IN
func (h *Handler) SignInHandler(w http.ResponseWriter, r *http.Request) {
	if r.Context().Value("authorised") == true {
		h.APIResponse(w, http.StatusForbidden, entity.Response{ErrorMessage: "Forbidden"}, "templates/errors.html")
		return
	}
	switch r.Method {
	case http.MethodGet:
		h.getSignIn(w, r)
	case http.MethodPost:
		h.postSignIn(w, r)
	default:
		h.APIResponse(w, http.StatusMethodNotAllowed, entity.Response{ErrorMessage: "Invalid method"}, "templates/errors.html")
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
	ctx, cancel := getTimeout(r.Context())
	defer cancel()
	sessionChan := make(chan entity.SessionResult)
	var sessionRes entity.SessionResult

	go h.auUcase.SignIn(ctx, credentials, sessionChan)

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
			case entity.ErrNotFound:
				h.APIResponse(w, http.StatusBadRequest, entity.Response{ErrorMessage: "User with a given email doesn't exist"}, "templates/login.html")
			case entity.ErrInvalidPassword:
				h.APIResponse(w, http.StatusUnauthorized, entity.Response{ErrorMessage: "Invalid password"}, "templates/login.html")
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
	h.APIResponse(w, http.StatusInternalServerError, entity.Response{ErrorMessage: "Internal Server Error"}, "templates/errors.html")
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
			http.Redirect(w, r, "/posts", 303)
		default:
			h.APIResponse(w, http.StatusInternalServerError, entity.Response{ErrorMessage: "Internal Server Error"}, "templates/errors.go")
		}
	}
}
