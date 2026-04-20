# Trie Implementation: Nim vs Java

This document provides a line-by-line explanation of the Nim Trie implementation (`trie.nim`), comparing it with how you would write the equivalent code in Java.

## Type Definitions

```nim
1: type
2:   TrieNode* = ref object
3:     children: array[26, TrieNode]
4:     isEndOfWord: bool
5: 
6:   Trie* = object
7:     root: TrieNode
```

*   **`TrieNode`**: The `type` keyword starts a block of type definitions. `TrieNode*` defines an exported (the `*` means `public`) reference type (`ref object`, which is managed by the garbage collector, just like a standard Java class). It contains a fixed-size array of 26 `children` and a boolean `isEndOfWord`.
    *   **Java Equivalent**:
        ```java
        public class TrieNode {
            TrieNode[] children = new TrieNode[26];
            boolean isEndOfWord = false;
        }
        ```
*   **`Trie`**: `Trie* = object` defines an exported value type (struct). It holds a single reference to the `root` node.
    *   **Java Equivalent**:
        ```java
        public class Trie {
            private TrieNode root;
            // ... methods ...
        }
        ```

## Constructors & Getters

```nim
10: proc newTrieNode(): TrieNode =
11:   TrieNode(
12:     isEndOfWord: false
13:   )
```

*   **`newTrieNode`**: A private `proc` (procedure/method without `*`) that acts as a constructor for `TrieNode`. It creates and returns a `TrieNode` with `isEndOfWord` set to `false`. In Nim, object initialization `TrieNode(...)` implicitly sets all omitted fields (like the `children` array) to their default zero-values (which is `nil` for references).
    *   **Java Equivalent**:
        ```java
        private TrieNode newTrieNode() {
            return new TrieNode(); // isEndOfWord is false by default
        }
        ```

```nim
16: proc newTrie*(): Trie =
17:   result.root = newTrieNode()
```

*   **`newTrie`**: An exported constructor for `Trie`. Nim automatically provides an implicit variable named `result` in procedures that return a value. Here, we assign the root node and implicitly return it.
    *   **Java Equivalent**:
        ```java
        public Trie() {
            this.root = new TrieNode();
        }
        ```

```nim
20: proc getRoot*(t: Trie): TrieNode =
21:   t.root
```

*   **`getRoot`**: A simple getter method. Nim allows omitting the `return` keyword for the last expression.
    *   **Java Equivalent**: `public TrieNode getRoot() { return this.root; }`

## Insertion

```nim
24: proc insert*(t: Trie, word: string) =
25:   var current = t.root
```

*   **`insert`**: Defines the `insert` procedure taking a `Trie` and a `string`. `var current` declares a mutable variable starting at the root. Notice that `t` is passed, Nim uses Uniform Function Call Syntax (UFCS), meaning `t.insert("hello")` is syntactically translated to `insert(t, "hello")`.
    *   **Java Equivalent**: `public void insert(String word) { TrieNode current = this.root; ... }`

```nim
27:   for c in word:
28:     let index = ord(c) - ord('a')
```

*   **Loop**: Iterates over each character `c` in the word. `let index` defines an immutable variable. `ord()` gets the ASCII integer value of the character.
    *   **Java Equivalent**:
        ```java
        for (char c : word.toCharArray()) {
            int index = c - 'a';
        ```

```nim
30:     if current.children[index] == nil:
31:       current.children[index] = newTrieNode()
32: 
33:     current = current.children[index]
```

*   **Child Check**: Checks if the child at that index is `nil` (null). If it is, it instantiates a new node. Then it moves the `current` pointer down to that child.
    *   **Java Equivalent**:
        ```java
            if (current.children[index] == null) {
                current.children[index] = new TrieNode();
            }
            current = current.children[index];
        }
        ```

```nim
35:   current.isEndOfWord = true
```

*   **End of Word**: After the loop finishes traversing the word, it marks the last node as the end of a valid word.
    *   **Java Equivalent**: `current.isEndOfWord = true;`

## Searching

```nim
38: proc search*(t: Trie, word: string): bool =
39:   var current = t.root
40:   for c in word:
41:     let index = ord(c) - ord('a')
42:     if current.children[index] == nil:
43:       return false
44:     current = current.children[index]
45:   return current.isEndOfWord
```

*   **`search`**: Similar traversal as `insert`. If at any point it encounters a `nil` child while reading the characters, the word doesn't exist (`return false`). If it finishes the loop, it returns `isEndOfWord` to ensure it's actually a complete word and not just a prefix.
    *   **Java Equivalent**:
        ```java
        public boolean search(String word) {
            TrieNode current = this.root;
            for (char c : word.toCharArray()) {
                int index = c - 'a';
                if (current.children[index] == null) return false;
                current = current.children[index];
            }
            return current.isEndOfWord;
        }
        ```

```nim
52: proc startsWith*(t: Trie, prefix: string): bool =
53:   var current = t.root
54:   for c in prefix:
55:     let index = ord(c) - ord('a')
56:     if current.children[index] == nil:
57:       return false
58:     current = current.children[index]
59:   return true
```

*   **`startsWith`**: Exact same traversal logic as `search`, but here, if the loop finishes successfully without hitting `nil`, we know the prefix exists, so it simply returns `true`.
    *   **Java Equivalent**: Same as Java `search()`, but ends with `return true;`.

## Auto-completion (Find Starts With)

```nim
66: proc findStartsWith*(t: Trie, prefix: string): seq[string] =
... (traversal logic same as above) ...
72:     if current.children[index] == nil:
73:       return @[]
```

*   **`findStartsWith`**: `seq[string]` is Nim's dynamically resizable array (like `ArrayList<String>` or `List<String>` in Java). `return @[]` returns an empty sequence if the prefix doesn't exist.

```nim
77:   result = @[]
78:   var currentPrefix = prefix
79:   collectWords(current, currentPrefix, result)
```

*   **Initialization**: Initializes the `result` sequence. `var currentPrefix` creates a mutable copy of the string (unlike Java, Nim strings are mutable value types, like a `StringBuilder`). Then it calls the recursive `collectWords` helper.
    *   **Java Equivalent**:
        ```java
        List<String> result = new ArrayList<>();
        StringBuilder currentPrefix = new StringBuilder(prefix);
        collectWords(current, currentPrefix, result);
        return result;
        ```

```nim
81: proc collectWords(node: TrieNode, current: var string, results: var seq[string]) =
```

*   **`collectWords`**: Defines the helper procedure. The `var string` and `var seq[string]` parameters are very important: they mean these arguments are **passed by reference** and can be mutated by the function.
    *   **Java Equivalent**: `private void collectWords(TrieNode node, StringBuilder current, List<String> results)`

```nim
82:   if node.isEndOfWord:
83:     results.add(current)
```

*   **Add Match**: If the current node is a valid word, append the `current` string to the `results` sequence.
    *   **Java Equivalent**: `if (node.isEndOfWord) results.add(current.toString());`

```nim
85:   for i in 0..25:
86:     if node.children[i] != nil:
87:       let c = char(ord('a') + i)
88:       current.add(c)
89:       collectWords(node.children[i], current, results)
90:       current.setLen(current.len - 1)
```

*   **Backtracking Loop**: Iterates through all 26 possible children (`0..25` is an inclusive range). If a child exists:
    *   Calculates the character `c` from the index.
    *   Appends it to the mutable `current` string (`current.add(c)`).
    *   Recursively calls `collectWords`.
    *   **Backtracking step**: `current.setLen(current.len - 1)` removes the last character we just added so we can test the next sibling character.
    *   **Java Equivalent**:
        ```java
        for (int i = 0; i < 26; i++) {
            if (node.children[i] != null) {
                char c = (char) ('a' + i);
                current.append(c);
                collectWords(node.children[i], current, results);
                current.deleteCharAt(current.length() - 1); // backtrack
            }
        }
        ```

## Overloading ToString

```nim
93: proc `$`*(t: Trie): string =
94:   var sb = ""
95:   buildString(t.root, sb, "<root>", "", "")
96:   result = sb
```

*   **`$` Operator**: In Nim, the `$` operator is used to convert types to strings (used by `echo` and string concatenation). Enclosing it in backticks `` `$` `` allows you to define this operator for your custom `Trie` type. This is exactly equivalent to overriding `toString()` in Java. It initializes a mutable string `sb` and passes it to a tree-building helper.
    *   **Java Equivalent**:
        ```java
        @Override
        public String toString() {
            StringBuilder sb = new StringBuilder();
            buildString(this.root, sb, "<root>", "", "");
            return sb.toString();
        }
        ```

## Deletion

```nim
proc isEmpty(node: TrieNode): bool =
  for child in node.children:
    if child != nil: return false
  return true

proc deleteNode(node: TrieNode, word: string, depth: int): bool =
  if node == nil:
    return false

  if depth == word.len:
    if not node.isEndOfWord:
      return false
    node.isEndOfWord = false
    return isEmpty(node)

  let index = ord(word[depth]) - ord('a')
  if deleteNode(node.children[index], word, depth + 1):
    node.children[index] = nil
    return not node.isEndOfWord and isEmpty(node)

  return false

proc delete*(t: Trie, word: string) =
  discard deleteNode(t.root, word, 0)
```

*   **`isEmpty`**: An internal helper that iterates over the children array to determine if a node has any remaining paths.
    *   **Java Equivalent**:
        ```java
        private boolean isEmpty(TrieNode node) {
            for (TrieNode child : node.children) {
                if (child != null) return false;
            }
            return true;
        }
        ```
*   **`deleteNode`**: A recursive method that walks down to the end of the word, unmarks `isEndOfWord`, and as the recursion unwinds, checks if the current node has become empty (no other words share it as a prefix). If so, it nils out the reference in the parent's `children` array.
    *   **Java Equivalent**:
        ```java
        private boolean deleteNode(TrieNode node, String word, int depth) {
            if (node == null) return false;

            if (depth == word.length()) {
                if (!node.isEndOfWord) return false;
                node.isEndOfWord = false;
                return isEmpty(node);
            }

            int index = word.charAt(depth) - 'a';
            if (deleteNode(node.children[index], word, depth + 1)) {
                node.children[index] = null;
                return !node.isEndOfWord && isEmpty(node);
            }
            return false;
        }
        ```
*   **`delete*`**: The exported procedure that initiates the deletion process, ignoring the final boolean result with `discard` (since Nim requires handling return values explicitly).
    *   **Java Equivalent**:
        ```java
        public void delete(String word) {
            deleteNode(this.root, word, 0);
        }
        ```
