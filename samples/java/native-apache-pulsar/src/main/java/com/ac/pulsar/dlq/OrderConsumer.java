package com.ac.pulsar.dlq;

import com.ac.pulsar.pojo.OrderMessage;
import org.apache.pulsar.client.api.*;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.concurrent.TimeUnit;

import static com.ac.pulsar.Config.*;

public class OrderConsumer {
    private static final Logger log = LoggerFactory.getLogger(OrderConsumer.class);

    static void main() {
        log.info("Starting Apache Pulsar DLT Order Consumer...");

        try (PulsarClient client = PulsarClient.builder()
                .serviceUrl(PULSAR_SERVICE_URL)
                .build();

             Consumer<OrderMessage> consumer = client.newConsumer(Schema.JSON(OrderMessage.class))
                     .topic(TOPIC_ORDERS)
                     .subscriptionInitialPosition(SubscriptionInitialPosition.Earliest)
                     .subscriptionName(ORDERS_SUBSCRIPTION)
                     .subscriptionType(SubscriptionType.Shared)
                     .negativeAckRedeliveryDelay(1, TimeUnit.SECONDS)
                     .deadLetterPolicy(DeadLetterPolicy.builder()
                             .maxRedeliverCount(1)
                             .deadLetterTopic(TOPIC_ORDERS_DLQ)
                             .build())
                     .subscribe()) {

            log.info("Consumer subscribed. Listening for order messages...");

            while (!Thread.currentThread().isInterrupted()) {
                Message<OrderMessage> msg = consumer.receive();
                try {
                    OrderMessage order = msg.getValue();
                    int attempt = msg.getRedeliveryCount() + 1;
                    if (order.value() > 0) {
                        log.info("SUCCESS: Order ID {} processed successfully. Acknowledging.", order.id());
                        consumer.acknowledge(msg);
                    } else {
                        log.warn("REJECTED: Order ID {} has invalid value: ${}. Negatively acknowledging. (Attempt {} of 3)",
                                order.id(), order.value(), attempt);
                        consumer.negativeAcknowledge(msg);
                    }
                } catch (Exception e) {
                    log.error("Failed to process message: {}, negative acknowledging.", msg.getMessageId(), e);
                    consumer.negativeAcknowledge(msg);
                }
            }

        } catch (PulsarClientException e) {
            log.error("Pulsar Consumer Exception: ", e);
        }
    }
}
