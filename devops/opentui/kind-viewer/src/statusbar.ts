import {
    BoxRenderable,
    TextRenderable,
} from "@opentui/core";

import type { CliRenderer } from "@opentui/core";

export class StatusBar {
    private renderer: CliRenderer;
    private container: BoxRenderable;
    private statusLabel: TextRenderable;

    constructor(renderer: CliRenderer) {
        this.renderer = renderer;
    
        this.container  = new BoxRenderable(this.renderer, {
            id: "status-bar",
            height: 3,
            border: true,
            borderStyle: "rounded",
            borderColor: "#30363d",
            alignItems: "center",
            paddingLeft: 2,
            backgroundColor: "#0d1117",
        });
        this.statusLabel = new TextRenderable(this.renderer, {
            id: "status-label",
            content: "",
            fg: "#8b949e",
        });

        this.container.add(this.statusLabel);
    }

    getComponent() {
        return this.container;
    }

    update(content: string) {
       this.statusLabel.content = content;
    }

}