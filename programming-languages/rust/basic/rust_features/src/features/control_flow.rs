
pub fn run() {
    let number = 3;

    if number == 3 {
        println!("number was three");
    } else if number < 3 {
        println!("condition was less than 3");
    } else {
        println!("condition was greater than 3");
    }

    let result = if number == 3 { 5 } else { 6 };

    println!("The value of result is: {result}");
}
