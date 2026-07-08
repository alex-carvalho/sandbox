FROM adoptopenjdk/openjdk11:jre-11.0.2.9-alpine
LABEL author="Aléx Carvalho"
COPY build/libs/*.jar app.jar
CMD java ${JAVA_OPTS} -jar app.jar