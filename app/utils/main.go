package main

import (
	"bufio"
	"crypto/md5"
	"crypto/rc4"
	"encoding/hex"
	"fmt"
	"os"
)

func rc4Encode() {
	key := []byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff}
	rc4Enc, err := rc4.NewCipher(key)
	if err != nil {
		panic(err)
	}

	scan := bufio.NewReader(os.Stdin)
	input, _, _ := scan.ReadLine()

	var encsStr  = make([]byte, len(input))
	rc4Enc.XORKeyStream(encsStr, input)
	printByteArr(encsStr)
}

func printByteArr(data []byte) {
	for _, val := range data {
		fmt.Printf("%#02x, ", val)
	}
}

func md5Encode(src string) string {
	md5Cipher := md5.Sum([]byte(src))
	result := hex.EncodeToString(md5Cipher[:])
	return result
}

var password = []string{"12345","8726t","ajuhs","alojs","alijs",}

func main() {
	// for _, val := range password {
	// 	fmt.Println(md5Encode(val))
	// }

	rc4Encode()
}