
#build a tiny image
FROM alpine:latest

RUN mkdir /app

COPY mailerApp /app

CMD ["/app/mailerApp"]