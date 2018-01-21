FROM ubuntu:16.04

RUN apt-get update && apt-get install -y build-essential

WORKDIR /app
ADD . /app

VOLUME /app