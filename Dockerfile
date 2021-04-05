# Use an official golang runtime as a parent image
FROM golang:1.16 as builder
ENV GO111MODULE=on
ADD . /src
WORKDIR /src

ENV TZ=Europe/Minsk
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone
RUN apt update -y
RUN apt-get -y install clang llvm make golang libpcap-dev
RUN go get -u github.com/gobuffalo/packr/packr && packr
RUN GOOS=linux GOARCH=amd64 go build -v -gcflags='-N -l'

FROM golang:1.16
RUN  apt update -y
RUN  apt-get -y install clang llvm make golang libpcap-dev
WORKDIR /root/
COPY --from=builder /src/kube-knark .
CMD ["./kube-knark"]