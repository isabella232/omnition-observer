Omnition Observer consists of two components, `omnition-observer-init` and `omnition-observer`. Both components must be enabled in order for the observer to function correctly. Omnition Observer can be deployed as a sidecar to Kubernetes pods to enable automatic tracing between services without instrumenting any code. 

## Installing in Kubernetes

To add Omnition Observer to any service running Kubernetes, you'll need to add `omnition-observer-init` as an init container and `omnition-observer` as a regular container to the pod. All incoming and outgoing traffic will automatically be routed through the observer.

### Step1: Add observer init container

`omnition-observer-init` must be added to a k8s pod as an init container.  If there are more than one init containers then `omnition-observer-init` must be the last one in the list. Example YAML config:

```
      initContainers:
      - name: omnition-observer-init
        image: omnition/omnition-observer-init:0.3.0
        imagePullPolicy: Always
        securityContext:
          capabilities:
            add:
            - NET_ADMIN
          privileged: true
```

### Step2: Add observer

The observer must be added as a regular k8s container to the pod. Example YAML config:

```
      containers:
      - name: omnition-observer
        image: omnition/omnition-observer:0.3.0
        imagePullPolicy: Always
        ports:
          - containerPort: 9901
            name: proxy
        env:
          - name: TRACING_HOST
            value: "zipkin"
          - name: TRACING_PORT
            value: "9411"
          - name: SERVICE_NAME
            value: "my-service"
          - name: TRACING_TAG_HEADERS
            value: "x_tenant_id region version"
```

Observer exports traces in the zipkin format. It use the `TRACING_HOST` and `TRACING_PORT` environment variables to determine the address of a zipkin compatible agent or collector it should send all the traces to. 

You should also set the `SERVICE_NAME` environment variable so each service is tagged correctly in traces.

`TRACING_TAG_HEADERS` takes a space separated list of HTTP headers that will automatically be added to tracing spans as tags when found on requests.
