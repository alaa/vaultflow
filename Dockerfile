FROM golang:latest
ADD bin/vaultflow /usr/local/bin
WORKDIR /root
ENTRYPOINT ["/usr/local/bin/vaultflow"]
