package main

import (
	"fmt"

	"um6p.ma/finalproject/database"
)

func main() {
	fmt.Println("Starting ...")
	database.ConnectDatabase()
}
