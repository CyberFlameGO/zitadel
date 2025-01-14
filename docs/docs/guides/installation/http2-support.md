---
title: HTTP/2 Support
---

The ZITADEL console (prefix `/ui/console`) uses [gRPC-Web](https://github.com/grpc/grpc-web) for its API calls.
The ZITADEL backend service accepts gRPC-Web requests and translates them into real gRPC calls to itself.
Because ZITADEL accepts gRPC-Web and translates it to gRPC itself, your reverse proxy doesn't need to be able to support gRPC or gRPC-Web.
However, as gRPC requires HTTP/2, your reverse proxy is required to send and receive downstream and upstream HTTP/2 traffic.

Go to the [loadbalancing example with Traefik](./loadbalancing-example) for seeing a working example configuration.
