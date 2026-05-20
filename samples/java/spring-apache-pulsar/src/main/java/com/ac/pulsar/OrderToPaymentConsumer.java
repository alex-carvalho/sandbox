package com.ac.pulsar;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.core.annotation.Order;
import org.springframework.pulsar.annotation.PulsarListener;
import org.springframework.pulsar.core.PulsarTemplate;
import org.springframework.stereotype.Component;

 @Component
class OrderToPaymentConsumer {

    private final Logger log = LoggerFactory.getLogger(OrderToPaymentConsumer.class);
    private final PulsarTemplate<String> pulsarTemplate;

    private final String paymentTopic;

    OrderToPaymentConsumer(PulsarTemplate<String> pulsarTemplate, @Value("${app.pulsar.topics.payment:payment}") String paymentTopic) {
        this.pulsarTemplate = pulsarTemplate;
        this.paymentTopic = paymentTopic;
    }

    @PulsarListener(id = "external", topics = "${app.pulsar.topics.orders:orders}", subscriptionName = "${app.pulsar.subscriptions.orders-to-payment:orders-to-payment}")
    void handle(OrderCreatedEvent orderMessage) {
        log.info("Received order: {}", orderMessage);
        this.pulsarTemplate.send(this.paymentTopic, orderMessage.toString());
    }
}
