package main

import (
	. "github.com/spouk/goparser/library"

	"fmt"
	"os"
	"log"
)

func main() {
	db := NewDatabase("simple.db")
	err := db.OpenDbs()
	if err != nil {
		log.Printf("Error: %v\n", err)
		os.Exit(1)
	} else {
		err := db.WriteRecordTester()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	}

	fmt.Printf("%v\n", db.DB)
	os.Exit(1)

	p, err := NewParser("config.yaml")
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
	p.Log.Print(p)

	p.Run()

}
