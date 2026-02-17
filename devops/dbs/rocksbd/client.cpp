#include <iostream>
#include <string>
#include "rocksdb/db.h"
#include "rocksdb/options.h"
#include "rocksdb/write_batch.h"

using namespace rocksdb;

int main() {
    DB* db;
    Options options;
    options.create_if_missing = true;
    
    std::cout << "Opening RocksDB...\n";
    Status status = DB::Open(options, "/tmp/rocksdb", &db);
    if (!status.ok()) {
        std::cerr << "Error opening database: " << status.ToString() << std::endl;
        return 1;
    }

    std::cout << "\n1. Writing data...\n";
    db->Put(WriteOptions(), "user:1", "Alice");
    db->Put(WriteOptions(), "user:2", "Bob");
    db->Put(WriteOptions(), "user:3", "Charlie");
    std::cout << "   ✓ Saved 3 users\n";

    std::cout << "\n2. Reading data...\n";
    std::string value;
    for (const auto& key : {"user:1", "user:2", "user:3"}) {
        status = db->Get(ReadOptions(), key, &value);
        if (status.ok()) {
            std::cout << "   " << key << " = " << value << std::endl;
        }
    }

    std::cout << "\n3. Batch write...\n";
    WriteBatch batch;
    batch.Put("product:1", "Laptop");
    batch.Put("product:2", "Mouse");
    status = db->Write(WriteOptions(), &batch);
    std::cout << "   ✓ Saved 2 products\n";

    std::cout << "\n4. All keys in database:\n";
    Iterator* it = db->NewIterator(ReadOptions());
    for (it->SeekToFirst(); it->Valid(); it->Next()) {
        std::cout << "   " << it->key().ToString() << " = " 
                  << it->value().ToString() << std::endl;
    }
    delete it;

    std::cout << "\n5. Deleting user:2...\n";
    db->Delete(WriteOptions(), "user:2");
    std::cout << "   ✓ Deleted\n";

    status = db->Get(ReadOptions(), "user:2", &value);
    std::cout << "   Verification: user:2 exists = " 
              << (status.ok() ? "true" : "false") << std::endl;

    std::cout << "\n✓ Demo complete!\n";

    delete db;
    return 0;
}
