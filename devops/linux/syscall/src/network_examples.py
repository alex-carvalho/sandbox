import socket

def main():
    # Create socket (equivalent to socket() syscall)
    server = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    
    # Enable address reuse (equivalent to setsockopt())
    server.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
    
    # Bind and listen (equivalent to bind() and listen() syscalls)
    server.bind(('0.0.0.0', 8080))
    server.listen(3)
    
    print("Server listening on port 8080...")
    
    # Clean up (equivalent to close() syscall)
    server.close()

if __name__ == "__main__":
    main()