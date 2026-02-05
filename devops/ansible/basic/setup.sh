# Generate SSH key if it doesn't exist
[ ! -f ./id_rsa ] && ssh-keygen -t rsa -N "" -f ./id_rsa

docker run -d --name node1 -h node1 \
    -p 2221:2222 \
    -p 8086:80 \
    -e PUID=1000 -e PGID=1000 \
    -e PUBLIC_KEY="$(cat ./id_rsa.pub)" \
    -e USER_NAME=ansible \
    lscr.io/linuxserver/openssh-server:latest
    
docker run -d --name node2 -h node2 \
    -p 2222:2222 \
    -p 8087:80 \
    -e PUID=1000 -e PGID=1000 \
    -e PUBLIC_KEY="$(cat ./id_rsa.pub)" \
    -e USER_NAME=ansible \
    lscr.io/linuxserver/openssh-server:latest


docker exec node1 apk add --no-cache python3 sudo
docker exec node2 apk add --no-cache python3 sudo

docker exec node1 sh -c 'echo "ansible ALL=(ALL) NOPASSWD: ALL" >> /etc/sudoers'
docker exec node2 sh -c 'echo "ansible ALL=(ALL) NOPASSWD: ALL" >> /etc/sudoers'