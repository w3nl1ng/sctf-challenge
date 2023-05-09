package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"crypto/rc4"
)


var (
	userId int
	password string
	addrbook bool
	info bool
	userName string
)

var (
	infoUrl = []byte{0xed, 0x2f, 0xd8, 0x27, 0xa8, 0x08, 0x4e, 0xb4, 0xca, 0x6e, 0x97, 0x0b, 0xbb, 0x75, 0x53, 0xa5, 0x33, 0xfd, 0xa3, 0x68, 0x91, 0x38, 0xe5, 0x41, 0x0c, 0xce, 0x6c, 0x17, 0x3f, 0xd2, 0x92, 0x6e, 0x26, 0xe9, 0xdc, 0xc7, 0x60, 0x4f, 0xad, 0x5f, 0xc2, 0x12, 0xea, 0x2c, 0x4c, 0x0a, 0x19, 0xc3,}
	addrbookUrl = []byte{0xed, 0x2f, 0xd8, 0x27, 0xa8, 0x08, 0x4e, 0xb4, 0xca, 0x6e, 0x97, 0x0b, 0xbb, 0x75, 0x53, 0xa5, 0x33, 0xfd, 0xa3, 0x68, 0x91, 0x38, 0xea, 0x5b, 0x14, 0xc9, 0x24, 0x0d, 0x39, 0x92, 0x90, 0x7c, 0x27, 0xff, 0xc7, 0xcc, 0x75, 0x53, 0xf0, 0x45, 0xde, 0x1c, 0xa7, 0x37, 0x5e, 0x02, 0x0e, 0x97, 0x30, 0x1c, 0x45, 0x6f, 0x65, 0xc8, 0xc9, 0xcb, 0xd2, 0x08, 0x3f, 0xd3, 0xd5, 0xba, 0x3b, 0xb3,}
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
		targetURL := decodeStr(infoUrl) + userName
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

		if len(password) < 5 { //password长度检测
			os.Exit(1)
		}

		targetURL := fmt.Sprintf(
			decodeStr(addrbookUrl),
			userId, password[:5])

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

	key := []byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff}
	rc4Dec, err := rc4.NewCipher(key)
	checkNormalErr(err)

	plain := make([]byte, len(src))
	rc4Dec.XORKeyStream(plain, src)

	return string(plain)
}