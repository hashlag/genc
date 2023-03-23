package cli

import (
	"cn-oxford/cryptoshield"
	"fmt"

	"github.com/Songmu/prompter"
)

func FileDecryptionHandler(target, out, password string) {
	if password == "" {
		password = prompter.Password("Password")
	}

	dec := cryptoshield.NewDecryptor()

	if err := dec.DecryptFile(target, out, password); err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		return
	}

	fmt.Printf("\nDecrypted.\n%s --> %s\n", target, out)
}
