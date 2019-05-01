FROM golang:1.12-alpine as builder
RUN apk add git
COPY . /go/src/shuCourse
ENV GO111MODULE on
WORKDIR /go/src/shuCourse/cli
RUN go get && go build
WORKDIR /go/src/shuCourse/web
RUN go get && go build

FROM alpine
MAINTAINER longfangsong@icloud.com
COPY --from=builder /go/src/shuCourse/web/web /
COPY --from=builder /go/src/shuCourse/cli/cli /
WORKDIR /
CMD ./web
ENV PORT 8000
EXPOSE 8000