FROM docker.io/cirss/repro-template

COPY .repro .repro

USER root

ENV GO_VERSION       1.16
ENV GO_DOWNLOADS_URL https://dl.google.com/go
ENV GO_ARCHIVE       go${GO_VERSION}.linux-amd64.tar.gz

RUN echo '****** Install Go development tools *****'                        \
    && wget ${GO_DOWNLOADS_URL}/${GO_ARCHIVE} -O /tmp/${GO_ARCHIVE}         \
    && tar -xzf /tmp/${GO_ARCHIVE} -C /usr/local

USER repro

# URLs for packages delivered as CIRSS GitHub releases
ENV CIRSS_RELEASES 'https://github.com/cirss/${1}/releases/download/v${2}/'

# install required repro modules
RUN repro.require blazegraph-service 0.2.6 ${CIRSS_RELEASES}
RUN repro.require blaze exported ${CIRSS_RELEASES}

RUN repro.setenv GOPATH '${REPRO_MNT}/.gopath'

RUN repro.prefixpath /usr/local/go/bin
RUN repro.prefixpath '${GOPATH}/bin'
RUN repro.prefixpath '${REPRO_MNT}/.repro/exported'

RUN repro.atstart start-blazegraph

CMD  /bin/bash -il

