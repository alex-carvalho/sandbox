package com.ac;

import org.apache.kafka.clients.consumer.ConsumerConfig;
import org.apache.kafka.clients.consumer.KafkaConsumer;
import org.apache.kafka.clients.consumer.ConsumerRecords;
import org.apache.kafka.clients.consumer.ConsumerRecord;
import org.apache.kafka.clients.producer.KafkaProducer;
import org.apache.kafka.clients.producer.ProducerConfig;
import org.apache.kafka.clients.producer.ProducerRecord;
import org.apache.kafka.common.serialization.StringDeserializer;
import org.apache.kafka.common.serialization.StringSerializer;
import org.junit.jupiter.api.*;
import org.testcontainers.kafka.ConfluentKafkaContainer;
import org.testcontainers.utility.DockerImageName;

import java.time.Duration;
import java.util.Collections;
import java.util.Properties;

public class MessageProcessorServiceTest {

    static ConfluentKafkaContainer kafka;

    @BeforeAll
    public static void setup() {
        kafka = new ConfluentKafkaContainer(DockerImageName.parse("confluentinc/cp-kafka:latest"));
        kafka.start();
        System.setProperty("kafka.bootstrap.servers", kafka.getBootstrapServers());
        System.setProperty("kafka.input.topic", "tc-source");
        System.setProperty("kafka.output.topic", "tc-dest");
    }

    @AfterAll
    public static void teardown() {
        if (kafka != null) kafka.stop();
    }

    @Test
    public void testRelay() throws Exception {
        // Start relay in separate thread (it will subscribe to tc-source)
        Thread relayThread = new Thread(() -> {
            MessageProcessorService relay = new MessageProcessorService(new KafkaConsumerService(), new KafkaProducerService());
            relay.start();
        }, "relay-thread");
        relayThread.setDaemon(true);
        relayThread.start();

        // Give relay a moment to start
        Thread.sleep(2000);

        // Produce a message to input topic
        Properties prodProps = new Properties();
        prodProps.put(ProducerConfig.BOOTSTRAP_SERVERS_CONFIG, kafka.getBootstrapServers());
        prodProps.put(ProducerConfig.KEY_SERIALIZER_CLASS_CONFIG, StringSerializer.class.getName());
        prodProps.put(ProducerConfig.VALUE_SERIALIZER_CLASS_CONFIG, StringSerializer.class.getName());
        KafkaProducer<String, String> producer = new KafkaProducer<>(prodProps);
        producer.send(new ProducerRecord<>("tc-source", "k1", "hello-test")).get();

        // Consume from output topic
        Properties consProps = new Properties();
        consProps.put(ConsumerConfig.BOOTSTRAP_SERVERS_CONFIG, kafka.getBootstrapServers());
        consProps.put(ConsumerConfig.GROUP_ID_CONFIG, "test-group");
        consProps.put(ConsumerConfig.KEY_DESERIALIZER_CLASS_CONFIG, StringDeserializer.class.getName());
        consProps.put(ConsumerConfig.VALUE_DESERIALIZER_CLASS_CONFIG, StringDeserializer.class.getName());
        consProps.put(ConsumerConfig.AUTO_OFFSET_RESET_CONFIG, "earliest");

        KafkaConsumer<String, String> consumer = new KafkaConsumer<>(consProps);
        consumer.subscribe(Collections.singletonList("tc-dest"));

        boolean found = false;
        long deadline = System.currentTimeMillis() + 10_000;
        while (System.currentTimeMillis() < deadline && !found) {
            ConsumerRecords<String, String> records = consumer.poll(Duration.ofMillis(500));
            for (ConsumerRecord<String, String> r : records) {
                if ("k1".equals(r.key()) && "hello-test".equals(r.value())) {
                    found = true;
                    break;
                }
            }
        }
        consumer.close();
        producer.close();

        Assertions.assertTrue(found, "Message should be relayed to destination topic");
    }
}
