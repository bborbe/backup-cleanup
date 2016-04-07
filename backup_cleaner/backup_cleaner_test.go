package backup_cleaner

import (
	"testing"

	. "github.com/bborbe/assert"
)

func TestImplementsBackupCleaner(t *testing.T) {
	c := New()
	var i *BackupCleaner
	err := AssertThat(c, Implements(i))
	if err != nil {
		t.Fatal(err)
	}
}
