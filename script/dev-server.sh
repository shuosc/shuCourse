#!/usr/bin/env bash
export DB_ADDRESS="postgres://postgres@localhost:5432/postgres?sslmode=disable"
export PROXY_AUTH_ADDRESS="http://cloud.shu.xn--io0a7i:30080/auth/login/proxy/course"
export PROXY_ADDRESS="http://cloud.shu.xn--io0a7i:30080/proxy/course/"
export PORT="8001"
gin -p 8000 run main.go