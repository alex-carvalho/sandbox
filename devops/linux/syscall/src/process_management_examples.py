import os
import sys

def main():
    print(f"Main process PID: {os.getpid()} Parent process PID: {os.getppid()}")

    pid = os.fork() # syscall fork()
    
    if pid < 0:
        print("fork failed", file=sys.stderr)
        sys.exit(1)
    elif pid == 0:
        print(f"Child process (PID: {os.getpid()}, Parent PID: {os.getppid()})")
        
        print("Child executing 'ls -l'...")
        os.execl('/bin/ls', 'ls', '-l') # syscall exec()

        # Child Process Timeline:
        # 1. Running Python code
        # 2. Reaches os.execl('/bin/ls', 'ls', '-l')
        # 3. Process memory is replaced with 'ls' program
        # 4. 'ls' runs and shows directory contents
        # 5. 'ls' finishes and the process exits
        
        # If execl returns, it failed
        print("execl failed", file=sys.stderr)
        sys.exit(1)
    else:
        # Parent process
        print(f"Parent process (PID: {os.getpid()}) waiting for child (PID: {pid})")
        
        pid, status = os.wait() # wait for all child processes finish
        if os.WIFEXITED(status):
            print(f"Child {pid} exited with status {os.WEXITSTATUS(status)}")

if __name__ == "__main__":
    main()