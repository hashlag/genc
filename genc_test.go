package main

import (
	"crypto/rand"
	"encoding/hex"
	"genc/cryptoshield"
	"io"
	"os"
	"strings"
	"testing"
)

func CreateTestFile() (string, []byte, error) {
	file, err := os.CreateTemp("./", "testdata_*.bin")
	if err != nil {
		return "", nil, err
	}

	defer file.Close()

	randomData := make([]byte, 1024)

	_, err = io.ReadFull(rand.Reader, randomData)
	if err != nil {
		return "", nil, err
	}

	_, err = file.Write(randomData)
	if err != nil {
		return "", nil, err
	}

	return file.Name(), randomData, nil
}

func TestEncryptFileAndDecryptFile(t *testing.T) {
	filename, expectedData, err := CreateTestFile()
	if err != nil {
		t.Error(err)
		return
	}

	enc := cryptoshield.NewEncryptor()

	err = enc.EncryptFile(filename, filename+".genc", "1234567890", true)
	if err != nil {
		t.Error(err)
		return
	}

	dec := cryptoshield.NewDecryptor()

	err = dec.DecryptFile(filename+".genc", filename+".genc.decrypted", "1234567890", true)
	if err != nil {
		t.Error(err)
		return
	}

	decryptedFile, err := os.Open(filename + ".genc.decrypted")
	if err != nil {
		t.Error(err)
		return
	}

	defer func() {
		decryptedFile.Close()
		os.Remove(filename + ".genc.decrypted")
	}()

	decryptedData, err := io.ReadAll(decryptedFile)
	if err != nil {
		t.Error(err)
		return
	}

	if hex.EncodeToString(decryptedData) != hex.EncodeToString(expectedData) {
		t.Error("The data has been corrupted during encryption/decryption.")
	}
}

func TestCorruptionDetection(t *testing.T) {
	filename, _, err := CreateTestFile()
	if err != nil {
		t.Error(err)
		return
	}

	enc := cryptoshield.NewEncryptor()

	err = enc.EncryptFile(filename, filename+".genc", "1234567890", true)
	if err != nil {
		t.Error(err)
		return
	}

	encryptedFile, err := os.OpenFile(filename+".genc", os.O_WRONLY, 0644)
	if err != nil {
		t.Error(err)
		return
	}

	defer func() {
		os.Remove(filename + ".genc")
		os.Remove(filename + ".genc.decrypted")
	}()

	_, err = encryptedFile.Seek(-200, 2)
	if err != nil {
		t.Error(err)
		return
	}

	randomData, err := io.ReadAll(io.LimitReader(rand.Reader, 10))
	if err != nil {
		t.Error(err)
		return
	}

	_, err = encryptedFile.Write(randomData)
	if err != nil {
		t.Error(err)
		return
	}

	encryptedFile.Close()

	dec := cryptoshield.NewDecryptor()

	err = dec.DecryptFile(filename+".genc", filename+".genc.decrypted", "1234567890", true)
	if err != nil {
		if !strings.Contains(err.Error(), "CORRUPTED/TAMPERED") {
			t.Error(err)
			return
		}
	} else {
		t.Error("corruption of the file was not detected")
		return
	}
}

func BenchmarkEncryptFile(b *testing.B) {
	enc := cryptoshield.NewEncryptor()

	for n := 0; n < b.N; n++ {
		enc.EncryptFile("localres/testdata.bin", "localres/testdata.enc.bin", "00000000000000000000", false)
	}
}
