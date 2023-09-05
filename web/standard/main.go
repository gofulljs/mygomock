package main

import "net/http"

func HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, world!"))
}

func InitMux() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", HelloWorldHandler)
	return mux
}

func main() {
	panic(http.ListenAndServe(":8080", InitMux()))
	// server := &http.Server{Addr: ":8080", Handler: Hellowo}
}
