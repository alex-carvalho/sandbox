package com.ac.batch.csvtomongo.job.mongo;

import com.ac.batch.csvtomongo.model.Sale;
import com.ac.batch.csvtomongo.repository.SaleRepository;
import org.springframework.batch.item.ItemWriter;
import org.springframework.batch.item.data.MongoItemWriter;
import org.springframework.data.mongodb.core.MongoTemplate;
import org.springframework.stereotype.Component;

@Component
public class MongoWriterFactory {

    private final SaleRepository saleRepository;
    private final MongoTemplate mongoTemplate;

    public MongoWriterFactory(SaleRepository saleRepository, MongoTemplate mongoTemplate) {
        this.saleRepository = saleRepository;
        this.mongoTemplate = mongoTemplate;
    }

    public ItemWriter<Sale> getWriter(MongoOperationType operationType) {
        switch (operationType) {
            case SINGLE:
                return getSingleItemWriter();
            case BULK:
            default:
                return getBulkWriter();
        }
    }

    private ItemWriter<Sale> getSingleItemWriter() {
        MongoItemWriter<Sale> mongoItemWriter = new MongoItemWriter<>();
        mongoItemWriter.setTemplate(mongoTemplate);
        return mongoItemWriter;
    }

    private ItemWriter<Sale> getBulkWriter() {
        return saleRepository::insertInBulk;
    }
}
