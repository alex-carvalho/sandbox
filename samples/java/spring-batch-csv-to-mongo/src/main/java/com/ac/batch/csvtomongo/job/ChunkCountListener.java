package com.ac.batch.csvtomongo.job;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.batch.core.ChunkListener;
import org.springframework.batch.core.scope.context.ChunkContext;

import java.util.concurrent.TimeUnit;

public class ChunkCountListener implements ChunkListener {

    private static final Logger LOGGER = LoggerFactory.getLogger(ChunkCountListener.class);

    private long actualTime = System.currentTimeMillis();

    @Override
    public void afterChunk(ChunkContext context) {
        if (TimeUnit.MILLISECONDS.toSeconds(System.currentTimeMillis() - actualTime) >= 2) {
            int count = context.getStepContext().getStepExecution().getReadCount();
            LOGGER.info("{} items processed", count);
            actualTime = System.currentTimeMillis();
        }
    }

    @Override
    public void beforeChunk(ChunkContext context) {
    }

    @Override
    public void afterChunkError(ChunkContext context) {
    }
}