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

ENTRYPOINT ["java", \
    "-javaagent:/app/apm-agent.jar", \
    "-Delastic.apm.service_name=elastic-stack-demo", \
    "-Delastic.apm.server_urls=http://apm-server:8200", \
    "-Delastic.apm.environment=testing", \
    "-Delastic.apm.transaction_sample_rate=1.0", \
    "-Delastic.apm.log_correlation=true", \
    "-Dspring.profiles.active=docker", \
    "-Dlogging.file.name=/app/logs/app.log", \
    "-jar", \
    "app.jar"]
