ARG arch

FROM golang:1.16-alpine as builder
ENV CGO_ENABLED=0
COPY . /bot
RUN apk add make
RUN cd /bot && if [[ "$arch" == "x86_64" ]]; then make build-amd64; else make build-arm64 ; fi


FROM alpine
RUN apk add curl ffmpeg youtube-dl
COPY --from=builder /bot/bot /root
ADD https://raw.githubusercontent.com/FINCTIVE/download-videos/master/download-videos.sh /root
RUN cd /bin && \
	if [[ "$arch" == "x86_64" ]] ; then \
		curl -OL "https://github.com/iawia002/annie/releases/download/0.10.3/annie_0.10.3_Linux_64-bit.tar.gz" && \
		tar -xzf annie_0.10.3_Linux_64-bit.tar.gz && \
		rm annie_0.10.3_Linux_64-bit.tar.gz ; \
	else \
		curl -OL "https://github.com/iawia002/annie/releases/download/0.10.3/annie_0.10.3_Linux_ARM64.tar.gz" && \
		tar -xzf annie_0.10.3_Linux_ARM64.tar.gz && \
		rm annie_0.10.3_Linux_ARM64.tar.gz ; \
	fi && \
	chmod +x /root/download-videos.sh
# config
VOLUME /root/.bot/
# output
VOLUME /root/downloads/
WORKDIR /root
ENTRYPOINT ["/root/bot"]
