FROM envoyproxy/envoy:v1.7.0

RUN apt update && apt-get install -y \
    iproute2 \
    iptables \
    gettext-base \
&& rm -rf /var/lib/apt/lists/*

ADD scripts/*.sh /usr/local/bin/
RUN chmod +x /usr/local/bin/*.sh

ADD envoy_tmpl.yaml /etc/envoy_tmpl.yaml

EXPOSE 15001
EXPOSE 9901

ENTRYPOINT ["/usr/local/bin/start_omnition_proxy.sh"]
CMD ["run"]
