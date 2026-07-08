
const age = 30

if age > 50:
  echo "Greater than 50"
elif age == 42:
  echo "Exactly 42"
else:
  echo "Less than 42"


let option = "b"
case option
of "a": echo "A selected"
of "b": echo "B selected"
else: echo "Other selected"


let score = 85
case score:
of 0..49:
  echo "Fail"
of 50..79:
  echo "Pass"
of 80..100:
  echo "Excellent"
else:
  echo "Invalid"


let x = 10
let result = if x > 5: "big" else: "small"

echo result

when defined(windows):
  echo "Running on Windows"
else:
  echo "Not Windows"