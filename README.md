# Cron for cleanup Backups

## Install

```
go get github.com/bborbe/backup-cleanup
```

## Run Backup

One time

```
backup-cleanup \
-logtostderr \
-v=2 \
-lock=/backup/backup-cleanup.lock \
-dir=/backup \
-match='backup_.*tar.gz' \
-keep=5 \
-one-time
```

Cron

```
backup-cleanup \
-logtostderr \
-v=2 \
-lock=/backup/backup-cleanup.lock \
-dir=/backup \
-match='backup_.*tar.gz' \
-keep=5 \
-wait=1h
```
