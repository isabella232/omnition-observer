FROM envoyproxy/envoy:latest

RUN apt update && apt-get install -y \
    iproute2 \
    iptables \
    gettext-base \
&& rm -rf /var/lib/apt/lists/*

ADD bin/envoy /usr/local/bin/
ADD bin/observer /usr/local/bin/
ADD scripts/*.sh /usr/local/bin/
RUN chmod +x /usr/local/bin/*.sh

# ADD envoy_tmpl.yaml /etc/envoy_tmpl.yaml

ENTRYPOINT ["/usr/local/bin/start_omnition_proxy.sh"]
CMD ["run"]
