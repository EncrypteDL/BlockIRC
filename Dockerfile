# Build Stage
FROM golang:alpine AS build

ARG TAG

ENV APP BlockIRC
ENV REPO EDL/$APP

RUN apk add --update git make build-base && \
    rm -rf /var/cache/apk/*

WORKDIR /go/src/github.com/$REPO
COPY . /go/src/github.com/$REPO
RUN make TAG=$TAG build

# Runtime Stage
FROM alpine

ENV APP BlockIRC
ENV REPO EDL/$APP

LABEL blockirc.app main

COPY --from=build /go/src/github.com/${REPO}/${APP} /${APP}

EXPOSE 6667/tcp 6697/tcp

ENTRYPOINT ["/BlockIRC"]
CMD [""]