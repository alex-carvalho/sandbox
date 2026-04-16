use std::fs::{self, OpenOptions};
use std::io::Write;

pub trait TaskRepository {
    fn add(&mut self, task: &str) -> Result<(), String>;
    fn list(&self) -> Result<Vec<String>, String>;
}

pub struct FileRepository {
    path: String,
}

impl FileRepository {
    pub fn new(path: &str) -> Self {
        Self { path: path.to_string() }
    }
}

impl TaskRepository for FileRepository {
    fn add(&mut self, task: &str) -> Result<(), String> {
        let mut file = OpenOptions::new()
            .create(true)
            .append(true)
            .open(&self.path)
            .map_err(|e| e.to_string())?;
        writeln!(file, "{}", task).map_err(|e| e.to_string())
    }

    fn list(&self) -> Result<Vec<String>, String> {
        let tasks = fs::read_to_string(&self.path)
            .unwrap_or_default()
            .lines()
            .filter(|l| !l.is_empty())
            .map(|l| l.to_string())
            .collect();
        Ok(tasks)
    }
}
