FROM golang:1-bullseye

RUN apt-get update && apt-get install -y ruby=1:2.7*

RUN mkdir -p /code /godir /root/data
VOLUME ["/code", "/godir"]

COPY dotfiles/* /root/data/
COPY scripts/copy_dotfiles.rb /root/
RUN ruby /root/copy_dotfiles.rb

WORKDIR /code
ENV GOPATH=/godir
