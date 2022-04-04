package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/nats-io/stan.go"
	"sync"
)

const (
	channel = "channel"
	client  = "clientid"
)

func main() {

	sc, _ := stan.Connect("test-cluster", client)
	sub, _ := sc.Subscribe(channel, func(m *stan.Msg) {
		text := string(m.Data)
		fmt.Printf("Received a message: %s\n", text)
		psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", "localhost", 5432, "postgres", "1234", "nats1")

		// open database
		db, err := sql.Open("postgres", psqlconn)
		if err != nil {
			panic(err)
		}

		defer db.Close()

		insert, err := db.Query("INSERT INTO nats (chan, client, texts) VALUES($1, $2, $3)",
			channel,
			client,
			text)
		if err != nil {
			panic(err)
		}
		defer insert.Close()
	}, stan.DeliverAllAvailable())

	wg := sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()
	// connection string

	sub.Unsubscribe()

	sc.Close()
}
