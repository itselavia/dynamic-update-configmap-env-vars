package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/fsnotify/fsnotify"
)

// Handler function to return the environment value of the key given in URL paramater 'var'
func getEnvValue(w http.ResponseWriter, req *http.Request) {
	envVar, ok := req.URL.Query()["var"]
	if !ok {
		errMessage := "variable " + envVar[0] + " not found"
		fmt.Fprintf(w, errMessage)
	}

	value := os.Getenv(envVar[0])
	fmt.Fprintf(w, value+"\n")
}

// reloads the enviroment variables from the mounted /config path
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
				continue
			} else {
				os.Setenv(key, string(value))
			}
		}
	}
}

func main() {

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Println("cannot initialize Watcher ", err)
	}
	defer watcher.Close()

	// watcher will monitor the files in a background goroutine
	go func() {
		for {
			select {
			// reload the environment variables whenever changes are made in the /config directory
			case _, ok := <-watcher.Events:
				if !ok {
					return
				}
				reloadEnvVars()
			case err, ok := <-watcher.Errors:
				if !ok {
					fmt.Println("error from Watcher: ", err)
					return
				}
			}
		}
	}()

	// monitor the /config directory
	err = watcher.Add("/config/")
	if err != nil {
		fmt.Println("error adding directory to Watcher", err)
	}

	http.HandleFunc("/getEnvValue", getEnvValue)

	http.ListenAndServe(":8080", nil)
}
