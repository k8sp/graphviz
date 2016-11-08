// graphviz is an HTTP server which calls GraphViz's dot command to
// visualize a .dot file given in the `dot` parameter.  For example,
// we can run this program as follows
//
//   go run graphviz.go
//
// and access it from the Web browser with URL:
//
//   http://localhost:8080/?dot=https://gist.githubusercontent.com/wangkuiyi/c4e0015211dd1b9bde2e20455a6cd38e/raw/4d5ec099f98a5f326cf6f108bcf510cadba1a0b4/ci-arch.dot
//
// then the visualization of ci-arch.dot should appear in the browser window.
//
package main

import (
	"crypto/md5"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"

	"github.com/topicai/candy"
)

func main() {
	dir := flag.String("dir", "/tmp", "The cache directory")
	addr := flag.String("addr", ":8080", "The listening address")
	flag.Parse()

	http.HandleFunc("/",
		makeSafeHandler(func(w http.ResponseWriter, r *http.Request) {
			if source := r.FormValue("dot"); len(source) > 0 {
				dot, e := candy.HTTPGet(source, 0)
				candy.Must(e)

				id := fmt.Sprintf("%015x", md5.Sum(dot))
				dotFile := path.Join(*dir, id) + ".dot"
				pngFile := path.Join(*dir, id) + ".png"

				if _, e := os.Stat(pngFile); os.IsNotExist(e) {
					// TODO(yi): Here we adopt an unlimted-size disk-based cache.  We should
					// either introduce LRU and limit the size, or periodically clean it.
					log.Printf("Update cache for %s", pngFile)
					candy.Must(ioutil.WriteFile(dotFile, dot, 0755))
					png, e := exec.Command("dot", "-Tpng", dotFile).Output()
					candy.Must(e)
					candy.Must(ioutil.WriteFile(pngFile, png, 0755))
				}

				png, e := ioutil.ReadFile(pngFile)
				candy.Must(e)
				_, e = w.Write(png)
				candy.Must(e)
			}
		}))

	candy.Must(http.ListenAndServe(*addr, nil))
}

func makeSafeHandler(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if e := recover(); e != nil {
				http.Error(w, fmt.Sprint(e), http.StatusInternalServerError)
			}
		}()
		h(w, r)
	}
}
