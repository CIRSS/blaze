FROM cirss/repro-parent:latest

COPY exports /repro/exports

ADD ${REPRO_DIST}/boot-setup /repro/dist/
RUN bash /repro/dist/boot-setup

USER repro

# install required repro modules
RUN repro.require blaze exports --code --demo
RUN repro.require blazegraph-service master ${CIRSS}

CMD  /bin/bash -il
