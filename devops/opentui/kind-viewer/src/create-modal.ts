import {
    BoxRenderable,
    TextRenderable,
    InputRenderable,
    InputRenderableEvents,
    type KeyEvent,
} from "@opentui/core";
import type { CliRenderer } from "@opentui/core";

export class CreateModal {
    private renderer: CliRenderer;
    private overlay: BoxRenderable;
    private inputRow: BoxRenderable;
    private input: InputRenderable;
    private keyHandler: ((key: KeyEvent) => void) | null = null;

    constructor(renderer: CliRenderer) {
        this.renderer = renderer;

        this.overlay = new BoxRenderable(renderer, {
            id: "overlay",
            position: "absolute",
            width: "100%",
            height: "100%",
            left: 0,
            top: 0,
            justifyContent: "center",
            alignItems: "center",
            zIndex: 100,
        });

        const modal = new BoxRenderable(renderer, {
            id: "modal",
            width: 44,
            height: 8,
            border: true,
            borderStyle: "rounded",
            borderColor: "#58a6ff",
            title: " Create New Cluster ",
            titleAlignment: "center",
            backgroundColor: "#161b22",
            flexDirection: "column",
            paddingLeft: 2,
            paddingRight: 2,
            paddingTop: 1,
            paddingBottom: 1,
            gap: 1,
        });
        this.overlay.add(modal);

        this.inputRow = new BoxRenderable(renderer, {
            id: "modal-row",
            flexDirection: "row",
            alignItems: "center",
            gap: 1,
        });
        modal.add(this.inputRow);

        this.inputRow.add(new TextRenderable(renderer, {
            id: "modal-label",
            content: "Name:",
            fg: "#c9d1d9",
        }));

        this.input = new InputRenderable(renderer, {
            id: "modal-input",
            width: 26,
            placeholder: "e.g. my-cluster",
            backgroundColor: "#0d1117",
            textColor: "#c9d1d9",
            cursorColor: "#58a6ff",
            focusedBackgroundColor: "#0d1117",
        });

        this.inputRow.add(this.input);

        modal.add(new TextRenderable(renderer, {
            id: "modal-hint",
            content: "[Enter] Create  [ESC] Cancel",
            fg: "#6e7681",
        }));

    }

    getComponent() {
        return this.overlay;
    }

    onChange(onChange: (name: string) => void) {
        this.input.on(InputRenderableEvents.CHANGE, onChange);
    }

    show(onSubmit: (name: string) => void, onCancel: () => void) {
        this.input.value = "";
        this.renderer.root.add(this.overlay);

        this.keyHandler = (key: KeyEvent) => {
            if (key.name === "escape") {
                this.hide();
                onCancel();
            } else if (key.name === "return" || key.name === "enter") {
                const name = this.input.value.trim();
                if (!name) return;
                this.hide();
                onSubmit(name);
            }
        };

        // if add it on the same loop the character `n` is added to the input
        setTimeout(() => {
            this.input.value = "";
            this.input.focus();
            this.renderer.keyInput.on("keypress", this.keyHandler!);
        }, 0);
    }

    private hide() {
        if (this.keyHandler) {
            this.renderer.keyInput.off("keypress", this.keyHandler);
            this.keyHandler = null;
        }
        this.renderer.root.remove("overlay");
    }
}
