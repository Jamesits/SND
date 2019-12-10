# SND

Minimal authoritative PTR (rDNS, reverse DNS) resolver with automatic generation of records.

**WARNING**: This piece of software is at a very early development stage and should be considered experimental. It only implements a minimal feature set to pass the RIPE automated DNS test and is by no means standard-compliant. Please don't run it in production. 

[![Build Status](https://dev.azure.com/nekomimiswitch/General/_apis/build/status/SND?branchName=master)](https://dev.azure.com/nekomimiswitch/General/_build/latest?definitionId=71&branchName=master)
[![](https://images.microbadger.com/badges/image/jamesits/snd.svg)](https://microbadger.com/images/jamesits/snd "Get your own image badge on microbadger.com")

## Motivation

Say you have a large collection of IP addresses (thousands of IPv4 /24 blocks, or one IPv6 /32 block), and you want to have PTR records on all of your IPs. Writing a zonefile and hosting it using any traditional authoritative DNS server is unrealistic: the zonefile will be of multiple GBs and you need an enormous amount of memory to even load it. 

SND provides you a simple alternative option: you name a base domain, and SND generates PTR records for you on the fly based on a set of pre-defined rules. 

```
1.1.168.192.in-addr.arpa.                                                 1000 IN PTR 192.168.1.1.ptr.example.com.
1.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.1.0.0.0.0.0.d.f.ip6.arpa. 1000 IN PTR fd00.1.0.0.0.0.0.1.ptr.example.com.
```

## Requirements

### Hardware

SND can run on very little processing power (Raspberry Pis are fine) and a very small memory footprint (a few MBs) although the performance will be not optimal.

As of version 0.1.2, on a 4-core Intel E5-2670 VM with more than 2GiB memory, SND can process around 25K RPS.

### Software

Officially supported OS: 

* Linux (kernel 4.19+ with glibc)
* Windows (Windows Server 2016 or later, Windows 10 Desktop 1809 or later)

Other OSes are not currently tested because of the lack of resources available to me.

## Usage

### Configure SND

Copy over the self-documented [example config](examples/config.toml) and tweak it for your own need. Please do not leave any `example.com` things in your own config. Remove what you don't need. 

Currently no strict config file format checking is implemented -- you might crash the program if some important keys are missing. 

### Set up SND

In most cases you are going to need 2 servers (or one server with 2 different IP addresses if you don't care about availability issues). Copy the exact same config file to both servers and launch SND on both of them:

Download the pre-compiled binary from [releases](https://github.com/Jamesits/SND/releases) to your server and run it:

```shell
chmod +x ./snd
./snd -config path/to/config.toml
```

Or, if you prefer Docker:

```shell
docker run --rm -p 53:53 -p 53:53/udp -v path/to/config.toml:/etc/snd/config.toml:ro snd:latest
```

Run a simple test using dig:

```shell
$ dig @localhost -x 192.0.2.1

; <<>> DiG 9.11.5-P4-5.1-Debian <<>> @localhost -x 192.0.2.1
; (1 server found)
;; global options: +cmd
;; Got answer:
;; ->>HEADER<<- opcode: QUERY, status: NOERROR, id: 50924
;; flags: qr aa rd; QUERY: 1, ANSWER: 1, AUTHORITY: 0, ADDITIONAL: 0
;; WARNING: recursion requested but not available

;; QUESTION SECTION:
;1.2.0.192.in-addr.arpa.      IN      PTR

;; ANSWER SECTION:
1.2.0.192.in-addr.arpa. 3600  IN      PTR     192.0.2.1.example.com.
```

### Set up DNS records

You need at least 2 `A` or `AAAA` records pointing to each of your SND servers. You might need to set them up as glue records based on your actual config.

```
ns1.example.com.	3600	IN	A	192.0.2.1
ns2.example.com.	3600	IN	A	192.0.2.2
```

### Set up PTR record delegation

Set up a `domain` object at your RIR like this. 

```
domain:   <zone name>
descr:    <description>
admin-c:  <nic-handle for administrative contact>
tech-c:   <nic-handle for technical contact>
zone-c:   <nic-handle for zone contact>
nserver:  ns1.example.com
nserver:  ns2.example.com
mnt-by:   <your maintainer>
source:   RIPE
```

Detailed instructions are provided per RIR:

* [AfriNIC](https://afrinic.net/support/requesting-reverse-delegation)
* [ARIN](https://www.arin.net/resources/manage/reverse/)
* [APNIC](https://www.apnic.net/manage-ip/manage-resources/reverse-dns/)
* [LACNIC](https://www.lacnic.net/685/2/lacnic/5-delegation-of-reverse-resolution)
* [RIPE NCC](https://www.ripe.net/manage-ips-and-asns/db/support/configuring-reverse-dns)

Notes:

* The smallest IP block sizes available for delegation differ
* Only RIPE NCC is currently tested because I cannot afford IP blocks from the other RIRs

## Compilation

Golang 1.13.5 or later is officially supported. Before starting, make sure the `GOROOT` and `GOPATH` environment variables are set correctly and there is a `go` binary is in your `PATH`.

### Linux

```shell
git clone https://github.com/Jamesits/SND.git
cd SND
./build.sh
```

Collect the binary in the `build` directory.

### Other OSes

#### Cross Compilation from Linux

Set the `GOOS` and `GOARCH` to the desired [values](https://gist.github.com/asukakenji/f15ba7e588ac42795f421b48b8aede63) and run `build.sh`.

#### Native Build

You need to figure out how to build on these platforms on yourself if you want a native build.