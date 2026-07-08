package com.ac.pulsar;

public class Config {

    public static final String PULSAR_SERVICE_URL = "pulsar://localhost:6651";
    public static final String TOPIC_USER_UPDATES = "persistent://public/default/user-updates-topic";
    public static final String TOPIC_ENRICHED_USER_UPDATES = "persistent://public/default/enriched-user-updates-topic";
    public static final String TOPIC_ORDERS = "persistent://public/default/orders";
    public static final String TOPIC_ORDERS_DLQ = "persistent://public/default/orders-dlq";
    public static final String ORDERS_SUBSCRIPTION = "orders-subscription";
    public static final String ORDERS_DLQ_SUBSCRIPTION = "orders-dlq-subscription";


}
