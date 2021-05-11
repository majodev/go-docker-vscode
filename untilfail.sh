#!/bin/bash
# $@ > error.log
while [ $? -eq 0 ]; do
    echo "$(date) Trying..."
    $@ &> /tmp/error.log
done

cat /tmp/error.log
cp -f /tmp/error.log error.log