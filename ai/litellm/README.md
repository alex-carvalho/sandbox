## Simple Sample using litellm


```shell
   docker compose up -d
```

Test it with a request:
   ```shell
   curl -X POST 'http://localhost:4000/chat/completions' \
   -H 'Content-Type: application/json' \
   -H 'Authorization: Bearer sk-xxxx' \
   -d '{
       "model": "llama3",
       "messages": [{"role": "user", "content": "Hello!"}]
   }'
   ```
