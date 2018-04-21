package main

import (
	"fmt"
	"time"

	"runtime"

	"context"
	"github.com/bborbe/backup-cleanup-cron/backup_cleaner"
	"github.com/bborbe/cron"
	flag "github.com/bborbe/flagenv"
	"github.com/bborbe/lock"
	"github.com/golang/glog"
)

const (
	defaultKeepAmount   = 5
	lockName            = "/var/run/backup-cleanup-cron.lock"
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

	if err := do(); err != nil {
		glog.Exit(err)
	}
}

func do() error {
	dir := *targetDirPtr
	match := *matchPtr
	keepAmount := *keepAmountPtr
	wait := *waitPtr
	oneTime := *oneTimePtr
	lockName := *lockPtr

	l := lock.NewLock(lockName)
	if err := l.Lock(); err != nil {
		return err
	}
	defer func() {
		if err := l.Unlock(); err != nil {
			glog.Warningf("unlock failed: %v", err)
		}
	}()

	glog.V(0).Info("backup cleanup cron started")
	defer glog.V(0).Info("backup cleanup cron finished")

	if len(dir) == 0 {
		return fmt.Errorf("parameter %s missing", parameterDirectory)
	}
	if len(match) == 0 {
		return fmt.Errorf("parameter %s missing", parameterMatch)
	}
	if keepAmount <= 0 {
		return fmt.Errorf("parameter %s missing", parameterKeepAmount)
	}

	backupCleaner := backup_cleaner.New()
	glog.V(1).Infof("dir: %s, match: %s, keepAmount %d, wait: %v, oneTime: %v, lockName: %s", dir, match, keepAmount, wait, oneTime, lockName)

	action := func(ctx context.Context) error {
		return backupCleaner.CleanupBackup(dir, match, keepAmount)
	}

	var c cron.Cron
	if *oneTimePtr {
		c = cron.NewOneTimeCron(action)
	} else {
		c = cron.NewWaitCron(
			*waitPtr,
			action,
		)
	}
	return c.Run(context.Background())
}
