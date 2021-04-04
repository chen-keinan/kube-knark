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
#RUN git clone https://github.com/go-delve/delve.git $GOPATH/src/github.com/go-delve/delve
#WORKDIR $GOPATH/src/github.com/go-delve/delve
#RUN make install
### export dlv bin path
#RUN export PATH=$PATH:/home/vagrant/go/bin >> ~/.bashrc
#RUN export PATH=$PATH:/root/go/bin >> ~/.bashrc
WORKDIR /root/
#EXPOSE 2345:2345
COPY --from=builder /src/kube-knark .
#CMD ["dlv","--listen=:2345", "--headless=true", "--api-version=2", "--accept-multiclient", "exec" ,"./kube-knark"]
CMD ["./kube-knark"]