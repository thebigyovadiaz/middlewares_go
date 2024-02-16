package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"time"
)

func getEnvConfig(key string) string {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}

// Basic Middlewares
func logging(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.URL.Path)
		f(w, r)
	}
}

func foo(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprintln(w, "foo")

	if err != nil {
		fmt.Printf("Server Error with FOO, error: %v\n", err.Error())
		return
	}

	fmt.Println("FOO is OK!")
}

func bar(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprintln(w, "bar")

	if err != nil {
		fmt.Printf("Server Error with BAR, error: %v\n", err.Error())
		return
	}

	fmt.Println("BAR is OK!")
}

// Advanced Middleware
type Middleware func(http.HandlerFunc) http.HandlerFunc

func Logging() Middleware {
	// New Middleware
	return func(handlerFunc http.HandlerFunc) http.HandlerFunc {

		// Define HandlerFunc
		return func(w http.ResponseWriter, r *http.Request) {

			// Do middleware things
			start := time.Now()
			fmt.Printf("Logging middleware - time: %v\n", start)

			defer func() {
				log.Println(r.URL.Path, time.Since(start))
			}()

			// Call the next middleware/handler in chain
			handlerFunc(w, r)
		}
	}
}

func Method(m string) Middleware {
	// New Middleware
	return func(handlerFunc http.HandlerFunc) http.HandlerFunc {

		// Define HandlerFunc
		return func(w http.ResponseWriter, r *http.Request) {

			// Do middleware things
			if r.Method != m {
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			} else {
				fmt.Printf("Method middleware - method param: %s\n", m)

				// Call the next middleware/handler in chain
				handlerFunc(w, r)
			}
		}
	}
}

func Chain(f http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
	for _, m := range middlewares {
		f = m(f)
	}

	return f
}

func HelloWorld(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprint(w, "Hello World with Middleware!")
}

func main() {
	port := getEnvConfig("PORT")
	fmt.Println(port)
	http.HandleFunc("/foo", logging(foo))
	http.HandleFunc("/bar", logging(bar))

	http.HandleFunc("/hello", Chain(HelloWorld, Method("GET"), Logging()))

	err := http.ListenAndServe(":"+port, nil)
	if err == nil {
		fmt.Printf("Server starting on port %s\n", port)
	} else {
		log.Fatalf("Server Error started, error: %v\n", err.Error())
	}
}
