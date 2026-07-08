use serde::{Deserialize, Serialize};
use std::fs;

#[derive(Serialize, Deserialize, Clone)]
pub struct Task {
    pub id: usize,
    pub description: String,
    pub completed: bool,
}

pub trait TaskRepository {
    fn add(&mut self, description: &str) -> Result<(), String>;
    fn list(&self) -> Result<Vec<Task>, String>;
    fn complete(&mut self, id: usize) -> Result<(), String>;
    fn delete(&mut self, id: usize) -> Result<(), String>;
}

pub struct FileRepository {
    path: String,
}

impl FileRepository {
    pub fn new(path: &str) -> Self {
        Self { path: path.to_string() }
    }

    fn load(&self) -> Result<Vec<Task>, String> {
        let content = fs::read_to_string(&self.path).unwrap_or_else(|_| "[]".to_string());
        serde_json::from_str(&content).map_err(|e| e.to_string())
    }

    fn save(&self, tasks: &[Task]) -> Result<(), String> {
        let content = serde_json::to_string_pretty(tasks).map_err(|e| e.to_string())?;
        fs::write(&self.path, content).map_err(|e| e.to_string())
    }
}

impl TaskRepository for FileRepository {
    fn add(&mut self, description: &str) -> Result<(), String> {
        let mut tasks = self.load()?;
        let id = tasks.iter().map(|t| t.id).max().unwrap_or(0) + 1;
        tasks.push(Task { id, description: description.to_string(), completed: false });
        self.save(&tasks)
    }

    fn list(&self) -> Result<Vec<Task>, String> {
        self.load()
    }

    fn complete(&mut self, id: usize) -> Result<(), String> {
        let mut tasks = self.load()?;
        let task = tasks.iter_mut().find(|t| t.id == id).ok_or(format!("Task {} not found", id))?;
        task.completed = true;
        self.save(&tasks)
    }

    fn delete(&mut self, id: usize) -> Result<(), String> {
        let mut tasks = self.load()?;
        let len_before = tasks.len();
        tasks.retain(|t| t.id != id);
        if tasks.len() == len_before {
            return Err(format!("Task {} not found", id));
        }
        self.save(&tasks)
    }
}
