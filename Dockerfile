FROM cirss/repro-parent:latest

COPY exports /repro/exports

ADD ${REPRO_DIST}/setup /repro/dist/
RUN bash /repro/dist/setup

USER repro

# install required repro modules
RUN repro.require repro master ${REPROS_DEV}
RUN repro.require blaze exports --code
RUN repro.require blazegraph-service master ${CIRSS}

RUN repro.atstart blazegraph-service.start

CMD  /bin/bash -il
