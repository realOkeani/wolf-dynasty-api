package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/gorilla/mux"

	route "github.com/realOkeani/wolf-dynasty-api/http"
)

func main() {
	/*
		Begin by formatting a string that will include all of the environment variables,
		needed to connect to a GCP mySQL database. To do this you use the sprintf function.
		func Sprintf(format string, a ...interface{}) string

		Sprintf formats according to a format specifier and returns the resulting string.

		This means that you use the % symbol and the letter s, to make %s in the places where you would like to inject variables
		into a string. Then for the each %s you add the variables. An example of this is

			const name, age = "Kim", 22
			s := fmt.Sprintf("%s is %d years old.\n", name, age)

		In the case of this project we are injecting local environment variables.

	*/
	// dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true",
	// 	os.Getenv("DB_USER"),
	// 	os.Getenv("DB_PASSWORD"),
	// 	os.Getenv("DB_HOST"),
	// 	os.Getenv("DB_NAME"))

	// db, err := sqlx.Open("mysql", dsn)
	// if err != nil {
	// 	panic(err)
	// }

	// sqlClient := leagueSQL.NewOwnersClient(db)

	// s := wolf.Services{
	// 	SQLClient: sqlClient,
	// }

	/*
		Package mux implments a request router and dispatcher. The name mux stands
		for "HTTP request multiplexer". Like the standard http.ServeMux, mux.Router
		matches incoming requests against a list of registered routes and calls a
		handler for the route that matches the URL or other conditions.
		In telecommunications and computer networks, multiplexing (sometimes contracted to muxing)
		is a method by which multiple analog or digital signals are combined into one signal over a
		shared medium. The aim is to share a scarce resource. A device that performs the multiplexing
		is called a multiplexer (MUX).
		The main features are:
		- Requests can be matched based on URL host, path, path prefix, schemes, header
			and query values, HTTP methods or using custom matchers.
		- URL hosts, paths and query values can have variables with an optional regular expression.
		- Registered URLs can be built, or "reveresed". which helps maintaining references to resources.
		- Routes can be used as subrouters: nested routes are only tested if the parent route matches.
			This is useful to define groups of routes that share common conditions like a host, a path
			prefix or other repeated attributes. As a bonus, this optimizes request matching.
		- It implements the http.Handler interface so it is compatible with the standard http.ServeMux.
	*/

	r := mux.NewRouter()
	route.AddRoutes(r)
	r.Use(route.CorsHandler)

	srv := &http.Server{
		Addr: "0.0.0.0:8080",
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r, // Pass our instance of gorilla/mux in.
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Println("shutting down")
	os.Exit(0)
}
