FROM golang AS builder

WORKDIR /src
COPY ./src .
RUN go build -o /src/ascii-art-web-dockerize

FROM debian:bookworm-slim

LABEL description="A Dockerized ASCII Art web application" 
LABEL version="1.0"
LABEL maintainer="Allen Lee (ylee) & Othmane Afilali (oafilali)"
LABEL vcs-url="https://01.gritlab.ax/git/ylee/ascii-art-web.git"

WORKDIR /app
COPY --from=builder /src/ascii-art-web-dockerize ./ascii-art-web-dockerize
COPY /app .

EXPOSE 8080
CMD ["/app/ascii-art-web-dockerize"]