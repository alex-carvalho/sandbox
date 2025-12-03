package com.elastic.test;

import io.micrometer.core.instrument.Counter;
import io.micrometer.core.instrument.MeterRegistry;
import io.micrometer.core.instrument.Timer;
import org.slf4j.Logger;
import org.springframework.stereotype.Component;

@Component
public class MetricsConfiguration {
    private final Logger logger = org.slf4j.LoggerFactory.getLogger(MetricsConfiguration.class);
    private final MeterRegistry meterRegistry;
    private final Counter requestCounter;
    private final Counter errorCounter;
    private final Timer responseTimeTimer;

    public MetricsConfiguration(MeterRegistry meterRegistry) {
        this.meterRegistry = meterRegistry;
        this.requestCounter = Counter.builder("app.requests.total")
                .description("Total number of requests processed")
                .register(meterRegistry);
        this.errorCounter = Counter.builder("app.errors.total")
                .description("Total number of errors encountered")
                .register(meterRegistry);
        this.responseTimeTimer = Timer.builder("app.response.time")
                .description("Response time in milliseconds")
                .publishPercentiles(0.5, 0.95, 0.99)
                .register(meterRegistry);
        logger.info("MetricsConfiguration initialized");
    }

    public void incrementRequestCounter() {
        requestCounter.increment();
    }

    public void incrementErrorCounter() {
        errorCounter.increment();
    }

    public Timer.Sample startTimer() {
        return Timer.start(meterRegistry);
    }

    public void stopTimer(Timer.Sample sample) {
        sample.stop(responseTimeTimer);
    }

    public MeterRegistry getMeterRegistry() {
        return meterRegistry;
    }
}
