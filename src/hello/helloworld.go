package main

import (
	"fmt"

	"gopkg.in/mgo.v2"
)

func main() {
	fmt.Println("Hello world! My lucky number is", add(10, 20))
}

type UserService struct {
	collection *mgo.Collection
}
