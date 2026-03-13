pub mod data_types;
pub mod loops;
pub mod ownership;
pub mod strings;

pub fn run_examples() {
    data_types::run();
    loops::run();
    ownership::run();
    strings::run();
}
