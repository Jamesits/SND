#!/usr/bin/env bash
set -Eeuo pipefail

cd "$( dirname "${BASH_SOURCE[0]}" )"/../..

DNSPERF_DATA_FILE="build/dnsperf.txt"

echo "" > "$DNSPERF_DATA_FILE"

# useful requests
# it takes a while to generate them -- be patient with bash
for i in `seq 1 20000`; do
  LOW_64_BITS=$(printf '%016x\n' $i | rev | fold -b1 | paste -sd'.' -)
  echo "$LOW_64_BITS.0.0.0.0.0.0.0.0.1.0.0.0.0.0.d.f.ip6.arpa. PTR" >> "$DNSPERF_DATA_FILE"
done

# errornous requests
# https://www.techietown.info/2017/03/load-testing-dns-using-dnsperf/
for i in `seq 1 2000000`; do
  echo "$i.example.com A" >> "$DNSPERF_DATA_FILE"
done

# download and compile https://github.com/cobblau/dnsperf
# need around 5mins for 5000000 requests
dnsperf -s 127.0.0.1 -d build/dnsperf.txt -c 1000 -Q 5000000
