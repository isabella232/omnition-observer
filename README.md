## A collection of [Envoy](envoyproxy.io) proxy filter examples

### License

This work is is licensed in [MIT](LICENSE). While Envoy is licensed in [Apache License 2.0](https://github.com/envoyproxy/envoy/blob/master/LICENSE).




### Using git proxy

- Add alias for the misbehaving repository to `~/.gitconfig`

```
[url "http://localhost:8000/google/boringssl"]
        insteadOf = https://github.com/google/boringssl
```

- ____mention how to isntall and setup git server


```
```

- Fetch shallow clones

```
git clone github.com/google/boringssl --shallow
```


- Convert shallow to full clone

```
cd boringssl
git fetch --unshallow
git config remote.origin.fetch "+refs/heads/*:refs/remotes/origin/*"
git fetch origin
git pull origin '*:*'
```

- Push to local "proxy" server

```
git push http://localhost:8000/google/boringssl '*:*'
git push http://localhost:8000/google/boringssl --all
```



