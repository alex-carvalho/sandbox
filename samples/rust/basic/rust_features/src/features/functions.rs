pub fn run() {
    another_function();
    parameter_function(42);
    print_labeled_measurement(5, 'h');
    print!("Return value from return_value function: {}", return_value());
}

fn another_function() {
    println!("Another function.");
}

fn parameter_function(x: i32) {
    println!("The value of x is: {x}");
}

fn print_labeled_measurement(value: i32, unit_label: char) {
    println!("The measurement is: {value}{unit_label}");
}

fn return_value() -> i32 {
    5
}