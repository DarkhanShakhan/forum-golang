package entity

import (
	"fmt"
	"net/http"
	"regexp"
	"unicode"
)

type Credentials struct {
	Name     string `json:"name,omitempty"`
	Login    string `json:"login,omitempty"`
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

type CredentialsResult struct {
	Credentials Credentials
	Err         error
}

func GetCredentials(r *http.Request) Credentials {
	return Credentials{
		Name:     r.FormValue("name"),
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}
}

func (c Credentials) ValidateSignUp(confirm_password string) (bool, string) {
	if !validName(c.Name) {
		return false, "Invalid name format\nName should be at least 5 symbols long and shouldn't contain empty space"
	} else if !validEmail(c.Email) {
		return false, "Invalid email format"
	} else if !validPassword(c.Password) {
		return false, "Invalid password format\nPassword should contain at least one number, one uppercase letter, one lowercase letter, one symbol or punctuation and at least 8 symbols"
	} else if c.Password != confirm_password {
		fmt.Println(c.Password)
		fmt.Println(confirm_password)
		return false, "Passwords don't match"
	}
	return true, ""
}

func (c Credentials) ValidateSignIn() (bool, string) {
	if !validEmail(c.Email) {
		return false, "Invalid email format"
	} else if !validPassword(c.Password) {
		return false, "Invalid password format\nPassword should contain at least one number, one uppercase letter, one lowercase letter, one symbol or punctuation and at least 8 symbols"
	}
	return true, ""
}

func validName(name string) bool {
	if len(name) < 5 {
		return false
	}
	for _, ch := range name {
		if ch == ' ' {
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
