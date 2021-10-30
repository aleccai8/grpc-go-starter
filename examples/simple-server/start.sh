#!/bin/bash
go build .
export app='one' server="user_server"
./simple-server