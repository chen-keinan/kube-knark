#!/bin/bash

echo "Provisioning virtual machine..."

### install clang llvm make
echo "install clang llvm make golang"
sudo apt-get update -y
sudo apt-get -y install clang llvm make golang

### Install dlv pkg
echo "Install dlv pkg"
git clone https://github.com/go-delve/delve.git $GOPATH/src/github.com/go-delve/delve
cd $GOPATH/src/github.com/go-delve/delve
make install

### export dlv bin path
echo export PATH=$PATH:/home/vagrant/go/bin >> ~/.bashrc
echo export PATH=$PATH:/root/go/bin >> ~/.bashrc
echo "Finished provisioning."
