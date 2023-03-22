package cryptoshield

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha512"
	"io"
	"os"

	"golang.org/x/crypto/pbkdf2"
)

type Encryptor struct{}

func NewEncryptor() *Encryptor {
	return &Encryptor{}
}

func (e *Encryptor) EncryptFile(targetPath, outPath, password string) error {
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

	key, salt, err := e.NewKey(password)
	if err != nil {
		return err
	}

	if _, err := out.Write(salt); err != nil {
		return err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return err
	}

	if _, err := out.Write(iv); err != nil {
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

func (e *Encryptor) NewKey(password string) (key, salt []byte, err error) {
	saltBuf := make([]byte, 48)

	if _, err := io.ReadFull(rand.Reader, saltBuf); err != nil {
		return nil, nil, err
	}

	return pbkdf2.Key([]byte(password), saltBuf, 1_048_576, 32, sha512.New), saltBuf, nil
}
