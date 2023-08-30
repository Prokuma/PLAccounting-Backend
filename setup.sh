#!/bin/sh
# This script is used to setup the environment for the project
git clone https://github.com/Prokuma/PLAccounting-WebFrontend.git
cp .env.docker PLAccounting-WebFrontend/.env
cd PLAccounting-WebFrontend
git checkout tags/alpha-v0.1.0
npm install
npm run build
mv ./build ../html
cd ..
rm -rf PLAccounting-WebFrontend
if [ ! -d "jwt_keys" ]; then
    mkdir jwt_keys
fi
ssh-keygen -m PKCS8 -t rsa -b 2048 -N "" -f ./jwt_keys/jwt_rsa
ssh-keygen -i -f ./jwt_keys/jwt_rsa.pub > ./jwt_keys/jwt_rsa.pub.pkcs8