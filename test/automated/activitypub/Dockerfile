FROM golang:1.26-bookworm

RUN apt-get update && apt-get install -y \
    gcc make \
    curl libcurl4-openssl-dev libssl-dev \
    jq sqlite3 ffmpeg \
    fonts-dejavu-core \
    libnss3-tools \
    lsof procps \
    && rm -rf /var/lib/apt/lists/*

# Install mkcert (architecture-aware)
RUN ARCH="$(dpkg --print-architecture)" && \
    curl -sL "https://github.com/FiloSottile/mkcert/releases/download/v1.4.4/mkcert-v1.4.4-linux-${ARCH}" -o /usr/local/bin/mkcert && \
    chmod +x /usr/local/bin/mkcert

# Install Caddy (architecture-aware)
RUN ARCH="$(dpkg --print-architecture)" && \
    curl -sL "https://github.com/caddyserver/caddy/releases/download/v2.8.4/caddy_2.8.4_linux_${ARCH}.tar.gz" | tar -xz -C /usr/local/bin caddy

# Build and install snac2
RUN git clone --depth 1 https://codeberg.org/grunfink/snac2.git /tmp/snac2-src \
    && cd /tmp/snac2-src && make && cp snac /usr/local/bin/snac \
    && rm -rf /tmp/snac2-src

# Install mkcert CA into system trust store so snac2 trusts the test certs
RUN mkcert -install

# Allow git operations on the mounted repo
RUN git config --global --add safe.directory /owncast

WORKDIR /owncast

COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]
