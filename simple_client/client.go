package main

import (
	"fmt"
	"io"
	"net/http"
)

func main() {
	//resp, err := http.Get("https://helloworld-icq63pqnqq-uc.a.run.app")
	resp, err := http.Get("http://localhost:7777")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	fmt.Println(string(body))

}
