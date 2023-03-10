FROM ubuntu:latest

RUN apt update -o Acquire::Check-Valid-Until=false -o Acquire::Check-Date=false
RUN apt install -y golang
RUN apt install ca-certificates -y
RUN update-ca-certificates
ADD . .
ENTRYPOINT go run .