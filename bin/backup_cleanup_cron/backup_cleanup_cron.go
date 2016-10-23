package main

import (
	"fmt"
	"time"

	"runtime"

	"github.com/bborbe/backup_cleanup_cron/backup_cleaner"
	flag "github.com/bborbe/flagenv"
	"github.com/bborbe/lock"
	"github.com/golang/glog"
)

const (
	defaultKeepAmount   = 5
	lockName            = "/var/run/backup_cleanup_cron.lock"
	parameterKeepAmount = "keep"
	parameterDirectory  = "dir"
	parameterMatch      = "match"
	parameterWait       = "wait"
	parameterOneTime    = "one-time"
	parameterLock       = "lock"
)

var (
	targetDirPtr  = flag.String(parameterDirectory, "", "target directory")
	matchPtr      = flag.String(parameterMatch, "", "match")
	keepAmountPtr = flag.Int(parameterKeepAmount, defaultKeepAmount, "keep amount")
	waitPtr       = flag.Duration(parameterWait, time.Minute*60, "wait")
	oneTimePtr    = flag.Bool(parameterOneTime, false, "exit after first backup")
	lockPtr       = flag.String(parameterLock, lockName, "lock")
)

type CleanupBackup func(directory string, match string, keepAmount int) error

func main() {
	defer glog.Flush()
	glog.CopyStandardLogTo("info")
	flag.Parse()
	runtime.GOMAXPROCS(runtime.NumCPU())

	backupCleaner := backup_cleaner.New()
	err := do(
		backupCleaner.CleanupBackup,
		*targetDirPtr,
		*matchPtr,
		*keepAmountPtr,
		*waitPtr,
		*oneTimePtr,
		*lockPtr,
	)
	if err != nil {
		glog.Exit(err)
	}
}

func do(
	cleanupBackup CleanupBackup,
	dir string,
	match string,
	keepAmount int,
	wait time.Duration,
	oneTime bool,
	lockName string,
) error {
	l := lock.NewLock(lockName)
	if err := l.Lock(); err != nil {
		return err
	}
	defer func() {
		if err := l.Unlock(); err != nil {
			glog.Warningf("unlock failed: %v", err)
		}
	}()

	glog.V(2).Info("backup cleanup cron started")
	defer glog.V(2).Info("backup cleanup cron finished")

	if len(dir) == 0 {
		return fmt.Errorf("parameter %s missing", parameterDirectory)
	}
	if len(match) == 0 {
		return fmt.Errorf("parameter %s missing", parameterMatch)
	}
	if keepAmount <= 0 {
		return fmt.Errorf("parameter %s missing", parameterKeepAmount)
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
