package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/your-username/utils"
)

var (
	version   = "1.0.0"
	buildTime = "2023-02-20T14:30:00Z"
)

func main() {
	flag.Parse()

	fmt.Printf("Starting application %s (version %s) built at %s\n", filepath.Base(os.Args[0]), version, buildTime)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Welcome to the main application!")
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
