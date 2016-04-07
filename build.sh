#!/bin/sh

SOURCEDIRECTORY="github.com/bborbe/backup_cleanup_cron"
VERSION="1.0.1-b${BUILD_NUMBER}"
NAME="backup-cleanup-cron"

################################################################################

export GOROOT=/opt/go
export PATH=/opt/go2xunit/bin/:/opt/utils/bin/:/opt/aptly_utils/bin/:/opt/aptly/bin/:/opt/debian_utils/bin/:/opt/debian/bin/:$GOROOT/bin:$PATH
export GOPATH=${WORKSPACE}
export REPORT_DIR=${WORKSPACE}/test-reports
INSTALLS=`cd src && find $SOURCEDIRECTORY/bin -name "*.go" | dirof | unique`
DEB="${NAME}_${VERSION}.deb"
rm -rf $REPORT_DIR ${WORKSPACE}/*.deb ${WORKSPACE}/pkg
mkdir -p $REPORT_DIR
PACKAGES=`cd src && find $SOURCEDIRECTORY -name "*_test.go" | dirof | unique`
FAILED=false
for PACKAGE in $PACKAGES
do
  XML=$REPORT_DIR/`pkg2xmlname $PACKAGE`
  OUT=$XML.out
  go test -i $PACKAGE
  go test -v $PACKAGE | tee $OUT
  cat $OUT
  go2xunit -fail=true -input $OUT -output $XML
  rc=$?
  if [ $rc != 0 ]
  then
    echo "Tests failed for package $PACKAGE"
    FAILED=true
  fi
done

if $FAILED
then
  echo "Tests failed => skip install"
  exit 1
else
  echo "Tests success"
fi

echo "Tests completed, install to $GOPATH"

go install $INSTALLS

echo "Install completed, create debian package"

create_debian_package \
-loglevel=DEBUG \
-version=$VERSION \
-config=src/$SOURCEDIRECTORY/create_debian_package_config.json || exit 1

echo "Create debian package completed, start upload to aptly"

aptly_upload \
-loglevel=DEBUG \
-url=https://www.benjamin-borbe.de/aptly \
-username=api \
-passwordfile=$HOME/aptly_api_password \
-file=$DEB \
-repo=unstable

echo "Upload completed"