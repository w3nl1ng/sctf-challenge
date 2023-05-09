package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
)


var (
	userId int
	password string
	addrbook bool
	info bool
	userName string
)

func init() {
	flag.StringVar(&password, "p", "", "set password")
	flag.IntVar(&userId, "id", 0, "set the user id")
	flag.StringVar(&userName, "username", "", "set the username")
	flag.BoolVar(&addrbook, "addrbook", false, "show address book")
	flag.BoolVar(&info, "info", false, "show user info")
}

func main() {
	flag.Parse()

	if info { //查询用户信息，需要提供username
		targetURL := "http://127.0.0.1:8080/userinfo?username=" + userName
		resp, err := http.Get(targetURL)
		if err != nil {
			fmt.Println(err)
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println(string(body))
	}
}