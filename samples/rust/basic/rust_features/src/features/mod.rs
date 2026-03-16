pub mod data_types;
pub mod loops;
pub mod ownership;
pub mod strings;
pub mod functions;

pub fn run_examples() {
    data_types::run();
    loops::run();
    ownership::run();
    strings::run();
    functions::run();
}
