import {
  BoxRenderable,
  TextRenderable,
  TextAttributes,
} from "@opentui/core";
import type { CliRenderer } from "@opentui/core";

export class Header {
  private container: BoxRenderable;
  private colimaStatusText: TextRenderable;

  constructor(renderer: CliRenderer) {
    this.container = new BoxRenderable(renderer, {
      id: "header",
      height: 3,
      border: true,
      borderStyle: "rounded",
      borderColor: "#30363d",
      flexDirection: "row",
      justifyContent: "space-between",
      alignItems: "center",
      paddingLeft: 1,
      paddingRight: 1,
      backgroundColor: "#0d1117",
    });

    this.container.add(new TextRenderable(renderer, {
      id: "header-title",
      content: "☸  Kind Cluster Manager",
      fg: "#58a6ff",
      attributes: TextAttributes.BOLD,
    }));

    this.colimaStatusText = new TextRenderable(renderer, {
      id: "colima-status",
      content: "○ Colima",
      fg: "#8b949e",
    });
    this.container.add(this.colimaStatusText);
  }

  getComponent() {
    return this.container;
  }

  updateColimaStatus(running: boolean) {
    this.colimaStatusText.content = running ? "● Colima" : "○ Colima";
    this.colimaStatusText.fg = running ? "#3fb950" : "#f85149";
  }
}
