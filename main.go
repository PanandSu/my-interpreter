package main

import (
	"fmt"
	"log"
	"my-interpreter/repl"
	"os"
	"os/user"
)

func main() {
	usr, err := user.Current()
	//MATEBOOK14S\35895,pansu
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Hello %s!\n", usr.Username)
	fmt.Println("I'm in Juejin")
	repl.Start(os.Stdin, os.Stdout)
}
