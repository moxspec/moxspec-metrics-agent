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
                                Prometheus Remote Write
                                              ┌────────────┐
                                              │            │
                  Metrics                 ┌──►│ Backend DB │
┌─────────────────┐    ┌───────────────┐  │   │            │
│                 │    │               │  │   └────────────┘
│ Hardware/Kernel ├───►│ Metrics Agent ├──┤
│                 │    │               │  │   ┌────────────┐
└─────────────────┘    └───────────────┘  │   │            │
                                          └──►│ Backend DB │
                                              │            │
                                              └────────────┘
                                  JSON over HTTP POST
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

# Metrics

```
- Processor
  - Physical core
    - cpu_throttlecount
    - cpu_scaling_cur_freq
    - cpu_scaling_max_freq
    - cpu_scaling_min_freq
    - cpu_temp
  - Power consumption by socket
    - rapl_package_energy_uj
    - rapl_package_power_limit
    - rapl_package_max_power
    - rapl_package_time_window
- Memory
  - ECC
    - edac_mc_size
    - edac_mc_ce_count
    - edac_mc_ce_noinfo_count
    - edac_mc_ue_count
    - edac_mc_ue_noinfo_count
    - edac_csrow_size
    - edac_csrow_ce_count
    - edac_csrow_ue_count
    - edac_ch_ce_count
- Storage
  - SATA/SAS drive
    - disk_cur_temp
    - disk_max_temp
    - disk_min_temp
    - disk_rotation
    - disk_byte_read
    - disk_byte_written
    - disk_neg_speed
    - disk_sig_speed
    - disk_power_cycle_count
    - disk_power_on_hours
    - disk_unsafe_shutdown_count
  - NVMe drive
    - disk_cur_temp
    - disk_warn_temp
    - disk_crit_temp
    - disk_byte_read
    - disk_byte_written
    - disk_size
    - disk_power_cycle_count
    - disk_power_on_hours
    - disk_unsafe_shutdown_count
    - drive_namespace_size
    - drive_namespace_phy_block_size
    - drive_namespace_log_block_size
  - RAID
    - Broadcom MegaRAID
      - SAS/SATA physical drive
        - disk_size
        - disk_error_count
        - disk_byte_written
        - disk_byte_read
        - disk_power_cycle_count
        - disk_power_on_hours
        - disk_unsafe_shutdown_count
    - Broadcom Fusion MPT (WIP)
    - HP Smart Array / Smart HBA (WIP)
- Network Interface
  - nw_rx_packets
  - nw_tx_packets
  - nw_rx_bytes
  - nw_tx_bytes
  - nw_rx_errors
  - nw_tx_errors
  - nw_rx_dropped
  - nw_tx_dropped
```
