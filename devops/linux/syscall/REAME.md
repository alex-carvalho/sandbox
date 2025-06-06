

# Linux System Calls Reference

Here are the most commonly used Linux syscalls and their purposes:

## Process Management
- `fork()` - Creates a new process by duplicating the calling process
- `exec()` - Replaces current process with a new program
- `exit()` - Terminates the calling process
- `wait()` - Waits for child process to terminate
- `kill()` - Sends a signal to a process

## File Operations
- `open()` - Opens a file or creates it if it doesn't exist
- `close()` - Closes a file descriptor
- `read()` - Reads from a file descriptor
- `write()` - Writes to a file descriptor
- `lseek()` - Repositions read/write file offset
- `unlink()` - Deletes a name from the filesystem

## Memory Management
- `mmap()` - Maps files or devices into memory
- `munmap()` - Unmaps files or devices from memory
- `brk()` - Changes data segment size
- `mprotect()` - Sets protection on a region of memory

## File System Operations
- `mkdir()` - Creates a directory
- `rmdir()` - Removes a directory
- `link()` - Creates a hard link to a file
- `symlink()` - Creates a symbolic link
- `chmod()` - Changes permissions of a file
- `chown()` - Changes ownership of a file

## Network Operations
- `socket()` - Creates an endpoint for communication
- `bind()` - Binds a name to a socket
- `connect()` - Initiates a connection on a socket
- `listen()` - Listens for connections on a socket
- `accept()` - Accepts a connection on a socket
- `send()` - Sends data on a socket
- `recv()` - Receives data from a socket

## Inter-Process Communication
- `pipe()` - Creates a pipe
- `msgget()` - Gets a message queue identifier
- `semget()` - Gets a semaphore set identifier
- `shmget()` - Gets shared memory segment

## System Information
- `uname()` - Gets system information
- `getpid()` - Gets process ID
- `getuid()` - Gets user ID
- `gettimeofday()` - Gets time of day

## Security
- `chmod()` - Changes file permissions
- `chown()` - Changes file ownership
- `getcwd()` - Gets current working directory
- `chroot()` - Changes root directory

Note: This is a subset of the most commonly used syscalls. Linux has over 300 system calls in total.


``` shell
gcc network_examples.c -o bin/network_examples
./bin/network_examples
python3 network_examples.py 



gcc process_management_examples.c -o bin/process_management_examples
./bin/process_management_examples
python3 process_management_examples.py 
```