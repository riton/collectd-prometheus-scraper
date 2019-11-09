package collectd

import (
	"context"
	"os"

	"collectd.org/api"
	"collectd.org/format"
)

// FileWriter allows to write collectd ValueList
// to a file. It is primarilly intended to help
// for debugging
type FileWriter struct {
	FileName   string
	file       *os.File
	putvalMeta *format.PutvalWithMeta
}

func NewFileWriter(fileName string) *FileWriter {
	fw := FileWriter{
		FileName: fileName,
	}

	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	fw.file = file
	fw.putvalMeta = format.NewPutvalWithMeta(fw.file)

	return &fw
}

func (fw *FileWriter) Write(ctx context.Context, vl *api.ValueList) error {
	return fw.putvalMeta.Write(ctx, vl)
}
