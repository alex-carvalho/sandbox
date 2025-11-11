# Build image using Gradle and Java 25
FROM eclipse-temurin:25-jdk as build

WORKDIR /app
# Copy Gradle files and source
COPY . /app

# Use Gradle wrapper if present; otherwise assume gradle is available.
RUN if [ -f ./gradlew ]; then ./gradlew --no-daemon clean build -x test; else gradle clean build -x test; fi

# Runtime image
FROM eclipse-temurin:25-jre
WORKDIR /app
COPY --from=build /app/build/libs/kafka-to-kafka-1.0.0.jar /app/app.jar

EXPOSE 8080
ENTRYPOINT ["java","-jar","/app/app.jar"]
