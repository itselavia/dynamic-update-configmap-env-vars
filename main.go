package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/fsnotify/fsnotify"
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

func reloadEnvVars() {

	configDir := "/config/"
	files, err := ioutil.ReadDir(configDir)
	if err != nil {
		fmt.Println("cannot read config dir ", err)
		return
	}
	for _, file := range files {
		key := file.Name()
		if !strings.HasPrefix(key, "..") {
			filename := configDir + key
			value, err := ioutil.ReadFile(filename)
			if err != nil || string(value) == "" {
				fmt.Println("Unable to read env variable: ", filename)
				continue
			}
			fmt.Println("Reloading key: ", key, " with value: ", string(value))
			os.Setenv(key, string(value))
		}

	}
}

func main() {

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Println("cannot initialize Watcher ", err)
	}
	defer watcher.Close()

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				fmt.Println("event:", event)
				reloadEnvVars()
			case err, ok := <-watcher.Errors:
				if !ok {
					fmt.Println("error from Watcher: ", err)
					return
				}
			}
		}
	}()

	err = watcher.Add("/config/")
	if err != nil {
		fmt.Println("error adding directory to Watcher", err)
	}

	http.HandleFunc("/getEnvValue", getEnvValue)

	http.ListenAndServe(":8080", nil)
}
