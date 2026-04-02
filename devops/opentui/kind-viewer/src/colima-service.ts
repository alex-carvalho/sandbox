import { logger } from "./logger";

export async function getColimaStatus(): Promise<boolean> {
  try {
    const result = await Bun.$`colima status`.quiet().nothrow();
    return result.exitCode === 0;
  } catch {
    return false;
  }
}

export async function startColima(): Promise<string | null> {
  logger.info("Starting Colima...");
  const result = await Bun.$`colima start`.quiet().nothrow();
  if (result.exitCode !== 0) {
    const message = result.stderr.toString().trim() || `exit code ${result.exitCode}`;
    logger.error(`Failed to start Colima: ${message}`);
    return message;
  }
  return null;
}

export async function stopColima(): Promise<string | null> {
  logger.info("Stopping Colima...");
  const result = await Bun.$`colima stop`.quiet().nothrow();
  if (result.exitCode !== 0) {
    const message = result.stderr.toString().trim() || `exit code ${result.exitCode}`;
    logger.error(`Failed to stop Colima: ${message}`);
    return message;
  }
  return null;
}
