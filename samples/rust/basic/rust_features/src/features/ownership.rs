pub fn run() {
    let message = String::from("ownership moved this value");
    takes_ownership(message);
    // println!("message after ownership moved: {message}"); // This will cause a compile error

    let message2 =  String::from("my_string");
    takes_ownership(message2.clone());
    println!("original message2: {message2}");
    println!("memory address: {:p}", &message2);
}

fn takes_ownership(value: String) {
    println!("ownership demo: {value}");
    println!("memory address: {:p}", &value);
}
