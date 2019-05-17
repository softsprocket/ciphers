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

var usageStrings = map[string]string{
	"app":            "usage: %s <command> <args>",
	"help":           "usage: %s help <command>\ncommands: gen-keys encrypt-string decrypt-string",
	"gen-keys":       "usage: %s gen-keys --private <filename> --public <filename>",
	"encrypt-string": "usage: %s encrypt-string <string>",
	"decrypt-string": "usage: %s decrypt-string <string>",
}

func printUsageAndExit(command string, exitValue int) {
	fmt.Printf(usageStrings[command], filepath.Base(os.Args[0]))
	os.Exit(exitValue)
}

func main() {
	arguments := os.Args

	if len(arguments) < 2 {
		printUsageAndExit("app", -1)
	}

	command := arguments[1]

	switch command {
	case "help":
		if len(arguments) < 3 {
			printUsageAndExit("help", -1)
		}

		printUsageAndExit(arguments[2], 0)

	case "gen-keys":
		if len(arguments) < 6 {
			printUsageAndExit("gen-keys", -1)
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
			printUsageAndExit("gen-keys", -1)
		}

		genKeyPair(files["private"], files["public"])

	case "encrypt-string":
		if len(arguments) < 3 {
			printUsageAndExit("encrypt-string", -1)
		}
	case "decrypt-string":
		if len(arguments) < 3 {
			printUsageAndExit("encrypt-string", -1)
		}
	}
}
