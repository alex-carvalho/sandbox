
package com.ac;

import java.util.ArrayList;
import java.util.Collections;
import java.util.List;

public class Trie {

    static class TrieNode {
        public TrieNode[] children;
        public boolean isEndOfWord;
        public char code; // only for debug purposes

        public TrieNode(char code) {
            children = new TrieNode[26];
            isEndOfWord = false;
            this.code = code;
        }

        @Override
        public String toString() {
            return "" + code;
        }
    }

    private final TrieNode root;

    public Trie() {
        root = new TrieNode('0');
    }

    public TrieNode getRoot() {
        return root;
    }

    public void insert(String word) {
        TrieNode currentNode = root;
        for (char c : word.toCharArray()) {
            int index = c - 'a';
            if (currentNode.children[index] == null) {
                currentNode.children[index] = new TrieNode(c);
            }
            currentNode = currentNode.children[index];
        }
        currentNode.isEndOfWord = true;
    }

    public boolean search(String word) {
        TrieNode currentNode = root;
        for (char c : word.toCharArray()) {
            int index = c - 'a';
            if (currentNode.children[index] == null) {
                return false;
            }
            currentNode = currentNode.children[index];
        }
        return currentNode.isEndOfWord;
    }

    public boolean startsWith(String prefix) {
        TrieNode current = root;
        for (int i = 0; i < prefix.length(); i++) {
            char c = prefix.charAt(i);
            int index = c - 'a';
            if (current.children[index] == null) {
                return false;
            }

            current = current.children[index];
        }
        return true;
    }

    public List<String> findStartsWith(String prefix) {
        TrieNode current = root;
        for (char c : prefix.toCharArray()) {
            int index = c - 'a';
            if (current.children[index] == null) {
                return Collections.emptyList();
            }
            current = current.children[index];
        }
        List<String> results = new ArrayList<>();
        collectWords(current, new StringBuilder(prefix), results);
        return results;
    }

    private void collectWords(TrieNode node, StringBuilder current, List<String> results) {
        if (node.isEndOfWord) {
            results.add(current.toString());
        }
        for (int i = 0; i < 26; i++) {
            if (node.children[i] != null) {
                current.append((char) ('a' + i));
                collectWords(node.children[i], current, results);
                current.deleteCharAt(current.length() - 1);
            }
        }
    }

    @Override
    public String toString() {
        StringBuilder sb = new StringBuilder();
        buildString(root, sb, "", "");
        return sb.toString();
    }

    private void buildString(TrieNode node, StringBuilder sb, String prefix, String childPrefix) {
        sb.append(prefix).append(node.code).append(node.isEndOfWord ? "*" : "").append("\n");
        for (int k = 0; k < node.children.length; k++) {
            var child = node.children[k];
            if (child != null) {
                buildString(child, sb, childPrefix + "├── ", childPrefix + "│   ");
            }
        }
    }

}
