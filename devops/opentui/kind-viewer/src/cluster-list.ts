import {
    BoxRenderable,
    SelectRenderable,
    SelectRenderableEvents,
} from "@opentui/core";
import type { CliRenderer } from "@opentui/core";

export class ClusterList {
    private renderer: CliRenderer;
    private container: BoxRenderable;
    private select: SelectRenderable | null = null;
    private onSelect: (cluster: string | null) => void = () => {};

    constructor(renderer: CliRenderer) {
        this.renderer = renderer;

        this.container = new BoxRenderable(renderer, {
            id: "clusters-panel",
            width: 30,
            border: true,
            borderStyle: "single",
            borderColor: "#30363d",
            title: " Clusters ",
            titleAlignment: "left",
            flexDirection: "column",
            backgroundColor: "#0d1117",
        });
    }

    getComponent() {
        return this.container;
    }

    onSelectionChange(fn: (cluster: string | null) => void) {
        this.onSelect = fn;
    }

    update(clusters: string[], selectedCluster: string | null) {
        if (this.select) {
            this.container.remove("cluster-list");
            this.select.destroy();
            this.select = null;
        }

        const options =
            clusters.length > 0
                ? clusters.map((n) => ({ name: n, description: "", value: n }))
                : [{ name: "(no clusters)", description: "Press [N] to create one", value: "" }];

        const initIdx = selectedCluster
            ? Math.max(clusters.indexOf(selectedCluster), 0)
            : 0;

        this.select = new SelectRenderable(this.renderer, {
            id: "cluster-list",
            flexGrow: 1,
            options,
            selectedIndex: initIdx,
        });

        this.select.on(SelectRenderableEvents.SELECTION_CHANGED, (index) => {
            if (clusters.length === 0) return;
            this.onSelect(clusters[index] ?? null);
        });

        this.container.add(this.select);
        this.select.focus();
    }

    focus() {
        this.select?.focus();
    }

    blur() {
        this.select?.blur();
    }
}
