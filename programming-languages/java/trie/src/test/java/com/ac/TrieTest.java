
package com.ac;

import org.junit.jupiter.api.Test;
import static org.junit.jupiter.api.Assertions.*;
import java.util.List;

public class TrieTest {

    @Test
    public void testTrie() {
        Trie trie = new Trie();
        trie.insert("apple");
        assertTrue(trie.search("apple"));
        assertFalse(trie.search("app"));
        assertTrue(trie.startsWith("app"));
        trie.insert("app");
        assertTrue(trie.search("app"));
    }

    @Test
    public void testFindStartsWith() {
        Trie trie = new Trie();
        trie.insert("apple");
        trie.insert("app");
        trie.insert("airbus");
        trie.insert("air");
        trie.insert("bat");

        List<String> results = trie.findStartsWith("ap");
        assertEquals(2, results.size());
        assertTrue(results.contains("app"));
        assertTrue(results.contains("apple"));

        List<String> airResults = trie.findStartsWith("air");
        assertEquals(2, airResults.size());
        assertTrue(airResults.contains("air"));
        assertTrue(airResults.contains("airbus"));

        assertTrue(trie.findStartsWith("xyz").isEmpty());
    }

}