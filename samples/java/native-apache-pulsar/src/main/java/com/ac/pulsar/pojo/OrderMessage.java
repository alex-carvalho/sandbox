package com.ac.pulsar.pojo;

public record OrderMessage(Integer id, String item, double value, long timestamp) {
}
