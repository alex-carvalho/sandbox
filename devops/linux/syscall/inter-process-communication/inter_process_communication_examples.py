import os
import sys
import multiprocessing
import queue
from multiprocessing import Process, Value, Array
import ctypes

# --- PIPE EXAMPLE: Parent writes, child reads ---
def pipe_example():
    r, w = os.pipe()  # r: read end, w: write end

    pid = os.fork()
    if pid == 0:
        # Child process
        os.close(w)  # Close write end in child
        msg = os.read(r, 32).decode()
        print(f"[Pipe] Child received: {msg}")
        os.close(r)
        sys.exit(0)
    else:
        # Parent process
        os.close(r)  # Close read end in parent
        msg = "Hello from parent!"
        os.write(w, msg.encode())
        os.close(w)
        os.wait()

# --- MESSAGE QUEUE EXAMPLE: Send and receive a message ---
def message_queue_example():
    q = multiprocessing.Queue()
    def sender(q):
        q.put("Hello via message queue!")
    def receiver(q):
        try:
            msg = q.get(timeout=2)
            print(f"[MsgQueue] Received: {msg}")
        except queue.Empty:
            print("No message received.")

    p_send = multiprocessing.Process(target=sender, args=(q,))
    p_recv = multiprocessing.Process(target=receiver, args=(q,))
    p_send.start()
    p_recv.start()
    p_send.join()
    p_recv.join()

# --- SEMAPHORE EXAMPLE: Lock and unlock ---
def semaphore_example():
    sem = multiprocessing.Semaphore(1)  # Initial value 1 (unlocked)
    def critical_section(sem):
        sem.acquire()
        print("[Semaphore] Locked critical section")
        sem.release()
        print("[Semaphore] Unlocked critical section")
    p = multiprocessing.Process(target=critical_section, args=(sem,))
    p.start()
    p.join()

# --- SHARED MEMORY EXAMPLE: Parent writes, child reads ---
def shared_memory_example():
    shm = multiprocessing.Array(ctypes.c_char, 100)
    def child(shm):
        print(f"[SharedMemory] Child read: {bytes(shm[:]).decode().rstrip(chr(0))}")

    p = Process(target=child, args=(shm,))
    msg = b"Hello from shared memory!"
    shm[:len(msg)] = msg
    p.start()
    p.join()

if __name__ == "__main__":
    pipe_example()
    message_queue_example()
    semaphore_example()
    shared_memory_example()