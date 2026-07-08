pub fn run() {
    let message = String::from("ownership moved this value");
    takes_ownership(message);
    // println!("message after ownership moved: {message}"); // This will cause a compile error

    let message2 =  String::from("my_string");
    takes_ownership(message2.clone());
    println!("original message2: {message2}");
    println!("memory address: {:p}", &message2);

    let message3 = String::from("ownership returned this value");
    let returned_message = return_ownership(message3);
    println!("returned message: {returned_message}");

    let number = 42;
    let second_number = number; // This is a copy, not a move, because i32 implements the Copy trait
    println!("number: {number} second number: {second_number}");

    let s1 = String::from("hello");
    let (s2, len) = calculate_length(s1);
    println!("The length of '{s2}' is {len}.");
}

fn takes_ownership(value: String) {
    println!("ownership demo: {value}");
    println!("memory address: {:p}", &value);
}

fn return_ownership(message: String) -> String {
    print!("returning ownership of message: {message}");
    message
}
fn calculate_length(s: String) -> (String, usize) {
    let length = s.len(); // len() returns the length of a String

    (s, length)
}