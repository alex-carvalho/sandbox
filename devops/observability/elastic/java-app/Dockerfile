FROM gradle:9.2.1-jdk25-ubi as builder

WORKDIR /workspace
COPY build.gradle gradle.properties ./
COPY src src

RUN gradle clean bootJar -DskipTests

RUN ls -lah build/libs/

FROM eclipse-temurin:25-jre-alpine

WORKDIR /app

COPY --from=builder /workspace/build/libs/elastic-stack-demo.jar app.jar

COPY --from=builder /workspace/build/elastic-apm-agent-*.jar apm-agent.jar

RUN mkdir -p /app/logs && chmod 777 /app/logs

EXPOSE 8080

ENV ELASTIC_APM_SERVER_URLS=http://apm-server:8200
ENV ELASTIC_APM_SERVICE_NAME=elastic-stack-demo
ENV ELASTIC_APM_ENVIRONMENT=testing

ENTRYPOINT ["java", \
    "-javaagent:/app/apm-agent.jar", \
    "-Delastic.apm.transaction_sample_rate=1.0", \
    "-Delastic.apm.log_ecs_reformatting=SHADE", \
    "-Dspring.profiles.active=docker", \
    "-Dlogging.file.name=/app/logs/app.log", \
    "-jar", \
    "app.jar"]
