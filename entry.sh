#/bin/bash
echo -- start git pull
cd /go/src/github.com/chongzii6/cloud-provider-baremetal

git config --global user.email "chenjun@molitv.cn"
git config --global user.name "chongzii6"

git pull
rm -rf vendor
tar xf v2.tar.gz

echo -- go build
go build

/go/bin/dlv debug --headless --listen=:2345 --log -- --cloud-provider=htnm --cloud-config=config/htnm.cfg
