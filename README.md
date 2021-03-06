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

With DASH streaming, videos uploaded to **cast** are first re-encoded by the Transcoder nodes (`cast-is`) to multiple resolutions (240p, 360p, 480p, 720p, 1080p). The video player selects the most suitable resolution, and transitions between them seamlessly, based on the client's available bandwidth and preferences.

### RTMP

RTMP streaming allows the users to stream live to their viewers via a direct uplink to **cast**'s servers. Users can use clients such as [OBS Studio](https://obsproject.com/) or [Streamlabs](https://streamlabs.com/).

### GPU Hardware Acceleration

Uploaded videos are ingested by **cast**'s transcoding nodes (`cast-is`) powered with NVIDIA CUDA GPU hardware acceleration. These transcoders are highly scalable and can be deployed with a high number of replica either on- or off-premise.

With GPU acceleration, `h264_nvenc` encoder is used. On environments without GPU, `cast-is` can also be started to use CPU encoding only using `libx264` encoder.

## Test Stream

You can use FFmpeg to create a test livestream. Use the following command:

```shell
$ ffmpeg -f lavfi -re -i testsrc2=s=1920x1080:r=60,format=yuv420p -f lavfi -i sine=f=440:b=4 -ac 2 -c:a pcm_s16le -c:v libx264 -f flv rtmp://cast.daystram.com/live/STREAM_KEY
```

This will create a sample 1080p 60 FPS stream with a 440 Hz sine wave sound. Ensure the stream key is provided correctly and stream window has been opened in the Dashboard.

## Services

The application comes in three parts:

| Name       | Code Name | Stack                                                                                                                                                               |
| ---------- | :-------: | ------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Back-end   | `cast-be` | [Go](https://golang.org/), [BeeGo](https://beego.me/), [MongoDB](https://www.mongodb.com/), [RabbitMQ](https://www.rabbitmq.com/), [S3](https://aws.amazon.com/s3/) |
| Transcoder | `cast-is` | [Go](https://golang.org/), [FFmpeg](https://ffmpeg.org/), [RabbitMQ](https://www.rabbitmq.com/), [S3](https://aws.amazon.com/s3/)                                   |
| Front-end  | `cast-fe` | JavaScript, [React](https://reactjs.org/)                                                                                                                           |

## Develop

### cast-be

`cast-be` uses [Go Modules](https://blog.golang.org/using-go-modules) module/dependency manager, hence at least Go 1.11 is required. BeeGo provides [Bee](https://beego.me/docs/install/bee.md) development tool which live-reloads the application. Install the tool as documented.

To begin developing, simply enter the sub-directory and run the development server:

```shell
$ cd cast-be
$ go mod tidy
$ bee run
```

### cast-is

`cast-is` uses [Go Modules](https://blog.golang.org/using-go-modules) module/dependency manager, hence at least Go 1.11 is required. To ease development, [comstrek/air](https://github.com/cosmtrek/air) is used to live-reload the application. Install the tool as documented.

To begin developing, simply enter the sub-directory and run the development server:

```shell
$ cd cast-is
$ go mod tidy
$ air
```

### cast-fe

Populate `.env.development` with the required credentials.

To begin developing, simply enter the sub-directory and run the development server:

```shell
$ cd cast-fe
$ yarn
$ yarn serve
```

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

Any S3 storage provider such as [AWS S3](https://aws.amazon.com/s3/) are supported. For this particular deployment for [cast.daystram.com](https://cast.daystram.com/), a self-hosted [MinIO](https://min.io/) is used.

### Helm Chart

To deploy to a Kubernetes cluster, Helm charts could be used. Add the [repository](https://charts.daystram.com):

```shell
$ helm repo add daystram https://charts.daystram.com
$ helm repo update
```

Ensure you have the secrets created for `cast-be` and `cast-is` by providing the secret name in `values.yaml`, or creating the secret from a populated `.env` file (make sure it is on the same namespace as `cast` installation):

```shell
$ kubectl create secret generic secret-cast-be --from-env-file=.cast-be.env
$ kubectl create secret generic secret-cast-is --from-env-file=.cast-is.env
```

And install `cast`:

```shell
$ helm install cast daystram/cast
```

You can override the chart values by providing a `values.yaml` file via the `--values` flag.

Pre-release and development charts are accessible using the `--devel` flag. To isntall the development chart, provide the `--set image.tag=dev` flag, as development images are deployed with the suffix `dev`.

### Docker Compose

For ease of deployment, the following `docker-compose.yml` file can be used to orchestrate the stack deployment:

```yaml
version: "3"
services:
  cast-be:
    image: daystram/cast:be
    ports:
      - "8080:8080"
      - "1935:1935"
    env_file:
      - /path_to_env_file/.env
    restart: unless-stopped
  cast-is: # no attached GPU
    image: daystram/cast:is
    env_file:
      - /path_to_env_file/.env
    restart: unless-stopped
  cast-fe:
    image: daystram/cast:fe
    ports:
      - "80:80"
    restart: unless-stopped
  mongodb:
    image: mongo:4.4-bionic
    environment:
      MONGO_INITDB_ROOT_USERNAME: MONGODB_USER
      MONGO_INITDB_ROOT_PASSWORD: MONGODB_PASS
    expose:
      - 27017
    volumes:
      - /path_to_mongo_data:/data/db
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

### MongoDB Indexes

For features to work properly, some indexes needs to be created in the MongoDB instance. Use the following command in `mongo` CLI to create indexes for the collated `video` collection:

```
use MONGODB_NAME;
db.createCollection("video", {collation: {locale: "en", strength: 2}})
db.video.createIndex({title: "text", description: "text"}, {collation: {locale: "simple"}});
db.video.createIndex({hash: "hashed"});
```

## License

This project is licensed under the [MIT License](https://github.com/daystram/cast/blob/master/LICENSE).
