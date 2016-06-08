package main

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func CalcMetadata(file string, writer io.Writer) error {
	fileReader, err := os.Open(file)
	if err != nil {
		return err
	}
	defer fileReader.Close()

	hasher := sha1.New()
	nBytes, err := io.Copy(hasher, fileReader)
	if err != nil {
		return err
	}
	sha1sum := hasher.Sum(nil)
	filename := filepath.Base(file)
	fmt.Fprintf(writer, "File: %s\n", filename)
	fmt.Fprintf(writer, "Size: %d\n", nBytes)
	fmt.Fprintf(writer, "SHA1: %s\n", hex.EncodeToString(sha1sum))
	fmt.Fprintln(writer)
	return nil
}
