import asyncdispatch, os

echo "--- Nim Concurrency (Async/Await) ---"

proc asyncTask(id: int, delayMs: int): Future[void] {.async.} =
  echo "Task ", id, " started."
  await sleepAsync(delayMs)
  echo "Task ", id, " finished."

proc main() {.async.} =
  echo "Starting main async execution"
  # Run concurrently
  let t1 = asyncTask(1, 50)
  let t2 = asyncTask(2, 30)
  
  await all(t1, t2)
  echo "All tasks completed."

waitFor main()
