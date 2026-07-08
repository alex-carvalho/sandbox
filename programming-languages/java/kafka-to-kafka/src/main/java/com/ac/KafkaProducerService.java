package com.ac;

import org.apache.kafka.clients.producer.KafkaProducer;
import org.apache.kafka.clients.producer.ProducerConfig;
import org.apache.kafka.clients.producer.ProducerRecord;
import org.apache.kafka.clients.producer.RecordMetadata;
import org.apache.kafka.common.serialization.StringSerializer;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.Properties;
import java.util.concurrent.Future;

public class KafkaProducerService {
    private static final Logger logger = LoggerFactory.getLogger(KafkaProducerService.class);

    private final KafkaProducer<String, String> producer;
    private final String outputTopic;

    public KafkaProducerService() {
        Properties props = new Properties();
        props.put(ProducerConfig.BOOTSTRAP_SERVERS_CONFIG, System.getProperty("kafka.bootstrap.servers",
                System.getenv().getOrDefault("KAFKA_BOOTSTRAP_SERVERS", "localhost:9092")));
        props.put(ProducerConfig.KEY_SERIALIZER_CLASS_CONFIG, StringSerializer.class.getName());
        props.put(ProducerConfig.VALUE_SERIALIZER_CLASS_CONFIG, StringSerializer.class.getName());
        props.put(ProducerConfig.ACKS_CONFIG, "all");

        this.producer = new KafkaProducer<>(props);
        this.outputTopic = System.getProperty("kafka.output.topic",
                System.getenv().getOrDefault("KAFKA_OUTPUT_TOPIC", "destination-topic"));
        logger.info("Producer configured for topic {}", outputTopic);
    }

    public Future<RecordMetadata> publish(String key, String value) {
        ProducerRecord<String, String> record = new ProducerRecord<>(outputTopic, key, value);
        Future<RecordMetadata> f = producer.send(record, (metadata, exception) -> {
            if (exception != null) {
                logger.error("Error sending record", exception);
            } else {
                logger.info("Sent record to topic={} partition={} offset={}", metadata.topic(), metadata.partition(), metadata.offset());
            }
        });
        return f;
    }

    public void close() {
        logger.info("Closing producer...");
        producer.close();
    }
}
