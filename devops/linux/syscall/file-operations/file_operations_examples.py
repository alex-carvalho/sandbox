import os
import sys
import mmap

filename = "example.txt"
msg = b"Hello, file operations with syscalls!\n"

# Open file for writing (os.O_CREAT | os.O_WRONLY | os.O_TRUNC)
fd = os.open(filename, os.O_CREAT | os.O_WRONLY | os.O_TRUNC, 0o644)
if fd < 0:
    print("open failed")
    sys.exit(1)
print(f"File '{filename}' opened with file descriptor {fd}")

# Write to the file
written = os.write(fd, msg)
if written < 0:
    print("write failed")
    os.close(fd)
    sys.exit(1)
print(f"Wrote {written} bytes to file.")

# Move file offset to position 7 (like lseek)
offset = os.lseek(fd, 7, os.SEEK_SET)
if offset == -1:
    print("lseek failed")
    os.close(fd)
    sys.exit(1)
print(f"File offset repositioned to {offset} (position 7).")

# Write again at position 7 (overwrite)
overwrite_msg = b"Overwritten! "
written = os.write(fd, overwrite_msg)
if written < 0:
    print("write failed (overwrite)")
    os.close(fd)
    sys.exit(1)
print(f"Overwrote {written} bytes at the 7 position of the file.")

# Close file after writing
os.close(fd)
print("File closed after writing.")

# Reopen file for reading
fd = os.open(filename, os.O_RDONLY)
if fd < 0:
    print("open failed (read)")
    sys.exit(1)
print(f"File '{filename}' reopened for reading with file descriptor {fd}")

# Read from the file
# 127 bytes is the maximum size for a single read operation
buffer = os.read(fd, 127)
if buffer is None:
    print("read failed")
    os.close(fd)
    sys.exit(1)
print(f"Read {len(buffer)} bytes: {buffer.decode()}")

# Close file after reading
os.close(fd)
print("File closed after reading.")

# Remove the file (unlink)
try:
    os.unlink(filename)
    print(f"File '{filename}' deleted from filesystem.")
except Exception as e:
    print("unlink failed:")