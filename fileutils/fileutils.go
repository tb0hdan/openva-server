package fileutils

import (
	"io"
	"os"

	"github.com/pkg/errors"
)

// MoveFile - actually move file across directories. https://stackoverflow.com/a/50741908
func MoveFile(sourcePath, destPath string) error {
	inputFile, err := os.Open(sourcePath)
	if err != nil {
		return errors.Wrap(err, "couldn't open source file")
	}
	outputFile, err := os.Create(destPath)
	if err != nil {
		_ = inputFile.Close()
		return errors.Wrap(err, "couldn't open dest file")
	}
	defer outputFile.Close()
	_, err = io.Copy(outputFile, inputFile)
	_ = inputFile.Close()
	if err != nil {
		return errors.Wrap(err, "writing to output file failed")
	}
	// The copy was successful, so now delete the original file
	err = os.Remove(sourcePath)
	if err != nil {
		return errors.Wrap(err, "failed to remove original file")
	}
	return nil
}
