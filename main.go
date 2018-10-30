// gobackgammond, a daemon that plays backgammon
package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/chandler37/gobackgammond/handlers"
)

var port = flag.Int64(
	"port",
	-1,
	"Port to listen on")

var seed = flag.Int64(
	"seed",
	time.Now().UnixNano(),
	"Port to listen on")

func main() {
	flag.Parse()
	fmt.Printf("rand.Seed(%v)\n", *seed)
	rand.Seed(*seed)
	if *port < 0 {
		panic("bad arg -port")
	}
	http.HandleFunc("/", handlers.RootHandler)
	http.HandleFunc("/game", handlers.GameHandler)
	http.HandleFunc("/game.svg", handlers.SvgHandler)
	p := fmt.Sprintf(":%d", *port)
	fmt.Printf("Listening on\nhttp://localhost%v/", p)
	log.Fatal(http.ListenAndServe(p, nil))
}
