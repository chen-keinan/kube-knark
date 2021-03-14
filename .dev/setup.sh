#!/bin/bash

echo "Provisioning virtual machine..."

### install clang llvm make
echo "install clang llvm make"
sudo apt-get update -y
sudo apt-get -y install clang llvm make

### install golang pkg
echo "install golang pkg"
sudo yum install -y golang

### Install dlv pkg
echo "Install dlv pkg"
git clone https://github.com/go-delve/delve.git $GOPATH/src/github.com/go-delve/delve
cd $GOPATH/src/github.com/go-delve/delve
make install

### export dlv bin path
export PATH=$PATH:/home/vagrant/go/bin
export PATH=$PATH:/root/go/bin
echo "Finished provisioning."
