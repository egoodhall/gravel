FROM golang:1.19 as build

WORKDIR /build

COPY ../ /build/

RUN go run . build

FROM gcr.io/distroless/base:debug

COPY --from=build /build/gravel/bin/gravel /bin/gravel

ENTRYPOINT [ "gravel" ]
