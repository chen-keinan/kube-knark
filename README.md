[![Go Report Card](https://goreportcard.com/badge/github.com/chen-keinan/kube-knark)](https://goreportcard.com/report/github.com/chen-keinan/kube-knark)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://github.com/chen-keinan/beacon/blob/main/LICENSE)
[![Build Status](https://travis-ci.com/chen-keinan/kube-knark.svg?branch=master)](https://travis-ci.com/chen-keinan/kube-knark)
[![Coverage Status](https://coveralls.io/repos/github/chen-keinan/kube-knark/badge.svg?branch=master)](https://coveralls.io/github/chen-keinan/kube-knark?branch=master)
[![Gitter](https://badges.gitter.im/kube-knark/community.svg)](https://gitter.im/kube-knark/community?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge)
<br><img src="./pkg/images/kube_krank.png" width="300" alt="kube-krank logo"><br>
# Kube-Knark Project
###  Trace your kubernetes runtime !!
Kube-Knark is an open source audit tracer (via ebpf technology) who perform runtime  tracing on a deployed kubernetes cluster looking for malware code execution.

### Requirements
- Go 1.10+
- Linux Kernel 4.15+
- Clang < 10
- LLVM
- Kernel Headers
