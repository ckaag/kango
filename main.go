package main

import (
	"context"
	"crypto/sha256"
	"crypto/subtle"
	"database/sql"
	"fmt"
	"kanco/kango/tutorial"
	"net/http"
	"strconv"

	_ "embed"
	"log"
	"reflect"

	_ "github.com/mattn/go-sqlite3"
)

//go:embed schema.sql
var ddl string

func main() {
	app := new(application)
	// to test basic auth:
	app.auth.username = "admin"
	app.auth.password = "test"

	run()

	http.HandleFunc("/", app.basicAuth(func(w http.ResponseWriter, r *http.Request) {
		name := r.Header.Get("username")
		counter, _ := getCounter(name)
		projects, _ := getProjects()
		mainPage(name, strconv.FormatInt(counter, 10), projects).Render(r.Context(), w)
	}))

	http.HandleFunc("/counter/increment", app.basicAuth(func(w http.ResponseWriter, r *http.Request) {
		name := r.Header.Get("username")
		counter, _ := incrementCounter(name)
		fmt.Println("Incrementing counter")
		clickCounter(strconv.FormatInt(counter, 10)).Render(r.Context(), w)
	}))

	http.HandleFunc("/projects/add", app.basicAuth(func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		name := r.Form.Get("title")
		fmt.Println("Adding Project: " + name)
		err := InsertProject(name)
		if err != nil {
			fmt.Println("Error 1", err)
		}
		projects, _ := getProjects()
		projectCrud(projects).Render(r.Context(), w)
	}))

	http.HandleFunc("/projects/{id}/delete", app.basicAuth(func(w http.ResponseWriter, r *http.Request) {
		idString := r.PathValue("id")
		fmt.Println("Deleting Project: " + idString)
		id, err := strconv.ParseInt(idString, 10, 0)
		if err != nil {
			fmt.Println("Error 1", err)
		} else {
			err := DeleteProject(id)
			if err != nil {
				fmt.Println("Error 1", err)
			}
		}
		projects, _ := getProjects()
		projectCrud(projects).Render(r.Context(), w)
	}))

	http.ListenAndServe("localhost:8080", nil)
}

func incrementCounter(name string) (int64, error) {
	ctx := context.Background()

	fmt.Println("Getting counter")
	oldValue, err := getCounter(name)
	if err != nil {
		return 0, err
	}

	fmt.Println("Got Counter")

	db, err := sql.Open("sqlite3", "file:.development.local.db")
	if err != nil {
		return 0, err
	}

	// create tables
	if _, err := db.ExecContext(ctx, ddl); err != nil {
		return 0, err
	}

	fmt.Println("Tables Created")

	queries := tutorial.New(db)

	fmt.Println("Before: ", oldValue)

	var ins tutorial.IncrementCounterAndReturnParams = tutorial.IncrementCounterAndReturnParams{Name: name, Counter: oldValue + 1}

	fmt.Printf("%+v\n", ins)

	count, err := queries.IncrementCounterAndReturn(ctx, ins)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func getProjects() ([]tutorial.Project, error) {
	ctx := context.Background()

	db, err := sql.Open("sqlite3", "file:.development.local.db")
	if err != nil {
		fmt.Println("Error 1", err)
		return nil, err
	}

	// create tables
	if _, err := db.ExecContext(ctx, ddl); err != nil {
		fmt.Println("Error 2", err)
		return nil, err
	}

	queries := tutorial.New(db)

	// list all authors
	count, err := queries.GetListOfProjects(ctx)
	if err != nil {
		fmt.Println("Error 3", err)
		return nil, err
	}
	return count, nil
}

func InsertProject(name string) error {
	ctx := context.Background()

	db, err := sql.Open("sqlite3", "file:.development.local.db")
	if err != nil {
		fmt.Println("Error 1", err)
		return err
	}

	// create tables
	if _, err := db.ExecContext(ctx, ddl); err != nil {
		fmt.Println("Error 2", err)
		return err
	}

	queries := tutorial.New(db)

	// list all authors
	err = queries.AddProject(ctx, name)
	if err != nil {
		fmt.Println("Error 3", err)
		return err
	}
	return nil
}

func DeleteProject(id int64) error {
	ctx := context.Background()

	db, err := sql.Open("sqlite3", "file:.development.local.db")
	if err != nil {
		fmt.Println("Error 1", err)
		return err
	}

	// create tables
	if _, err := db.ExecContext(ctx, ddl); err != nil {
		fmt.Println("Error 2", err)
		return err
	}

	queries := tutorial.New(db)

	// list all authors
	err = queries.DeleteProject(ctx, id)
	if err != nil {
		fmt.Println("Error 3", err)
		return err
	}
	return nil
}

func getCounter(name string) (int64, error) {
	ctx := context.Background()

	db, err := sql.Open("sqlite3", "file:.development.local.db")
	if err != nil {
		fmt.Println("Error 1", err)
		return 0, err
	}

	// create tables
	if _, err := db.ExecContext(ctx, ddl); err != nil {
		fmt.Println("Error 2", err)
		return 0, err
	}

	queries := tutorial.New(db)

	// list all authors
	count, err := queries.GetCounter(ctx, name)
	if err != nil {
		fmt.Println("Error 3", err)
		return 0, err
	}
	return count, nil
}

func run() error {
	ctx := context.Background()

	db, err := sql.Open("sqlite3", "file:.development.local.db")
	if err != nil {
		fmt.Println("Error I", err)
		return err
	}

	// create tables
	if _, err := db.ExecContext(ctx, ddl); err != nil {
		fmt.Println("Error I", err)
		return err
	}

	queries := tutorial.New(db)

	// list all authors
	authors, err := queries.ListAuthors(ctx)
	if err != nil {
		fmt.Println("Error I", err)
		return err
	}
	log.Println(authors)

	// create an author
	insertedAuthor, err := queries.CreateAuthor(ctx, tutorial.CreateAuthorParams{
		Name: "Brian Kernighan",
		Bio:  sql.NullString{String: "Co-author of The C Programming Language and The Go Programming Language", Valid: true},
	})
	if err != nil {
		fmt.Println("Error I", err)
		return err
	}
	log.Println(insertedAuthor)

	// get the author we just inserted
	fetchedAuthor, err := queries.GetAuthor(ctx, insertedAuthor.ID)
	if err != nil {
		fmt.Println("Error I", err)
		return err
	}

	// prints true
	log.Println(reflect.DeepEqual(insertedAuthor, fetchedAuthor))
	return nil
}

type application struct {
	auth struct {
		username string
		password string
	}
}

// Based on: https://www.alexedwards.net/blog/basic-authentication-in-go
func (app *application) basicAuth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if ok {
			nameHash := sha256.Sum256([]byte(username))
			secretHash := sha256.Sum256([]byte(password))
			expectedName := sha256.Sum256([]byte(app.auth.username))
			expectedSecret := sha256.Sum256([]byte(app.auth.password))

			usernameMatch := (subtle.ConstantTimeCompare(nameHash[:], expectedName[:]) == 1)
			passwordMatch := (subtle.ConstantTimeCompare(secretHash[:], expectedSecret[:]) == 1)

			if usernameMatch && passwordMatch {
				r.Header.Set("username", username)
				next.ServeHTTP(w, r)
				return
			}
		}

		w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	})
}
