type
  TrieNode* = ref object
    children: array[26, TrieNode]
    isEndOfWord: bool

  Trie* = object
    root: TrieNode


proc newTrieNode(): TrieNode =
  TrieNode(
    isEndOfWord: false
  )


proc newTrie*(): Trie =
  result.root = newTrieNode()


proc getRoot*(t: Trie): TrieNode =
  t.root


proc insert*(t: Trie, word: string) =
  var current = t.root

  for c in word:
    let index = ord(c) - ord('a')

    if current.children[index] == nil:
      current.children[index] = newTrieNode()

    current = current.children[index]

  current.isEndOfWord = true


proc search*(t: Trie, word: string): bool =
  var current = t.root

  for c in word:
    let index = ord(c) - ord('a')

    if current.children[index] == nil:
      return false

    current = current.children[index]

  return current.isEndOfWord


proc startsWith*(t: Trie, prefix: string): bool =
  var current = t.root

  for c in prefix:
    let index = ord(c) - ord('a')

    if current.children[index] == nil:
      return false

    current = current.children[index]

  return true


proc findStartsWith*(t: Trie, prefix: string): seq[string] =
  var current = t.root

  for c in prefix:
    let index = ord(c) - ord('a')

    if current.children[index] == nil:
      return @[]

    current = current.children[index]

  result = @[]
  var currentPrefix = prefix
  collectWords(current, currentPrefix, result)

proc collectWords(node: TrieNode, current: var string, results: var seq[string]) =
  if node.isEndOfWord:
    results.add(current)

  for i in 0..25:
    if node.children[i] != nil:
      let c = char(ord('a') + i)
      current.add(c)
      collectWords(node.children[i], current, results)
      current.setLen(current.len - 1)


proc `$`*(t: Trie): string =
  var sb = ""
  buildString(t.root, sb, "<root>", "", "")
  result = sb

proc buildString(node: TrieNode, sb: var string, label: string, prefix: string, childPrefix: string) =
  sb.add(prefix & label & (if node.isEndOfWord: "*" else: "") & "\n")

  for i, child in node.children:
    if child != nil:
      let childLabel = $char(ord('a') + i)
      buildString(child, sb, childLabel, childPrefix & "├── ", childPrefix & "│   ")
