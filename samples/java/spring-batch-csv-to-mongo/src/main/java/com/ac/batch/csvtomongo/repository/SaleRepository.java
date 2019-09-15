package com.ac.batch.csvtomongo.repository;

import com.ac.batch.csvtomongo.model.Sale;
import com.mongodb.BasicDBObject;
import org.springframework.data.mongodb.core.BulkOperations;
import org.springframework.data.mongodb.core.MongoTemplate;
import org.springframework.data.mongodb.core.aggregation.Aggregation;
import org.springframework.stereotype.Component;

import java.util.List;

import static org.springframework.data.mongodb.core.aggregation.Aggregation.group;
import static org.springframework.data.mongodb.core.aggregation.Aggregation.newAggregation;
import static org.springframework.data.mongodb.core.aggregation.Aggregation.out;

@Component
public class SaleRepository {

    private final MongoTemplate mongoTemplate;

    public SaleRepository(MongoTemplate mongoTemplate) {
        this.mongoTemplate = mongoTemplate;
    }

    public void runAggregationSalesByRegion() {
        Aggregation agg = newAggregation(
                group("region", "country")
                        .sum("$unitsSold").as("totalUnitsSoldByCity"),
                group("_id.region")
                        .sum("totalUnitsSoldByCity").as("totalUnitsSoldByRegion")
                        .push(new BasicDBObject()
                                .append("name", "$_id.country")
                                .append("totalUnitsSoldByCity", "$totalUnitsSoldByCity")).as("countries"),
                out("sales_by_region")
        );

        mongoTemplate.aggregate(agg, Sale.class, Object.class);
    }

    public void insertInBulk(List<? extends Sale> sales) {
        mongoTemplate
                .bulkOps(BulkOperations.BulkMode.ORDERED, Sale.class)
                .insert(sales)
                .execute();
    }

}
