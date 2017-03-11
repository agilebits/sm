package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/AgileBits/sm/secrets"

	"encoding/json"
)

const (
	cmdEncrypt = "encrypt"
	cmdDecrypt = "decrypt"
)

func main() {
	encryptCommand := flag.NewFlagSet(cmdEncrypt, flag.ExitOnError)
	decryptCommand := flag.NewFlagSet(cmdDecrypt, flag.ExitOnError)

	if len(os.Args) < 2 {
		mainUsage()
		return
	}

	switch os.Args[1] {
	case cmdEncrypt:
		encrypt(encryptCommand, os.Args[2:])
	case cmdDecrypt:
		decrypt(decryptCommand, os.Args[2:])
	default:
		fmt.Printf("Unknown command: %s\n", os.Args[1])
		mainUsage()
		return
	}
}

func mainUsage() {
	fmt.Printf("Usage: sm [%s|%s]\n", cmdEncrypt, cmdDecrypt)
}

func encrypt(cmd *flag.FlagSet, args []string) {
	env := cmd.String("env", "dev", "environment type: 'dev' or 'aws'")
	masterKeyID := cmd.String("master", "", "master key identifier")
	region := cmd.String("region", "", "AWS region, for example 'us-east-1'")
	cmd.Parse(args)

	reader := bufio.NewReader(os.Stdin)
	message, err := ioutil.ReadAll(reader)
	if err != nil {
		log.Fatal("failed to read:", err)
	}

	envelope, err := secrets.EncryptEnvelope(*env, *region, *masterKeyID, message)
	if err != nil {
		log.Fatal("failed to encrypt", err)
	}

	buf, err := json.Marshal(envelope)
	if err != nil {
		log.Fatal("failed to Marshal:", err)
	}

	fmt.Println(string(buf))
}

func decrypt(cmd *flag.FlagSet, args []string) {
	reader := bufio.NewReader(os.Stdin)
	message, err := ioutil.ReadAll(reader)
	if err != nil {
		log.Fatal("failed to read:", err)
	}

	envelope := &secrets.Envelope{}
	if err := json.Unmarshal(message, &envelope); err != nil {
		log.Fatal("failed to Unmarshal:", err)
	}

	result, err := secrets.DecryptEnvelope(envelope)
	if err != nil {
		log.Fatal("failed to DecryptEnvelope:", err)
	}

	fmt.Println(string(result))
}
