import os
import sys

def main():
    print(f"Main process PID: {os.getpid()}")
    print(f"Parent process PID: {os.getppid()}")

    pid = os.fork() # Create child process
    
    if pid < 0:
        print("fork failed", file=sys.stderr)
        sys.exit(1)
    elif pid == 0:
        print(f"Child process (PID: {os.getpid()}, Parent PID: {os.getppid()})")
        
        print("Child executing 'ls -l'...")
        os.execl('/bin/ls', 'ls', '-l')
        
        # If execl returns, it failed
        print("execl failed", file=sys.stderr)
        sys.exit(1)
    else:
        # Parent process
        print(f"Parent process (PID: {os.getpid()}) waiting for child (PID: {pid})")
        
        pid, status = os.wait()
        if os.WIFEXITED(status):
            print(f"Child {pid} exited with status {os.WEXITSTATUS(status)}")

if __name__ == "__main__":
    main()