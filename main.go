package main

import (
	"context"
	"fmt"
	"github.com/Keith1039/dbvg_usecase/templates"
	"github.com/a-h/templ"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"net/http"
	"os"
)

var (
	dbpool *pgxpool.Pool
	ctx    = context.Background()
	name   string
)

func init() {
	var err error
	err = os.Setenv("DATABASE_URL", "postgres://postgres:localDB12@localhost:5432/testgres?sslmode=disable")
	if err != nil {
		return
	}
	dbpool, err = pgxpool.New(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v\n", err)
	}

	fmt.Println("The database is connected")
}

func main() {
	defer dbpool.Close()
	http.Handle("/", templ.Handler(templates.Welcome()))
	http.Handle("/home", templ.Handler(templates.Home(&name)))
	http.Handle("/signup", templ.Handler(templates.Signup()))
	http.Handle("/login", templ.Handler(templates.Login()))
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))

	http.HandleFunc("/validate-login/", HandleLogin)
	http.HandleFunc("/validate-signup/", HandleSignup)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("listening on port 8080...")
}

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	username := r.PostFormValue("username")
	password := r.PostFormValue("password")
	query := dbpool.QueryRow(ctx, "SELECT USERNAME FROM USERS WHERE USERNAME=$1 AND PASSWORD=$2", username, password)
	err := query.Scan(&name)
	if err != nil { // probably no error row, but I should add the explicit condition
		component := templates.FailedLogin(username, password)
		err = Render(w, r, component)
	} else {
		err = hxRedirect(w, r, "/home")
		if err != nil {
			log.Fatal(err)
		}
	}

}

func HandleSignup(w http.ResponseWriter, r *http.Request) {
	username := r.PostFormValue("username")
	password := r.PostFormValue("password")
	query := dbpool.QueryRow(ctx, "SELECT USERNAME FROM USERS WHERE USERNAME=$1")
	err := query.Scan(&name)
	if err != nil { // new username
		_, err = dbpool.Exec(ctx, "INSERT INTO USERS (USERNAME, PASSWORD) VALUES ($1, $2)", username, password)
		if err != nil {
			log.Fatal(err)
		}
		name = username
		err = hxRedirect(w, r, "/home")
		if err != nil {
			log.Fatal(err)
		}

	} else {
		component := templates.FailedUsername(username)
		err = Render(w, r, component)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func hxRedirect(w http.ResponseWriter, r *http.Request, url string) error {
	if len(r.Header.Get("HX-Request")) > 0 {
		w.Header().Set("HX-Redirect", url)
		w.WriteHeader(http.StatusSeeOther)
		return nil
	}
	http.Redirect(w, r, url, http.StatusSeeOther)
	return nil
}

func Render(w http.ResponseWriter, r *http.Request, c templ.Component) error {
	return c.Render(r.Context(), w)
}
