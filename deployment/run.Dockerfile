FROM debian

WORKDIR /driver

COPY --from=build:develop /driver_build/app /driver/app

CMD ["./app", "-cfg", "config.yaml"]