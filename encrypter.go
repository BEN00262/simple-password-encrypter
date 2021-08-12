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

func encryptFile(filename string, fileRead []byte, key string) {
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

func decryptFile(filename string, key string) ([]byte, error) {

	fileRead := readFile(filename)

	block, err := aes.NewCipher(hashKey(key))

	if err != nil {
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(block)

	if err != nil {
		return nil, err
	}

	nonceSize := aesGCM.NonceSize()

	nonce, ciphertext := fileRead[:nonceSize], fileRead[nonceSize:]

	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)

	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

var (
	filename = flag.String("filename", "", "<filename> of the file to encrypt or decrypt")
	key      = flag.String("key", "", "<key>  master secret to be used for encryption or decryption")

	// save read
	mode = flag.String("mode", "r", "<mode> either r (reading) or s (saving) or d (deleting)")

	// used for the saving and retrieving of the passwords
	target   = flag.String("target", "", "<target> target of the password to be saved")
	password = flag.String("password", "", "<password> the password to save")
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
	case "r":
		{
			fileData, err := decryptFile(*filename, *key)

			if err != nil {
				fmt.Println(err.Error())
				return
			}

			credentials := parseJSONFile(fileData)

			if *target == "" {
				for i := 0; i < len(credentials); i++ {
					fmt.Printf("Target: %s\n", credentials[i].Target)
				}

				fmt.Println()
				return
			}

			foundCredential := searchCredential(*target, credentials)

			if foundCredential != "" {
				err := CopyToClipboard(foundCredential)

				if err != nil {
					fmt.Println(err.Error())
					return
				}

				fmt.Println("Copied to clipboard")
				return
			}

			fmt.Println("Password for the given Target not found :(")
		}
	case "s":
		{
			if *target == "" || *password == "" {
				fmt.Println("Cannot save the credential. Pass the required arguments")
				return
			}

			var credentials []Credential = []Credential{}

			_, err := os.Stat(*filename)

			if !os.IsNotExist(err) {
				fileData, err := decryptFile(*filename, *key)

				if err != nil {
					fmt.Println(err.Error())
					return
				}

				credentials = parseJSONFile(fileData)
			}

			credentials = insertCredential(Credential{
				Target:   *target,
				Password: *password,
			}, credentials)

			jsonBytes, err := dumpToJSONFile(&credentials)

			if err != nil {
				fmt.Println("Failed to save password")
				return
			}

			encryptFile(*filename, jsonBytes, *key)
		}
	case "d":
		{

			if *target == "" {
				fmt.Println("Target should not be empty :(")
				return
			}

			fileData, err := decryptFile(*filename, *key)

			if err != nil {
				fmt.Println(err.Error())
				return
			}

			credentials := parseJSONFile(fileData)
			credentials = removeCredential(*target, credentials)

			dumpedJSON, err := dumpToJSONFile(&credentials)

			if err != nil {
				fmt.Println("Failed to delete credentials for target")
				return
			}

			encryptFile(*filename, dumpedJSON, *key)
		}
	default:
		fmt.Println("Invalid mode")
	}

	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()
}
