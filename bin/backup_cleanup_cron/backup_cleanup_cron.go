package main

import (
	"fmt"
	"io"
	"os"
	"time"

	flag "github.com/bborbe/flagenv"
	"github.com/bborbe/lock"
	"github.com/bborbe/log"
	"github.com/bborbe/backup_cleanup_cron/backup_cleaner"
)

var logger = log.DefaultLogger

const (
	DEFAULT_KEEP_AMOUNT = 5
	LOCK_NAME = "/var/run/backup_cleanup_cron.lock"
	PARAMETER_LOGLEVEL = "loglevel"
	PARAMETER_KEEP_AMOUNT = "keep"
	PARAMETER_DIRECTORY = "dir"
	PARAMETER_PREFIX = "prefix"
	PARAMETER_WAIT = "wait"
	PARAMETER_ONE_TIME = "one-time"
	PARAMETER_LOCK = "lock"
)

type CleanupBackup func(directory string, prefix string, keepAmount int) error

func main() {
	defer logger.Close()
	logLevelPtr := flag.String(PARAMETER_LOGLEVEL, log.INFO_STRING, "one of OFF,TRACE,DEBUG,INFO,WARN,ERROR")
	targetDirPtr := flag.String(PARAMETER_DIRECTORY, "", "target directory")
	prefixPtr := flag.String(PARAMETER_PREFIX, "", "prefix")
	keepAmountPtr := flag.Int(PARAMETER_KEEP_AMOUNT, DEFAULT_KEEP_AMOUNT, "keep amount")
	waitPtr := flag.Duration(PARAMETER_WAIT, time.Minute * 60, "wait")
	oneTimePtr := flag.Bool(PARAMETER_ONE_TIME, false, "exit after first backup")
	lockPtr := flag.String(PARAMETER_LOCK, LOCK_NAME, "lock")

	flag.Parse()
	logger.SetLevelThreshold(log.LogStringToLevel(*logLevelPtr))
	logger.Debugf("set log level to %s", *logLevelPtr)

	backupCleaner := backup_cleaner.New()

	writer := os.Stdout
	err := do(writer, backupCleaner.CleanupBackup, *targetDirPtr, *prefixPtr, *keepAmountPtr, *waitPtr, *oneTimePtr, *lockPtr)
	if err != nil {
		logger.Fatal(err)
		logger.Close()
		os.Exit(1)
	}
}

func do(writer io.Writer, cleanupBackup CleanupBackup, dir string, prefix string, keepAmount int, wait time.Duration, oneTime bool, lockName string) error {
	l := lock.NewLock(lockName)
	if err := l.Lock(); err != nil {
		return err
	}
	defer l.Unlock()
	logger.Debug("start")
	defer logger.Debug("done")

	if len(dir) == 0 {
		return fmt.Errorf("parameter %s missing", PARAMETER_DIRECTORY)
	}
	if len(prefix) == 0 {
		return fmt.Errorf("parameter %s missing", PARAMETER_PREFIX)
	}
	if keepAmount <= 0 {
		return fmt.Errorf("parameter %s missing", PARAMETER_KEEP_AMOUNT)
	}

	logger.Debugf("dir: %s, keepAmount %s, wait: %v, oneTime: %v, lockName: %s", dir, keepAmount, wait, oneTime, lockName)

	for {
		logger.Debugf("backup started")
		if err := cleanupBackup(dir, prefix, keepAmount); err != nil {
			return err
		}
		logger.Debugf("backup completed")

		if oneTime {
			return nil
		}

		logger.Debugf("wait %v", wait)
		time.Sleep(wait)
		logger.Debugf("done")
	}
}
