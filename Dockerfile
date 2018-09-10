FROM golang:1.10

COPY entry.sh /go/bin/entry.sh
RUN go get -u github.com/derekparker/delve/cmd/dlv \
  && mkdir -p /go/src/github.com/chongzii6/cloud-provider-baremetal \
  && chmod a+x /go/bin/entry.sh

#COPY . /go/src/github.com/chongzii6/cloud-provider-baremetal
#RUN cd /go/src/github.com/chongzii6/cloud-provider-baremetal \
#  && go build

EXPOSE 2345
WORKDIR /go/src/github.com/chongzii6/cloud-provider-baremetal

# CMD /go/bin/dlv debug --headless --listen=:2345 --log
ENTRYPOINT /go/bin/entry.sh

