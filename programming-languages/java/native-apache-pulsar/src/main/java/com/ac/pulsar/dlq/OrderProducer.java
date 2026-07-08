package com.ac.pulsar.dlq;

import com.ac.pulsar.pojo.OrderMessage;
import org.apache.pulsar.client.api.*;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import static com.ac.pulsar.Config.*;

public class OrderProducer {
    private static final Logger log = LoggerFactory.getLogger(OrderProducer.class);

    static void main() {
        log.info("Starting Apache Pulsar DLT Order Producer...");

        try (PulsarClient client = PulsarClient.builder()
                .serviceUrl(PULSAR_SERVICE_URL)
                .build();
             Producer<OrderMessage> producer = client.newProducer(Schema.JSON(OrderMessage.class))
                     .topic(TOPIC_ORDERS)
                     .producerName("dlq-order-producer")
                     .create()) {

            OrderMessage[] orders = {
                    new OrderMessage(101, "Laptop (Valid)", 1200.50, System.currentTimeMillis()),
                    new OrderMessage(102, "Promo Discount (Invalid)", -25.00, System.currentTimeMillis()),
                    new OrderMessage(103, "Smartphone (Valid)", 850.00, System.currentTimeMillis()),
                    new OrderMessage(104, "Free Gift Error (Invalid)", 0.00, System.currentTimeMillis()),
                    new OrderMessage(105, "Headphones (Valid)", 120.00, System.currentTimeMillis())
            };

            for (OrderMessage order : orders) {
                MessageId msgId = producer.newMessage()
                        .value(order)
                        .key(String.valueOf(order.id()))
                        .send();
                log.info("SENT: Order ID {} ({}) - Value: ${} - Message ID: {}",
                        order.id(), order.item(), order.value(), msgId);
            }

        } catch (PulsarClientException e) {
            log.error("Pulsar Producer Exception: ", e);
        }
    }
}
