#!/usr/bin/env bash
export DB_ADDRESS="postgres://localhost:5432/postgres?sslmode=disable"
export PROXY_AUTH_ADDRESS="https://cloud.shuosc.com/auth/login/shu-course-proxy"
export PROXY_ADDRESS="https://cloud.shuosc.com/api/shu-course-proxy/"
export PORT="8001"
gin -p 8000 run web/main.go
