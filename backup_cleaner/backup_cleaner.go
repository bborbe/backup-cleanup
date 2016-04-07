package backup_cleaner

import (
	"io/ioutil"

	"github.com/bborbe/log"
	"strings"
	"os"
"sort"
)

var logger = log.DefaultLogger

type backupCleaner struct {
}

type BackupCleaner interface {
	CleanupBackup(directory string, prefix string, keepAmount int) error
}

func New() *backupCleaner {
	return new(backupCleaner)
}

func (b *backupCleaner) CleanupBackup(directory string, prefix string, keepAmount int) error {
	logger.Debugf("backup cleanup started")

	allBackups, err := listBackups(directory, prefix)
	if err != nil {
		return err
	}

	sort.Sort(FileInfoByName(allBackups))


	logger.Debugf("backup cleanup finished")
	return nil
}

func listBackups(directory string, prefix string) ([]os.FileInfo, error) {
	files, err := ioutil.ReadDir(directory)
	if err != nil {
		return nil, err
	}
	list := make([]os.FileInfo, 0)
	for _, f := range files {
		if f.IsDir() {
			continue
		}

		if strings.Index(f.Name(), prefix) != -1 {
			list = append(list, f)
		}
	}
	return list, nil
}