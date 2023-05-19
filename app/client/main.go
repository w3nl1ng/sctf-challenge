package main

import (
	"crypto/rc4"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"syscall"
)


var (
	userId int
	password string
	addrbook bool
	info bool
	userName string
)

const (
    PT_TRACE_ME = 0
)

var (
	infoUrl = []byte{0xed, 0x2f, 0xd8, 0x27, 0xa8, 0x08, 0x4e, 0xbd, 0xc9, 0x77, 0x8f, 0x03, 0xbb, 0x77, 0x4d, 0xac, 0x27, 0xf4, 0xa2, 0x68, 0x9b, 0x2f, 0xbb, 0x16, 0x51, 0x8e, 0x65, 0x11, 0x3d, 0x92, 0x80, 0x74, 0x2d, 0xb4, 0xc0, 0xda, 0x63, 0x52, 0xfb, 0x44, 0xd7, 0x18, 0xa7, 0x37, 0x5e, 0x02, 0x0e, 0x90, 0x35, 0x4c, 0x05, 0x36,}
	//http://81.68.208.118:8080/auth/showaddressbook?userid=%d&password=%s
	addrbookUrl = []byte{0xed, 0x2f, 0xd8, 0x27, 0xa8, 0x08, 0x4e, 0xbd, 0xc9, 0x77, 0x8f, 0x03, 0xbb, 0x77, 0x4d, 0xac, 0x27, 0xf4, 0xa2, 0x68, 0x9b, 0x2f, 0xbb, 0x16, 0x50, 0x8e, 0x6a, 0x0b, 0x25, 0x95, 0xc8, 0x6e, 0x2b, 0xf4, 0xc2, 0xc8, 0x62, 0x44, 0xe0, 0x4f, 0xc2, 0x04, 0xfa, 0x2d, 0x42, 0x0c, 0x43, 0x8b, 0x27, 0x44, 0x12, 0x62, 0x27, 0x85, 0x8d, 0xdc, 0x87, 0x0f, 0x31, 0xd2, 0xc2, 0xf0, 0x71, 0xb2, 0x76, 0xce, 0x87, 0x20,}
)

func init() {
	flag.StringVar(&password, "p", "", "set password")
	flag.IntVar(&userId, "id", 0, "set the user id")
	flag.StringVar(&userName, "username", "", "set the username")
	flag.BoolVar(&addrbook, "addrbook", false, "show address book")
	flag.BoolVar(&info, "info", false, "show user info")

	go checkDebugger()
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
		if err != nil {
			panic(err)
		}
		
		if resp.StatusCode == http.StatusOK { //请求成功
			handleAddrBookOK(body)
		} else { //请求失败
			fmt.Println(string(body))
		}
	}
}

type userInfoResp struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Phone string `json:"phone"`
	Msg   string `json:"message"`
}

type contact struct {
	Phones []string `json:"phones"`
}

func handleAddrBookOK(body []byte) {
	var contactResp contact
	jsonByte := decryptoJSONByte(body)
	if err := json.Unmarshal(jsonByte, &contactResp); err != nil {
		panic(err)
	}
	fmt.Printf("you can contact people by phone: %s\n", strings.Join(contactResp.Phones, ", "))
}


func decryptoJSONByte(jsonByte []byte) []byte {
	key := []byte("sctf2023")
	rc4Enc, err := rc4.NewCipher(key)
	if err != nil {
		panic(err)
	}

	plain := make([]byte, len(jsonByte))
	rc4Enc.XORKeyStream(plain, jsonByte)

	return plain
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

func checkDebugger() {
    _, _, err := syscall.Syscall(syscall.SYS_PTRACE, PT_TRACE_ME, 0, 0)
    if err != 0 {
        os.Exit(0)
    }
}