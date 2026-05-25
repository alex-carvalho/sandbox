package com.ac.pulsar.functions;

import com.ac.pulsar.pojo.UserMessage;
import com.ac.pulsar.pojo.UserWithEmailMessage;
import org.apache.pulsar.functions.api.Context;
import org.apache.pulsar.functions.api.Function;
import org.slf4j.Logger;

import java.time.Instant;

public class UserAddEmailFunction implements Function<UserMessage, UserWithEmailMessage> {

    @Override
    public UserWithEmailMessage process(UserMessage input, Context context) throws Exception {
        Logger log = context.getLogger();
        
        if (input == null) {
            log.warn("Received a null user update message");
            return null;
        }

        log.info("Processing user update - ID: {}, Name: {}", input.id(), input.name());

        UserWithEmailMessage enriched = new UserWithEmailMessage(
                input.id(),
                 input.name(),
                 input.name().toLowerCase().replace("_" , "") + "@company.com" ,
                input.timestamp(),
                Instant.now().toString()
        );

        log.info("Successfully enriched user ID {}: Name: {}, Email: {}", 
                enriched.id(), enriched.name(), enriched.email());

        return enriched;
    }
}
