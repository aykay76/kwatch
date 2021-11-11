# kwatch

## Overview

A simple microservice that watches for changes in Kubernetes internal state and publishes events on a Redis Stream for others to consume.
On its own it may not seem particularly useful but as part of an event driven architecture it will be an essential componentof a larger system.

My motivation is perhaps better documented in my `autograf` repository which combines this with other microservices to automate the creation of dashboards in Grafana. Just one use case that also might be of use to others.

## To use..

You can run this program from inside Kubernetes (by deploying it there) or externally by just running it locally. It will detect using environment variables specific to Kubernetes where it is running.
Note: if running outside a cluster you will need a `~/.kube/config` file that allows you to connect to the cluster you wish to watch.

On starting it will attach to the namespace API and watch for changes. As changes occur it will publish messages to a Redis stream. (I chose Redis Streams because it's GA in Dapr, a CNCF framework).
Of course this could be extended to watch other watchable 'things' in k8s.
(It's worth noting that an `added` notification happens for every existing namespaces on startup, so handle that in an idempotent way - or do what I do and just push those events on and let the downstream consumers worry about it `:D`)
I settled on a single stream/topic in Redis called `kubernetes` so other metadta could be used to identify the resource type and action that happened.

You can configure the Redis endpoint to connect to with the `REDIS_ADDR` environment variable. I have tried to adhere to 12-factor principles and ensure the service is confiured through environment variables. For local testing it can also be configured with command line flags.

And that's it - microdocumentation for a microservice that does one thing.

Enjoy, star, follow `;)`

---

Basic build size: 39MB
Build size with Redis: 40MB
Build size with Kafka: 40MB
Build size with both: 41MB

Confluent Kafka Golang client wouldn't work on Apple M1 because:
```
ld: warning: ignoring file ../../../go/pkg/mod/gopkg.in/confluentinc/confluent-kafka-go.v1@v1.7.0/kafka/librdkafka_vendor/librdkafka_darwin.a, building for macOS-arm64 but attempting to link with file built for macOS-x86_64
```
Using segmentio Kafka instead.