import trie

proc main*() =
  var trie = newTrie()

  trie.insert("apple")
  trie.insert("app")
  trie.insert("banana")

  echo trie.search("apple")
  echo trie.search("ap")
  echo trie.startsWith("app")
  echo trie.findStartsWith("app")

  echo trie

when isMainModule:
  main()
