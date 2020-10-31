package main

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/gocraft/dbr/dialect"
	"github.com/gocraft/dbr/v2"
	"github.com/joho/godotenv"
	"github.com/muhammadisa/godbconn"
	"github.com/muhammadisa/goredisku/goredisku"
)

// Person struct
type Person struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Ages      int    `json:"age"`
}

func main() {

	// Create client from Connecting
	client := goredisku.RedisCred{
		Address:  "localhost:6379",
		Password: "",
		DB:       0,
		Debug:    true,
	}.Connect()

	// connect to db
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	// Load database credential env and use it
	db, err := godbconn.DBCred{
		DBDriver:   "mysql",
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),
	}.Connect()
	if err != nil {
		log.Fatal(err)
	}
	conn := &dbr.Connection{
		DB:            db,
		EventReceiver: &dbr.NullEventReceiver{},
		Dialect:       dialect.MySQL,
	}
	conn.SetMaxOpenConns(10)
	session := conn.NewSession(nil)
	session.Begin()

	// Create command from Client
	command := goredisku.NewGrkClient(
		client,
		time.Duration(60000)*time.Millisecond,
	)

	person1 := Person{
		FirstName: "Muhammad Isa",
		LastName:  "Wijaya Kusuma 1",
		Ages:      21,
	}

	person2 := Person{
		FirstName: "Olav",
		LastName:  "Alan Walker 2",
		Ages:      28,
	}

	command.WT(
		"person_1",
		person1,
		func() {
			fmt.Println("Write through db")
			_, err := session.InsertInto("persons").
				Columns("id", "first_name", "last_name", "ages").
				Record(&person1).
				Exec()
			if err != nil {
				log.Fatal(err)
			}
		},
	)
	command.WB(
		"person_2",
		person2,
		func(wg *sync.WaitGroup, mtx *sync.Mutex) {
			mtx.Lock()
			fmt.Println("Write back db")
			_, err := session.InsertInto("persons").
				Columns("id", "first_name", "last_name", "ages").
				Record(&person2).
				Exec()
			if err != nil {
				log.Fatal(err)
			}
			mtx.Unlock()
			wg.Done()
		},
	)

	// // Insert or Set
	// person1 := Person{
	// 	FirstName: "Muhammad Isa",
	// 	LastName:  "Wijaya Kusuma",
	// 	Age:       21,
	// }
	// inserted, err := command.Set("person_1", &person1)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(inserted)
	// // End of Insert or Set

	// // Retrieve or Get
	// var cachedPerson1 Person
	// err = command.Get("person_1", &cachedPerson1)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(cachedPerson1)
	// // End of Retrieve or Get

	// // Remove or Delete
	// affected, err := command.Del("person_1")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// if affected != 0 {
	// 	fmt.Println(affected)
	// }
	// // End of Remove or Delete
}
