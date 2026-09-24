package file

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/gwaylib/log/logger"
	"github.com/gwaylib/log/proto"
)

func TestPut(t *testing.T) {
	path := filepath.Join(t.TempDir(), "file.log")
	// use a large size limit so that no rotation happens.
	log := logger.NewDefaultLogger("testing", NewFile(path, 10, 1024*1024))
	// output all the levels.
	log.SetOutputLevel(proto.LevelFatal)
	log.Debug("debug")
	log.Debug([]byte{0xff})
	log.Info("info")
	log.Warn("warn")
	log.Error("error")
	//log.Fatal("fatal")
	//log.Exit(0, "exit")
	log.Close()

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), "info") {
		t.Fatalf("expect the log file contains the message, got: %q", string(data))
	}
}

func TestRotate(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "rotate.log")

	const backups = 3
	const fileSize = 128
	log := logger.NewDefaultLogger("testing", NewFile(path, backups, fileSize))
	log.SetOutputLevel(proto.LevelFatal)
	for i := 0; i < 200; i++ {
		log.Infof("rotate message number %d", i)
	}
	log.Close()

	// the current file and its backups should be no more than backups+1.
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) > backups+1 {
		t.Fatalf("expect at most %d files, got %d", backups+1, len(entries))
	}
	if len(entries) < 2 {
		t.Fatalf("expect the file to be rotated, got %d files", len(entries))
	}
	// check the expected files exist.
	for i := 0; i <= backups; i++ {
		name := path
		if i > 0 {
			name = filepath.Join(dir, "rotate.log."+strconv.Itoa(i))
		}
		if _, err := os.Stat(name); err != nil {
			t.Fatalf("expect the file %s exists: %v", name, err)
		}
	}
}

// the adapter should not output to the stdio nor panic when the file can
// not be opened.
func TestFileError(t *testing.T) {
	dir := t.TempDir()
	// use a regular file as the parent directory, so that the log file can
	// not be created.
	blocker := filepath.Join(dir, "blocker")
	if err := os.WriteFile(blocker, []byte("x"), 0644); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(blocker, "sub", "file.log")

	log := logger.NewDefaultLogger("testing", NewFile(path, 10, 1024))
	log.SetOutputLevel(proto.LevelFatal)
	log.Info("should not be written")
	log.Close()

	if _, err := os.Stat(path); err == nil {
		t.Fatal("expect the log file not created")
	}
}
