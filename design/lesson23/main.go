package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
)

type Address struct {
	StreetAddress, City, Country string
}

func (a *Address) DeepCopy() *Address {
	return &Address{
		StreetAddress: a.StreetAddress,
		City:          a.City,
		Country:       a.Country,
	}
}

type Person struct {
	Name    string
	Address *Address
	Friends []string
}

func (p *Person) DeepCopy() *Person {
	b := bytes.Buffer{}
	e := gob.NewEncoder(&b)
	_ = e.Encode(p)

	fmt.Println(string(b.Bytes()))

	d := gob.NewDecoder(&b)
	result := Person{}
	_ = d.Decode(&result)
	return &result
}

func main() {
	john := &Person{
		Name: "John",
		Address: &Address{
			StreetAddress: "123 London Rd",
			City:          "London",
			Country:       "UK",
		},
		Friends: []string{"Chris", "Matt"},
	}

	jane := john.DeepCopy()
	jane.Name = "Jane"
	jane.Address.StreetAddress = "321 Baker St"
	jane.Friends = append(jane.Friends, "Angela")
	fmt.Println(john, john.Address)
	fmt.Println(jane, jane.Address)
}
