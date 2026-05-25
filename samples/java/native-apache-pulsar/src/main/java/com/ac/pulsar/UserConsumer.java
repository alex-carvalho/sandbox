package com.ac.pulsar;

import com.ac.pulsar.pojo.UserMessage;
import org.apache.pulsar.client.api.*;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import static com.ac.pulsar.Config.*;

public class UserConsumer {
    private static final Logger log = LoggerFactory.getLogger(UserConsumer.class);
    private static final String SUBSCRIPTION_NAME = "user-updates-subscription";

    static void main() {
        log.info("Starting Apache Pulsar Consumer (listening forever)...");

        try (PulsarClient client = PulsarClient.builder()
                .serviceUrl(PULSAR_SERVICE_URL)
                .build();

             Consumer<UserMessage> consumer = client.newConsumer(Schema.JSON(UserMessage.class))
                .topic(TOPIC_USER_UPDATES)
                 .consumerName("java-PulsarConsumer")
                .subscriptionName(SUBSCRIPTION_NAME)
                .subscriptionType(SubscriptionType.Exclusive)
//                .subscriptionType(SubscriptionType.Key_Shared)
//                .subscriptionType(SubscriptionType.Shared)
//                     .deadLetterPolicy()
                .subscribe()) {

            while (!Thread.currentThread().isInterrupted()) {
                Message<UserMessage> msg = consumer.receive();
                try {
                    UserMessage value = msg.getValue();
                    log.info("ID: {}, Name: {}, Timestamp: {}",
                            value.id(), value.name(), value.timestamp());
                    consumer.acknowledge(msg);
                } catch (Exception e) {
                    log.error("Failed to process message: {}", msg.getMessageId(), e);
                    consumer.negativeAcknowledge(msg);
                }
            }

        } catch (PulsarClientException e) {
            log.error("Pulsar Consumer Exception: ", e);
        }
    }
}
