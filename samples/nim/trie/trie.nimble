# Package

version       = "0.1.0"
author        = "Alex Carvalho"
description   = "A new awesome nimble package"
license       = "MIT"
srcDir        = "src"
bin           = @["maintrie"]


# Dependencies

requires "nim >= 2.2.8"

task test, "Run the trie test suite":
  exec "nim c -r --path:src --nimcache:build/nimcache tests/test_trie.nim"
