package cryptoshield

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha512"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"golang.org/x/crypto/pbkdf2"
)

type Encryptor struct{}

func NewEncryptor() *Encryptor {
	return &Encryptor{}
}

func (e *Encryptor) EncryptDir(targetPath, password string, deleteTargets bool) (int, []error) {
	var (
		filesFound int
		encErrors  []error
	)

	filepath.Walk(targetPath, func(path string, info fs.FileInfo, err error) error {
		if !info.IsDir() {
			filesFound++

			err := e.EncryptFile(path, path+".genc", password, deleteTargets)
			if err != nil {
				encErrors = append(encErrors, err)
			}
		}
		return nil
	})

	return filesFound, encErrors
}

func (e *Encryptor) EncryptFile(targetPath, outPath, password string, deleteTarget bool) error {
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

	buf := make([]byte, 102400)

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
