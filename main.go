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
	"os"
	"strconv"

	"github.com/hashicorp/vault/shamir"
)

func main() {
	//fmt.Println("Hello World!!")

	if len(os.Args) != 2 {
		fmt.Println("Please provide a filepath")
		return
	}
	filepath := os.Args[1]
	/*file, err := os.Open(filepath)
	if err != nil {
		fmt.Println("Error opening %s: %s", filepath, err)
	}*/

	pass := "Hello, World!!"
	shares := 5
	threshold := 2

	c := sha256.New()

	c.Write([]byte(pass))
	key := c.Sum(nil)

	text, err := os.ReadFile(filepath)
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println(string(text))

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

	err = ioutil.WriteFile("files/secret.bin", secret, 0777)
	if err != nil {
		log.Fatal(err)
	}

	n, err := shamir.Split(secret, shares, threshold)
	//fmt.Printf("%t", n[0][0])
	// fmt.Println(n)
	if err != nil {
		log.Fatal(err)
	}

	//err = ioutil.WriteFile("files/share.txt", n, 0777)
	file, err := os.Create("files/shares.txt")
	if err != nil {
		log.Fatal("os.Create", err)
	}

	for l := 0; l < shares; l++ {
		z := strconv.Itoa(l)
		fmt.Fprintf(file, z)
		fmt.Fprintln(file, n[l])
	}
	fmt.Println("Shares created Successfully")
	a := len(n[0])

	var j int

	fmt.Print("Enter the Number of Secret Shares you want to enter: ")
	fmt.Scanf("%d", &j)
	if j < 2 {
		fmt.Println("Please enter the minimum number of Shares i.e. 2")
		return
	} else if j > shares {
		fmt.Println("Exceeded the Number of Shares!!")
		return
	}
	var parts [60][60]byte
	for i := 0; i < j; i++ {
		fmt.Println("Enter the Secret Share: ")
		for x := 0; x < a; x++ {
			fmt.Scanf("%d", &parts[i][x])
		}

	}
	fmt.Println("DONE")

	//for _, i := range n {
	//	fmt.Println(i)
	//}

	//fmt.Println(secret)
	//fmt.Println(key)
	//fmt.Println(text)
}
