package backup_cleaner

import (
	"testing"

	"os"

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
	list, err := listBackups("/tmp", "notexistingbackupprefix")
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(len(list), Is(0)); err != nil {
		t.Fatal(err)
	}
}

func TestGetBackupsToDelete(t *testing.T) {
	list := getBackupsToDelete(make([]os.FileInfo, 10), 5)
	if err := AssertThat(len(list), Is(5)); err != nil {
		t.Fatal(err)
	}
}

func TestGetBackupsToDeleteLessBackups(t *testing.T) {
	list := getBackupsToDelete(make([]os.FileInfo, 4), 5)
	if err := AssertThat(len(list), Is(0)); err != nil {
		t.Fatal(err)
	}
}

func TestGetBackupsToDeleteSub(t *testing.T) {
	list := getBackupsToDelete(make([]os.FileInfo, 6), 5)
	if err := AssertThat(len(list), Is(1)); err != nil {
		t.Fatal(err)
	}
}
