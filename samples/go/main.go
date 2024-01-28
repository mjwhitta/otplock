package main

import (
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
)

func main() {
	var e error
	var encHex string = "ENCHEX"
	var key []byte
	var otpURL string = "OTPURL"
	var r *http.Response
	var tmp []byte

	// Convert payload to byte array for decryption
	if tmp, e = hex.DecodeString(encHex); e != nil {
		fmt.Println(e.Error())
		os.Exit(1)
	}

	// Fetch decryption key
	if r, e = http.Get(otpURL); e != nil {
		fmt.Println(e.Error())
		os.Exit(2)
	}

	if key, e = io.ReadAll(r.Body); e != nil {
		fmt.Println(e.Error())
		os.Exit(3)
	}

	// Exit if key length doesn't match
	if len(key) != len(tmp) {
		return
	}

	// Decrypt
	for i, b := range key {
		tmp[i] ^= b
	}

	// Do stuff with decrypted payload
	fmt.Println(string(tmp))
}
