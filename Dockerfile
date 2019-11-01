FROM alpine:3.10

RUN apk add -U wireguard-tools
RUN apk add -U --virtual .build-deps build-base libmnl-dev git gcc wget 

RUN cd /tmp/ \
    && git clone https://github.com/WireGuard/wg-dynamic.git \
    && cd wg-dynamic \
    && sed -i 's/install: wg/install: all/g' Makefile \
    && make \
    && make install || true

RUN apk del .build-deps

COPY dist/wgrest-linux-amd64 /usr/bin/wgrest

ENTRYPOINT ["/usr/bin/wgrest"]
CMD ["--scheme", "http", "--host", "0.0.0.0", "--port", "8000"]