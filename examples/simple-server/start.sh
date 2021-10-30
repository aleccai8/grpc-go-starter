#!/bin/bash
go build .
export app='examples' server="simple-server"
./simple-server