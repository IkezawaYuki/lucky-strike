package main

import "fmt"

type Employee struct {
	Name, Position string
	AnnualIncome   int
}

const (
	Developer = iota
	Manager
)

func NewEmployee(role int) *Employee {
	switch role {
	case Developer:
		return &Employee{
			Name:         "",
			Position:     "Developer",
			AnnualIncome: 60000,
		}
	case Manager:
		return &Employee{
			Name:         "",
			Position:     "Manager",
			AnnualIncome: 80000,
		}
	default:
		panic("unsupported role")
	}
}

func main() {
	m := NewEmployee(Manager)
	m.Name = "Sam"
	fmt.Println(m)
}
