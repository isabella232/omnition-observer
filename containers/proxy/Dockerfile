FROM envoyproxy/envoy:v1.13.1

RUN apt update && apt-get install -y \
    iproute2 \
    iptables \
    gettext-base \
    curl \
&& rm -rf /var/lib/apt/lists/*

RUN curl -Lo /usr/local/lib/libjaegertracing_plugin.so https://github.com/jaegertracing/jaeger-client-cpp/releases/download/v0.4.2/libjaegertracing_plugin.linux_amd64.so

COPY bin/observer /usr/local/bin/
COPY scripts/*.sh /usr/local/bin/
RUN chmod +x /usr/local/bin/*.sh

ENTRYPOINT ["/usr/local/bin/start_omnition_proxy.sh"]
CMD ["run"]
#CMD ["show-config"]
