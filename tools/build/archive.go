package main

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
)

type ZipWorker struct {
	zipWriter *zip.Writer
	root      string
}

func NewZipWorker(zipFile io.Writer, root string) *ZipWorker {
	return &ZipWorker{
		zipWriter: zip.NewWriter(zipFile),
		root:      root,
	}
}

func (worker *ZipWorker) run() error {
	defer worker.close()
	return filepath.Walk(worker.root, worker.zipAllFiles)
}

func (worker *ZipWorker) zipAllFiles(path string, info os.FileInfo, err error) error {
	if info.IsDir() {
		return nil
	}
	fileWriter, err := worker.zipWriter.Create(path)
	if err != nil {
		return err
	}
	fileReader, err := os.Open(path)
	if err != nil {
		return err
	}
	_, err = io.Copy(fileWriter, fileReader)
	if err != nil {
		return err
	}
	return nil
}

func (worker *ZipWorker) close() {
	worker.zipWriter.Close()
}

func zipFolder(folder string, file string) error {
	if _, err := os.Stat(file); err == nil {
		os.Remove(file)
	}

	zipFile, err := os.Create(file)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	return NewZipWorker(zipFile, folder).run()
}
