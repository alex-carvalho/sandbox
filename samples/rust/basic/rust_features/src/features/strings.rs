pub fn run() {
    // stored on the program's binary, immutable, and has a fixed size
    let simple = "hello"; // type is &str
    println!("simple string literal: {simple}");

    let mut simple_mut = "hello"; // this is still a string literal, but we can reassign the variable to point to a different string literal
    if simple_mut != "" {
        simple_mut = "hello rust";
    }
    println!("simple mut string literal: {simple_mut}");

    // stored on the heap, mutable, and can grow in size
    let mut greeting = String::from("hello"); // type is String
    greeting.push_str(" rust");
    println!("string demo: {greeting}");

    // any diferente way to create a new string?
    let name = "world".to_string(); // this is a method on the &str type that creates a String
    println!("name: {name}");

}