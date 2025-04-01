package main

import (
	"fmt"
	"log"
	"temp/helpers"
)

func main() {
	password := "test123"
	hashed, err := helpers.HashPassword(password)
	if err != nil {
		log.Fatal("Error hashing password:", err)
	}
	fmt.Println("Hashed password:", hashed)
}
