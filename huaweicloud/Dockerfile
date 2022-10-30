FROM golang:latest as BUILDER

MAINTAINER zengchen1024<chenzeng765@gmail.com>

# build binary
WORKDIR /go/src/github.com/opensourceways/xihe-training-center
COPY . .
RUN cd huaweicloud && GO111MODULE=on CGO_ENABLED=0 go build -a -o xihe-training-center .
RUN tar -xf ./huaweicloud/trainingimpl/tools/obsutil.tar.gz

# copy binary config and utils
FROM alpine:3.14
RUN apk update && apk add --no-cache \
        git \
        bash \
        libc6-compat
COPY --from=BUILDER /go/src/github.com/opensourceways/xihe-training-center/huaweicloud/xihe-training-center /opt/app/xihe-training-center
COPY --from=BUILDER /go/src/github.com/opensourceways/xihe-training-center/obsutil /opt/app/obsutil
COPY --from=BUILDER /go/src/github.com/opensourceways/xihe-training-center/huaweicloud/trainingimpl/tools/sync_files.sh /opt/app/sync_file.sh
COPY --from=BUILDER /go/src/github.com/opensourceways/xihe-training-center/huaweicloud/trainingimpl/tools/upload_folder.sh /opt/app/upload_folder.sh

ENTRYPOINT ["/opt/app/xihe-training-center"]
