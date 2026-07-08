package com.ac.pulsar.api;

import com.ac.pulsar.OrderCreatedEvent;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.pulsar.core.PulsarTemplate;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

@RestController
@RequestMapping("/order")
public class OrderApi {

    private final PulsarTemplate<OrderCreatedEvent> pulsarTemplate;
    private final String topic;

    public OrderApi(PulsarTemplate<OrderCreatedEvent> pulsarTemplate, @Value("${app.pulsar.topics.orders:orders}") String topic) {
        this.pulsarTemplate = pulsarTemplate;
        this.topic = topic;

    }

    @PostMapping
    public void createOrder() {
        OrderCreatedEvent event = new OrderCreatedEvent((int)(Math.random() * 100) + 1, 10);

        pulsarTemplate.send(topic, event);
    }
}
