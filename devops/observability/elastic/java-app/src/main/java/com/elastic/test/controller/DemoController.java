package com.elastic.test.controller;

import com.elastic.test.MetricsConfiguration;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.util.HashMap;
import java.util.Map;
import java.util.UUID;

@RestController
@RequestMapping("/api")
public class DemoController {

    private final Logger logger = LoggerFactory.getLogger(DemoController.class);
    private final MetricsConfiguration metrics;

    public DemoController(MetricsConfiguration metrics) {
        this.metrics = metrics;
    }

    @GetMapping("/hello")
    public ResponseEntity<Map<String, Object>> hello(@RequestParam(defaultValue = "World") String name) {
        logger.info("Received request for hello endpoint with name: {}", name);
        metrics.incrementRequestCounter();

        var sample = metrics.startTimer();
        try {
            // Simulate some processing
            Thread.sleep((long) (Math.random() * 100));

            Map<String, Object> response = new HashMap<>();
            response.put("message", "Hello, " + name + "!");
            response.put("timestamp", System.currentTimeMillis());
            response.put("requestId", UUID.randomUUID().toString());

            logger.debug("Successfully processed hello request for: {}", name);
            return ResponseEntity.ok(response);
        } catch (InterruptedException e) {
            logger.error("Interrupted while processing hello request", e);
            metrics.incrementErrorCounter();
            Thread.currentThread().interrupt();
            return ResponseEntity.status(500).build();
        } finally {
            metrics.stopTimer(sample);
        }
    }

    @GetMapping("/health-check")
    public ResponseEntity<Map<String, String>> healthCheck() {
        logger.info("Health check endpoint called");
        metrics.incrementRequestCounter();

        return ResponseEntity.ok(Map.of(
                "status", "UP",
                "service", "elastic-stack-demo",
                "timestamp", String.valueOf(System.currentTimeMillis())
        ));
    }

    @PostMapping("/test-error")
    public ResponseEntity<Map<String, String>> testError() {
        logger.warn("Test error endpoint triggered");
        metrics.incrementRequestCounter();
        metrics.incrementErrorCounter();

        logger.error("Simulating application error for testing error tracking");
        return ResponseEntity.status(500).body(Map.of(
                "error", "This is a test error for error tracking",
                "timestamp", String.valueOf(System.currentTimeMillis())
        ));
    }

    @PostMapping("/test-exception")
    public ResponseEntity<String> testException() {
        logger.warn("Test exception endpoint triggered");
        metrics.incrementRequestCounter();
        metrics.incrementErrorCounter();

        throw new RuntimeException("This is a test exception for APM tracing");
    }

    @GetMapping("/slow-endpoint")
    public ResponseEntity<Map<String, Object>> slowEndpoint(@RequestParam(defaultValue = "2000") long delayMs) {
        logger.info("Slow endpoint called with delay: {} ms", delayMs);
        metrics.incrementRequestCounter();

        var sample = metrics.startTimer();
        try {
            Thread.sleep(Math.min(delayMs, 5000)); // Cap at 5 seconds to prevent issues

            Map<String, Object> response = new HashMap<>();
            response.put("message", "This was a slow request");
            response.put("delayMs", delayMs);
            response.put("timestamp", System.currentTimeMillis());

            logger.debug("Slow endpoint completed after {} ms", delayMs);
            return ResponseEntity.ok(response);
        } catch (InterruptedException e) {
            logger.error("Slow endpoint interrupted", e);
            metrics.incrementErrorCounter();
            Thread.currentThread().interrupt();
            return ResponseEntity.status(500).build();
        } finally {
            metrics.stopTimer(sample);
        }
    }

    @GetMapping("/metrics-info")
    public ResponseEntity<Map<String, Object>> getMetricsInfo() {
        logger.info("Metrics info endpoint called");
        metrics.incrementRequestCounter();

        Map<String, Object> info = new HashMap<>();
        info.put("service", "elastic-stack-demo");
        info.put("version", "1.0.0");
        info.put("environment", "testing");
        info.put("timestamp", System.currentTimeMillis());

        return ResponseEntity.ok(info);
    }
}
