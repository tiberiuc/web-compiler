FROM phusion/baseimage:latest

RUN apt-get update
RUN apt-get -y install git
RUN curl https://godeb.s3.amazonaws.com/godeb-amd64.tar.gz > godeb.tar.gz
RUN tar -xzvf godeb.tar.gz
RUN ./godeb  install 1.4.2

VOLUMES [".:/web-compiler"]
WORKDIR /web-compiler
