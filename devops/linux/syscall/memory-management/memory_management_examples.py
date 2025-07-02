import mmap
import ctypes
import os

# Allocate one memory page (usually 4096 bytes)
length = 4096

# Create an anonymous memory map (similar to mmap syscall in C)
mm = mmap.mmap(-1, length, prot=mmap.PROT_READ | mmap.PROT_WRITE)

print(f"Memory mapped at address {ctypes.addressof(ctypes.c_char.from_buffer(mm)):#x}")

# Write a string to the allocated memory
msg = b"Hello, memory management with syscalls!"
mm[:len(msg)] = msg
print(f"Content: {mm[:len(msg)].decode()}")

# Change memory protection to read-only
# This requires using mprotect via ctypes
libc = ctypes.CDLL("libc.so.6")
addr = ctypes.addressof(ctypes.c_char.from_buffer(mm))
if libc.mprotect(ctypes.c_void_p(addr), ctypes.c_size_t(length), mmap.PROT_READ) != 0:
    print("mprotect failed")
    mm.close()
    exit(1)
print("Memory protection changed to read-only.")

# Uncommenting the next line would cause a TypeError in Python (cannot write to read-only buffer)
# mm[:4] = b"Fail"

# Unmap the memory region
mm.close()
print("Memory unmapped successfully.")