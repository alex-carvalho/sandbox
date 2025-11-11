package com.ac;

public class App {
    static void main() {
        new MessageProcessorService(new KafkaConsumerService(), new KafkaProducerService()).start();
    }
}
