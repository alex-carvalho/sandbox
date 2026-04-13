use clap::{Parser, Subcommand};
use std::process::Command;

#[derive(Parser)]
#[command(name = "simple-cli", about = "A simple CLI example")]
struct Cli {
    #[command(subcommand)]
    command: Commands,
}

#[derive(Subcommand)]
enum Commands {
    /// List directory contents using ls
    Ls {
        /// Directory path to list (defaults to current directory)
        path: Option<String>,
    },
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_ls_no_path() {
        let cli = Cli::parse_from(["simple-cli", "ls"]);
        match cli.command {
            Commands::Ls { path } => assert!(path.is_none()),
        }
    }

    #[test]
    fn test_ls_with_path() {
        let cli = Cli::parse_from(["simple-cli", "ls", "/tmp"]);
        match cli.command {
            Commands::Ls { path } => assert_eq!(path.as_deref(), Some("/tmp")),
        }
    }

    #[test]
    fn test_missing_subcommand_fails() {
        let result = Cli::try_parse_from(["simple-cli"]);
        assert!(result.is_err());
    }
}

fn main() {
    let cli = Cli::parse();

    match cli.command {
        Commands::Ls { path } => {
            let dir = path.as_deref().unwrap_or(".");
            let status = Command::new("ls")
                .arg(dir)
                .status()
                .expect("failed to execute ls");
            std::process::exit(status.code().unwrap_or(1));
        }
    }
}
