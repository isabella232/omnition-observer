# Protocol Annotator

This is an HTTP filter that annotates upstream request with request protocol
string as a header entry with a specified header name.

The possible values are:

- `HTTP/1.0`
- `HTTP/1.1`
- `HTTP/2`
- `-` (if no valid protocol detected)

A configuration example:

```yaml
http_filters:
  - name: io.omnition.envoy.http.protocol_annotator
    config:
      header_name: x-envoy-annotated-protocol
```

This can be used as a matching value against Route's [HeaderMatcher](https://github.com/envoyproxy/envoy/blob/master/api/envoy/api/v2/route/route.proto#L293)
config.

## License

SEE LICENSE at [LICENSE](LICENSE).