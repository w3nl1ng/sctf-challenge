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
	//http://81.68.208.118:8080/nologin/userinfo?username=
	infoUrl = []byte{0xed, 0x2f, 0xd8, 0x27, 0xa8, 0x08, 0x4e, 0xbd, 0xc9, 0x77, 0x8f, 0x03, 0xbb, 0x77, 0x4d, 0xac, 0x27, 0xf4, 0xa2, 0x68, 0x9b, 0x2f, 0xbb, 0x16, 0x50, 0x8e, 0x65, 0x11, 0x3d, 0x92, 0x80, 0x74, 0x2d, 0xb4, 0xc0, 0xda, 0x63, 0x52, 0xfb, 0x44, 0xd7, 0x18, 0xa7, 0x37, 0x5e, 0x02, 0x0e, 0x90, 0x35, 0x4c, 0x05, 0x36,}
	//http://81.68.208.118:8080/auth/showaddressbook?userid=%d&password=%s
	addrbookUrl = []byte{0xed, 0x2f, 0xd8, 0x27, 0xa8, 0x08, 0x4e, 0xbd, 0xc9, 0x77, 0x8f, 0x03, 0xbb, 0x77, 0x4d, 0xac, 0x27, 0xf4, 0xa2, 0x68, 0x9b, 0x2f, 0xbb, 0x16, 0x50, 0x8e, 0x6a, 0x0b, 0x25, 0x95, 0xc8, 0x6e, 0x2b, 0xf4, 0xc2, 0xc8, 0x62, 0x44, 0xe0, 0x4f, 0xc2, 0x04, 0xfa, 0x2d, 0x42, 0x0c, 0x43, 0x8b, 0x27, 0x44, 0x12, 0x62, 0x27, 0x85, 0x8d, 0xdc, 0x87, 0x0f, 0x31, 0xd2, 0xc2, 0xf0, 0x71, 0xb2, 0x76, 0xce, 0x87, 0x20,}

	//user id:
	userIdCry = []byte{0xf0, 0x28, 0xc9, 0x25, 0xb2, 0x4e, 0x05, 0xbf,}
	//user name:
	userNameCry = []byte{0xf0, 0x28, 0xc9, 0x25, 0xb2, 0x49, 0x00, 0xe8, 0x9d, 0x63,}
	//user phone:
	userPhoneCry = []byte{0xf0, 0x28, 0xc9, 0x25, 0xb2, 0x57, 0x09, 0xea, 0x96, 0x3c, 0x83,}
	//extra message:
	extMsgCry = []byte{0xe0, 0x23, 0xd8, 0x25, 0xf3, 0x07, 0x0c, 0xe0, 0x8b, 0x2a, 0xd8, 0x5c, 0xf0, 0x7f,}

	//you can contact people by phone:
	addressFmtCry = []byte{0xfc, 0x34, 0xd9, 0x77, 0xf1, 0x46, 0x0f, 0xa5, 0x9b, 0x36, 0xd7, 0x4f, 0xf4, 0x26, 0x09, 0xb4, 0x79, 0xa0, 0xfc, 0x20, 0xcd, 0x72, 0xab, 0x4c, 0x19, 0x81, 0x7b, 0x16, 0x3e, 0x93, 0x82, 0x27,}

	//sctf2023
	keyCry = []byte{0xf6, 0x38, 0xd8, 0x31, 0xa0, 0x17, 0x53, 0xb6,}
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
		if err != nil {
			panic(err)
		}

		if resp.StatusCode == http.StatusOK {
			handleInfoOK(body)
		} else {
			fmt.Println(string(body))
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
			// fmt.Println("hello")
		}
	}
}

type userInfoResp struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Phone string `json:"phone"`
	Msg   string `json:"message"`
}

func handleInfoOK (body []byte) {
	var infoResp userInfoResp
	bodyJSON := decryptoJSONByte(body)
	if err := json.Unmarshal(bodyJSON, &infoResp); err != nil {
		panic(err)
	}
	fmt.Println(decodeStr(userIdCry), infoResp.ID)
	fmt.Println(decodeStr(userNameCry), infoResp.Name)
	fmt.Println(decodeStr(userPhoneCry), infoResp.Phone)

	if infoResp.Msg != "" {
		fmt.Println(decodeStr(extMsgCry), infoResp.Msg)
	}
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
	fmt.Println(decodeStr(addressFmtCry), strings.Join(contactResp.Phones, ", "))
}


func decryptoJSONByte(jsonByte []byte) []byte {
	key := []byte(decodeStr(keyCry))
	rc4Enc, err := rc4.NewCipher(key)
	if err != nil {
		panic(err)
	}

	plain := make([]byte, len(jsonByte))
	rc4Enc.XORKeyStream(plain, jsonByte)

	return plain
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