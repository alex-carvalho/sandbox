#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
#include <sys/wait.h>

int main() {
    printf("Main process PID: %d\n", getpid());
    printf("Parent process PID: %d\n", getppid());

    pid_t pid = fork(); // Create child process
    
    if (pid < 0) {
        perror("fork failed");
        exit(1);
    }
    else if (pid == 0) {
        printf("Child process (PID: %d, Parent PID: %d)\n", getpid(), getppid());
        
        printf("Child executing 'ls -l'...\n");
        execl("/bin/ls", "ls", "-l", NULL);
        
        // If execl returns, it failed
        perror("execl failed");
        exit(1);
    } 
    else {
        printf("Parent process (PID: %d) waiting for child (PID: %d)\n", getpid(), pid);
        
        int status;
        pid_t wait_pid = wait(&status);
        
        if (WIFEXITED(status)) {
            printf("Child %d exited with status %d\n", wait_pid, WEXITSTATUS(status));
        }
    }

    return 0;
}