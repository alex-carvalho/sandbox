pub fn run() {
    let age: i32 = 18; // types can be inferred let age = 18
    let height: f64 = 1.75;
    let is_cool: bool = true;
    let letter: char = 'R';
    let tuple: (i32, f64, char) = (500, 6.4, 'z');
    let array: [i32; 3] = [1, 2, 3];


    println!("integer: {age}");
    println!("float: {height}");
    println!("boolean: {is_cool}");
    println!("char: {letter}");
    println!("tuple second value: {}", tuple.1);
    println!("array first value: {}", array[0]);
}
