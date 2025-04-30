package file

import (
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gwaylib/errors"
	"github.com/gwaylib/log/behavior"
)

type fileInfo struct {
	fileIdx  int
	fileDir  string
	fileName string
}

type logFile struct {
	file *os.File
	size int64
}

type Adapter struct {
	filePath   string // file path contain the base name
	backups    int    // backup files
	fileSize   int64  // size of each file
	file       *logFile
	buffer     chan *behavior.Event
	closing    bool
	closeEvent chan bool
	doneEvent  chan bool
}

func (a *Adapter) closeFile() {
	if a.file != nil {
		if err := a.file.file.Close(); err != nil {
			println(errors.As(err))
		}
	}
	a.file = nil
}

func (a *Adapter) reBackup() error {
	fileDir := filepath.Dir(a.filePath)
	dirs, err := ioutil.ReadDir(fileDir)
	if err != nil {
		return errors.As(err)
	}
	baseName := filepath.Base(a.filePath)
	names := []fileInfo{}
	for _, dir := range dirs {
		dirName := dir.Name()
		if !strings.HasPrefix(dirName, baseName) {
			continue
		}

		dotNames := strings.Split(dirName, ".")
		idx, err := strconv.Atoi(dotNames[len(dotNames)-1])
		if err != nil {
			if baseName != dirName {
				println(errors.As(err))
				continue
			}
		}
		names = append(names, fileInfo{fileDir: fileDir, fileName: dirName, fileIdx: idx})
	}
	if len(names) == 0 {
		return nil
	}
	sort.Slice(names, func(i, j int) bool {
		return names[i].fileIdx < names[j].fileIdx
	})
	for i := len(names) - 1; i > -1; i-- {
		info := names[i]
		filePath := filepath.Join(info.fileDir, info.fileName)
		// remove file when over backups
		if info.fileIdx+1 >= a.backups {
			if err := os.Remove(filePath); err != nil {
				println(errors.As(err))
			}
			continue
		}
		// rename to new index
		if err := os.Rename(filePath, fmt.Sprintf("%s.%d", a.filePath, info.fileIdx+1)); err != nil {
			println(errors.As(err))
		}
	}
	return nil
}

func (a *Adapter) getFile() *logFile {
	if a.file != nil {
		return a.file
	}
	file, err := os.OpenFile(a.filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		println(errors.As(err))
		return nil
	}
	stat, err := file.Stat()
	if err != nil {
		println(errors.As(err))
		file.Close()
		return nil
	}
	if stat.Size() > a.fileSize {
		file.Close()
		if err := a.reBackup(); err != nil {
			println(errors.As(err))
			return nil
		}
		return a.getFile()
	}
	a.file = &logFile{
		file: file,
		size: stat.Size(),
	}
	return a.file
}

func (a *Adapter) Close() {
	a.closeEvent <- true
	<-a.doneEvent
	a.closeFile()
}

func (a *Adapter) Put(p *behavior.Event) {
	if a.closing {
		println(errors.New("adapter is closing"))
		return
	}
	a.buffer <- p
}

func (a *Adapter) run() {
	for {
		select {
		case p := <-a.buffer:
		reWrite:
			output := string(p.ToJson()) + "\n"
			file := a.getFile()
			if file == nil {
				println("file adapter not work, dupm to stderr")
				println(string(output))
				continue
			}
			data := []byte(output)
			w, err := file.file.Write(data)
			if err != nil {
				println(errors.As(err))
				a.closeFile()
				time.Sleep(1e9)
				goto reWrite
			}
			file.size += int64(w)
			if file.size >= a.fileSize {
				a.closeFile()
			}
		case <-a.closeEvent:
			a.closing = true
			if len(a.buffer) > 0 {
				a.closeEvent <- true
				continue
			}
			a.doneEvent <- true
			return
		}
	}
}

func New(path string) behavior.Client {
	return NewFile(path, math.MaxInt, math.MaxInt)
}

func NewFile(path string, backups int, fileSize int64) behavior.Client {
	a := &Adapter{
		filePath:   path,
		backups:    backups,
		fileSize:   fileSize,
		buffer:     make(chan *behavior.Event, 100),
		closeEvent: make(chan bool, 1),
		doneEvent:  make(chan bool, 1),
	}
	go a.run()
	return a
}
