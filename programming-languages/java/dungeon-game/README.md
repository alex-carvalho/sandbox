# Dungeon Game API

A Java 23 Spring Boot application that solves the LeetCode Dungeon Game problem (Problem 174) with REST API, PostgreSQL persistence, Docker container and docker-compose

## Running the application

1.  **Build the project:**
    ```bash
    mvn clean install
    ```

2.  **Run the application:**
   ```bash
    mvn clean spring-boot:run
    ```
    Using Docker Compose
    ```bash
    docker-compose up --build
    ```

3.  **Send a POST request to the API:**
    ```bash
    curl -X POST http://localhost:8080/game \
         -H "Content-Type: application/json" \
         -d '{"dungeon": [[-2, -3, 3], [-5, -10, 1], [10, 30, -5]]}'
    ```
