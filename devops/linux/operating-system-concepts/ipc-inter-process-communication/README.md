Inter-Process Communication (IPC) 



gcc -std=c11 -Wall -Wextra -o bin/shared-memory-producer shared-memory-producer.c

./bin/shared-memory-producer

gcc -std=c11 -Wall -Wextra -o bin/shared-memory-consumer shared-memory-consumer.c

./bin/shared-memory-consumer