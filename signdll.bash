#!/bin/bash

function die() {
  echo $*
  exit 1
}

if [ -z "$BNS_CERT" ]
then
    die "$0: Please set BNS_CERT to the bns signing certificate for windows"
fi

if [ -z "$BNS_CERT_PASS" ]
then
    die "$0: Please set BNS_CERT_PASS to the password for the $BNS_CERT signing key"
fi

which osslsigncode > /dev/null
if [ $? -ne 0 ]
then
    echo "Installing osslsigncode"
    brew install osslsigncode || die "Could not install osslsigncode"
fi
osslsigncode sign -pkcs12 "$BNS_CERT" -pass "$BNS_CERT_PASS" -in dll/systray386.dll_unsigned -out dll/systray386.dll || die "Could not sign windows 386 dll"
osslsigncode sign -pkcs12 "$BNS_CERT" -pass "$BNS_CERT_PASS" -in dll/systrayamd64.dll_unsigned -out dll/systrayamd64.dll || die "Could not sign windows 386 dll"
