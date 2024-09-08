
FROM ubuntu:24.04

ENV DEBIAN_FRONTEND=noninteractive
ENV NVM_DIR=/root/.nvm
ENV NODE_VERSION=20.16.0
ENV CHIA_ROOT=/root/.chia/testnet
ENV PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin

RUN apt-get update && apt-get install -y \
    git \
    curl \
    sudo \
    python3-venv \
    build-essential \
    coreutils \
    procps \
    grep \
    sed \
    findutils \
    gnupg \
    ca-certificates \
    bsdmainutils \
    gawk \
    netcat-openbsd \
    lsof \
    psmisc \
    dos2unix \
    apt-utils \
    less \
    file \
    man-db \
    && apt-get clean && \
    apt-get update && apt-get upgrade -y

RUN command -v curl && command -v apt-get && command -v less && command -v groups

RUN echo "PATH is set to: $PATH"

COPY ./chia-blockchain /chia-blockchain

RUN dos2unix /chia-blockchain/install.sh

RUN cd /chia-blockchain && \
    bash install.sh && \
    /bin/bash -c "source ./activate && chia init && chia init --fix-ssl-permissions"

RUN curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.35.3/install.sh | bash && \
    . $NVM_DIR/nvm.sh && nvm install $NODE_VERSION && nvm alias default $NODE_VERSION && nvm use default && \
    echo "Node version is now: $(node -v)"

RUN git clone -b main https://github.com/Chia-Network/cadt.git && \
    cd cadt && \
    /bin/bash -c ". $NVM_DIR/nvm.sh && nvm use $NODE_VERSION && npm install"

COPY modify_config.sh /modify_config.sh
RUN chmod +x /modify_config.sh

COPY initialize_chia.sh /initialize_chia.sh
RUN chmod +x /initialize_chia.sh

COPY initialize_cadt.sh /initialize_cadt.sh
RUN chmod +x /initialize_cadt.sh

CMD ["/bin/bash", "-c", "echo 'Starting Chia and CADT...' && /initialize_chia.sh && /initialize_cadt.sh"]
