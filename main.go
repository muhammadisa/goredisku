package main

import (
	"fmt"
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

	command.WT(
		func() {
			fmt.Println("Write through db")
		},
	)
	command.WB(
		func() {
			fmt.Println("Write back db")
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
