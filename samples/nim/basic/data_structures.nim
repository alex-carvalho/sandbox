echo "--- Nim Data Structures ---"

var arr: array[3, int] = [1, 2, 3]
echo "Array: ", arr

var seq1 = @[10, 20, 30]
seq1.add(40)
echo "Sequence: ", seq1

var person = (name: "Alice", age: 30)
echo "Tuple: ", person
echo "Name: ", person.name

var mySet: set[char] = {'a', 'b', 'c'}
mySet.incl('d')
echo "Set: ", mySet
echo "Contains 'b'? ", 'b' in mySet

var str = "Nim is "
str.add("awesome!")
echo "String: ", str
