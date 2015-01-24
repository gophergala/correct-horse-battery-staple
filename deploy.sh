#!/bin/bash

echo "Building..."
GOOS=linux go build

echo "Stopping service"
ssh root@where-you-at.com "service whereyouat stop"

echo "Uploading..."
scp main root@where-you-at.com:/home/whereyouat/
scp -r assets/. root@where-you-at.com:/home/whereyouat/assets/
ssh root@where-you-at.com "chown -R whereyouat:whereyouat /home/whereyouat/"

echo "Starting service"
ssh root@where-you-at.com "service whereyouat start"

