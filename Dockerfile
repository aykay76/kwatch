FROM golang:alpine AS builder
RUN apk add --no-cache gcc libc-dev
WORKDIR /workspace
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a

FROM alpine AS final
ARG ENV
WORKDIR /
COPY --from=builder /workspace/kwatch .
CMD [ "./kwatch" ]

LABEL description="Watches Kubernetes for namespaces changes, then publishes the changes to an event queue."
LABEL maintainer="alan.kelly.london@gmail.com"
LABEL source="https://github.com/aykay76/kwatch"
LABEL labels="kubernetes,event driven,redis,kafka,alpine"
LABEL org.opencontainers.image.documentation="https://github.com/"