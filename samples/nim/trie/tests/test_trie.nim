import std/unittest
import trie

suite "Trie":
  test "testTrie":
    var trie = newTrie()

    trie.insert("apple")

    check trie.search("apple")
    check not trie.search("app")
    check trie.startsWith("app")

    trie.insert("app")

    check trie.search("app")

  test "testFindStartsWith":
    var trie = newTrie()

    trie.insert("apple")
    trie.insert("app")
    trie.insert("airbus")
    trie.insert("air")
    trie.insert("bat")

    let results = trie.findStartsWith("ap")
    check results.len == 2
    check "app" in results
    check "apple" in results

    let airResults = trie.findStartsWith("air")
    check airResults.len == 2
    check "air" in airResults
    check "airbus" in airResults

    check trie.findStartsWith("xyz").len == 0
