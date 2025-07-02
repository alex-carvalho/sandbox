#include <stdio.h>      // Used for printf(), perror()
#include <unistd.h>     // Used for syscall()
#include <sys/syscall.h> // Used for SYS_mmap, SYS_munmap, SYS_mprotect
#include <sys/mman.h>   // Used for PROT_* and MAP_* flags
#include <string.h>     // Used for strcpy()
#include <stdint.h>     // For uintptr_t

int main() {
    size_t length = 4096; // Allocate one memory page (usually 4096 bytes)
    // mmap syscall: void *mmap(void *addr, size_t length, int prot, int flags, int fd, off_t offset);
    // PROT_READ | PROT_WRITE: Pages may be read and written
    // MAP_PRIVATE | MAP_ANONYMOUS: Mapping is not backed by any file; changes are private
    void *addr = (void *)syscall(SYS_mmap, NULL, length, PROT_READ | PROT_WRITE,
                                 MAP_PRIVATE | MAP_ANONYMOUS, -1, 0);
    if (addr == MAP_FAILED) {
        perror("mmap"); // Print error if mapping fails
        return 1;
    }

    printf("Memory mapped at address %p\n", addr);

    // Write a string to the allocated memory
    strcpy((char *)addr, "Hello, memory management with syscalls!");
    printf("Content: %s\n", (char *)addr);

    // Change memory protection to read-only
    // mprotect syscall: int mprotect(void *addr, size_t len, int prot);
    if (syscall(SYS_mprotect, addr, length, PROT_READ) != 0) {
        perror("mprotect"); // Print error if protection change fails
        syscall(SYS_munmap, addr, length); // Clean up
        return 1;
    }
    printf("Memory protection changed to read-only.\n");

    // Uncommenting the next line would cause a segmentation fault:
    // strcpy((char *)addr, "This will fail!");

    // Unmap the memory region
    // munmap syscall: int munmap(void *addr, size_t length);
    if (syscall(SYS_munmap, addr, length) != 0) {
        perror("munmap"); // Print error if unmapping fails
        return 1;
    }
    printf("Memory unmapped successfully.\n");

    return 0;
}