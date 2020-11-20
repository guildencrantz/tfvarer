FROM golang:1.13 AS build

COPY . /tfvarser
WORKDIR /tfvarser

RUN go build .

FROM scratch

COPY --from=build /tfvarser/tfvarser /tfvarser

ENTRYPOINT ["/tfvarser"]
