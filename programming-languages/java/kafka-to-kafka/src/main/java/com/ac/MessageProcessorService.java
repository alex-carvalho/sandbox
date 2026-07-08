package com.ac;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

public class MessageProcessorService {
    private static final Logger logger = LoggerFactory.getLogger(MessageProcessorService.class);

    private final KafkaConsumerService consumerService;
    private final KafkaProducerService producerService;

    public MessageProcessorService(KafkaConsumerService consumerService, KafkaProducerService producerService) {
        this.consumerService = consumerService;
        this.producerService = producerService;
    }

    public void start() {
        logger.info("Starting Service...");
        Runtime.getRuntime().addShutdownHook(new Thread(() -> {
            logger.info("Shutting down Service...");
            consumerService.close();
            producerService.close();
        }));

        consumerService.consume(record -> {
            logger.info("Received record key={} value={}", record.key(), record.value());
            producerService.publish(record.key(), record.value());
        });
    }
}
