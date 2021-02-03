FROM nvidia/cuda:11.1.1-base-ubuntu20.04 AS builder
ENV DEBIAN_FRONTEND=noninteractive
ENV FFMPEG_EXECUTABLE=/usr/bin/ffmpeg
ENV MP4BOX_EXECUTABLE=/usr/local/bin/MP4Box
ENV GPAC_EXECUTABLE=/usr/local/bin/gpac
RUN apt-get update
RUN apt-get install -y build-essential pkg-config git zlib1g-dev
RUN apt-get install -y ffmpeg
RUN git clone https://github.com/gpac/gpac /gpac_public
WORKDIR /gpac_public
RUN ./configure --static-mp4box
RUN make && make install
WORKDIR /
RUN rm -rf /gpac_public
RUN apt-get remove -y build-essential pkg-config git
RUN apt-get autoremove -y --purge
