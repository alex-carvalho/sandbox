import { logger } from "./logger";

export async function getClusters(): Promise<string[]> {
  try {
    logger.info("Loading clusters...");
    const output = await Bun.$`kind get clusters`.quiet().text();
    return output.trim().split("\n").filter(Boolean);
  } catch (error) {
    logger.error(`Failed to load clusters: ${error}`);
    return [];
  }
}

export async function getNodes(cluster: string): Promise<string[]> {
  try {
    logger.info(`Loading nodes for cluster "${cluster}"...`);
    const output = await Bun.$`kind get nodes --name ${cluster}`.quiet().text();
    return output.trim().split("\n").filter(Boolean);
  } catch (error) {
    logger.error(`Failed to load nodes for cluster "${cluster}": ${error}`);
    return [];
  }
}

export async function runCreate(name: string): Promise<string | null> {
  logger.info(`Creating cluster "${name}"...`);
  const result = await Bun.$`kind create cluster --name ${name}`.quiet().nothrow();
  if (result.exitCode !== 0) {
    const message = result.stderr.toString().trim() || `exit code ${result.exitCode}`;
    logger.error(`Failed to create cluster "${name}": ${message}`);
    return message;
  }
  return null;
}

export async function runDelete(name: string): Promise<string | null> {
  logger.info(`Deleting cluster "${name}"...`);
  const result = await Bun.$`kind delete cluster --name ${name}`.quiet().nothrow();
  if (result.exitCode !== 0) {
    const message = result.stderr.toString().trim() || `exit code ${result.exitCode}`;
    logger.error(`Failed to delete cluster "${name}": ${message}`);
    return message;
  }
  return null;
}