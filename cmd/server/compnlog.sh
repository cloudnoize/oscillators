#!/bin/bash

GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -a -o  ./osc  ./main.go

scp -i ~/awskey/ireland_cy.pem ./osc  ubuntu@ec2-34-246-193-176.eu-west-1.compute.amazonaws.com:~

ssh  -i ~/awskey/ireland_cy.pem -v ubuntu@ec2-34-246-193-176.eu-west-1.compute.amazonaws.com
