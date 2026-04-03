import {
  createCliRenderer,
  BoxRenderable,
} from "@opentui/core";
import { getClusters, getNodes, runCreate, runDelete } from "./kind-service";
import { getColimaStatus, startColima, stopColima } from "./colima-service";
import { ClusterDetailPanel } from "./cluster-detail";
import { ClusterList } from "./cluster-list";
import { Header } from "./header";
import { StatusBar } from "./statusbar";
import { CreateModal } from "./create-modal";
import { logger } from "./logger";

type AppMode = "list" | "creating" | "deleting" | "loading";

let clusters: string[] = [];
let selectedCluster: string | null = null;
let clusterNodes: string[] = [];
let mode: AppMode = "loading";
let colimaRunning = false;


const renderer = await createCliRenderer({ exitOnCtrlC: true, useConsole: true });

const rootBox = new BoxRenderable(renderer, {
  id: "root",
  flexDirection: "column",
  flexGrow: 1,
});
renderer.root.add(rootBox);

const header = new Header(renderer);
rootBox.add(header.getComponent());

const contentRow = new BoxRenderable(renderer, {
  id: "content",
  flexDirection: "row",
  flexGrow: 1,
});
rootBox.add(contentRow);

const clusterListPanel = new ClusterList(renderer);
clusterListPanel.onSelectionChange(async (cluster) => {
  if (mode !== "list") return;
  selectedCluster = cluster;
  detailPanel.update(selectedCluster, await getNodeList());
});

const nodesCache: Record<string, string[]> = {};

async function getNodeList(): Promise<string[]> {
  if (selectedCluster) {
    if (nodesCache[selectedCluster]) {
      return nodesCache[selectedCluster]!;
    }
    const nodes = await getNodes(selectedCluster);
    nodesCache[selectedCluster] = nodes;
    return nodes;
  }
  return [];
}

contentRow.add(clusterListPanel.getComponent());

const rightPanel = new BoxRenderable(renderer, {
  id: "details-panel",
  flexGrow: 1,
  border: true,
  borderStyle: "single",
  borderColor: "#30363d",
  title: " Details ",
  titleAlignment: "left",
  flexDirection: "column",
  padding: 1,
  backgroundColor: "#0d1117",
});
contentRow.add(rightPanel);


const statusBar = new StatusBar(renderer);
rootBox.add(statusBar.getComponent());


const detailPanel = new ClusterDetailPanel(renderer);
rightPanel.add(detailPanel.getComponent());


function updateStatus(newMode: AppMode, override?: string) {
  mode = newMode;
  if (override) {
    statusBar.update(`  ${override}`);
    return;
  }
  const map: Record<AppMode, string> = {
    loading: "Loading...",
    list: "[N] New  [D] Delete  [R] Refresh  [C] Colima  [↑↓] Navigate",
    deleting: `Delete "${selectedCluster}"?  [Y] Confirm  [N / ESC] Cancel`,
    creating: "[Enter] Create  [ESC] Cancel",
  };
  statusBar.update(`  ${map[newMode]}`);
}

function updateColimaDisplay(running: boolean) {
  colimaRunning = running;
  header.updateColimaStatus(running);
}

// ── Load / refresh ────────────────────────────────────────────────────────────

async function loadClusters(preferCluster?: string | null) {
  updateStatus("loading");
  clusters = await getClusters();
  const pref = preferCluster ?? selectedCluster;
  selectedCluster = pref && clusters.includes(pref) ? pref : (clusters[0] ?? null);
  clusterNodes = await getNodeList();
  clusterListPanel.update(clusters, selectedCluster);
  detailPanel.update(selectedCluster, clusterNodes);
  updateStatus("list");
}


const createModal = new CreateModal(renderer);

function showModalCreate() {

  createModal.show(
    async (name) => {
      updateStatus("loading", `Creating cluster "${name}"...`);
      const err = await runCreate(name);
      if (err) {
        clusterListPanel.focus();
        updateStatus("list", `Error: ${err}`);
        setTimeout(() => updateStatus(mode), 6000);
        return;
      }
      await loadClusters(name);
    },
    () => {
      clusterListPanel.focus();
      updateStatus("list");
    },
  );
  updateStatus("creating");
}



// ── Keyboard handling ─────────────────────────────────────────────────────────

renderer.keyInput.on("keypress", async (key) => {
  if (mode === "loading") return;

  if (mode === "deleting") {
    switch (key.name) {
      case "y": {
        const name = selectedCluster!;
        updateStatus("loading", `Deleting cluster "${name}"...`);
        const err = await runDelete(name);
        if (err) {
          clusterListPanel.focus();
          updateStatus("list", `Error: ${err}`);
          setTimeout(() => updateStatus(mode), 6000);
          return;
        }
        await loadClusters(null);
        break;
      }
      case "n":
      case "escape":
        clusterListPanel.focus();
        updateStatus("list");
        break;
    }
    return;
  }

  if (mode !== "list") {
    return;
  }

  switch (key.name) {
    case "n":
      showModalCreate();
      break;
    case "d":
      if (selectedCluster && clusters.length > 0) updateStatus("deleting");
      break;
    case "r":
      await loadClusters();
      break;
    case "c":
      await toggleColima();
      break;
  }
});

async function toggleColima() {
  if (colimaRunning) {
      updateStatus("loading", "Stopping Colima...");
      const err = await stopColima();
      if (err) {
        updateStatus("list", `Colima error: ${err}`);
        setTimeout(() => updateStatus(mode), 6000);
      } else {
        updateColimaDisplay(false);
        updateStatus("list");
      }
    } else {
      updateStatus("loading", "Starting Colima...");
      const err = await startColima();
      if (err) {
        updateStatus("list", `Colima error: ${err}`);
        setTimeout(() => updateStatus(mode), 6000);
      } else {
        updateColimaDisplay(true);
        await loadClusters();
      }
    }
}

await loadClusters();
getColimaStatus().then(updateColimaDisplay);
