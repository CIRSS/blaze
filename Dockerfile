FROM docker.io/cirss/repros-template

COPY .repro .repro

USER repro

# install required repro modules
RUN repro.require blaze exported
RUN repro.require blazegraph-service 0.2.6 ${CIRSS_RELEASE}

RUN repro.prefixpath '${REPRO_MNT}/.repro/exported'

RUN repro.atstart start-blazegraph

CMD  /bin/bash -il

