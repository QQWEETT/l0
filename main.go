package main

import (
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/patrickmn/go-cache"
	"html/template"
	"net/http"
	"strconv"
	"time"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "username"
	password = "password"
	dbname   = "dbname"
)

type Example struct {
	Id      int
	Channel string
	Client  string

	Text string
}

var showPost = Example{}

func db(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	t, err := template.ParseFiles("templates/show.html")

	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	res, err := db.Query("select id, chan, client, texts from nats")
	if err != nil {
		panic(err)
	}
	showPost = Example{}
	for res.Next() {
		var example Example
		err = res.Scan(&example.Id, &example.Channel, &example.Client, &example.Text)
		if err != nil {
			panic(err)
		}

		c := cache.New(5*time.Minute, 10*time.Minute)
		c.Set(strconv.Itoa(example.Id), example, cache.DefaultExpiration)
		foo, found := c.Get(vars["id"])
		if found {
			showPost = example
			fmt.Println(foo)
			t.ExecuteTemplate(w, "show", foo)

		}
	}

}

func handleFunc() {
	rtr := mux.NewRouter()
	rtr.HandleFunc("/post/{id:[0-9]+}", db).Methods("GET")
	http.Handle("/", rtr)
	http.ListenAndServe(":8080", nil)
}

func main() {
	handleFunc()
}
