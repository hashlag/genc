package cli

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

func Root() {
	mode := flag.String("m", "", "Mode: encryption (-m E) or decryption (-m D). Required.")
	target := flag.String("t", "", "Path to the target file or directory. Required.")
	out := flag.String("o", "", "Path to the output file. Ignored when target is a directory. Optional flag.")
	password := flag.String("p", "", "Allows user to set password without additional prompting. Use with caution as the password may remain in the command line history. Optional flag.")
	deleteTargets := flag.Bool("d", false, "Delete target files after encryption/decryption. Optional flag.")

	flag.Parse()

	if *mode == "" {
		fmt.Println("Missing required mode flag (-m).")
		return
	}

	if *target == "" {
		fmt.Println("Missing required target flag (-t).")
		return
	}

	targetStat, err := os.Stat(*target)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		return
	}

	switch strings.ToUpper(*mode) {
	case "E":
		if targetStat.IsDir() {
			DirEncryptionHandler(*target, *password, *deleteTargets)
		} else {
			SingleFileEncryptionHandler(*target, *out, *password, *deleteTargets)
		}
	case "D":
		if targetStat.IsDir() {
			DirDecryptionHandler(*target, *password, *deleteTargets)
		} else {
			SingleFileDecryptionHandler(*target, *out, *password, *deleteTargets)
		}
	default:
		fmt.Println("Unknown mode (-m).")
		return
	}
}
