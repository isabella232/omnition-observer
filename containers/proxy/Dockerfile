FROM envoyproxy/envoy:v1.13.0

RUN apt update && apt-get install -y \
    iproute2 \
    iptables \
    gettext-base \
&& rm -rf /var/lib/apt/lists/*

COPY bin/observer /usr/local/bin/
COPY scripts/*.sh /usr/local/bin/
RUN chmod +x /usr/local/bin/*.sh

ENTRYPOINT ["/usr/local/bin/start_omnition_proxy.sh"]
CMD ["run"]
#CMD ["show-config"]
