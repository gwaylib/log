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
	"github.com/gwaylib/log/proto"
)

type fileInfo struct {
	fileIdx  int
	fileName string
}

type logFile struct {
	file *os.File
	size int64
}

type Adapter struct {
	fileName   string // 文件路径及名称
	backups    int    // 保留的文件数
	fileSize   int64  // 每个文件大小
	file       *logFile
	buffer     chan *proto.Proto
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
	dirs, err := ioutil.ReadDir(".")
	if err != nil {
		return errors.As(err)
	}
	baseName := filepath.Base(a.fileName)
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
		names = append(names, fileInfo{fileName: dirName, fileIdx: idx})
	}
	if len(names) == 0 {
		return nil
	}
	sort.Slice(names, func(i, j int) bool {
		return names[i].fileIdx < names[j].fileIdx
	})
	for i := len(names) - 1; i > -1; i-- {
		info := names[i]
		// remove file when over backups
		if info.fileIdx+1 >= a.backups {
			if err := os.Remove(info.fileName); err != nil {
				println(errors.As(err))
			}
			continue
		}
		// rename to new index
		if err := os.Rename(info.fileName, fmt.Sprintf("%s.%d", a.fileName, info.fileIdx+1)); err != nil {
			println(errors.As(err))
		}
	}
	return nil
}

func (a *Adapter) getFile() *logFile {
	if a.file != nil {
		return a.file
	}
	file, err := os.OpenFile(a.fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
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

func (a *Adapter) Put(p *proto.Proto) {
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
			for _, val := range p.Data {
			reWrite:
				date := val.Date.Format("2006-01-02T15:04:05.000Z07:00")
				output := fmt.Sprintf("%s [%s] (%s) %s",
					date,
					val.Level.String(),
					val.Logger,
					string(val.Msg),
				)

				file := a.getFile()
				if file == nil {
					println("file adapter not work, dupm to stderr")
					println(output)
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
				print(string(data))
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

func New(fileName string) proto.Adapter {
	return NewFile(fileName, math.MaxInt, math.MaxInt)
}

func NewFile(fileName string, backups int, fileSize int64) proto.Adapter {
	a := &Adapter{
		fileName:   fileName,
		backups:    backups,
		fileSize:   fileSize,
		buffer:     make(chan *proto.Proto, 100),
		closeEvent: make(chan bool, 1),
		doneEvent:  make(chan bool, 1),
	}
	go a.run()
	return a
}
