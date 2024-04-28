package main

import (
	"context"
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

	run()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		name := "Max Mustermann"
		counter, _ := getCounter(name)
		mainPage(name, strconv.FormatInt(counter, 10)).Render(r.Context(), w)
	})

	http.HandleFunc("/counter/increment", func(w http.ResponseWriter, r *http.Request) {
		name := "Max Mustermann"
		counter, _ := incrementCounter(name)
		fmt.Println("Incrementing counter")
		clickCounter(strconv.FormatInt(counter, 10)).Render(r.Context(), w)
	})

	//fmt.Println("Listening on http://localhost:8080")
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
		return err
	}

	// create tables
	if _, err := db.ExecContext(ctx, ddl); err != nil {
		return err
	}

	queries := tutorial.New(db)

	// list all authors
	authors, err := queries.ListAuthors(ctx)
	if err != nil {
		return err
	}
	log.Println(authors)

	// create an author
	insertedAuthor, err := queries.CreateAuthor(ctx, tutorial.CreateAuthorParams{
		Name: "Brian Kernighan",
		Bio:  sql.NullString{String: "Co-author of The C Programming Language and The Go Programming Language", Valid: true},
	})
	if err != nil {
		return err
	}
	log.Println(insertedAuthor)

	// get the author we just inserted
	fetchedAuthor, err := queries.GetAuthor(ctx, insertedAuthor.ID)
	if err != nil {
		return err
	}

	// prints true
	log.Println(reflect.DeepEqual(insertedAuthor, fetchedAuthor))
	return nil
}
