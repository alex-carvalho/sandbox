mod storage;

use clap::{Parser, Subcommand};
use storage::{FileRepository, TaskRepository};

#[derive(Parser)]
#[command(name = "todo", about = "A simple todo list CLI")]
struct Cli {
    #[command(subcommand)]
    command: Commands,
}

#[derive(Subcommand)]
enum Commands {
    /// Add new task, e.g. `todo add "Buy groceries"`
    Add { task: String },
    /// List all tasks
    List,
}

fn main() {
    let cli = Cli::parse();
    let mut repo = FileRepository::new("tasks.txt");

    match cli.command {
        Commands::Add { task } => {
            repo.add(&task).expect("Failed to add task");
            println!("Added: {}", task);
        }
        Commands::List => {
            let tasks = repo.list().expect("Failed to load tasks");
            if tasks.is_empty() {
                println!("No tasks.");
            } else {
                for (i, task) in tasks.iter().enumerate() {
                    println!("{}: {}", i + 1, task);
                }
            }
        }
    }
}
