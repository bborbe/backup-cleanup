package backup_cleaner

import (
	"io/ioutil"

	"fmt"
	"os"
	"regexp"
	"sort"

	"github.com/bborbe/log"
)

var logger = log.DefaultLogger

type backupCleaner struct {
}

type BackupCleaner interface {
	CleanupBackup(directory string, match string, keepAmount int) error
}

func New() *backupCleaner {
	return new(backupCleaner)
}

func (b *backupCleaner) CleanupBackup(directory string, match string, keepAmount int) error {
	logger.Debugf("backup cleanup started")

	allBackups, err := listBackups(directory, match)
	if err != nil {
		return err
	}
	logger.Debugf("found %d backups", len(allBackups))

	sort.Sort(FileInfoByName(allBackups))

	toDeleteBackups := getBackupsToDelete(allBackups, keepAmount)
	logger.Debugf("found %d backups to delete", toDeleteBackups)

	if err = deleteBackups(directory, toDeleteBackups); err != nil {
		return err
	}

	logger.Debugf("backup cleanup finished")
	return nil
}

func getBackupsToDelete(allBackups []os.FileInfo, keepAmount int) []os.FileInfo {
	emptyBackups := emptyFiles(allBackups)
	logger.Debugf("found %d empty backups to delete", len(emptyBackups))
	notEmptyBackups := notEmptyFiles(allBackups)
	logger.Debugf("found %d not empty backups", len(notEmptyBackups))
	pos := len(notEmptyBackups) - keepAmount
	if pos < 0 {
		logger.Debugf("nothing to delete => return only empty backups")
		return emptyBackups
	}
	return append(notEmptyBackups[0:pos], emptyBackups...)
}

func emptyFiles(files []os.FileInfo) []os.FileInfo {
	return filterFiles(files, func(file os.FileInfo) bool {
		return file.Size() == 0
	})
}

func notEmptyFiles(files []os.FileInfo) []os.FileInfo {
	return filterFiles(files, func(file os.FileInfo) bool {
		return file.Size() != 0
	})
}

func filterFiles(files []os.FileInfo, filter func(file os.FileInfo) bool) []os.FileInfo {
	result := make([]os.FileInfo, 0)
	for _, file := range files {
		if filter(file) {
			result = append(result, file)
		}
	}
	return result
}

func deleteBackups(directory string, files []os.FileInfo) error {
	for _, file := range files {
		logger.Debugf("delete backup %s", file.Name())
		if err := os.Remove(fmt.Sprintf("%s/%s", directory, file.Name())); err != nil {
			return err
		}
	}
	return nil
}

func listBackups(directory string, match string) ([]os.FileInfo, error) {
	re, err := regexp.Compile(match)
	if err != nil {
		return nil, err
	}
	files, err := ioutil.ReadDir(directory)
	if err != nil {
		return nil, err
	}
	list := make([]os.FileInfo, 0)
	for _, f := range files {
		logger.Tracef("found %s", f.Name())
		if f.IsDir() {
			logger.Tracef("skip directory %s", f.Name())
			continue
		}
		if re.MatchString(f.Name()) {
			logger.Tracef("%s matches %s", f.Name(), match)
			logger.Debugf("add backup %s", f.Name())
			list = append(list, f)
		} else {
			logger.Tracef("%s mismatches %s", f.Name(), match)
		}
	}
	return list, nil
}
