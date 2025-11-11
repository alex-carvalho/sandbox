package com.ac;

import org.apache.kafka.clients.consumer.ConsumerConfig;
import org.apache.kafka.clients.consumer.ConsumerRecord;
import org.apache.kafka.clients.consumer.ConsumerRecords;
import org.apache.kafka.clients.consumer.KafkaConsumer;
import org.apache.kafka.common.serialization.StringDeserializer;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.time.Duration;
import java.util.Collections;
import java.util.Properties;
import java.util.function.Consumer;

public class KafkaConsumerService {
    private static final Logger logger = LoggerFactory.getLogger(KafkaConsumerService.class);



    private final KafkaConsumer<String, String> consumer;

    public KafkaConsumerService() {
        Properties props = new Properties();
        props.put(ConsumerConfig.BOOTSTRAP_SERVERS_CONFIG, System.getProperty("kafka.bootstrap.servers",
                System.getenv().getOrDefault("KAFKA_BOOTSTRAP_SERVERS", "localhost:9092")));
        props.put(ConsumerConfig.GROUP_ID_CONFIG, System.getProperty("kafka.group.id",
                System.getenv().getOrDefault("KAFKA_GROUP_ID", "ktk-group")));
        props.put(ConsumerConfig.KEY_DESERIALIZER_CLASS_CONFIG, StringDeserializer.class.getName());
        props.put(ConsumerConfig.VALUE_DESERIALIZER_CLASS_CONFIG, StringDeserializer.class.getName());
        props.put(ConsumerConfig.AUTO_OFFSET_RESET_CONFIG, "earliest");
        props.put(ConsumerConfig.ENABLE_AUTO_COMMIT_CONFIG, "true");

        this.consumer = new KafkaConsumer<>(props);
        String topic = System.getProperty("kafka.input.topic",
                System.getenv().getOrDefault("KAFKA_INPUT_TOPIC", "source-topic"));
        consumer.subscribe(Collections.singletonList(topic));
        logger.info("Subscribed to topic {}", topic);
    }

    public void consume(Consumer<ConsumerRecord<String, String>> handler) {
        while (true) {
            ConsumerRecords<String, String> records = consumer.poll(Duration.ofMillis(200));
            for (ConsumerRecord<String, String> record : records) {
                try {
                    handler.accept(record);
                } catch (Exception e) {
                    logger.error("Error handling record", e);
                }
            }
        }
    }

    public void close() {
        logger.info("Closing consumer...");
        consumer.close();
    }
}
