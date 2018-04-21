# Cron for cleanup Backups

## Install

```
go get github.com/bborbe/backup-cleanup-cron
```

## Run Backup

One time

```
backup-cleanup-cron \
-logtostderr \
-v=2 \
-lock=/backup/backup-cleanup-cron.lock \
-dir=/backup \
-match='backup_.*tar.gz' \
-keep=5 \
-one-time
```

Cron

```
backup-cleanup-cron \
-logtostderr \
-v=2 \
-lock=/backup/backup-cleanup-cron.lock \
-dir=/backup \
-match='backup_.*tar.gz' \
-keep=5 \
-wait=1h
```
