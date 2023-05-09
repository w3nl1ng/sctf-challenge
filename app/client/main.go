package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
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
		targetURL := "http://127.0.0.1:8080/nologin/userinfo?username=" + userName
		resp, err := http.Get(targetURL)
		checkNormalErr(err)
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		checkNormalErr(err)
		checkBodyErr(body)

		var data map[string]interface{}
		err = json.Unmarshal(body, &data)
		checkNormalErr(err)

		fmt.Printf("user name: %s\nuser id: %d\nuser phone: %s\n", 
					data["name"], int(data["id"].(float64)), data["phone"])

		if message, ok := data["message"]; ok {
			fmt.Printf("extra message: %s\n", message)
		}
	}

	if addrbook {
		targetURL := fmt.Sprintf(
			"http://127.0.0.1:8080/auth/showaddressbook?userid=%d&password=%s",
			userId, password)

		resp, err := http.Get(targetURL)
		checkNormalErr(err)

		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		checkNormalErr(err)
		checkBodyErr(body)

		fmt.Println(string(body))
	}
}

func checkBodyErr(body []byte) {
		var errMsg map[string]interface{}
		if err1 := json.Unmarshal(body, &errMsg); err1 != nil { //body不是json格式
			return
		} else { //json解析成功
			if msg, ok := errMsg["error"]; ok {
				fmt.Println(msg)
				os.Exit(1)
			} else { //未设置err键
				return
			}
		}
}

func checkNormalErr(err error) {
	if err != nil {
		panic(err)
	}
}

func decodeStr(src []byte) string {
	
}