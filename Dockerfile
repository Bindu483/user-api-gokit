FROM alpine:3.13.5

COPY main /app/user-api

EXPOSE 8500

CMD /app/user-api