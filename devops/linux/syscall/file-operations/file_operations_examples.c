#include <stdio.h>      // Used for printf(), perror()
#include <unistd.h>     // Used for syscall()
#include <sys/syscall.h> // Used for SYS_open, SYS_write, SYS_read, SYS_close, SYS_lseek
#include <fcntl.h>      // Used for O_* flags
#include <string.h>     // Used for strlen>

int main() {
    const char *filename = "example.txt";
    const char *msg = "Hello, file operations with syscalls!\n";
    char buffer[128];

    // open syscall: int open(const char *pathname, int flags, mode_t mode);
    // O_CREAT: Create file if it does not exist
    // O_WRONLY: Open for writing only
    // 0644: File permissions (rw-r--r--)
    int fd = syscall(SYS_open, filename, O_CREAT | O_WRONLY | O_TRUNC, 0644);
    if (fd < 0) {
        perror("open failed");
        return 1;
    }
    printf("File '%s' opened with file descriptor %d\n", filename, fd);

    // write syscall: ssize_t write(int fd, const void *buf, size_t count);
    ssize_t written = syscall(SYS_write, fd, msg, strlen(msg));
    if (written < 0) {
        perror("write failed");
        syscall(SYS_close, fd); // Clean up
        return 1;
    }
    printf("Wrote %zd bytes to file.\n", written);

    // lseek syscall: off_t lseek(int fd, off_t offset, int whence);
    // Move file offset to the position 7
    off_t offset = syscall(SYS_lseek, fd, 7, SEEK_SET);
    if (offset == (off_t)-1) {
        perror("lseek failed");
        syscall(SYS_close, fd);
        return 1;
    }
    printf("File offset repositioned to %ld (position 7).\n", (long)offset);

    // write again at the position 7 (overwrite)
    const char *overwrite_msg = "Overwritten! ";
    written = syscall(SYS_write, fd, overwrite_msg, strlen(overwrite_msg));
    if (written < 0) {
        perror("write failed (overwrite)");
        syscall(SYS_close, fd);
        return 1;
    }
    printf("Overwrote %zd bytes at the 7 position of the file.\n", written);

    // close syscall: int close(int fd);
    if (syscall(SYS_close, fd) != 0) {
        perror("close failed"); 
        return 1;
    }
    printf("File closed after writing.\n");

    // Reopen file for reading
    fd = syscall(SYS_open, filename, O_RDONLY);
    if (fd < 0) {
        perror("open failed (read)");
        return 1;
    }
    printf("File '%s' reopened for reading with file descriptor %d\n", filename, fd);

    // read syscall: ssize_t read(int fd, void *buf, size_t count);
    ssize_t read_bytes = syscall(SYS_read, fd, buffer, sizeof(buffer) - 1);
    if (read_bytes < 0) {
        perror("read failed");
        syscall(SYS_close, fd); // Clean up
        return 1;
    }
    buffer[read_bytes] = '\0'; // Null-terminate the buffer
    printf("Read %zd bytes: %s", read_bytes, buffer);

    // close syscall: int close(int fd);
    if (syscall(SYS_close, fd) != 0) {
        perror("close fail (read)");
        return 1;
    }
    printf("File closed after reading.\n");
    
    // unlink syscall: int unlink(const char *pathname);
    // Removes the file from the filesystem
    if (syscall(SYS_unlink, filename) != 0) {
        perror("unlink failed");
        return 1;
    }
    printf("File '%s' deleted from filesystem.\n", filename);

    return 0;
}