package cli

import (
	"flag"
	"fmt"
	"strings"
)

func Root() {
	mode := flag.String("m", "", "Mode: encryption (-m E) or decryption (-m D). Required.")
	target := flag.String("t", "", "Path to the target file. Required.")
	out := flag.String("o", "", "Path to the output file. Required.")
	password := flag.String("p", "", "Allows user to set password without additional prompting. Use with caution as the password may remain in the command line history. Optional flag.")

	flag.Parse()

	if *target == "" {
		fmt.Println("Missing required target flag (-t).")
		return
	}

	if *out == "" {
		fmt.Println("Missing required output flag (-o).")
		return
	}

	switch strings.ToUpper(*mode) {
	case "E":
		FileEncryptionHandler(*target, *out, *password)
	case "D":
		FileDecryptionHandler(*target, *out, *password)
	default:
		fmt.Println("Unknown mode (-m).")
		return
	}
}
