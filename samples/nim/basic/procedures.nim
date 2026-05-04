# Basic procedure
proc greet(name: string): string =
  result = "Hello, " & name

echo greet("World")

proc power(base: int, exp: int = 2): int =
  var res = 1
  for _ in 1..exp:
    res *= base
  return res

echo "3^2 = ", power(3)
echo "3^3 = ", power(3, 3)

echo "2^4 = ", power(exp=4, base=2)
