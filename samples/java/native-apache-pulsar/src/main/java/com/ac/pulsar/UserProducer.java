package com.ac.pulsar;

import com.ac.pulsar.pojo.UserMessage;
import org.apache.pulsar.client.api.*;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;


import java.util.Optional;
import java.util.function.Supplier;

import static com.ac.pulsar.Config.*;

public class UserProducer {
    private static final Logger log = LoggerFactory.getLogger(UserProducer.class);

    static void main() {
        log.info("Starting Apache Pulsar Producer...");

        try (PulsarClient client = PulsarClient.builder()
                .serviceUrl(PULSAR_SERVICE_URL)
                .build();
             Producer<UserMessage> producer = client.newProducer(Schema.JSON(UserMessage.class))
                .topic(TOPIC_USER_UPDATES)
                .create()) {

            Supplier<Integer> idGenerator = () ->  (int) (Math.random() * 1000) + 1;

            UserMessage[] messages = {
                new UserMessage(idGenerator.get(), "Alice Smith", System.currentTimeMillis()),
                new UserMessage(idGenerator.get(), "Bob Jones", System.currentTimeMillis()),
                new UserMessage(idGenerator.get(), "Charlie Brown", System.currentTimeMillis())
            };

            for (UserMessage msg : messages) {
                MessageId msgId = producer.newMessage()
                        .value(msg)
                        .key(String.valueOf(msg.id()))
                        .send();
                log.info("SUCCESS: Sent message! Name: '{}', MsgId: {}", msg.name(), msgId);
            }

        } catch (PulsarClientException e) {
            log.error("Pulsar Producer Exception: ", e);
        }
    }
}
