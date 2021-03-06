package main

import (
	"encoding/json"
	"fmt"
	"html"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

type vcapService struct {
	Name        string                 `json:"name"`
	Credentials map[string]interface{} `json:"credentials"`
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Path: %s", html.EscapeString(r.URL.Path))
	})

	http.HandleFunc("/env", func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprintln(w, strings.Join(os.Environ(), "\n"))
	})

	http.HandleFunc("/run", func(w http.ResponseWriter, r *http.Request) {
		for name, arg := range r.URL.Query() {
			out, err := exec.Command(name, arg...).CombinedOutput()
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, "%s\n", err)
				return
			}
			fmt.Fprintf(w, "%s", out)
		}
	})

	http.HandleFunc("/services", func(w http.ResponseWriter, r *http.Request) {
		vcapServices := map[string][]vcapService{}
		if err := json.Unmarshal([]byte(os.Getenv("VCAP_SERVICES")), &vcapServices); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "%s\n", err)
			return
		}
		for _, services := range vcapServices {
			for _, service := range services {
				uri := service.Credentials["uri"].(string)
				fmt.Fprintf(w, "Name: %s\nURI: %s\n", service.Name, uri)
				req, err := http.NewRequest("GET", uri, nil)
				if err != nil {
					fmt.Fprintf(w, "Error: %s\n\n", err)
					continue
				}
				req.Host = service.Credentials["host_header"].(string)
				resp, err := http.DefaultClient.Do(req)
				if err != nil {
					fmt.Fprintf(w, "Error: %s\n\n", err)
					continue
				}
				body, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					fmt.Fprintf(w, "Error: %s\n\n", err)
					continue
				}
				fmt.Fprintf(w, "Response: %s\n\n", body)
			}
		}
	})

	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), nil))
}
