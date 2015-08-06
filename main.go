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

type Database interface {
	Store(k string, v string)
	Read(k string) (v string)
}

type RiakDatabase struct{}

func (r RiakDatabase) Store(k string, v string) {
	obj, _ := riak.NewObjectIn("urls", k)
	obj.ContentType = "text/plain"
	obj.Data = []byte(v)
	obj.Store()
}

func (r RiakDatabase) Read(k string) (v string) {
	obj, _ := riak.GetFrom("urls", k)
	v = string(obj.Data)
	return v
}

func redirectHandler(db Database) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		k := r.URL.Query().Get(":shorturl")
		v := db.Read(k)
		http.Redirect(w, r, "http://"+v, http.StatusFound)
	})
}

func shortenHandler(db Database) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		val := r.URL.Query().Get(":longurl")
		key := generateKey()
		db.Store(key, val)
		// goal: go.url/s23Fs
		io.WriteString(w, "localhost:12345/"+key)
	})
}

func generateKey() (k string) {
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, 5)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func newDatabase() Database {
	if err := riak.ConnectClient("127.0.0.1:8087"); err != nil {
		fmt.Println("Cannot connect, is Riak running?")
	}

	riak.NewBucket("urls")
	return RiakDatabase{}
}

func main() {
	m := pat.New()
	db := newDatabase()
	m.Get("/:shorturl", redirectHandler(db))
	m.Post("/:longurl", shortenHandler(db))
	http.Handle("/", m)
	if err := http.ListenAndServe(":12345", nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
