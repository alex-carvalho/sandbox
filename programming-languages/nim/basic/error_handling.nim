echo "--- Nim Error Handling ---"

type
  MyCustomError = object of CatchableError

proc riskyOperation(fail: bool) =
  echo "Starting risky operation..."
  defer:
    echo "It will run in the end even if an error occurs."
  if fail:
    echo "The value is true, raising an error."
    raise newException(MyCustomError, "My custom error occurred!")
  echo "Operation succeeded."

try:
  riskyOperation(true)
except MyCustomError as e:
  echo "Caught custom error: ", e.msg
except Exception as e:
  echo "Caught general error: ", e.msg
finally:
  echo "Cleanup: this always runs."


proc useDefer() =
  echo "Starting scope"
  defer: echo "Deferred cleanup"
  echo "Ending scope"

useDefer()
