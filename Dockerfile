FROM docker.io/cirss/repro-template

COPY .repro .repro

USER repro

# install required repro modules
RUN repro.require blaze exported --dev
RUN repro.require blazegraph-service master ${CIRSS_BRANCH}

RUN repro.atstart blazegraph-service.start

CMD  /bin/bash -il

