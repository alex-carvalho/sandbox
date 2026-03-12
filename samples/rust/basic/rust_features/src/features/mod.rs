pub mod loops;
pub mod ownership;
pub mod strings;

pub fn run_examples() {
    loops::run();
    ownership::run();
    strings::run();
}
