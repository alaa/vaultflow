FROM golang:latest
ADD vaultflow /usr/local/bin
WORKDIR /root
ENTRYPOINT ["/usr/local/bin/vaultflow"]
