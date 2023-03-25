package cli

import (
	"cn-oxford/cryptoshield"
	"fmt"
	"strings"

	"github.com/Songmu/prompter"
)

func DirDecryptionHandler(target, password string, deleteTargets bool) {
	if password == "" {
		password = prompter.Password("Password")
	}

	dec := cryptoshield.NewDecryptor()

	files, decErrors := dec.DecryptDir(target, password, deleteTargets)

	for _, err := range decErrors {
		fmt.Printf("Error: %s\n", err.Error())
	}

	fmt.Printf("\nSuccessfully decrypted %d/%d files found.\n", files-len(decErrors), files)
}

func SingleFileDecryptionHandler(target, out, password string, deleteTarget bool) {
	if out == "" {
		if strings.HasSuffix(target, ".genc") {
			out = target[:len(target)-5]
		} else {
			fmt.Println("The target file does not have .genc extension.\nOutput flag (-o) is required.")
			return
		}
	}

	if password == "" {
		password = prompter.Password("Password")
	}

	dec := cryptoshield.NewDecryptor()

	if err := dec.DecryptFile(target, out, password, deleteTarget); err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		return
	}

	fmt.Printf("\nDecrypted.\n%s --> %s\n", target, out)
}
