package backup_cleaner

import (
	"io/ioutil"

	"fmt"
	"os"
	"regexp"
	"sort"

	"github.com/golang/glog"
)

type backupCleaner struct {
}

type BackupCleaner interface {
	CleanupBackup(directory string, match string, keepAmount int) error
}

func New() *backupCleaner {
	return new(backupCleaner)
}

func (b *backupCleaner) CleanupBackup(directory string, match string, keepAmount int) error {
	glog.V(2).Infof("backup cleanup started")

	allBackups, err := listBackups(directory, match)
	if err != nil {
		return err
	}
	glog.V(2).Infof("found %d backups", len(allBackups))

	sort.Sort(FileInfoByName(allBackups))

	toDeleteBackups := getBackupsToDelete(allBackups, keepAmount)
	glog.V(2).Infof("found %d backups to delete", toDeleteBackups)

	if err = deleteBackups(directory, toDeleteBackups); err != nil {
		return err
	}

	glog.V(2).Infof("backup cleanup finished")
	return nil
}

func getBackupsToDelete(allBackups []os.FileInfo, keepAmount int) []os.FileInfo {
	emptyBackups := emptyFiles(allBackups)
	glog.V(2).Infof("found %d empty backups to delete", len(emptyBackups))
	notEmptyBackups := notEmptyFiles(allBackups)
	glog.V(2).Infof("found %d not empty backups", len(notEmptyBackups))
	pos := len(notEmptyBackups) - keepAmount
	if pos < 0 {
		glog.V(2).Infof("nothing to delete => return only empty backups")
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
		glog.V(2).Infof("delete backup %s", file.Name())
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
		glog.V(4).Infof("found %s", f.Name())
		if f.IsDir() {
			glog.V(4).Infof("skip directory %s", f.Name())
			continue
		}
		if re.MatchString(f.Name()) {
			glog.V(4).Infof("%s matches %s", f.Name(), match)
			glog.V(2).Infof("add backup %s", f.Name())
			list = append(list, f)
		} else {
			glog.V(4).Infof("%s mismatches %s", f.Name(), match)
		}
	}
	return list, nil
}
