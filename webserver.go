package main

import (
	"flag"
	"fmt"
	"net/http"
)

var (
	webport string
)
func init() {
	flag.StringVar(&webport, "webport", "127.0.0.1:80", "webserver listening port <:port>")
}

func main() {
	flag.Parse()
	fmt.Printf("webserver running on %v\n", webport)
	//	http.Handle("/", http.FileServer(http.Dir("./vscode")))
	http.Handle("/", http.FileServer(http.Dir(".")))
	err:= http.ListenAndServe(webport, nil)
	if err!=nil { panic(err) }
}