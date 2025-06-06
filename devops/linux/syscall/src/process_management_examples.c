#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
#include <sys/wait.h>

int main() {
    printf("Main process PID: %d Parent process PID: %d\n", getpid(), getppid());

    pid_t pid = fork(); // syscall fork()
    
    if (pid < 0) {
        perror("fork failed");
        exit(1);
    }
    else if (pid == 0) {
        printf(">Child process (PID: %d, Parent PID: %d)\n", getpid(), getppid());
        
        printf(">Child executing 'ls -l'...\n");
        execl("/bin/ls", "ls", "-l", NULL); // syscall exec()

        // Child Process Timeline:
        // 1. Running C code
        // 2. Reaches execl('/bin/ls', 'ls', '-l')
        // 3. Process memory is replaced with 'ls' program
        // 4. 'ls' runs and shows directory contents
        // 5. 'ls' finishes and the process exits
        
        // If execl returns, it failed
        perror("execl failed");
        exit(1);
    } 
    else {
        printf("Parent process (PID: %d) waiting for child (PID: %d)\n", getpid(), pid);
        
        int status;
        pid_t wait_pid = wait(&status); // wait for all child processes finish
        
        if (WIFEXITED(status)) {
            printf("Child %d exited with status %d\n", wait_pid, WEXITSTATUS(status));
        }
    }

    return 0;
}