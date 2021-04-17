package register

import (
	"database/sql"
	"encoding/json"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/sessions"
	"github.com/kyleparisi/expiration.dev/src/app"
	"github.com/kyleparisi/expiration.dev/src/framework"
	"github.com/ucarion/urlpath"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"os"
	"io"
	"log"
	"net/mail"
	"strings"
)

type Register struct {
	Email    string
	Password string
}

type RegisterError struct {
	Email  string `json:email`
	Errors struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	} `json:"errors"`
}

func GetHandler(db *sql.DB, session *sessions.Session) func(_ urlpath.Match) framework.Response {
	return func (_ urlpath.Match) framework.Response {
		t, err := template.ParseFiles(os.Getenv("APP_DIR") + "/views/register.gohtml")
		if err != nil {
			panic(err)
		}
		return framework.Response{StatusCode: 200, Template: t}
	}
}

func PostHandler(db *sql.DB, session *sessions.Session) func(_ urlpath.Match, body io.Reader) framework.Response {
	return func(_ urlpath.Match, body io.Reader) framework.Response {
		register := Register{}
		t, err := template.ParseFiles(os.Getenv("APP_DIR") + "/views/register.gohtml")
		if err != nil {
			panic(err)
		}
		err = json.NewDecoder(body).Decode(&register)
		if err != nil {
			panic(err.Error())
		}
		hasEmail := register.Email != ""
		hasPassword := register.Password != ""
		// Input validation
		if !hasEmail || !hasPassword {
			registerError := RegisterError{Errors: struct {
				Email    string `json:"email"`
				Password string `json:"password"`
			}{Email: "", Password: ""}}
			if !hasEmail {
				registerError.Errors.Email = "Please provide an email address"
			}
			if !hasPassword {
				registerError.Errors.Password = "Please provide a password"
			}
			log.Printf("RegisterHandler: %+v", registerError)
			return framework.Response{StatusCode: 400, Data: registerError, Template: t}
		}
		// Email validation
		_, err = mail.ParseAddress(register.Email)
		if err != nil {
			registerError := RegisterError{
				Email: register.Email,
				Errors: struct {
					Email    string `json:"email"`
					Password string `json:"password"`
				}{Email: "Not a valid email address", Password: ""}}
			log.Printf("RegisterHandler: %+v", registerError)
			return framework.Response{StatusCode: 400, Data: registerError, Template: t}
		}
		// Check for existing user
		var appUser app.AppUser
		err = db.QueryRow("SELECT id, `name`, email, password FROM `app_user` where email = ?", strings.ToLower(register.Email)).Scan(&appUser.Id, &appUser.Name, &appUser.Email, &appUser.Password)
		switch {
		case err == sql.ErrNoRows:
			log.Printf("RegisterHandler: no user with email: %s\n", register.Email)
		case err != nil:
			log.Fatalf("RegisterHandler: query error: %v\n", err)
		default:
			log.Printf("RegisterHandler: user already registered: %s\n", register.Email)
			registerError := RegisterError{
				Email: register.Email,
				Errors: struct {
					Email    string `json:"email"`
					Password string `json:"password"`
				}{Email: "User already exists", Password: ""}}
			return framework.Response{StatusCode: 400, Data: registerError, Template: t}
		}

		hash, hashError := bcrypt.GenerateFromPassword([]byte(register.Password), bcrypt.DefaultCost)
		if hashError != nil {
			log.Printf("RegisterHandler: failed to make hash password: %s\n", err.Error())
			return framework.Response{StatusCode: 500}
		}
		res, err := db.Exec("INSERT INTO app_user(`name`, email, password) VALUES('', ?, ?)", strings.ToLower(register.Email), string(hash))
		if err != nil {
			log.Printf("RegisterHandler: failed to register user: %s\n", err.Error())
			return framework.Response{StatusCode: 500}
		}
		log.Printf("RegisterHandler: new user registered: %s\n", register.Email)
		lastId, err := res.LastInsertId()
		if err != nil {
			panic(err)
		}
		err = db.QueryRow("SELECT id, `name`, email from `app_user` where id = ?", lastId).Scan(&appUser.Id, &appUser.Name, &appUser.Email)
		if err != nil {
			panic(err)
		}
		session.Values["Id"] = appUser.Id
		session.Values["Name"] = appUser.Name
		session.Values["Email"] = appUser.Email

		return framework.Response{Redirect: "/dashboard"}
	}
}
