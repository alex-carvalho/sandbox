#include <stdio.h>      // Used for printf(), perror()
#include <unistd.h>     // Used for syscall()
#include <sys/syscall.h> // Used for SYS_mkdir, SYS_rmdir, SYS_link, SYS_symlink, SYS_chmod, SYS_chown
#include <sys/stat.h>   // Used for mode constants
#include <fcntl.h>      // Used for file permissions
#include <string.h>     // Used for strerror()
#include <errno.h>      // Used for errno

int main() {
    // mkdir syscall: int mkdir(const char *pathname, mode_t mode);
    const char *dir_name = "example_dir";
    if (syscall(SYS_mkdir, dir_name, 0755) != 0) {
        perror("mkdir"); // Print error if directory creation fails
    } else {
        printf("Directory '%s' created.\n", dir_name);
    }

    // Create a file inside the directory to demonstrate link and symlink
    char file_name[256];
    snprintf(file_name, sizeof(file_name), "%s/%s", dir_name, "example_file.txt");
    int fd = syscall(SYS_open, file_name, O_CREAT | O_WRONLY, 0644);
    if (fd < 0) {
        perror("open");
    } else {
        write(fd, "Hello, file system syscalls!\n", 29);
        close(fd);
        printf("File '%s' created.\n", file_name);
    }

    // link syscall: int link(const char *oldpath, const char *newpath);
    char hard_link[256];
    snprintf(hard_link, sizeof(hard_link), "%s/%s", dir_name, "example_file_hardlink.txt");
    if (syscall(SYS_link, file_name, hard_link) != 0) {
        perror("link"); // Print error if hard link creation fails
    } else {
        printf("Hard link '%s' created.\n", hard_link);
    }

    // symlink syscall: int symlink(const char *target, const char *linkpath);
    char sym_link[256];
    snprintf(sym_link, sizeof(sym_link), "%s/%s", dir_name, "example_file_symlink.txt");
    if (syscall(SYS_symlink, "example_file.txt", sym_link) != 0) {
        perror("symlink"); // Print error if symbolic link creation fails
    } else {
        printf("Symbolic link '%s' created (points to 'example_file.txt').\n", sym_link);
    }

    // chmod syscall: int chmod(const char *pathname, mode_t mode);
    if (syscall(SYS_chmod, file_name, 0600) != 0) {
        perror("chmod"); // Print error if permission change fails
    } else {
        printf("Permissions of '%s' changed to 0600.\n", file_name);
    }

    // chown syscall: int chown(const char *pathname, uid_t owner, gid_t group);
    // Set owner and group to current user (getuid(), getgid())
    if (syscall(SYS_chown, file_name, getuid(), getgid()) != 0) {
        perror("chown"); // Print error if ownership change fails
    } else {
        printf("Ownership of '%s' changed to current user.\n", file_name);
    }

    // Remove all files and links inside the directory before removing the directory itself
    // Unlink (remove) the regular file
    if (syscall(SYS_unlink, file_name) != 0) {
        perror("unlink file"); // Print error if file removal fails
    } else {
        printf("File '%s' removed.\n", file_name);
    }

    // Unlink (remove) the hard link
    if (syscall(SYS_unlink, hard_link) != 0) {
        perror("unlink hard link"); // Print error if hard link removal fails
    } else {
        printf("Hard link '%s' removed.\n", hard_link);
    }

    // Unlink (remove) the symbolic link
    if (syscall(SYS_unlink, sym_link) != 0) {
        perror("unlink symbolic link"); // Print error if symbolic link removal fails
    } else {
        printf("Symbolic link '%s' removed.\n", sym_link);
    }

    // Now remove the directory itself
    if (syscall(SYS_rmdir, dir_name) != 0) {
        perror("rmdir"); // Print error if directory removal fails
    } else {
        printf("Directory '%s' removed.\n", dir_name);
    }

    return 0;
}