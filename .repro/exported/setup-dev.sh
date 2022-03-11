sudo apt -y install gcc

GO_VERSION=1.16
GO_DOWNLOADS_URL=https://dl.google.com/go
GO_ARCHIVE=go${GO_VERSION}.linux-amd64.tar.gz

sudo wget ${GO_DOWNLOADS_URL}/${GO_ARCHIVE} -O /tmp/${GO_ARCHIVE}
sudo tar -xzf /tmp/${GO_ARCHIVE} -C /usr/local

repro.setenv GOPATH '${REPRO_MNT}/.gopath'
repro.prefixpath /usr/local/go/bin
repro.prefixpath '${GOPATH}/bin'
