#include <stdio.h> // Used for printf()
#include <unistd.h>  // Used for close()
#include <sys/socket.h> // Used for socket(), bind(), listen()
#include <netinet/in.h> // Used for struct sockaddr_in

int main() {
    // AF_INET: Address Family Internet - IPv4 (AF_INET6: IPv6, AF_UNIX Unix domain sockets, etc.)
    // SOCK_STREAM: TCP (SOCK_DGRAM: UDP, etc)
    // IPPROTO_TCP: Protocol for TCP (IPPROTO_UDP for UDP, etc.)
    // int socket(int domain, int type, int protocol);
    int server_fd = socket(AF_INET, SOCK_STREAM, IPPROTO_TCP);
    struct sockaddr_in address;
    int opt = 1; // bool
    // SO_REUSEADDR: Allows the socket to bind to an address that is already in use
    setsockopt(server_fd, SOL_SOCKET, SO_REUSEADDR, &opt, sizeof(opt));
    
    address.sin_family = AF_INET;
    address.sin_addr.s_addr = INADDR_ANY; // INADDR_ANY allows the server to accept connections on any IP address
    address.sin_port = htons(8080); // htons converts port number to network byte order
    
    // bind(int sockfd, const struct sockaddr *addr, socklen_t addrlen);
    bind(server_fd, (struct sockaddr *)&address, sizeof(address)); // Bind the socket to the address and port
    //int listen(int sockfd, int backlog); backlog is the maximum length of the queue of pending connections
    listen(server_fd, 3); // Marks socket as passive with backlog queue of 3 connections
    
    printf("Server listening on port 8080...\n");
    
    // Clean up
    close(server_fd);
    printf("Exiting server...\n");
    return 0;
}