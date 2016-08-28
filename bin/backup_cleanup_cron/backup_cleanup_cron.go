package main

import (
	"fmt"
	"io"
	"os"
	"time"

	"runtime"

	"github.com/bborbe/backup_cleanup_cron/backup_cleaner"
	flag "github.com/bborbe/flagenv"
	"github.com/bborbe/lock"
	"github.com/golang/glog"
)

const (
	DEFAULT_KEEP_AMOUNT   = 5
	LOCK_NAME             = "/var/run/backup_cleanup_cron.lock"
	PARAMETER_KEEP_AMOUNT = "keep"
	PARAMETER_DIRECTORY   = "dir"
	PARAMETER_MATCH       = "match"
	PARAMETER_WAIT        = "wait"
	PARAMETER_ONE_TIME    = "one-time"
	PARAMETER_LOCK        = "lock"
)

var (
	targetDirPtr  = flag.String(PARAMETER_DIRECTORY, "", "target directory")
	matchPtr      = flag.String(PARAMETER_MATCH, "", "match")
	keepAmountPtr = flag.Int(PARAMETER_KEEP_AMOUNT, DEFAULT_KEEP_AMOUNT, "keep amount")
	waitPtr       = flag.Duration(PARAMETER_WAIT, time.Minute*60, "wait")
	oneTimePtr    = flag.Bool(PARAMETER_ONE_TIME, false, "exit after first backup")
	lockPtr       = flag.String(PARAMETER_LOCK, LOCK_NAME, "lock")
)

type CleanupBackup func(directory string, match string, keepAmount int) error

func main() {
	defer glog.Flush()
	glog.CopyStandardLogTo("info")
	flag.Parse()
	runtime.GOMAXPROCS(runtime.NumCPU())

	backupCleaner := backup_cleaner.New()
	writer := os.Stdout
	err := do(writer, backupCleaner.CleanupBackup, *targetDirPtr, *matchPtr, *keepAmountPtr, *waitPtr, *oneTimePtr, *lockPtr)
	if err != nil {
		glog.Exit(err)
	}
}

func do(writer io.Writer, cleanupBackup CleanupBackup, dir string, match string, keepAmount int, wait time.Duration, oneTime bool, lockName string) error {
	l := lock.NewLock(lockName)
	if err := l.Lock(); err != nil {
		return err
	}
	defer l.Unlock()
	glog.V(2).Info("backup cleanup cron started")
	defer glog.V(2).Info("backup cleanup cron finished")

	if len(dir) == 0 {
		return fmt.Errorf("parameter %s missing", PARAMETER_DIRECTORY)
	}
	if len(match) == 0 {
		return fmt.Errorf("parameter %s missing", PARAMETER_MATCH)
	}
	if keepAmount <= 0 {
		return fmt.Errorf("parameter %s missing", PARAMETER_KEEP_AMOUNT)
	}

	glog.V(2).Infof("dir: %s, match: %s, keepAmount %d, wait: %v, oneTime: %v, lockName: %s", dir, match, keepAmount, wait, oneTime, lockName)

	for {
		glog.V(2).Infof("backup started")
		if err := cleanupBackup(dir, match, keepAmount); err != nil {
			return err
		}
		glog.V(2).Infof("backup completed")

		if oneTime {
			return nil
		}

		glog.V(2).Infof("wait %v", wait)
		time.Sleep(wait)
		glog.V(2).Infof("sleep done")
	}
}
