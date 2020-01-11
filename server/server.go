// package main

// import (
// 	"flag"
// 	"log"
// 	"net/http"
// )

package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os/exec"

	"github.com/jaschaephraim/lrserver"
	"gopkg.in/fsnotify.v1"
)

func server() {
	port := flag.String("p", "2021", "port to serve on")
	directory := flag.String("d", ".", "the directory of static file to host")
	flag.Parse()

	http.Handle("/", http.FileServer(http.Dir(*directory)))

	log.Printf("Serving %s on HTTP port: %s\n", *directory, *port)
	log.Fatal(http.ListenAndServe(":"+*port, nil))
}

const html = `<!doctype html>
<html>
<head>
<script src="http://localhost:35729/livereload.js"></script>
</head>
</html>`

func main() {
	// Create file watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatalln(err)
	}
	defer watcher.Close()

	// Add dir to watcher
	err = watcher.Add(".")
	if err != nil {
		log.Fatalln(err)
	}

	// Create and start LiveReload server
	lr := lrserver.New(lrserver.DefaultName, lrserver.DefaultPort)
	go lr.ListenAndServe()

	// Start goroutine that requests reload upon watcher event
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				fmt.Println(event)
				cmd := exec.Command("make")
				cmd.Run()
				lr.Reload("index.html")
			case err := <-watcher.Errors:
				log.Println(err)
			}
		}
	}()
	// server()
	http.HandleFunc("/", func(rw http.ResponseWriter, req *http.Request) {
		rw.Write([]byte(html))
	})
	http.ListenAndServe(":3000", nil)

}
