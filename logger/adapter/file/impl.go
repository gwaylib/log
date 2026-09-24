package file

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"time"

	"github.com/gwaylib/errors"
	"github.com/gwaylib/log/proto"
)

const (
	// Retry times when writing failed, so that the log can survive a
	// momentary file system error.
	writeRetryTimes = 3
	// Delay between write retry.
	writeRetryDelay = 200 * time.Millisecond
)

// logFile is a rolling file writer.
//
// The rolling rule is referenced from gwaysys/supd/internal/filebak:
// when the current file reaches maxSize, it is moved to "<name>.1", the
// existing backups are shifted up ("<name>.i" -> "<name>.i+1"), and the
// oldest backup ("<name>.<backups>") is removed. So there is at most
// `backups` backup files kept beside the current one.
type logFile struct {
	file    *os.File
	name    string
	curSize int64
	maxSize int64
	backups int
}

func openLogFile(name string, maxSize int64, backups int) (*logFile, error) {
	f, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		if strings.HasSuffix(err.Error(), "no such file or directory") {
			// make the parent directory, then retry.
			if err := os.MkdirAll(filepath.Dir(name), 0755); err != nil {
				return nil, errors.As(err, name)
			}
			return openLogFile(name, maxSize, backups)
		}
		return nil, errors.As(err, name)
	}
	stat, err := f.Stat()
	if err != nil {
		f.Close()
		return nil, errors.As(err, name)
	}
	return &logFile{
		file:    f,
		name:    name,
		curSize: stat.Size(),
		maxSize: maxSize,
		backups: backups,
	}, nil
}

// move renames a file, ignoring the case while the source is not exist.
func (lf *logFile) move(src, dst string) error {
	if err := os.Rename(src, dst); err != nil {
		if !strings.HasSuffix(err.Error(), "no such file or directory") {
			return errors.As(err, src)
		}
	}
	return nil
}

// rotate closes the current file, shifts the backups, and opens a new
// current file.
func (lf *logFile) rotate() error {
	if err := lf.file.Close(); err != nil {
		return errors.As(err, lf.name)
	}
	lf.file = nil

	name := lf.name
	if lf.backups > 0 {
		// remove the oldest backup.
		if err := os.RemoveAll(fmt.Sprintf("%s.%d", name, lf.backups)); err != nil {
			return errors.As(err, name)
		}
		// shift the backups.
		for i := lf.backups - 1; i > 0; i-- {
			if err := lf.move(fmt.Sprintf("%s.%d", name, i), fmt.Sprintf("%s.%d", name, i+1)); err != nil {
				return err
			}
		}
		// backup the current file.
		if err := lf.move(name, fmt.Sprintf("%s.1", name)); err != nil {
			return err
		}
	}

	// open a new empty current file.
	f, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return errors.As(err, name)
	}
	lf.file = f
	lf.curSize = 0
	return nil
}

func (lf *logFile) write(data []byte) error {
	n, err := lf.file.Write(data)
	if err != nil {
		return errors.As(err, lf.name)
	}
	lf.curSize += int64(n)
	if lf.maxSize > 0 && lf.curSize >= lf.maxSize {
		return lf.rotate()
	}
	return nil
}

func (lf *logFile) close() error {
	if lf.file == nil {
		return nil
	}
	err := lf.file.Close()
	lf.file = nil
	return err
}

type Adapter struct {
	filePath   string // file path contain the base name
	backups    int    // backup files
	fileSize   int64  // size of each file
	file       *logFile
	buffer     chan *proto.Proto
	closing    int32
	closeEvent chan bool
	doneEvent  chan bool
}

func (a *Adapter) closeFile() {
	if a.file != nil {
		// keep the output in file only, ignore the close error.
		_ = a.file.close()
		a.file = nil
	}
}

func (a *Adapter) getFile() *logFile {
	if a.file != nil {
		return a.file
	}
	f, err := openLogFile(a.filePath, a.fileSize, a.backups)
	if err != nil {
		return nil
	}
	a.file = f
	return a.file
}

// write writes one log line to the file.
//
// It retries on failure by reopening the file, and drops the log when it
// still can not be written, never output to the stdio.
func (a *Adapter) write(data []byte) {
	for i := 0; i < writeRetryTimes; i++ {
		f := a.getFile()
		if f == nil {
			// can not open the file, wait for a while and retry.
			time.Sleep(writeRetryDelay)
			continue
		}
		if err := f.write(data); err != nil {
			// reopen the file, wait for a while and retry.
			a.closeFile()
			time.Sleep(writeRetryDelay)
			continue
		}
		return
	}
}

func (a *Adapter) Close() {
	a.closeEvent <- true
	<-a.doneEvent
	a.closeFile()
}

// Put puts a log protocol to the log queue.
func (a *Adapter) Put(p *proto.Proto) {
	if p == nil {
		panic("argument is nil")
	}
	if atomic.LoadInt32(&a.closing) != 0 {
		// the adapter is closing, drop the log.
		return
	}
	a.buffer <- p
}

func (a *Adapter) run() {
	for {
		select {
		case p := <-a.buffer:
			for _, val := range p.Data {
				date := val.Date.Format("2006-01-02T15:04:05.000Z07:00")
				output := fmt.Sprintf("%s [%s] (%s) %s",
					date,
					val.Level.String(),
					val.Logger,
					string(val.Msg),
				)
				a.write([]byte(output))
			}
		case <-a.closeEvent:
			atomic.StoreInt32(&a.closing, 1)
			if len(a.buffer) > 0 {
				a.closeEvent <- true
				continue
			}
			a.doneEvent <- true
			return
		}
	}
}

func New(path string) proto.Adapter {
	return NewFile(path, math.MaxInt, math.MaxInt)
}

// NewFile creates a file adapter.
//
// Params
//
//	path     -- the log file path, the parent directory is created if not exist.
//	backups  -- the max count of the backup files.
//	fileSize -- the max size of a single log file in bytes, the file is
//	            rotated when it reached the size.
func NewFile(path string, backups int, fileSize int64) proto.Adapter {
	a := &Adapter{
		filePath:   path,
		backups:    backups,
		fileSize:   fileSize,
		buffer:     make(chan *proto.Proto, 100),
		closeEvent: make(chan bool, 1),
		doneEvent:  make(chan bool, 1),
	}
	go a.run()
	return a
}
