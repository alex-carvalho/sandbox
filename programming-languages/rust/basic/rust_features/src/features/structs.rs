use std::fmt;


fn struct_user() {
    struct User {
        active: bool,
        username: String,
        email: String,
        sign_in_count: u64,
    }

    // like toString() on java
    impl fmt::Display for User {
        fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
            write!(f, "User {{ active: {}, username: {}, email: {}, sign_in_count: {} }}\n", self.active, self.username, self.email, self.sign_in_count)
        }
    }

    let mut user = User {
        active: true,
        username: String::from("someusername123"),
        email: String::from("someone@example.com"),
        sign_in_count: 1,
    };

    print!("user: {}", user);


    user.email = String::from("another@email.com");
    print!("user email changed: {}", user);
}

fn struct_tuple() {
    struct Point(i32, i32, i32);

    let origin = Point(0, 0, 0);
    println!("origin: ({}, {}, {})", origin.0, origin.1, origin.2);
}

fn struct_vs_tuple_vs_parameters() {
    struct Rectangle {
        width: u32,
        height: u32,
    }

    fn area_simple(width: u32, height: u32) -> u32 {
        width * height
    }

    fn area_tuple(dimensions: (u32, u32)) -> u32 {
        dimensions.0 * dimensions.1
    }

    fn area_struct(rectangle: &Rectangle) -> u32 {
        rectangle.width * rectangle.height
    }

    println!("rect1 area: {}", area_simple(30, 50));
    println!("rect1 area: {}", area_tuple((30, 50)));
    println!("rect2 area: {}", area_struct(&Rectangle { width: 30, height: 50 }));
}

pub fn run() {
    struct_user();
    struct_tuple();
    struct_vs_tuple_vs_parameters();
}

