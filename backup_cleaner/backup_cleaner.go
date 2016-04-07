package backup_cleaner

import (
	"io/ioutil"

	"os"
	"sort"
	"strings"

	"github.com/bborbe/log"
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
	logger.Debugf("found %d backups", len(allBackups))

	sort.Sort(FileInfoByName(allBackups))

	toDeleteBackups := getBackupsToDelete(allBackups, keepAmount)
	logger.Debugf("found %d backups to delete", toDeleteBackups)

	if err = deleteBackups(toDeleteBackups); err != nil {
		return err
	}

	logger.Debugf("backup cleanup finished")
	return nil
}

func getBackupsToDelete(allBackups []os.FileInfo, keepAmount int) []os.FileInfo {
	pos := len(allBackups) - keepAmount
	if pos < 0 {
		logger.Debugf("nothing to delete")
		return nil
	}
	return allBackups[0:pos]
}

func deleteBackups(files []os.FileInfo) error {
	for _, file := range files {
		logger.Debugf("delete backup %s", file.Name())
	}
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
			logger.Debugf("add backup %s", f.Name())
			list = append(list, f)
		}
	}
	return list, nil
}
