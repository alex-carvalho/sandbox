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
    /// Add a new task, e.g. `todo add "Buy groceries"`
    Add { task: String },
    /// List all tasks
    List,
    /// Mark a task as complete by ID
    Complete { id: usize },
    /// Delete a task by ID
    Delete { id: usize },
}

fn main() {
    let cli = Cli::parse();
    let mut repo = FileRepository::new("tasks.json");

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
                for task in &tasks {
                    let status = if task.completed { "[x]" } else { "[ ]" };
                    println!("{} {} {}", status, task.id, task.description);
                }
            }
        }
        Commands::Complete { id } => {
            repo.complete(id).expect("Failed to complete task");
            println!("Task {} marked as complete.", id);
        }
        Commands::Delete { id } => {
            repo.delete(id).expect("Failed to delete task");
            println!("Task {} deleted.", id);
        }
    }
}
