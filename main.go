package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"
	"log"

	"github.com/hashicorp/vault/shamir"
)

func main() {
	//fmt.Println("Hello World!!")
	pass := "Hello, World!!"
	shares := 5
	threshold := 2

	c := sha256.New()

	c.Write([]byte(pass))
	key := c.Sum(nil)

	text, err := ioutil.ReadFile("store.txt")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(text))

	block, err := aes.NewCipher(key)
	if err != nil {
		log.Fatal(err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		log.Fatal(err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		log.Fatal(err)
	}

	var secret []byte = gcm.Seal(nonce, nonce, text, nil)

	err = ioutil.WriteFile("secret.bin", secret, 0777)
	if err != nil {
		log.Fatal(err)
	}

	n, err := shamir.Split(secret, shares, threshold)
	// fmt.Println(n)
	if err != nil {
		log.Fatal(err)
	}

	var parts []byte
	var j int

	fmt.Scanf("Enter the Number of Secret Shares you want to enter: ", &j)
	for i := 0; i < j; i++ {
		fmt.Scanf("Enter the Secret Share: ", &parts[i])
	}

	for _, i := range n {
		fmt.Println(i)
	}

	//fmt.Println(secret)
	//fmt.Println(key)
	//fmt.Println(text)
}
