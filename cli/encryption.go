package cli

import (
	"fmt"
	"genc/cryptoshield"

	"github.com/Songmu/prompter"
)

func DirEncryptionHandler(target, password string, deleteTargets bool) {
	if password == "" {
		password = prompter.Password("Set password")

		if prompter.Password("Repeat password") != password {
			fmt.Println("Passwords do not match.")
			return
		}
	}

	fmt.Println("\nEncrypting... Please wait...")

	enc := cryptoshield.NewEncryptor()

	files, encErrors := enc.EncryptDir(target, password, deleteTargets)

	for _, err := range encErrors {
		fmt.Printf("Error: %s\n", err.Error())
	}

	fmt.Printf("\nSuccessfully encrypted %d/%d files found.\n", files-len(encErrors), files)
}

func SingleFileEncryptionHandler(target, out, password string, deleteTarget bool) {
	if password == "" {
		password = prompter.Password("Set password")

		if prompter.Password("Repeat password") != password {
			fmt.Println("Passwords do not match.")
			return
		}
	}

	if out == "" {
		out = target + ".genc"
	}

	fmt.Println("\nEncrypting... Please wait...")

	enc := cryptoshield.NewEncryptor()

	if err := enc.EncryptFile(target, out, password, deleteTarget); err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		return
	}

	fmt.Printf("\nEncryption succeeded.\n%s --> %s\n", target, out)
}
