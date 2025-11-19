package main

import (
	"encoding/json"
	"fmt"
	"log"
)

type Address struct {
	City   string
	Street string
}

type Person struct {
	Name string
	Address
}

type Array []int

func (a *Array) Append(numbers ...int) {
	*a = Array(append([]int(*a), numbers...))
}

type Timer interface {
	Time(time int) int
}

func Factorial(x int) int {
	if x <= 1 {
		return 1
	}
	return x * Factorial(x-1)
}

func MarshalingTest() {
	fmt.Println("hi")
	names := make(map[string]int)
	names["Dima"]++
	fmt.Printf("%#v\n", names)

	client := Person{Name: "Mama", Address: Address{City: "Likino-Dulevo", Street: "Stepana Morozkina"}}
	var newClient Person

	jsonClient, err := json.MarshalIndent(client, "", "\t")

	if err != nil {
		log.Fatalf("%v\n", err)
	}
	fmt.Printf("%s\n", jsonClient)

	err = json.Unmarshal(jsonClient, &newClient)
	if err != nil {
		log.Fatalf("Error Unmarshaling : %v", err)
	}

	jsonNames, err := json.Marshal(names)
	if err != nil {
		log.Fatalf("Error Marshaling : %v", err)
	}

	fmt.Printf("%+s\n%+s\n", jsonNames, newClient)
}
