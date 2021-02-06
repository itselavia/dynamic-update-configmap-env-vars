package main

import (
	"fmt"
	"net/http"
	"os"
)

func getEnvValue(w http.ResponseWriter, req *http.Request) {
	envVar, ok := req.URL.Query()["var"]
	if !ok {
		errMessage := "variable " + envVar[0] + " not found"
		fmt.Fprintf(w, errMessage)
	}

	value := os.Getenv(envVar[0])

	fmt.Fprintf(w, value+"\n")
}

func main() {

	http.HandleFunc("/getEnvValue", getEnvValue)

	http.ListenAndServe(":8080", nil)
}
