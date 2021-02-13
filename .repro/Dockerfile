FROM ubuntu:20.04

ENV REPRO_NAME  blazegraph-cli
ENV REPRO_MNT   /mnt/${REPRO_NAME}
ENV REPRO_USER  repro
ENV REPRO_UID   1000
ENV REPRO_GID   1000

RUN echo '***** Update packages *****'                                      \
    && apt-get -y update

RUN echo '***** Set timezone noninteractively install JDK8 *****'
RUN DEBIAN_FRONTEND="noninteractive" TZ="America/Los_Angeles"               \
    apt -y install tzdata openjdk-8-jdk

RUN echo '***** Install packages required for creating this image *****'    \
    && apt -y install apt-utils wget curl makepasswd make git               \
                                                                            \
    && echo '***** Install command-line utility packages *****'             \
    && apt -y install sudo less file tree

RUN echo '***** Add the REPRO user and group *****'                         \
    && groupadd ${REPRO_USER} --gid ${REPRO_GID}                            \
    && useradd ${REPRO_USER} --uid ${REPRO_UID} --gid ${REPRO_GID}          \
        --shell /bin/bash                                                   \
        --create-home                                                       \
        -p `echo repro | makepasswd --crypt-md5 --clearfrom - | cut -b8-`   \
    && echo "${REPRO_USER} ALL=(ALL) NOPASSWD: ALL"                         \
            > /etc/sudoers.d/${REPRO_USER}                                  \
    && chmod 0440 /etc/sudoers.d/repro

ENV HOME /home/${REPRO_USER}
ENV HOME_BIN ${HOME}/bin
ENV PATH ${HOME_BIN}:${PATH}
ENV BASHRC ${HOME}/.bashrc
USER  ${REPRO_USER}
WORKDIR $HOME

ENV BLAZEGRAPH_VER RELEASE_2_1_5
ENV BLAZEGRAPH_RELEASES https://github.com/blazegraph/database/releases
ENV BLAZEGRAPH_DOWNLOAD_DIR ${BLAZEGRAPH_RELEASES}/download/BLAZEGRAPH_${BLAZEGRAPH_VER}
ENV BLAZEGRAPH_DOWNLOAD_JAR ${BLAZEGRAPH_DOWNLOAD_DIR}/blazegraph.jar
ENV BLAZEGRAPH_JAR $HOME/blazegraph-${BLAZEGRAPH_VER}.jar
RUN echo '***** Download Blazegraph jar *****'                              \
    && wget -O ${BLAZEGRAPH_JAR} ${BLAZEGRAPH_DOWNLOAD_JAR}

ENV BLAZEGRAPH_DOT_DIR ${REPRO_MNT}/.blazegraph
ENV BLAZEGRAPH_PROPERTY_FILE=${BLAZEGRAPH_DOT_DIR}/.properties
ENV BLAZEGRAPH_OPTIONS "-server -Xmx4g -Dbigdata.propertyFile=${BLAZEGRAPH_PROPERTY_FILE}"
ENV BLAZEGRAPH_CMD "java ${BLAZEGRAPH_OPTIONS} -jar ${BLAZEGRAPH_JAR}"
ENV BLAZEGRAPH_LOG ${BLAZEGRAPH_DOT_DIR}/blazegraph_`date +%s`.log

ENV GEIST_RELEASES https://github.com/CIRSS/geist/releases/download
ENV GEIST_VERSION 0.2.0
ENV BLAZEGRAPH_DOWNLOAD_BINARY blazegraph-cli-${GEIST_VERSION}-linux-x86-64
ENV BLAZEGRAPH_CLI_COMMAND ${HOME_BIN}/blazegraph

RUN echo '***** Install the Geist blazegraph command *****'                 \
    && mkdir -p ${HOME_BIN}                                                 \
    && wget -O ${BLAZEGRAPH_CLI_COMMAND}                                    \
        ${GEIST_RELEASES}/v${GEIST_VERSION}/${BLAZEGRAPH_DOWNLOAD_BINARY}                    \
    && chmod a+rx ${BLAZEGRAPH_CLI_COMMAND}

RUN echo '***** Start Blazegraph at login *****'                            \
    && echo "source .bashrc" >> .bash_profile                               \
    && echo "cd ${BLAZEGRAPH_DOT_DIR} && ${BLAZEGRAPH_CMD} 2>&1 > ${BLAZEGRAPH_LOG} &" >> .bash_profile

RUN echo "export IN_RUNNING_REPRO=${REPRO_NAME}" >> ${BASHRC}
RUN echo "cd ${REPRO_MNT}" >> ${BASHRC}

CMD  /bin/bash -il
