
package com.ac;

public class MainTrie {


    static void main() {
        Trie trie = new Trie();
        trie.insert("apple");
        trie.insert("app");
        trie.insert("airbus");
        trie.insert("access");
        trie.insert("air");
        System.out.println(trie);

        System.out.println(trie.findStartsWith("app"));
    }

}
