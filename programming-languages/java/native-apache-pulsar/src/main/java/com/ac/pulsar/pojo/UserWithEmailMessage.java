package com.ac.pulsar.pojo;

public record UserWithEmailMessage(
    Integer id,
    String name,
    String email,
    long timestamp,
    String processingTimestamp
) {}
