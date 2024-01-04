package main

import "fmt"

type User struct {
	name string
	age  int
}

func (user User) message() string {
	return fmt.Sprintf("Hi %s", user.name)
}

func main() {

	mark := User{
		name: "Mark",
		age:  19,
	}

	fmt.Println(mark.message())
}
