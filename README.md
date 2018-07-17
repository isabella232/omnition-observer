# Omnition L7 Inspector


## Installation 

L7 inspector can be installed as a sidecar to your application pods by adding a few lines to your apps' kubernetes yaml definitions.

### Example Pod Spec

```
spec:
  ...
  containers:
    ... your app container(s)

    # Omnition L7 Proxy Container
    - name: omnition-proxy
      image: 592865182265.dkr.ecr.us-west-2.amazonaws.com/omnition/proxy:latest
      env:
        - name: ZIPKIN_HOST
          value: "<ZIPKIN_HOST>"
        - name: ZIPKIN_PORT
          value: "<ZIPKIN_PORT>"
        - name: SERVICE_NAME
          value: "<SERVICE_NAME>"

  initContainers:
    ... your init container(s)

    # Omnition Init container
    - args:
      image: 592865182265.dkr.ecr.us-west-2.amazonaws.com/omnition/proxy-init:latest
      name: omnition-proxy-init
      securityContext:
        capabilities:
          add:
          - NET_ADMIN
        privileged: true
```

Here we injected two containers in the pod. An init-container and a proxy container. The init-container runs before everything else and sets up networking routing rules that route all traffic in and out of the pod through the proxy.

## Configuration

The containers can be configured through environment variables. Currently supported variables are:

#### ZIPKIN_HOST

Default value: `zipkin.opsoss.svc.cluster.local`

This address at which zipkin is listening for new traces.


#### ZIPKIN_PORT

Default value: `9411`

This port on which zipkin is listening for new traces.


#### SERVICE_NAME

This gives the service a label in zipkin/jaeger instead of using service's IP address.
