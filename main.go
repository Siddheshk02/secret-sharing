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
		//z := strconv.Itoa(l)
		//fmt.Fprintf(file, z)
		fmt.Fprintln(file, n[l])
	}
	fmt.Println("Shares created Successfully")
	a := len(n[0])
	//fmt.Println(a)

	var j int
label:
	fmt.Print("Enter the Number of Secret Shares you want to enter: ")
	fmt.Scanf("%d", &j)

	if j < 2 {
		fmt.Println("Please enter the minimum number of Shares i.e. 2")
		fmt.Println(" ")

		goto label
	} else if j > shares {
		fmt.Println("Exceeded the Number of Shares!!")

		goto label
	}

	var parts [10][51]byte

	for i := 0; i < j; i++ {
		fmt.Print("Enter the Secret Share: ")
		for x := 0; x < a; x++ {
			fmt.Scan(&parts[i][x])

		}
		fmt.Println(" ")

	}

	for h := 0; h < j; h++ {
		var con string
		if (parts[h][0] == parts[h+1][0]) && (parts[h][1] == parts[h+1][1]) {
			fmt.Println("Share ", h+1, " is repeated")
			fmt.Print("Do you Want to continue? Yes/No :")
			fmt.Scan(&con)
			fmt.Println(" ")
			if con == "Yes" {
				goto label
			} else {
				return
			}

		}
	}

	boolean := true
	var loc int

	for i := 0; i < j; i++ {
		boolean = true
		for x := 0; x < shares; x++ {
			if (parts[i][0] == n[x][0]) && (parts[i][1] == n[x][1]) {
				//var eg [60]byte = parts[i]
				//bool := bytes.Equal(eg, n[x])

				for p := 0; p < len(n[0]); p++ {
					if parts[i][p] != n[x][p] {
						boolean = false
						loc = i + 1
						goto label1
					}
				}
			}
		}
	}

label1:

	if boolean == false {
		var con string
		fmt.Println("Invalid Share ", loc)
		fmt.Print("Do you Want to continue? Yes/No :")
		fmt.Scan(&con)
		fmt.Println(" ")
		if con == "Yes" {
			goto label
		} else {
			return
		}
	}

	fileloc := "C:/Users/user/go/src/github.com/Siddheshk02/secret-sharing/files/secret.bin"

	ciphertext, err := os.ReadFile(fileloc)
	if err != nil {
		log.Fatal(err)
	}

	nonce1 := ciphertext[:gcm.NonceSize()]
	ciphertext = ciphertext[gcm.NonceSize():]
	plaintext, err := gcm.Open(nil, nonce1, ciphertext, nil)
	if err != nil {
		log.Panic(err)
	}

	err = ioutil.WriteFile("Users/encrypted.txt", plaintext, 0777)
	if err != nil {
		log.Panic(err)
	}
	fmt.Println("File Decrypted Successfully!!")

	//fmt.Println(parts[0])
	//fmt.Println(parts[1])
	//fmt.Println(n[0][1])

	//for _, i := range n {
	//	fmt.Println(i)
	//}

	//fmt.Println(secret)
	//fmt.Println(key)
	//fmt.Println(text)
}
