
#build a tiny image
FROM alpine:latest

RUN mkdir /app

COPY listenerApp /app

CMD ["/app/listenerApp"]