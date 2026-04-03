import { appendFileSync } from "fs";
import { join } from "path";

const LOG_FILE_ENABLED = process.env.LOG_FILE_ENABLED === "true";

const LOG_FILE = join(import.meta.dir, "..", "app.log");

function write(level: string, msg: string) {
  const line = `${new Date().toISOString()} [${level}] ${msg}\n`;
  if (LOG_FILE_ENABLED) {
    appendFileSync(LOG_FILE, line);
  }
}

export const logger = {
  info: (msg: string) => write("INFO", msg),
  warn: (msg: string) => write("WARN", msg),
  error: (msg: string) => write("ERROR", msg),
};
