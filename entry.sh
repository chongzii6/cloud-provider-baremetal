#/bin/bash
echo -- start git pull
cd /go/src/github.com/chongzii6/cloud-provider-baremetal
git pull

echo -- go build
go build

/go/bin/dlv debug --headless --listen=:2345 --log
