#[derive(Debug)] // it enables the {:?} format specifier for the enum and prints the debug representation
enum IpAddr {
        V4(u8, u8, u8, u8),
        V6(String),
}

enum Coin {
    Penny,
    Nickel,
    Dime,
    Quarter,
}

fn value_in_cents(coin: Coin) -> u8 {
    match coin {
        Coin::Penny => 1,
        Coin::Nickel => 5,
        Coin::Dime => 10,
        Coin::Quarter => 25,
    }
}


fn match_math_operations(operation: char, n1: i32, n2: i32) -> i32 {
    return match operation {
        '+' => n1 + n2,
        '-' => n1 - n2,
        '*' => n1 * n2,
        '/' => n1 / n2,
        _ => panic!("Unsupported operation: {}", operation),
    }
}

fn match_call_function(value: i32) {
    match value {
        1 => println!("one"),
        3 => println!("three"),
        5 => println!("five"),
        7 => println!("seven"),
        _ => println!("other"), // catch all pattern
    }
}

pub fn run() {
    let home = IpAddr::V4(127, 0, 0, 1);
    let loopback = IpAddr::V6(String::from("::1"));
    
    for (name, addr) in &[("home", &home), ("loopback", &loopback)] {
        match addr {
            IpAddr::V4(a, b, c, d) => println!("{} is an IPv4 address: {}.{}.{}.{}", name, a, b, c, d),
            IpAddr::V6(addr) => println!("{} is an IPv6 address: {}", name, addr),
        }
    }

    println!("Value in cents: {}", value_in_cents(Coin::Nickel));
    println!("Value in pennies: {}", value_in_cents(Coin::Penny));
    println!("Value in dimes: {}", value_in_cents(Coin::Dime));
    println!("Value in quarters: {}", value_in_cents(Coin::Quarter));

    println!("Math operation result: {}", match_math_operations('+', 5, 3));
    match_call_function(3);
}
