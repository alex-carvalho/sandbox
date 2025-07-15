#include <stdio.h>      // Used for printf(), perror()
#include <unistd.h>     // Used for syscall(), fork(), read(), write(), close()
#include <sys/syscall.h> // Used for SYS_* syscall numbers
#include <sys/ipc.h>    // Used for IPC_CREAT, IPC_PRIVATE
#include <sys/types.h>  // Used for key_t, pid_t
#include <sys/msg.h>    // Used for message queue flags and structures
#include <sys/sem.h>    // Used for semaphore flags and structures
#include <sys/shm.h>    // Used for shared memory flags and structures
#include <string.h>     // Used for strcpy(), strlen()
#include <stdlib.h>     // Used for exit()
#include <sys/wait.h>   // Used for wait()

// Message structure for message queue
struct msgbuf {
    long mtype;         // Message type, must be > 0
    char mtext[100];    // Message data
};

int main() {
    // --- PIPE EXAMPLE: Parent writes, child reads ---
    int pipefd[2];
    // pipe syscall: int pipe(int pipefd[2]);  pipefd[0]: read end, pipefd[1]: write end
    if (syscall(SYS_pipe, pipefd) == -1) {
        perror("error pipe creation fails");
    } else {
        pid_t pid = fork(); // Create a child process
        if (pid == 0) { // Child process
            close(pipefd[1]); // Close write end in child
            char buf[32];
            // Read data from the pipe
            read(pipefd[0], buf, sizeof(buf));
            printf("[Pipe] Child received: %s\n", buf);
            close(pipefd[0]); // Close read end
            exit(0);
        } else { // Parent process
            close(pipefd[0]); // Close read end in parent
            const char *msg = "Hello from parent!";
            // Write data to the pipe
            write(pipefd[1], msg, strlen(msg) + 1); // +1 to include null terminator
            close(pipefd[1]); // Close write end
            wait(NULL); // Wait for child to finish
        }
    }

    // --- MESSAGE QUEUE EXAMPLE: Send and receive a message ---
    // msgget syscall: int msgget(key_t key, int msgflg);
    // IPC_PRIVATE: create a new, unique message queue
    // 0666: permissions (read and write for everyone)
    // IPC_CREAT: create the queue if it does not exist
    int msqid = syscall(SYS_msgget, IPC_PRIVATE, 0666 | IPC_CREAT);
    if (msqid == -1) {
        perror("error msgget queue creation fails");
    } else {
        struct msgbuf msg = {1, "Hello via message queue!"};
        // msgsnd syscall: int msgsnd(int msqid, const void *msgp, size_t msgsz, int msgflg);
        // Send a message of type 1
        if (syscall(SYS_msgsnd, msqid, &msg, sizeof(msg.mtext), 0) == -1) {
            perror("error msgsnd fails");
        } else {
            struct msgbuf rcv;
            // msgrcv syscall: int msgrcv(int msqid, void *msgp, size_t msgsz, long msgtyp, int msgflg);
            // Receive a message of type 1
            if (syscall(SYS_msgrcv, msqid, &rcv, sizeof(rcv.mtext), 1, 0) == -1) {
                perror("error msgrcv fails");
            } else {
                printf("[MsgQueue] Received: %s\n", rcv.mtext);
            }
        }
        // msgctl syscall: int msgctl(int msqid, int cmd, struct msqid_ds *buf);
        // Remove the message queue
        syscall(SYS_msgctl, msqid, IPC_RMID, NULL);
    }

    // --- SEMAPHORE EXAMPLE: Lock and unlock ---
    // semget syscall: int semget(key_t key, int nsems, int semflg);
    // Create a set with 1 semaphore
    int semid = syscall(SYS_semget, IPC_PRIVATE, 1, 0666 | IPC_CREAT);
    if (semid == -1) {
        perror("error semget semaphore set creation fails");
    } else {
        struct sembuf lock = {0, -1, 0};   // P operation: decrement semaphore to lock
        struct sembuf unlock = {0, 1, 0};  // V operation: increment semaphore to unlock
        // semctl syscall: int semctl(int semid, int semnum, int cmd, ...);
        // Initialize semaphore to 1 (unlocked)
        union semun { int val; } arg;
        arg.val = 1;
        syscall(SYS_semctl, semid, 0, SETVAL, arg);
        // semop syscall: int semop(int semid, struct sembuf *sops, size_t nsops);
        // Lock (decrement)
        syscall(SYS_semop, semid, &lock, 1);
        printf("[Semaphore] Locked critical section\n");
        // Unlock (increment)
        syscall(SYS_semop, semid, &unlock, 1);
        printf("[Semaphore] Unlocked critical section\n");
        // Remove semaphore set
        syscall(SYS_semctl, semid, 0, IPC_RMID);
    }

    // --- SHARED MEMORY EXAMPLE: Parent writes, child reads ---
    // shmget syscall: int shmget(key_t key, size_t size, int shmflg);
    // IPC_PRIVATE: create a new, unique shared memory segment
    // 4096: size in bytes (one memory page)
    // 0666: permissions (read and write for everyone)
    // IPC_CREAT: create the segment if it does not exist
    int shmid = syscall(SYS_shmget, IPC_PRIVATE, 4096, 0666 | IPC_CREAT);
    if (shmid == -1) {
        perror("error shmget shared memory creation fails");
    } else {
        pid_t pid = fork(); // Create a child process
        if (pid == 0) { // Child process
            // shmat syscall: void *shmat(int shmid, const void *shmaddr, int shmflg);
            // Attach shared memory segment to address space
            char *shmaddr = (char *)syscall(SYS_shmat, shmid, NULL, 0);
            printf("[SharedMemory] Child read: %s\n", shmaddr);
            // shmdt syscall: int shmdt(const void *shmaddr);
            // Detach shared memory
            syscall(SYS_shmdt, shmaddr);
            exit(0);
        } else { // Parent process
            char *shmaddr = (char *)syscall(SYS_shmat, shmid, NULL, 0);
            // Write a string to shared memory
            strcpy(shmaddr, "Hello from shared memory!");
            syscall(SYS_shmdt, shmaddr); // Detach
            wait(NULL); // Wait for child to finish
            // shmctl syscall: int shmctl(int shmid, int cmd, struct shmid_ds *buf);
            // Remove the shared memory segment
            syscall(SYS_shmctl, shmid, IPC_RMID, NULL);
        }
    }

    return 0;
}