package com.elastic.test;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.context.annotation.Bean;
import io.micrometer.core.instrument.MeterRegistry;
import io.micrometer.core.instrument.Timer;

@SpringBootApplication
public class ElasticStackDemoApplication {

    public static void main(String[] args) {
        SpringApplication.run(ElasticStackDemoApplication.class, args);
    }

    @Bean
    public MetricsConfiguration metricsConfiguration(MeterRegistry meterRegistry) {
        return new MetricsConfiguration(meterRegistry);
    }
}
