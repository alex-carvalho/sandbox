package com.ac.pulsar.dlq;

import com.ac.pulsar.pojo.OrderMessage;
import org.apache.pulsar.client.api.*;
import org.apache.pulsar.common.api.proto.CommandSubscribe;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import static com.ac.pulsar.Config.*;

public class OrderDLQConsumer {
    private static final Logger log = LoggerFactory.getLogger(OrderDLQConsumer.class);

    static void main() {
        log.info("Starting Apache Pulsar Dead Letter Consumer...");

        try (PulsarClient client = PulsarClient.builder()
                .serviceUrl(PULSAR_SERVICE_URL)
                .build();

             Consumer<OrderMessage> consumer = client.newConsumer(Schema.JSON(OrderMessage.class))
                     .topic(TOPIC_ORDERS_DLQ)
                     .subscriptionInitialPosition(SubscriptionInitialPosition.Earliest)
                     .subscriptionName(ORDERS_DLQ_SUBSCRIPTION)
                     .subscriptionType(SubscriptionType.Exclusive)
                     .subscribe()) {

            while (!Thread.currentThread().isInterrupted()) {
                Message<OrderMessage> msg = consumer.receive();
                try {
                    OrderMessage order = msg.getValue();
                    log.error("OrderDLQConsumer - Topic: {} - Message ID: {} - Order ID: {} ({}) - Value: ${}",
                            msg.getTopicName(), msg.getMessageId(), order.id(), order.item(), order.value());
                    consumer.acknowledge(msg);
                } catch (Exception e) {
                    log.error("Failed to process dead letter message: {}", msg.getMessageId(), e);
                    consumer.negativeAcknowledge(msg);
                }
            }

        } catch (PulsarClientException e) {
            log.error("Pulsar Dead Letter Consumer Exception: ", e);
        }
    }
}
