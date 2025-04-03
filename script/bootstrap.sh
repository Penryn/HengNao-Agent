#!/bin/bash
CURDIR=$(cd $(dirname $0); pwd)
BinaryName=meeting_agent
echo "$CURDIR/bin/${BinaryName}"
exec $CURDIR/bin/${BinaryName}