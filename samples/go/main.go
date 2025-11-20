package main

import (
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
)

func main() {
	var decHex []byte
	var e error
	var encHex string = "ENCHEX"
	var key []byte
	var otpURL string = "OTPURL"
	var r *http.Response

	// Convert payload to byte array for decryption
	if decHex, e = hex.DecodeString(encHex); e != nil {
		fmt.Println(e.Error())
		return
	}

	// Fetch decryption key
	if r, e = http.Get(otpURL); e != nil {
		fmt.Println(e.Error())
		return
	}
	defer func() {
		if e := r.Body.Close(); e != nil {
			panic(e)
		}
	}()

	if key, e = io.ReadAll(r.Body); e != nil {
		fmt.Println(e.Error())
		return
	}

	// Exit if key length doesn't match
	if len(key) != len(decHex) {
		return
	}

	// Decrypt
	for i, b := range key {
		decHex[i] ^= b
	}

	// Do stuff with decrypted payload (probably not print it)
	fmt.Println(string(decHex))
}
