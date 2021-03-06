[![Go Report Card](https://goreportcard.com/badge/github.com/chen-keinan/kube-knark)](https://goreportcard.com/report/github.com/chen-keinan/kube-knark)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://github.com/chen-keinan/beacon/blob/main/LICENSE)
[![Build Status](https://travis-ci.com/chen-keinan/kube-knark.svg?branch=master)](https://travis-ci.com/chen-keinan/kube-knark)
<img src="./pkg/images/coverage_badge.png" alt="test coverage badge">
[![Gitter](https://badges.gitter.im/kube-knark/community.svg)](https://gitter.im/kube-knark/community?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge)
<br><img src="./pkg/images/kube_krank.png" width="300" alt="kube-krank logo"><br>
# Kube-Knark Project
###  Trace your kubernetes runtime !!
Kube-Knark is an open-source tracer that uses pcap & ebpf technology to perform runtime tracing on a deployed Kubernetes cluster. It is tracing the Kubernetes API execution and master node configuration files permission changes. The trace matching events are leveraged via go plugin webhooks

### kube-knark trace the following :
- The full [Kubernetes API specification](https://kubernetes.io/docs/reference/kubernetes-api/) execution calls especially mutation <br> 
- kubernetes master node configuration files permission changes [CIS Kubernetes Benchmark specification](https://www.cisecurity.org/benchmark/kubernetes/)

kube-knark tracing data are reported :
- Console dashboard
- Go Plugin hooks

kube-Knark console:
<br><img src="./pkg/images/knark-console.png" width="800" alt="kube-krank-console logo"><br>

* [Requirements](#requirements)
* [Installation](#installation)
* [Quick Start](#quick-start)
* [User Plugin Usage](#user-plugin-usage)
* [Supported Specs](#supported-specs)
* [Contribution](#Contribution)

## Requirements
- Go 1.13+
- Linux Kernel 4.15+
- Clang 10+
- LLVM
- Kernel Headers
- Pcap

## Installation
```shell
git clone https://github.com/chen-keinan/kube-knark
cd kube-knark
make build
```
## Quick Start
Execute kube-knark without plugins 
```shell
 ./kube-knark 
```
## User Plugin Usage (via go plugins)
The Kube-knark expose 2 hooks for user plugins [Example](https://github.com/chen-keinan/kube-knark/tree/master/examples/plugins) :
- **OnK8sAPICallHook** - this hook accepts k8s api call event with all details (http request /response ,matching API spec) 
- **OnK8sFileConfigChangeHook** - this hook accepts master file configuration change event with command details (chown or chmod ,args and matching file change spec)

##### Compile user plugin
```shell
go build -buildmode=plugin -o=~/<plugin folder>/<plugin>.so ~/<plugin folder>/<plugin>.go
```
##### Copy plugin to folder (.kube-knark folder is created on the 1st startup)
```shell
cp ~/<plugin folder>/<plugin>.so ~/.kube-knark/plugins/compile/<plugin>.so
```
## Supported Specs 
The Kube-knark support 2 specs and can be easily extended:
- The full k8s API spec [Kubernetes API specification](https://github.com/chen-keinan/kube-knark/tree/master/internal/spec/api)
- master config file change spec [Master Node Config](https://github.com/chen-keinan/kube-knark/tree/master/internal/spec/filesystem)

both specs can be easily extended by amended the spec files under ```~/.kube-knark/spec``` folder

## Contribution
- code contribution is welcome !! , contribution with tests and passing linter is more than welcome :) 
- /.dev folder include vagrantfile to be used for development : [Dev Instruction](https://github.com/chen-keinan/kube-knark/tree/master/.dev)
