#!/bin/bash - 
#===============================================================================
#
#          FILE: benchmark.sh
#
#         USAGE: ./benchmark.sh
#
#   DESCRIPTION: Runs a benchmark against the kubenab-Server
#
#       OPTIONS: ---
#  REQUIREMENTS: go, docker, openssl
#          BUGS: ---
#         NOTES: ---
#        AUTHOR: Francesco Emanuel Bennici <benniciemanuel78@gmail.com>
#  ORGANIZATION:
#       CREATED: 28.09.2019 00:32:08
#      REVISION:  ---
#===============================================================================

set -o nounset                              # Treat unset variables as an error

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
tmp=$(mktemp -d)
curr=$(pwd)

## ===> Compile `kubenab` <===

echo "[i] Generating Self-Signed Certificates"
openssl req -x509 -newkey rsa:4096 -sha256 -days 3650 -nodes \
  -keyout ${tmp}/tls.key -out ${tmp}/tls.crt -extensions san -config \
  <(echo "[req]"; 
    echo distinguished_name=req; 
    echo "[san]"; 
    echo subjectAltName=DNS:localhost,IP:127.0.0.1
    ) \
   -subj "/CN=localhost"


echo "[i] Compiling kubenab"
cd ${DIR}/../
docker build -t temp/build:kubenab .
id=$(docker run -p 8443:443 \
  -v ${tmp}:/etc/admission-controller/tls \
  -d --env "DOCKER_REGISTRY_URL=jfrog" --env "REPLACE_REGISTRY_URL=false" \
  temp/build:kubenab)
cd ${curr}

## ==> Benchmark <==#

echo "[i] Installing bombardier"
go get -u github.com/codesenberg/bombardier

echo "###########################"
echo "### Starting bombardier ###"
echo -e "###########################\n"

echo -e "==> Mutate Webhook\n\n"
bombardier -c 125 -n 10000000 --insecure --latencies \
  --fasthttp --body $(cat ${DIR}/mutate_body) \
  --print 'i,p,r' --method POST https://localhost:8443/mutate

echo -e "\n\n==> Validate Webhook\n\n"
bombardier -c 125 -n 10000000 --insecure --latencies \
  --fasthttp --body $(cat ${DIR}/validate_body) \
  --print 'i,p,r' --method POST https://localhost:8443/validate


## Cleaning Up
docker rm --force --volumes ${id}
rm -rf ${tmp}
