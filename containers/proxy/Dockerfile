FROM envoyproxy/envoy:latest

RUN apt update && apt-get install -y \
    iproute2 \
    iptables \
    gettext-base \
&& rm -rf /var/lib/apt/lists/*


COPY bin/envoy /usr/local/bin/
COPY bin/observer /usr/local/bin/

COPY scripts/*.sh /usr/local/bin/
RUN chmod +x /usr/local/bin/*.sh

COPY envoy.yaml /etc/envoy.yaml

ENTRYPOINT ["/usr/local/bin/start_omnition_proxy.sh"]
CMD ["run"]
