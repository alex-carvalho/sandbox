pub fn run() {
    for i in 1..=3 {
        println!("for i: {i}");
    }

    let mut count = 0;
    while count < 3 {
        count += 1;
        println!("while count: {count}");
    }

    let mut count = 0;

    loop {
        count += 1;

        println!("loop count: {count}");

        if count >= 3 {
            break;
        }
    }

    let result = return_value_from_loop();
    println!("The result from the loop is: {result}");

}

fn return_value_from_loop() -> i32 {
    let mut count = 0;

    let result = loop {
        count += 1;

        if count >= 10 {
            break count * 2;
        }
    };

    result
}