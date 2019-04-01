# TODO LIST API

## Spring web with PostgreSQL

### Api documentation: [localhost:8080/swagger-ui.html](http://localhost:8080/swagger-ui.html#/)

__docker image:__ alexcarvalhoac/spring-web-todolist-postgres 

---
> Start postgres container local: 
>```
>docker run --name postgres -e POSTGRES_PASSWORD=secret -e POSTGRES_DB=todolist -p 5432:5432 -d postgres
>```