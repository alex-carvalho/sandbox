import socket

def main():
    # AF_INET: Address Family Internet - IPv4 (AF_INET6: IPv6, AF_UNIX Unix domain sockets, etc.)
    # SOCK_STREAM: TCP (SOCK_DGRAM: UDP, etc)
    # IPPROTO_TCP: Protocol for TCP (IPPROTO_UDP for UDP, etc.)
    server_fd = socket.socket(socket.AF_INET, socket.SOCK_STREAM, socket.IPPROTO_TCP)
    
    # SO_REUSEADDR: Allows the socket to bind to an address that is already in use, last parameter is a boolean
    server_fd.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
    
    # Port 8080 is converted to network byte order automatically in Python
    address = ('0.0.0.0', 8080)  # '0.0.0.0' is equivalent to INADDR_ANY (allows the server to accept connections on any IP address)
    server_fd.bind(address)  # Bind the socket to the address and port
    
    # backlog is the maximum length of the queue of pending connections
    server_fd.listen(3)  # Marks socket as passive with backlog queue of 3 connections
    
    print("Server listening on port 8080...")
    
    # Clean up
    server_fd.close()
    print("Exiting server...")

if __name__ == "__main__":
    main()