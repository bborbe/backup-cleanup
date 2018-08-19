FROM golang:1.10 AS build
COPY . /go/src/github.com/bborbe/backup-cleanup
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-s" -a -installsuffix cgo -o /backup-cleanup ./src/github.com/bborbe/backup-cleanup
CMD ["/bin/bash"]

FROM scratch
MAINTAINER Benjamin Borbe <bborbe@rocketnews.de>

ENV LOCK /backup-cleanup.run
ENV WAIT 1h
ENV ONE_TIME false
ENV KEEP 5
ENV DIR /backup
ENV MATCH backup_.*sql

VOLUME ["/backup"]

COPY  --from=build backup-cleanup /
COPY files/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
ENTRYPOINT ["/backup-cleanup"]
