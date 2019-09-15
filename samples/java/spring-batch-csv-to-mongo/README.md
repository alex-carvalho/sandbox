# CSV to MongoDB

## Spring Batch load csv and save in MongoDB

This app load sales records from csv file using Spring Batch and store in MongoDB, 
after execute a mongo aggregation by country and city.

Sales file can be found in: http://eforexcel.com/wp/downloads-18-sample-csv-files-data-sets-for-testing-sales/


**JVM arguments:**

 Property               |  Description  |
|------------------------|---------------|
| spring.data.mongodb.uri| URI MongoDB. Default "mongodb://localhost/sales" |
| input.file.path        | Path for sales csv file. Default use file inside the project with 100 records |
| job.chunkSize          | Size of chunk. Default 100 |
| mongodb.operation.type | Define if use sequential insert ou bulk operation, values SINGLE or BULK. Default BULK |

