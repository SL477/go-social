package main

import (
	"fmt"
	"log"

	"github.com/SL477/go-social/internal/database"
)

func main() {
	c := database.NewClient("db.json")
	err := c.EnsureDB()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("database ensured!")
}