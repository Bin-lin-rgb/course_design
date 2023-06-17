FROM alpine:latest

RUN mkdir "/app"
RUN mkdir "/app/config"

WORKDIR "/app"
COPY course_design /app/
COPY config.json /app/config/

ENTRYPOINT ["./course_design"]