package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
)

func main() {
	//content, err := ioutil.ReadFile("sakai/ioutil/main.go")
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//fmt.Println(string(content))
	//
	//if err := ioutil.WriteFile("sakai/ioutil/main_copy.go", content, 0666); err != nil {
	//	log.Fatalln(err)
	//}

	r := bytes.NewBuffer([]byte("abc"))
	content, _ := ioutil.ReadAll(r)
	fmt.Println(string(content))
}
