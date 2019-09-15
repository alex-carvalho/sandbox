package com.ac.batch.csvtomongo.job;

import com.ac.batch.csvtomongo.job.mongo.MongoOperationType;
import com.ac.batch.csvtomongo.job.mongo.MongoWriterFactory;
import com.ac.batch.csvtomongo.model.Sale;
import com.ac.batch.csvtomongo.repository.SaleRepository;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.batch.core.Job;
import org.springframework.batch.core.Step;
import org.springframework.batch.core.configuration.annotation.EnableBatchProcessing;
import org.springframework.batch.core.configuration.annotation.JobBuilderFactory;
import org.springframework.batch.core.configuration.annotation.StepBuilderFactory;
import org.springframework.batch.item.ItemWriter;
import org.springframework.batch.item.file.FlatFileItemReader;
import org.springframework.batch.item.file.mapping.BeanWrapperFieldSetMapper;
import org.springframework.batch.item.file.mapping.DefaultLineMapper;
import org.springframework.batch.item.file.transform.DelimitedLineTokenizer;
import org.springframework.batch.repeat.RepeatStatus;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.core.io.ClassPathResource;
import org.springframework.core.io.FileSystemResource;
import org.springframework.core.io.Resource;

import java.beans.PropertyEditor;
import java.io.IOException;
import java.time.LocalDate;
import java.util.HashMap;

import static org.springframework.util.StringUtils.isEmpty;

@EnableBatchProcessing
@Configuration
public class JobConfig {

    private final Logger logger = LoggerFactory.getLogger(JobConfig.class);

    private final SaleRepository saleRepository;
    private final JobBuilderFactory jobBuilderFactory;
    private final StepBuilderFactory stepBuilderFactory;
    private final String filePath;
    private final String[] fileHeader;
    private final Integer chunkSize;
    private final ItemWriter<Sale> writer;

    public JobConfig(SaleRepository saleRepository, JobBuilderFactory jobBuilderFactory,
                     StepBuilderFactory stepBuilderFactory,
                     MongoWriterFactory mongoWriterFactory,
                     @Value("${input.file.path}") String filePath,
                     @Value("#{'${input.file.header}'.split(',')}") String[] fileHeader,
                     @Value("${job.chunkSize}") Integer chunkSize,
                     @Value("${mongodb.operation.type}") MongoOperationType mongoOperationType) {
        this.saleRepository = saleRepository;
        this.jobBuilderFactory = jobBuilderFactory;
        this.stepBuilderFactory = stepBuilderFactory;
        this.filePath = filePath;
        this.fileHeader = fileHeader;
        this.chunkSize = chunkSize;
        this.writer = mongoWriterFactory.getWriter(mongoOperationType);
    }

    @Bean
    public Job readCSVFile() {
        return jobBuilderFactory.get("readCSVFileAndSave")
                .start(processFile())
                .next(mongoAggregation())
                .build();
    }

    private Step processFile() {
        return stepBuilderFactory.get("processFile")
                .<Sale, Sale>chunk(chunkSize)
                .reader(reader())
                .writer(writer)
                .listener(new ChunkCountListener())
                .build();
    }

    private FlatFileItemReader<Sale> reader() {
        var tokenizer = new DelimitedLineTokenizer();
        tokenizer.setNames(fileHeader);

        var fieldSetMapper = new BeanWrapperFieldSetMapper<Sale>();
        var customEditors = new HashMap<Object, PropertyEditor>();
        customEditors.put(LocalDate.class, new LocalDatePropertyEditor("M/d/yyyy"));
        fieldSetMapper.setCustomEditors(customEditors);
        fieldSetMapper.setTargetType(Sale.class);

        var lineMapper = new DefaultLineMapper<Sale>();
        lineMapper.setLineTokenizer(tokenizer);
        lineMapper.setFieldSetMapper(fieldSetMapper);

        var reader = new FlatFileItemReader<Sale>();
        reader.setResource(getResource());
        reader.setLinesToSkip(1);
        reader.setLineMapper(lineMapper);
        return reader;
    }

    private Step mongoAggregation() {
        return this.stepBuilderFactory.get("mongoAggregation")
                .tasklet((contribution, chunkContext) -> {
                    saleRepository.runAggregationSalesByRegion();
                    return RepeatStatus.FINISHED;
                })
                .build();
    }

    private Resource getResource() {
        try {
            var resource = isEmpty(filePath)
                    ? new ClassPathResource("100-Sales-Records.csv")
                    : new FileSystemResource(filePath);
            logger.info("File path: " + resource.getFile().getPath());
            return resource;
        } catch (IOException e) {
            logger.error("Error on load csv file", e);
            throw new RuntimeException(e);
        }
    }
}
