package main

import (
	"errors"
	"fmt"
	"net/http"
)

func printPort (port string) (error, string) {
    var error error
    if port == "" {
        error = errors.New("Port is not defined")
        fmt.Println(error)
        return error, port 
    }
    fmt.Printf("This is listening on port %v", port)
    if error != nil {
        fmt.Println(error)
    }
    return error, port
}

func main() {

    http.HandleFunc("/authorize", func (w http.ResponseWriter, r *http.Request) {
        Authorize(w, r)
    })


    port := "80"
    go printPort(port)
    http.ListenAndServe(`:80`, nil)
}