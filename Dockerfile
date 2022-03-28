FROM docker.io/cirss/repro-template

COPY exports /repro/exports

USER repro

# install required repro modules
RUN repro.require blaze exports --code
RUN repro.require blazegraph-service master ${CIRSS_BRANCH}

RUN repro.atstart blazegraph-service.start

CMD  /bin/bash -il

