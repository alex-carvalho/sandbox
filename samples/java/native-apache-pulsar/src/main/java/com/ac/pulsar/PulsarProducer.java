package com.ac.pulsar;

import com.ac.pulsar.pojo.UserMessage;
import org.apache.pulsar.client.api.*;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.UUID;

import static com.ac.pulsar.Config.*;

public class PulsarProducer {
    private static final Logger log = LoggerFactory.getLogger(PulsarProducer.class);

    static void main() {
        log.info("Starting Apache Pulsar Producer...");

        try (PulsarClient client = PulsarClient.builder()
                .serviceUrl(PULSAR_SERVICE_URL)
                .listenerName("external")
                .build();
             Producer<UserMessage> producer = client.newProducer(Schema.JSON(UserMessage.class))
                .topic(TOPIC_USER_UPDATES)
                .create()) {

            UserMessage[] messages = {
                new UserMessage(UUID.randomUUID().toString(), "Alice Smith", "alice@example.com", System.currentTimeMillis()),
                new UserMessage(UUID.randomUUID().toString(), "Bob Jones", "bob@example.com", System.currentTimeMillis()),
                new UserMessage(UUID.randomUUID().toString(), "Charlie Brown", "charlie@example.com", System.currentTimeMillis())
            };

            for (UserMessage msg : messages) {
                MessageId msgId = producer.newMessage()
                        .value(msg)
                        .send();
                log.info("SUCCESS: Sent message! Name: '{}', MsgId: {}", msg.name(), msgId);
            }

        } catch (PulsarClientException e) {
            log.error("Pulsar Producer Exception: ", e);
        }
    }
}
