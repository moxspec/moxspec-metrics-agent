moxspec-metrics-agent
===

[![CircleCI](https://circleci.com/gh/moxspec/moxspec-metrics-agent.svg?style=shield&circle-token=951127200cfb6aaa6c85939e7344aba39b888bb7)](https://circleci.com/gh/moxspec/moxspec-metrics-agent)
[![Maintainability](https://api.codeclimate.com/v1/badges/aedbc4adcc2ffe91cea1/maintainability)](https://codeclimate.com/github/moxspec/moxspec-metrics-agent/maintainability)
[![Test Coverage](https://api.codeclimate.com/v1/badges/aedbc4adcc2ffe91cea1/test_coverage)](https://codeclimate.com/github/moxspec/moxspec-metrics-agent/test_coverage)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

***Early Prototype***

A hardware and kernel metrics collector. OS independent, easy to deploy.   
Any DB that implements Prometheus remote write interface can be used as a backend.

```
                  Metrics
┌─────────────────┐    ┌───────────────┐    ┌────────────┐
│                 │    │               │    │            │
│ Hardware/Kernel ├───►│ Metrics Agent ├───►│ Backend DB │
│                 │    │               │    │            │
└─────────────────┘    └───────────────┘    └────────────┘
                              Prometheus Remote Write
```

# Prerequisites

- x86_64 or ARMv8 (experimental)
- Intel(Westmere or later) or AMD (Zen or later) processor
- PCI Rev 3.0
- PCI Express Rev 3.0+
- SMBIOS v2.4+
- Linux kernel 2.6.32+

# Installation

Requirements: Docker
```
$ git clone https://github.com/moxspec/moxspec-metrics-agent.git
$ cd moxspec-metrics-agent
$ make bin
$ sudo bin/mox-metrics-agent
```

# Quick start

```
% bin/mox-metrics-agent -h
Usage of bin/mox-metrics-agent:
  -d    enable debug logging
  -e string
        a remote-write endpoint (default "http://localhost:3030/remote/write")
  -o string
        output selection (promRemote, jsonHttp, stdout) (default "promRemote")
  -r int
        default retrieval interval (sec) (default 5)
  -s int
        default submission interval (sec) (default 10)
```

# Example

<img alt="metrics-sample" src="https://user-images.githubusercontent.com/44848317/111316332-6797a900-8620-11eb-804e-758917d7fac2.png">


