package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/muhammadisa/goredisku/goredisku"
)

// Person struct
type Person struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Age       int    `json:"age"`
}

func main() {

	// Create client from Connecting
	client := goredisku.RedisCred{
		Address:  "localhost:6379",
		Password: "",
		DB:       0,
		Debug:    true,
	}.Connect()

	// Create command from Client
	command := goredisku.GoRedisKu{
		Client:       client,
		GlobalExpire: time.Duration(60000) * time.Millisecond,
	}

	person1 := Person{
		FirstName: "Muhammad Isa",
		LastName:  "Wijaya Kusuma 1",
		Age:       21,
	}

	person2 := Person{
		FirstName: "Olav",
		LastName:  "Alan Walker 2",
		Age:       28,
	}

	command.WT(
		"person_1",
		person1,
		func() {
			fmt.Println("Write through db")
		},
	)
	command.WB(
		"person_2",
		person2,
		func(wg *sync.WaitGroup, mtx *sync.Mutex) {
			mtx.Lock()
			fmt.Println("Write back db")
			time.Sleep(3 * time.Second)
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
