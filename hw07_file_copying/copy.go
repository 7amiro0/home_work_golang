package main

import (
	"errors"
	"github.com/cheggaaa/pb/v3"
	"io"
	"os"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	fromFile, errorOpen := os.Open(fromPath)
	outFile, errorCreate := os.Create(toPath)

	if errorOpen != nil {
		return ErrUnsupportedFile
	}
	if errorCreate != nil {
		return ErrUnsupportedFile
	}

	infoFile, _ := os.Stat(fromPath)
	if infoFile.Size() <= offset {
		return ErrOffsetExceedsFileSize
	}
	if _, err := fromFile.Seek(offset, io.SeekStart); err != nil {
		return err
	}

	if limit == 0 {
		limit = infoFile.Size() - offset
	}

	bar := pb.Full.Start64(limit)
	rd := bar.NewProxyReader(fromFile)
	io.CopyN(outFile, rd, limit)
	bar.Finish()

	return nil
}
