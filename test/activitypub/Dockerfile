FROM alpine:3.20 AS builder

RUN apk add openssl-dev curl-dev git make gcc openssl libc-dev

RUN git clone "https://codeberg.org/grunfink/snac2"
WORKDIR ./snac2
RUN make

FROM alpine:3.20

COPY --from=builder /snac2/snac /usr/local/bin/snac

RUN apk add libssl3 libcurl openssh jq

ENV SNAC_DATA_PATH="/var/local/snac-data"
COPY --chmod=744 ./*.sh /usr/local/bin/
ENTRYPOINT ["/usr/local/bin/run.sh"]

