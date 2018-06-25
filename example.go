package main

import (
	. "github.com/spouk/goparser/library"
	"fmt"
	"os"
)

func main() {
	p, err := NewParser("config.yaml")
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
	p.Log.Print(p)



}
