# Native Apache Pulsar Functions Example (Java)

A modern, complete demonstration of **Apache Pulsar Functions** implemented in Java. The project showcases how to build, run, and verify a lightweight stream processing pipeline using Pulsar's Java SDK.

---

## Features

- Simple produce consumer
- Functions


### Run the Pulsar Function (Localrun Mode)
Use Pulsar's `localrun` command to start the function locally on your machine. This runs the function process independently of the main broker environment, which is ideal for development and testing.

```bash
pulsar-admin functions localrun \
  --jar "$PWD/build/libs/native-apache-pulsar.jar" \
  --className com.ac.pulsar.functions.UserAddEmailFunction \
  --inputs persistent://public/default/user-updates-topic \
  --output persistent://public/default/enriched-user-updates-topic \
  --broker-service-url pulsar://localhost:6651
  
./gradlew runUserWithEmailConsumer
```


### Consumers and producers

```bash
./gradlew runUserConsumer
./gradlew runUserProducer
```
