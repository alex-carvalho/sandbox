var counter = 0
while counter < 3:
  echo "while count: ", counter
  inc(counter)

for j in 1..3:
  echo "for count (inclusive): ", j

for j in 1..<3:
  echo "for count (exclusive): ", j

let nums = [10, 20, 30]
for n in nums:
  echo n

let numsIndex = [10, 20, 30]
for i, n in numsIndex:
  echo "Index: ", i, " Value: ", n


# using iterator to create a custom loop
iterator countTo(n: int): int =
  for i in 0..n:
    yield i

for x in countTo(3):
  echo x