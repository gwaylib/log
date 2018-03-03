package file

import (
	"os"

	"github.com/gwaylib/log/logger"
	"github.com/gwaylib/log/logger/adapter/stdio"
)

type Adapter struct {
	file *os.File
	logger.Adapter
}

func (a *Adapter) Close() {
	if a.file != nil {
		a.file.Close()
	}
	a.Adapter.Close()
}

func New(fileName string) logger.Adapter {
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	return &Adapter{
		file:    file,
		Adapter: stdio.New(file),
	}
}
