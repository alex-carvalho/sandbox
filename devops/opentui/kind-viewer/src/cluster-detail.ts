import {
  BoxRenderable,
  TextRenderable,
  TextTableRenderable,
  TextAttributes,
  RGBA,
  type TextTableContent,
} from "@opentui/core";
import type { CliRenderer } from "@opentui/core";


type Destroyable = { destroy(): void };

export class ClusterDetailPanel {
  private renderer: CliRenderer;
  private parent: BoxRenderable;
  private items: Array<{ id: string; r: Destroyable }> = [];
  private idCounter = 0;

  constructor(renderer: CliRenderer) {
    this.renderer = renderer;
    this.parent = new BoxRenderable(renderer, {});
  }

  private nextId() {
    return `cd-${this.idCounter++}`;
  }

  private addRenderable<T extends Destroyable>(r: T, id: string): T {
    this.parent.add(r as unknown as BoxRenderable);
    this.items.push({ id, r });
    return r;
  }

  private text(content: string, fg = "#c9d1d9", attributes?: number) {
    const id = this.nextId();
    const r = new TextRenderable(this.renderer, {
      id,
      content,
      fg,
      ...(attributes !== undefined ? { attributes } : {}),
    });
    return this.addRenderable(r, id);
  }

  private box(opts: Record<string, unknown>): BoxRenderable {
    const id = this.nextId();
    const r = new BoxRenderable(this.renderer, { id, ...opts });
    return this.addRenderable(r, id);
  }

  clear() {
    for (const { id, r } of this.items) {
      this.parent.remove(id);
      r.destroy();
    }
    this.items = [];
  }

  getComponent() {
    return this.parent;
  }

  update(cluster: string | null, nodes: string[]) {
    this.clear();

    if (!cluster) {
      this.text("Select a cluster to view details.", "#6e7681");
      return;
    }

    const name = cluster;

    const headerRow = this.box({
      flexDirection: "row",
      alignItems: "center",
      marginBottom: 1,
    });

    headerRow.add(new TextRenderable(this.renderer, {
      id: `${this.idCounter}-name`,
      content: `☸  ${name}`,
      fg: "#58a6ff",
      attributes: TextAttributes.BOLD,
    }));

    const spacer = new BoxRenderable(this.renderer, {
      id: `${this.idCounter}-spacer`,
      flexGrow: 1,
    });
    headerRow.add(spacer);

    const badge = new TextRenderable(this.renderer, {
      id: `${this.idCounter}-badge`,
      content: `${nodes.length} node${nodes.length !== 1 ? "s" : ""}`,
      fg: "#3fb950",
    });
    headerRow.add(badge);

    this.text("─".repeat(40), "#30363d");

    this.text("NODES", "#8b949e", TextAttributes.BOLD);

    if (nodes.length === 0) {
      this.text("  No nodes found.", "#6e7681");
      return;
    }

    const header = [
      [{ __isChunk: true as const, text: "Name", fg: RGBA.fromHex("#c9d1d9") }],
      [{ __isChunk: true as const, text: "Role", fg: RGBA.fromHex("#c9d1d9") }],
    ];

    const rows = nodes.map((node) => {
      const role = node.includes("control-plane") ? "control-plane" : "worker";
      return [
        [{ __isChunk: true as const, text: node, fg: RGBA.fromHex("#8b949e") }],
        [{ __isChunk: true as const, text: role, fg: RGBA.fromHex(role === "control-plane" ? "#58a6ff" : "#3fb950") }],
      ];
    });

    const tableId = this.nextId();
    const table = new TextTableRenderable(this.renderer, {
      id: tableId,
      content: [header, ...rows] as unknown as TextTableContent,
      columnWidthMode: "content",
      border: true,
      outerBorder: true,
      borderStyle: "single",
      borderColor: "#30363d",
      backgroundColor: "#0d1117",
      cellPadding: 1,
      marginTop: 1,
    });
    this.addRenderable(table, tableId);
  }
}
