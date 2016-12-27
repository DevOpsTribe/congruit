package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"./congruit-go/libs"
)

func main() {

	fmt.Println("******************************************")
	fmt.Println("*              congruit!                 *")
	fmt.Println("******************************************")

	files, err := ioutil.ReadDir(".")
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		fmt.Println(file.Name())
	}
	mywork := new(congruit.Work)
	mywork.Name = "create file"
	mywork.DoAfter = "echo now"
	mywork.Command = "touch /tmp/foobar"
	mywork.Idempotency = "[ -e /tmp/foobar ]"
	mywork.DoWork()

}
