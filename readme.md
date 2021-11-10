A simple microservice that watches for Kubernetes changes and publishes events on a Redis Stream in response.

You can run this program from inside Kubernetes (by deploying it there) or "externally" / locally by just running it. It will detect using environment variables specific to Kubernetes where it is running.
Note: if running outside a cluster you will need a ~/.kube/config file that allows you to connect to the cluster you wish to watch.

On starting it will attach to the namespace API and watch for changes. As changes occur it will publish messages to a Redis stream. (I chose Redis Streams because it's GA in Dapr, a CNCF framework).

You can configure the Redis endpoint to connect to with the `REDIS_ADDR` environment variable. The stream and events are hard-coded currently but these could easily be made to be configurable.