package com.ac.pulsar;

import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.pulsar.core.PulsarTopic;
import org.springframework.pulsar.core.PulsarTopicBuilder;

@Configuration
class PulsarTopicConfiguration {

    @Bean
    public PulsarTopic ordersTopic(@Value("${app.pulsar.topics.orders:orders}") String ordersTopic, PulsarTopicBuilder topicBuilder) {
        return topicBuilder.name(ordersTopic).build();
    }

    @Bean
    public PulsarTopic paymentTopic( @Value("${app.pulsar.topics.payment:payment}") String paymentTopic, PulsarTopicBuilder topicBuilder) {
        return topicBuilder.name(paymentTopic)
                .numberOfPartitions(3)
                .build();
    }

}
