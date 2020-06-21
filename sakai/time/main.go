package main

import (
	"fmt"
	"regexp"
)

func main() {
	match, _ := regexp.MatchString("a([a-z]+)e", "apple")
	fmt.Println(match)

	r := regexp.MustCompile("a([a-z]+)e")
	ms := r.MatchString("apple")
	fmt.Println(ms)

	//s := "/view/test"
	r2 := regexp.MustCompile("^/(edit|save|view)/([a-zA-z0-9]+)$")

	fs := r2.FindString("/view/test")
	fmt.Println(fs)

	fss := r2.FindStringSubmatch("/view/test")
	fmt.Println(fss)
}
