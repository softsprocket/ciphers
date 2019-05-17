package main

import (
	"ciphers/cipher"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func genKeyPair(privateFileName, publicFileName string) {
	priv, pub := cipher.GenerateKeyPair(4096)
	err := ioutil.WriteFile(privateFileName, cipher.PrivateKeyToBytes(priv), 0644)

	if err != nil {
		log.Panic(err)
	}

	err = ioutil.WriteFile(publicFileName, cipher.PublicKeyToBytes(pub), 0644)

	if err != nil {
		log.Panic(err)
	}

}

func encryptString(str string, pubKey []byte) []byte {
	pub := cipher.BytesToPublicKey(pubKey)
	enbuf := cipher.EncryptWithPublicKey([]byte(str), pub)

	return enbuf
}

func decryptToString(buf []byte, privKey []byte) string {
	priv := cipher.BytesToPrivateKey(privKey)
	debuf := cipher.DecryptWithPrivateKey(buf, priv)

	return string(debuf)
}

func main() {
	fmt.Println("Hello, Ciphers")
	arguments := os.Args

	if len(arguments) < 2 {
		fmt.Printf("usage: %s <command>", filepath.Base(arguments[0]))
		fmt.Println("commands: gen-keys help")
		os.Exit(-1)
	}

	command := arguments[1]

	switch command {
	case "help":
		if len(arguments) < 3 {
			fmt.Printf("usage: %s help <command>", filepath.Base(arguments[0]))
			fmt.Println("commands: gen-keys")
			os.Exit(-1)
		}

		switch arguments[2] {
		case "gen-keys":
			fmt.Printf("usage: %s gen-keys --private <filename> --public <filename>", filepath.Base(arguments[0]))
		}

		os.Exit(0)
	case "gen-keys":
		if len(arguments) < 6 {
			fmt.Printf("usage: %s gen-keys --private <filename> --public <filename>", filepath.Base(arguments[0]))
			os.Exit(-1)
		}

		files := make(map[string]string)
		runes := []rune(arguments[2])
		first := string(runes[2:])
		runes = []rune(arguments[4])
		second := string(runes[2:])

		files[first] = arguments[3]
		files[second] = arguments[5]

		_, privOk := files["private"]
		_, pubOk := files["public"]
		if !privOk || !pubOk {

		}

		genKeyPair(files["private"], files["public"])

	}
}
