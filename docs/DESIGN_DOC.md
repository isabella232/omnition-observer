# Overview

## Problem Statement

While microservice architecture solves a lot of problems, it also comes with a whole set of new ones. Microservices are much harder to monitor and fix when things go wrong. Root cause identification is a major pain as is figuring out the right person/team to alert. Traditionally, teams have had ways to monitor and troubleshoot their own applications, but in a microservice architecture, the distributed nature leads to issues between microservices and makes it really hard to make sense of things. Instrumenting such systems provides a rich set a data that can be used to answer a lot of questions. We believe such data can power a daily use tool that can provide a lot of insight to teams, automate and aid root causes identification, and dramatically reduce the mean time to repair. 

Collecting and analyzing traces is still not very common and is one of the last things teams look into when it comes to observability. In addition to setting up things like Jaeger or Zipkin, teams also need to instrument their codebases in order to generate and collect traces. This can be quite an expensive process and can take weeks or even months to complete. 

To this end, we're building a very lightweight layer-7 proxy that can be deployed as a sidecar/companion to all services. The proxy will automatically generate traces (and other information) with no or very little modifications to the application code.


## Why sidecars or Out of process architecture?

Microservices also result in a lot of duplicated effort as teams often need to implement the same features across services, especially things like monitoring and integration with the cloud platform. This makes it very expensive to maintain such integrations especially when different services are written in different languages and maintained by different teams. There is a ton of work that every single service needs to handle on it's own. Things like metrics, logs, traces, networkig, TLS, authentication and authorization, service discovery, retries and exponential backoffs, circuit-break, etc. Most of this work implies tight integration with the underlying cloud platform, scheduling system, service discovery system, metrics system, etc. This results in coupling of application code with the platform it runs on and duplication of work across services, languages and teams. 

All of this work can be extracted out of the services into separate companion processes or containers more commonly known as sidecars. Each service can delegate any or all of these tasks to it's sidecar and only focus on the business logic itself. Any of the above mentioned components can be removed, replaced or modified across an entire cluster without having to make any changes to the application code. The sidecar pattern makes it easier define ownsership boundaries, allows both dev and ops to move faster and enables many times faster and safer deployment of things like distributed tracing. 


## Goals

* Automatic trace generation compatible with Jaeger and Zipkin.
* Easy installation with zero or minimal configuration.
* Low latency impact on traffic, and low resource consumption.
* Tracing support for most common protocols (HTTP, GRPC, HTTP/2).
* Support micro-services and monoliths in both contianerized and ochestrated workloads.
* Compatibility with existing ecosystem such as CNI plugin for kubernetes.

## Non-Goals

* This document does not focus on the trace collector component. Details about collector are tracked here: https://cloudnativelabs.atlassian.net/wiki/spaces/DEV/pages/43712522/Omnition+Collector
* Collect logs and other metrics from the services.
* Support for non UDP, ICMP protocols in first version. 
* Control-plane for proxy mesh. We aren't targetting remote, dynamic configuration management through a control-plane.
* Any other features usually found in service meshes and similar systems.


## Priorities

### P0

* Generate and traces in Zipkin and/or Jaeger comptaible formats.
* Support trace generation for HTTP, HTTP/2 and GRPC.
* Silent pass-through for unsupported protocols such Thrift, CQL, Kafka, etc.
* Support kubernetes.
* Support non-containerized ochestrated workloads.
* Support TLS between services.
* Provides health checks and basic monitoring for itself.
* Zero or very little configuration installation in all environments.
* Should work with network providers such as Calico and Flannel.
* Benchmark and understand resource consumption and scaling limit.

### P1

* Support service meshes like Istio and LinkerD (may build a different component).
* Support trace generation for more application protocols such as thrift, kafka, cql, etc.
* Provide an optional control-plane if it makes sense to support remote, dynamic configuration.

## Assumptions

* We assume customers will be OK with introduction of sub 10ms latency.
* We assume customers will be OK routing all traffic through a sidecar
* We assume that users will be able and willing to let us run a small script with root priviledges in order to set up iptables rules. We'll likely replace iptables with eBPF in future releases.
* We assume HTTP/1.1, HTTP/2 and GRPC to be the most common protocols customers will want to trace. 
* We assume customers will be be fine with not having tracing support for less common binary formats. 


## Risks

#### Reliability and stability
We'll expose statsD/prometheus compatible metrics that can be consumed by a local monitoring system or remote SaaS provider.

#### Iptables need to be set up with root priviledges
We'll need to gain users' trust in order to run the iptables script as root. The script is very small and can be easily audited but there are still some risks involved especially in non-containerized environments. 

#### Performance issues/slowdown services
Envoy claims to be very performant and is used by similar projects like Istio.
_todo_: run benchmarks to understand real time 

#### Conflicts with already existing iptables rules 
Users might be using iptables already inside their k8s pods or VMs and they can conflict with our rules since we capture all. We'll need to allow some customization/configuration of the proxy and iptables rules to let users work around such limitations. As a stretch goal, we'll also re-work our iptables rules to reduce chances of conflicts.

#### TLS support
We need to make sure we support TLS traffic. Envoy supports TLS termination so we can terminate encryption at proxy and talk to the local service on an unencrypted connection. This is how people usually deploy service behing nginx, haproxy or cloud load balancers so it should not be an issue.

#### Flexibility of envoy to allow custom use-cases
We might need to support complex use cases in future that Envoy does not support. Envoy supports custom C++ extensions that can act as L3/L4 and L7 filters. We've already implemented one custom filter to detect application protocols from data coming into a Linux socket.

#### Support unknown protocols
Envoy was designed to act as a proxy for known services and as a result it is not possible to support any protocol out of the box without forcing the users to put extra effort into configuring the L7 Inspector. So, we need to support transparent passthrough for unknown/unsupported protocols. We've built a prototye Envoy extension that acts as an L3 filter, detect known application protocol from the network stream and annotates the request. The request can then be matched easily with different L7 filters in Envoy.


### Compatibility with kubernetes CNI plugins
Users might be using CNI plugins in kubernetes such as Calico. This should not be an issue L7 Inspector will operate only inside kubernetes pods.

# Design

## Proposal

Our L7 Inspector will use iptables to capture all incoming and outgoing traffic and re-route it to the locally running Envoy instance. Envoy will be configured to transparently forward the traffic to the original destination and collect traces.

### Capture traffic with Iptables

Iptables will be used to capture all incoming and outgoing traffic and re-route it to locally running Envoy instance. Some things we need to consider when using iptables

#### Needs admin privileges

Iptables need to run as root and with NET_ADMIN capability in order to be able to inject/modify routing rules. The L7 Inspector component will consist of two sub-components, the init-script and the proxy. Init-script will be responsible for injecting iptables rules. Since it'll need to run as root, we'll keep the surface area as small as possible. The script will be implemented in bash for wide-compatibility and easy auditing.

#### Circular routing

If iptables is configured to capture all traffic, it'll also capture request initiated by Envoy to upstream servers and route them back to Envoy resulting in a circule route. Luckily, iptables supports filtering traffic by the group ID of the Linux processes. We leverage this by runnig Envoy as a special omnition-user and then configure iptables to ignore outgoing traffic from any process running as omnition-user.

#### Cluster vs internet traffic

We'll also need the ability to only capture service to service traffic and ignore traffic between services and the wider internet. We'll need to make this behaviour configurable as different users might want to handle it differently. We might need to support more granular controls in future.

#### Ignoring other traffic

We don't want to capture traffic generated by components other the main services like package managers, monitoring systems, etc. We'll need to automatically detect the ports services run on and allow users to configure it as needed.


### Envoy as proxy

We use Envoy proxy to do the heavy lifting for us. Envoy was built with very similar use-cases in mind. It supports almost all features we need and has extension points which we can use to add extra functionality. Envoy is also used by other more complex projects like Istio and Ambassador, has a thriving community, is performant and flexible enough to work for us. We can replace it in future if it turns out to be too limiting.

### Installation

#### Kubernetes

In Kubernetes, users will inject our L7 proxy as a sidecar for each running service. Installation will be quite straightforward. Users can inject omnition sidecar to their pods by adding 8-10 lines to code to their kubernetes deployment yaml files. (TODO: explore building a config injectio tool) In future releases we can leverage kubernetes' mutating webhook feature to autnomatically inject the sidecars during deployments.

#### VMs

In the first version, we'll ship a bash script that'll automate installation and setup of L7 Inspector. It'll install two CLI programs to the system. One that will be used to setup iptables rules (usually after each boot) and another that starts the envoy proxy.

#### Containers on VM

This will be similar to the VMs use-case but the envoy proxy program will ship as a docker container instead of native program.

### Configuration

We'll need to allow users to customize some behavior if needed. Users might want to change the port envoy runs on, the UID/GID of omnition-user, IP address range for which traffic is captured, etc.

#### Kubernetes and Docker

In containerized environments, all configuration options will be passed down to the containers as environment variables. If/when the number of configuration knobs increases too much, we can allow configuring using a YAML file that can be added to the containers at build time by the user. If it makes sense, we may consider a lightweight control plane in future.

### TLS

We'll support TLS by terminating encryption at the proxy and then establish an unencrypted connection to the proxies service. Such unencrypted traffic will never leave the pod/container/node the proxy and service will be running on. We'll need to allow users to easily configure the the inspector to use their TLS certificates.


#### TLS with outgoing connections

Unless the code is instrumented and modified to talk through the local Envoy proxy, we cannot intercept and modify outgoing TLS calls to other services. In such cases Envoy will only trace calls only on the receiving end/target service.

##### Tracing even with outgoing TLS connections

Some proxies such as the mitmproxy allow intercepting even TLS connections. mitmproxy manages to do this by generating fake TLS certificates on the fly and establishing a connection with the client using that certificate. It essentially pretends to be the server when talking to the client and pretends to be the client when talking to the server. In order to make this work, the local container or VM needs to trust mitmproxy's fake CA. If needed, we can try to implement similar system with Envoy or even use mitmproxy in future.

### Protocol detection

In order to support any possible protocol, we can ship configs for each protocol with the inspector and expose proxy for different protocols on different ports. However, this requires the users to configure the proxy before deploying it so they can map different services to corresponding proxy ports that support the protocol spoken by the service. It means they'll likely need to configure it for every service they want to cover. This is because Envoy does not support automatic application protocol detection right now and delegates this responsibility to the users as they have knowledge about their systems and can pre-configure Envoy to enable specific protocols for different services on different ports.

In order to have out of box support automatic tracing on known protocols and silent passthrough for unsupported protocols, we'll  extend Envoy with a custom network filter and ship a custom build.

The filter will receive network traffic, detect the application protocol and annotate the requet with the protocol name. Envoy can be then configured to have different routing and tracing rules for all known protocols, and a catch-all rule for all unknown protocols that'll just re-route the traffic at TCP level without tracing support.

#### Custom filter details

Envoy has a feature called listener filters. A listener filter is kind of a middleware that is activated whenever a new incoming connection is established. Such filter deal directly with low level sockets and usually annotate the requests with some additional information. For example, the `envoy.listener.tls_inspector` filter listens for new connections, detects if the connection is encrypted with TLS or not and annotates the request so routing rules (filter chains) can decide how to handle the request. Envoy also allows writing custom listener filters in C++ out of tree and build custom version of Envoy from that with the help of the Bazel build system. Such custom filters are not uncommon in the Envoy ecosystem and many companies build custom filters to implement custom usecases.

Out listener filter is activated by envoy every time envoy accepts a new connection. The filter also gets a reference to the socket representing the connection. We use the `recv` syscall with `MSG_PEEK` on the socket to inspect the first few kilo-bytes of the request. We then look for protocol headers in the read data and annotate the request with the protocol label if we find a well-known protocol.

Since we configure Envoy to have different routing rules for different protocols, it automatically activates the correct one based on the annotated application protocol. Envoy configuration looks something like the following where `application_protocols` value is set by our custom listener filter before Envoy matches a filter chain.

```
    - filter_chain_match: {application_protocols: "http/1.1"}
      filters:
        - name: envoy.http_connection_manager
          config: {} # http1 routing rules

    - filter_chain_match: {application_protocols: "http/2"}
      filters:
        - name: envoy.http_connection_manager
          config: {} # http2 routing rules

    - filters:
      - name: envoy.tcp_proxy
        config: {} # fallback TCP routing
```

#### Pros

* A custom filter can be very flexible and can help us support a lot of special use-cases.
* We can try to upstream the filter back to Envoy, build reputation and gain trust in the eco-system

#### Cons

* Possible increased surface area for bugs especially if filter remains in active development.
* Upgrading to newer version of Envoy will not be as straightforward as we'll need to update filter API usage if needed and rebuild Enovy.
* Will need to invest time and build expertise in C++ and low level networking as such filters can become performance bottlenecks if not built well. This can be consider a pro in the long run.


## Diagrams

## Alternatives

### eBPF vs Ipables

eBPF based routing can be more performant than iptables and apparently doesn't require root priviledges (need to confirm). Currently Envoy does not support eBPF fitered traffic when running as a transparent proxy. Iptables gets us up and running much faster and works now. We can evaluate using eBPF with Envoy in future or even build something custom on top of eBPF.

### mitmproxy vs envoy

#### Resource consumption
mitmproxy is implemented in Python3 so wont be as lightweight as Envoy. It'll have a significant memory overhead compared to Envoy. 

Winner: Envoy.

#### Performance
We have not run any benchmarks yet to compare to two but mitmproxy should have decent performance since it uses Python3's asyncio which is a non-blocking networking library like NodeJS. In any case Envoy is expected to outperform mitmproxy both because of their tech stack and also because of their intended use cases.

Winner: Envoy.

#### Customizability
mitmproxy exposes a python API for scripting so it is super easy to customize and extend. Envoy allows custom C++ L3 and L7 filters, and dynamically loadable tracing plugins. What Envoy offers is good enough for us but we'll need to build expertise in C++. 

Winner: Tie but tipping in favour of mitmproxy.

#### TLS support

mitmproxy shines when it comes to intercepting TLS traffic. mitm even stands for "man in the middle" so we could say intercepting traffic is the main usecase of mitm.

Winnder: mitmproxy


# Operations

## Telemetry

We'll add monitoring for each running Envoy proxy and report health information to Omnition backend. We'll also expose a statsd  and probably prometheus compatible stats endpoint for local monitoring support.
