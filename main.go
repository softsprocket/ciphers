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
	"encrypt-string": "usage: %s encrypt-string --public <filename> --in <filename> --out <filename>",
	"decrypt-string": "usage: %s decrypt-string --private <filename> --in <filename> --out <filename>",
}

func printUsageAndExit(command string, exitValue int) {
	fmt.Printf(usageStrings[command], filepath.Base(os.Args[0]))
	os.Exit(exitValue)
}

func parseArguments(arguments []string) map[string]string {
	var args = arguments[2:]
	argMap := make(map[string]string)

	for i := 0; i < len(args); i++ {
		runes := []rune(args[i])
		key := string(runes[2:])
		i++
		argMap[key] = args[i]
	}

	return argMap
}

func getKeyFileNames(command string, arguments []string) map[string]string {
	files := parseArguments(arguments)

	ok := false

	if command == "gen-keys" {
		if _, ok = files["private"]; ok {
			_, ok = files["public"]
		}
	} else {
		if command == "encrypt-string" {
			_, ok = files["public"]
		} else if command == "decrypt-string" {
			_, ok = files["private"]
		}

		if ok {
			if _, ok = files["in"]; ok {
				_, ok = files["out"]
			}
		}
	}
	if !ok {
		printUsageAndExit(command, -1)
	}

	return files

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

		files := getKeyFileNames("gen-keys", arguments)

		genKeyPair(files["private"], files["public"])

	case "encrypt-string":
		if len(arguments) < 8 {
			printUsageAndExit("encrypt-string", -1)
		}

		files := getKeyFileNames("encrypt-string", arguments)

		publicKey, err := ioutil.ReadFile(files["public"])

		if err != nil {

		}

		toEncrypt, err := ioutil.ReadFile(files["in"])

		encrypted := encryptString(string(toEncrypt), publicKey)

		ioutil.WriteFile(files["out"], encrypted, 0644)

	case "decrypt-string":
		if len(arguments) < 8 {
			printUsageAndExit("decrypt-string", -1)
		}

		files := getKeyFileNames("decrypt-string", arguments)

		privateKey, err := ioutil.ReadFile(files["private"])

		if err != nil {

		}

		toDecrypt, err := ioutil.ReadFile(files["in"])

		decrypted := decryptToString(toDecrypt, privateKey)

		ioutil.WriteFile(files["out"], []byte(decrypted), 0644)
	}
}
