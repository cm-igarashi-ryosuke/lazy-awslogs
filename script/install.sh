#!/bin/sh

package=`echo $PWD |sed -e "s,$GOPATH/src/,,g"`
vers=`git describe --dirty --always`

go install -a -ldflags "-X $package/cmd.appVersion=$vers"
