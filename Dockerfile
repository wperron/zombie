FROM golang:1.16 as build
WORKDIR ./app
COPY . .
RUN make release

FROM scratch
ENV PATH=/
COPY --from=build /go/app/bin/release/zombie /zombie
CMD [ "./zombie", "-config", "/zombie.yaml" ]