FROM ubuntu:xenial

RUN apt update && apt-get install -y \
    iproute2 \
    iptables \
&& rm -rf /var/lib/apt/lists/*

ADD scripts/setup_omnition.sh /usr/local/bin/
RUN chmod +x /usr/local/bin/setup_omnition.sh

ENTRYPOINT ["/usr/local/bin/setup_omnition.sh"]
