package framework

import (
	"database/sql"
	"fmt"
	"github.com/gorilla/sessions"
	"github.com/srinathgs/mysqlstore"
	"github.com/ucarion/urlpath"
	"html/template"
	"io"
	"net/http"
	"os"
)

type Router struct {
	GET  []map[string]Handler
	POST []map[string]PostHandler
}
type Redirect struct {
	location string
	code     int
}
type Response struct {
	StatusCode int
	Data       interface{}
	Template   *template.Template
	Redirect   interface{}
	Session    *sessions.Session
}
type Handler func(*sql.DB, *sessions.Session) func(urlpath.Match) Response
type PostHandler func(*sql.DB, *sessions.Session) func(urlpath.Match, io.Reader) Response

func NewRouter() *Router {
	router := new(Router)
	router.GET = make([]map[string]Handler, 0)
	router.POST = make([]map[string]PostHandler, 0)
	return router
}

func NewDatabaseConnection() *sql.DB {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/%s", os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_DATABASE")))
	if err != nil {
		panic(err.Error())
	}
	return db
}

func NewSessionStore() *mysqlstore.MySQLStore {
	dsn := fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/%s?parseTime=true&loc=Local", os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_DATABASE"))
	store, err := mysqlstore.NewMySQLStore(dsn, "session", "/", 3600, []byte(os.Getenv("APP_SECRET_KEY")))
	if err != nil {
		panic(err.Error())
	}
	return store
}

func (r *Router) Get(path string, f Handler) {
	r.GET = append(r.GET, map[string]Handler{path: f})
}

func (r *Router) Post(path string, f PostHandler) {
	r.POST = append(r.POST, map[string]PostHandler{path: f})
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	fmt.Println(req.Method, req.URL.Path)
	db := NewDatabaseConnection()
	defer db.Close()
	store := NewSessionStore()
	defer store.Close()

	session, _ := store.Get(req, "user")

	if req.Method == "GET" {
		for _, element := range r.GET {
			for path, handler := range element {
				route := urlpath.New(path)
				match, ok := route.Match(req.URL.Path)
				if !ok {
					break
				}
				response := handler(db, session)(match)
				// check for redirect first
				if redirect, ok := response.Redirect.(Redirect); ok {
					http.Redirect(w, req, redirect.location, redirect.code)
					goto Done
				}
				response.Template.Execute(w, response.Data)
				goto Done
			}
		}
	}

	if req.Method == "POST" {
		for _, element := range r.POST {
			for path, post := range element {
				route := urlpath.New(path)
				match, ok := route.Match(req.URL.Path)
				if !ok {
					break
				}
				response := post(db, session)(match, req.Body)
				if response.Session != nil {
					err := response.Session.Save(req, w)
					if err != nil {
						fmt.Println("problem saving session", err.Error())
					}
				}
				// check for redirect first
				if redirect, ok := response.Redirect.(Redirect); ok {
					http.Redirect(w, req, redirect.location, redirect.code)
					goto Done
				}
				if response.Template != nil {
					err := response.Template.Execute(w, response.Data)
					if err != nil {
						fmt.Println(err)
					}
				}
				goto Done
			}
		}
	}

	// 404
	http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
Done:
}
