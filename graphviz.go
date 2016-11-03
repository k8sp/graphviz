package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os/exec"
)

func main() {
	http.HandleFunc("/", makeSafeHandler(func(w http.ResponseWriter, r *http.Request) {
		source := "https:/" + r.URL.Path
		resp, e := http.Get(source)
		if e != nil {
			log.Panicf("Cannot retreive source %s: %v", source, e)
		}
		defer resp.Body.Close()

		cmd := exec.Command("dot", "-Tpng")
		w.Header().Set("Content-Type", "image/png")

		if stdin, e := cmd.StdinPipe(); e != nil {
			log.Panicf("Cannot pipe to stdin: %v", e)
		} else {
			go func() {
				if _, e := io.Copy(stdin, resp.Body); e != nil {
					log.Panicf("Failed copying source to graphviz: %v", e)
				}
			}()
		}

		if stdout, e := cmd.StdoutPipe(); e != nil {
			log.Panicf("Cannot read stdout: %v", e)
		} else {
			go func() {
				if _, e := io.Copy(w, stdout); e != nil {
					log.Panicf("Failed delivering output: %v", e)
				}
			}()
		}

		if e := cmd.Start(); e != nil {
			log.Panicf("%v", e)
		}
	}))

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func makeSafeHandler(in http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if e := recover(); e != nil {
				http.Error(w, fmt.Sprint(e), http.StatusInternalServerError)
			}
		}()
		in(w, r)
	}
}
