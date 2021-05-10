package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
)

func readFile(filename string) []byte {
	fileRead, error := ioutil.ReadFile(filename)

	if error != nil {
		panic(error.Error())
	}

	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()

	return fileRead
}

func hashKey(key string) []byte {
	hashedKey := md5.New()
	io.WriteString(hashedKey, key)
	return hashedKey.Sum(nil)
}

func encryptFile(filename string, key string) {
	fileRead := readFile(filename)

	hashedKey := hashKey(key)

	block, error := aes.NewCipher(hashedKey)

	if error != nil {
		panic(error.Error())
	}

	aesGCM, error := cipher.NewGCM(block)
	if error != nil {
		panic(error.Error())
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, error = io.ReadFull(rand.Reader, nonce); error != nil {
		panic(error.Error())
	}

	ciphertext := aesGCM.Seal(nonce, nonce, fileRead, nil)

	error = ioutil.WriteFile(filename, ciphertext, 0644)

	if error != nil {
		panic(error.Error())
	}

	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()
}

func decryptFile(filename string, key string) string {

	fileRead := readFile(filename)

	hashedKey := hashKey(key)

	block, err := aes.NewCipher(hashedKey)

	if err != nil {
		panic(err.Error())
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}

	nonceSize := aesGCM.NonceSize()

	nonce, ciphertext := fileRead[:nonceSize], fileRead[nonceSize:]

	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)

	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()

	if err != nil {
		panic(err.Error())
	}

	return fmt.Sprintf("%s", plaintext)
}

var (
	filename = flag.String("filename", "", "<filename> of the file to encrypt or decrypt")
	key      = flag.String("key", "", "<key>  master secret to be used for encryption or decryption")
	mode     = flag.String("mode", "e", "<mode> either e (encryption) or d (for decryption)")
)

func main() {

	flag.Parse()

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage %s:\n", os.Args[0])

		flag.VisitAll(func(f *flag.Flag) {
			fmt.Fprintf(os.Stderr, "    %v\n", f.Usage)
		})
	}

	if *filename == "" || *key == "" {
		flag.Usage()
		return
	}

	switch *mode {
	case "d":
		fmt.Println(decryptFile(*filename, *key))
	case "e":
		encryptFile(*filename, *key)
	default:
		panic("Invalid mode")
	}

	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()
}
