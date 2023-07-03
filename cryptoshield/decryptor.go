package cryptoshield

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/sha512"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/crypto/pbkdf2"
)

type Decryptor struct{}

func NewDecryptor() *Decryptor {
	return &Decryptor{}
}

func (d *Decryptor) DecryptDir(targetPath, password string, deleteTargets bool) (int, []error) {
	var (
		filesFound int
		decErrors  []error
	)

	filepath.Walk(targetPath, func(path string, info fs.FileInfo, err error) error {
		if !info.IsDir() && strings.HasSuffix(path, ".genc") {
			filesFound++

			err := d.DecryptFile(path, path[:len(path)-5], password, deleteTargets)
			if err != nil {
				decErrors = append(decErrors, err)
			}
		}
		return nil
	})

	return filesFound, decErrors
}

func (d *Decryptor) DecryptFile(targetPath, outPath, password string, deleteTarget bool) error {
	target, err := os.Open(targetPath)
	if err != nil {
		return err
	}

	defer func(p string, d bool) {
		target.Close()
		if d {
			os.Remove(targetPath)
		}
	}(targetPath, deleteTarget)

	out, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer out.Close()

	expectedMAC := make([]byte, 64)
	if _, err := io.ReadFull(target, expectedMAC); err != nil {
		return err
	}

	salt := make([]byte, 48)
	if _, err := io.ReadFull(target, salt); err != nil {
		return err
	}

	block, err := aes.NewCipher(d.DeriveKey(password, salt))
	if err != nil {
		return err
	}

	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(target, iv); err != nil {
		return err
	}

	stream := cipher.NewCTR(block, iv)

	dataMAC := hmac.New(sha512.New, pbkdf2.Key([]byte(password), salt, 524_288, 64, sha512.New))

	buf := make([]byte, 102400)

	for {
		n, err := target.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}

		if n > 0 {
			_, err := dataMAC.Write(buf[:n])
			if err != nil {
				return err
			}

			stream.XORKeyStream(buf[:n], buf[:n])

			if _, err := out.Write(buf[:n]); err != nil {
				return err
			}
		}

		if err == io.EOF {
			break
		}
	}

	if !hmac.Equal(dataMAC.Sum(nil), expectedMAC) {
		return fmt.Errorf("wrong password or %s has been CORRUPTED/TAMPERED", targetPath)
	}

	return nil
}

func (e *Decryptor) DeriveKey(password string, salt []byte) []byte {
	return pbkdf2.Key([]byte(password), salt, 1_048_576, 32, sha512.New)
}
