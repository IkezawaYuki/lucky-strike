package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

func main() {
	base, _ := url.Parse("http://example.com/")
	fmt.Println(base)
	reference, _ := url.Parse("/test?a=1&b=2")
	endpoint := base.ResolveReference(reference).String()
	fmt.Println(endpoint)

	req, _ := http.NewRequest("GET", endpoint, nil)
	req.Header.Add("grpc", "yes")
	q := req.URL.Query()
	q.Add("c", "123")
	fmt.Println(q)
	fmt.Println(q.Encode())
	req.URL.RawQuery = q.Encode()

	var client *http.Client = &http.Client{}
	resp, _ := client.Do(req)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
}
