package backup_cleaner

import (
	"testing"

	"os"

	"time"

	. "github.com/bborbe/assert"
)

func TestImplementsBackupCleaner(t *testing.T) {
	c := New()
	var i *BackupCleaner
	if err := AssertThat(c, Implements(i)); err != nil {
		t.Fatal(err)
	}
}

func TestListBackups(t *testing.T) {
	list, err := listBackups("/tmp", "notexistingbackupmatch")
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(len(list), Is(0)); err != nil {
		t.Fatal(err)
	}
}

func TestGetBackupsToDelete(t *testing.T) {
	list := getBackupsToDelete(createFileInfos(10, 0), 5)
	if err := AssertThat(len(list), Is(5)); err != nil {
		t.Fatal(err)
	}
}

func TestGetBackupsToDeleteLessBackups(t *testing.T) {
	list := getBackupsToDelete(createFileInfos(4, 0), 5)
	if err := AssertThat(len(list), Is(0)); err != nil {
		t.Fatal(err)
	}
}

func TestGetBackupsToDeleteSub(t *testing.T) {
	list := getBackupsToDelete(createFileInfos(6, 42), 5)
	if err := AssertThat(len(list), Is(1+42)); err != nil {
		t.Fatal(err)
	}
}

func TestGetBackupsToDeleteWithEmpty(t *testing.T) {
	list := getBackupsToDelete(createFileInfos(10, 42), 5)
	if err := AssertThat(len(list), Is(5+42)); err != nil {
		t.Fatal(err)
	}
}

func TestGetBackupsToDeleteLessBackupsWithEmpty(t *testing.T) {
	list := getBackupsToDelete(createFileInfos(4, 42), 5)
	if err := AssertThat(len(list), Is(0+42)); err != nil {
		t.Fatal(err)
	}
}

func TestGetBackupsToDeleteSubWithEmpty(t *testing.T) {
	list := getBackupsToDelete(createFileInfos(6, 42), 5)
	if err := AssertThat(len(list), Is(1+42)); err != nil {
		t.Fatal(err)
	}
}

func createFileInfos(notEmpty int, empty int) []os.FileInfo {
	result := make([]os.FileInfo, notEmpty+empty)
	for i := 0; i < notEmpty; i++ {
		result[i] = &fileInfo{
			size: 1337,
		}
	}
	for i := 0; i < empty; i++ {
		result[i+notEmpty] = &fileInfo{
			size: 0,
		}

	}
	return result
}

type fileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
	isDir   bool
	sys     interface{}
}

func (f *fileInfo) Name() string {
	return f.name
}

func (f *fileInfo) Size() int64 {
	return f.size
}

func (f *fileInfo) Mode() os.FileMode {
	return f.mode
}

func (f *fileInfo) ModTime() time.Time {
	return f.modTime
}

func (f *fileInfo) IsDir() bool {
	return f.isDir
}

func (f *fileInfo) Sys() interface{} {
	return f.sys
}
