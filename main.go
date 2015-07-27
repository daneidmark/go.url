package main

import (
	"fmt"
	"github.com/bmizerany/pat"
	"github.com/tpjg/goriakpbc"
	"io"
	"log"
	"math/rand"
	"net/http"
	"time"
)

var letters = []rune("1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func redirect(w http.ResponseWriter, req *http.Request) {
	k := req.URL.Query().Get(":shorturl")
	v := read(k)
	http.Redirect(w, req, "http://"+v, http.StatusFound)
}

func shorten(w http.ResponseWriter, req *http.Request) {
	val := req.URL.Query().Get(":longurl")
	key := generateKey()
	store(key, val)
	// goal: go.url/s23Fs
	io.WriteString(w, "localhost:12345/"+key)
}

func generateKey() (k string) {
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, 5)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func read(k string) (v string) {
	obj, _ := riak.GetFrom("urls", k)
	v = string(obj.Data)
	return v
}

func store(k string, v string) {
	obj, _ := riak.NewObjectIn("urls", k)
	obj.ContentType = "text/plain"
	obj.Data = []byte(v)
	obj.Store()
}

func setupDatabase() {
	if err := riak.ConnectClient("127.0.0.1:8087"); err != nil {
		fmt.Println("Cannot connect, is Riak running?")
		return
	}
	riak.NewBucket("urls")
}

func setupRest() {
	m := pat.New()
	m.Get("/:shorturl", http.HandlerFunc(redirect))
	m.Post("/:longurl", http.HandlerFunc(shorten))
	http.Handle("/", m)
	if err := http.ListenAndServe(":12345", nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func main() {
	setupDatabase()
	setupRest()
}
