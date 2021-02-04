<br/>
<img src="https://raw.githubusercontent.com/daystram/cast/master/cast-fe/src/components/logo.svg" alt="logo" width="200"/>

---

[![Gitlab Pipeline Status](https://img.shields.io/gitlab/pipeline/daystram/cast/master)](https://gitlab.com/daystram/cast/-/pipelines)
[![Docker Pulls](https://img.shields.io/docker/pulls/daystram/cast)](https://hub.docker.com/r/daystram/cast)
[![MIT License](https://img.shields.io/github/license/daystram/cast)](https://github.com/daystram/cast/blob/master/LICENSE)

DASH video-streaming and RTMP live-streaming platform.

## Features
- [DASH](https://en.wikipedia.org/wiki/Dynamic_Adaptive_Streaming_over_HTTP) Video Streaming
- [RTMP](https://en.wikipedia.org/wiki/Real-Time_Messaging_Protocol) Live Streaming
- Live Chat (WebSocket)
- Highly Scalable Transcoder Nodes
- GPU Accelerated Video Transcoding
- [Ratify](https://ratify.daystram.com/) Authentication

### DASH
With DASH streaming, videos uploaded to __cast__ are first re-encoded by the Transcoder nodes (`cast-is`) to multiple resolutions (240p, 360p, 480p, 720p, 1080p). The video player selects the most suitable resolution, and transitions between them seamlessly, based on the client's available bandwidth and preferences.

### RTMP
RTMP streaming allows the users to stream live to their viewers via a direct uplink to __cast__'s servers. Users can use clients such as [OBS Studio](https://obsproject.com/) or [Streamlabs](https://streamlabs.com/).

## Services
The application comes in three parts:

|Name|Code Name|Stack|
|----|:-------:|-----|
|Back-end|`cast-be`|[Go](https://golang.org/), [BeeGo](https://beego.me/), [MongoDB](https://www.mongodb.com/), [RabbitMQ](https://www.rabbitmq.com/), S3|
|Transcoder|`cast-is`|[Go](https://golang.org/), [FFMpeg](https://ffmpeg.org/), [RabbitMQ](https://www.rabbitmq.com/), S3|
|Front-end|`cast-fe`|JavaScript, [ReactJS](https://beego.me/)|

## Deploy
`cast-be`, `cast-is`, and `cast-fe` are containerized and pushed to [Docker Hub](https://hub.docker.com/r/daystram/cast). They are tagged based on their application name and version, e.g. `daystram/cast:be` or `daystram/cast:be-v2.0.1`.

To run `cast-be`, run the following:
```console
$ docker run --name cast-be --env-file ./.env -p 8080:8080 -d daystram/cast:be
```

To run `cast-is`, run the following:
```console
$ docker run --name cast-is --env-file ./.env -d daystram/cast:is
```

And `cast-fe` as follows:
```console
$ docker run --name cast-fe -p 80:80 -d daystram/cast:fe
```

### Dependencies
The following are required for `cast-be` to function properly:
- [MongoDB](https://www.mongodb.com/)
- [RabbitMQ](https://www.rabbitmq.com/)
- S3 storage provider

The following are required for `cast-is` to function properly:
- [RabbitMQ](https://www.rabbitmq.com/)
- S3 storage provider

Their credentials must be provided in their respective services' configuration file.

Any S3 storage provider such as [AWS S3](https://aws.amazon.com/s3/ are supported. For this particular deployment for [cast.daystram.com](https://cast.daystram.com/), a self-hosted [MinIO](https://min.io/) is used.

### Docker Compose
For ease of deployment, the following `docker-compose.yml` file can be used to orchestrate the stack deployment:
```yaml
version: "3"
services:
  cast-fe:
    image: daystram/cast:fe
    ports:
      - "80:80"
    restart: unless-stopped
  cast-is:  # no attached GPU
    image: daystram/cast:is
    env_file:
     - /path_to_env_file/.env
    restart: unless-stopped
  cast-be:
    image: daystram/cast:be
    ports:
      - "8080:8080"
    env_file:
      - /path_to_env_file/.env
    restart: unless-stopped
  mongodb:
    image: mongo:4.4-bionic
    environment:
      MONGO_INITDB_ROOT_USERNAME: MONGODB_USER
      MONGO_INITDB_ROOT_PASSWORD: MONGODB_PASS
    expose:
      - 27017
    volumes:
      - cast-mongodb:/data/db
    restart: unless-stopped
  rabbitmq:
    image: rabbitmq:3.8-alpine
    environment:
      RABBITMQ_DEFAULT_USER: RABBITMQ_USER
      RABBITMQ_DEFAULT_PASS: RABBITMQ_PASS
    expose:
      - 5672
    restart: unless-stopped
```

### GPU Containers
The Transcoder service `cast-is` is able to utilize GPU (NVIDIA graphics cards only) for harware accelerated video transcoding. To enable this, simply set the environment `USE_CUDA=true` when running the container.

Docker Engine also needs to be configured to allow GPU passthrough. Follow the steps provided [here](https://docs.nvidia.com/datacenter/cloud-native/container-toolkit/install-guide.html#docker) to enable NVIDIA's `container-toolkit`. Ensure that the host machine already have the GPU drivers installed.

To run `cast-is` with GPU attached, use the following:
```shell
$ docker run --name cast-is --env-file ./.env --gpus all -d daystram/cast:is
```

For `docker-compose.yml`, as noted [here](https://github.com/docker/compose/issues/6691#issuecomment-758460418), use the following:
```yml
version: "3"
services:
  cast-is:
    image: daystram/cast:is
    env_file:
      - /path_to_env_file/.env
    deploy:
      resources:
        reservations:
          devices:
            - capabilities:
              - gpu
    restart: unless-stopped
```

### Ingest-Base Image
`daystram/ingest-base` ([DockerHub](https://hub.docker.com/r/daystram/ingest-base)) is the base image used to build `cast-is`. This image contains the required tools for the Transcoder to properly re-encode the source videos. The tools required are:
- [FFmpeg](https://ffmpeg.org/)
- [MP4Box + gpac](https://github.com/gpac/gpac)

This image is built on top of [NVIDIA's CUDA images](https://hub.docker.com/r/nvidia/cuda/) to enable FFmpeg harware acceleration on supported hosts. MP4Box is built from source, as seen on the [Dockerfile](https://github.com/daystram/cast/blob/master/cast-is/ingest-base.Dockerfile).

## License
This project is licensed under the [MIT License](https://github.com/daystram/cast/blob/master/LICENSE).
