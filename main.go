package main

import (
	"fmt"
	"log"

	"github.com/Diyatt/GolangProject/database"
)

func main() {
	err := database.InitDb()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Hello, world!")

}
