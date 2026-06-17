package com.ac.chat.domain;

import java.util.Map;
import java.util.concurrent.ConcurrentHashMap;
import java.util.concurrent.atomic.AtomicLong;

public class IdGenerator {

    private static final Map<String, AtomicLong> IDS = new ConcurrentHashMap<>();

    public static Long getNextId(String identifier) {
        return IDS.computeIfAbsent(identifier, k -> new AtomicLong(0)).addAndGet(1);
    }

}
