FROM docker.io/library/eclipse-temurin:21-jdk

WORKDIR /app

COPY build/libs/*.jar app.jar

LABEL org.opencontainers.image.source=https://github.com/alex-carvalho/sandbox

EXPOSE 8080

ENTRYPOINT ["java", "-jar", "app.jar"]
