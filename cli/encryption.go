package cli

import (
	"cn-oxford/cryptoshield"
	"fmt"

	"github.com/Songmu/prompter"
)

func FileEncryptionHandler(target, out, password string) {
	if password == "" {
		password = prompter.Password("Set password")

		if prompter.Password("Repeat password") != password {
			fmt.Println("Passwords do not match.")
			return
		}
	}

	enc := cryptoshield.NewEncryptor()

	if err := enc.EncryptFile(target, out, password); err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		return
	}

	fmt.Printf("\nEncryption succeeded.\n%s --> %s\n", target, out)
}
