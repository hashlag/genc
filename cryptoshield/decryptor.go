package cryptoshield

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha512"
	"io"
	"os"

	"golang.org/x/crypto/pbkdf2"
)

type Decryptor struct{}

func NewDecryptor() *Decryptor {
	return &Decryptor{}
}

func (d *Decryptor) DecryptFile(targetPath, outPath, password string) error {
	target, err := os.Open(targetPath)
	if err != nil {
		return err
	}
	defer target.Close()

	out, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer out.Close()

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

	buf := make([]byte, 2048)

	for {
		n, err := target.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}

		if n > 0 {
			stream.XORKeyStream(buf[:n], buf[:n])
			if _, err := out.Write(buf[:n]); err != nil {
				return err
			}
		}

		if err == io.EOF {
			break
		}
	}

	return nil
}

func (e *Decryptor) DeriveKey(password string, salt []byte) []byte {
	return pbkdf2.Key([]byte(password), salt, 1_048_576, 32, sha512.New)
}
