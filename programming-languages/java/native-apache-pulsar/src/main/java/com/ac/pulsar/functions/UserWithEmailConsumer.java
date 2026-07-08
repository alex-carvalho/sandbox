package com.ac.pulsar.functions;

import com.ac.pulsar.pojo.UserWithEmailMessage;
import org.apache.pulsar.client.api.*;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import static com.ac.pulsar.Config.*;

public class UserWithEmailConsumer {
    private static final Logger log = LoggerFactory.getLogger(UserWithEmailConsumer.class);
    private static final String SUBSCRIPTION_NAME = "enriched-user-updates-subscription";

    public static void main(String[] args) {
        log.info("Starting Apache Pulsar Enriched Consumer (listening forever)...");

        try (PulsarClient client = PulsarClient.builder()
                .serviceUrl(PULSAR_SERVICE_URL)
                .listenerName("external")
                .build();

             Consumer<UserWithEmailMessage> consumer = client.newConsumer(Schema.JSON(UserWithEmailMessage.class))
                .topic(TOPIC_ENRICHED_USER_UPDATES)
                .consumerName("java-PulsarEnrichedConsumer")
                .subscriptionName(SUBSCRIPTION_NAME)
                .subscriptionType(SubscriptionType.Exclusive)
                .subscribe()) {

            while (!Thread.currentThread().isInterrupted()) {
                Message<UserWithEmailMessage> msg = consumer.receive();
                try {
                    UserWithEmailMessage value = msg.getValue();
                    log.info("ID: {}, Name: {}, Email: {}", value.id(), value.name(), value.email());
                    consumer.acknowledge(msg);
                } catch (Exception e) {
                    log.error("Failed to process enriched message: {}", msg.getMessageId(), e);
                    consumer.negativeAcknowledge(msg);
                }
            }

        } catch (PulsarClientException e) {
            log.error("Pulsar Enriched Consumer Exception: ", e);
        }
    }
}
