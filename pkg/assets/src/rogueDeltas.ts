import { parseId } from "@/lib/utils";
import { RogueEditor, Id } from "./rogueEditor";

const deleteModalBottomPadding = 4;
const dataDeltaStart = "data-delta-start";
const dataDeltaEnd = "data-delta-end";

export class RogueDeltasManager extends HTMLElement {
  editor: RogueEditor | null = null;
  openDelta: HTMLElement | null = null;
  openDeltaModal: HTMLElement | null = null;

  constructor() {
    super();
  }

  close() {
    if (this.openDeltaModal) {
      this.openDeltaModal.remove();
      this.openDeltaModal = null;
    }
    if (this.openDelta) {
      this.openDelta.classList.remove("show");
      this.openDelta = null;
    }
  }

  rewind(target: HTMLElement) {
    if (!this.editor) {
      console.error("no editor");
      return;
    }

    if (!target) {
      console.error("no target");
      return;
    }

    const changes = findAdjacentChanges(target);
    if (!changes || !changes.length) {
      console.error("no changes");
      return;
    }

    const ids = getDeltaIds(changes);
    if (!ids || !ids.length) {
      console.error("no ids");
      return;
    }

    const wrappingIds = findLargestAndSmallestIndex(this.editor, ids);
    if (!wrappingIds || !wrappingIds.smallestId || !wrappingIds.largestId) {
      console.error("missing wrapping id(s) " + JSON.stringify(wrappingIds));
      return;
    }

    this.close();

    this.editor?.rewind(wrappingIds.smallestId, wrappingIds.largestId);
  }

  show(target: HTMLElement) {
    if (this.openDelta) {
      this.close();
    }

    target.classList.add("show");
    this.openDelta = target;

    if (this.editor?.editorMode !== "diff") {
      return;
    }

    const rects = target.getClientRects();

    // Find the topmost rectangle (first line)
    const topRect = Array.from(rects).reduce((top, rect) =>
      rect.top < top.top ? rect : top,
    );

    const parentRect = this.getBoundingClientRect();

    const relativeX = topRect.left - parentRect.left;
    const relativeY = topRect.top - parentRect.top;

    const modal = new RogueDeltaModal();
    modal.manager = this;
    modal.target = target;
    this.append(modal);

    const modalRect = modal.getBoundingClientRect();

    modal.style.left = relativeX + "px";
    modal.style.top =
      relativeY - deleteModalBottomPadding - modalRect.height + "px";

    this.openDeltaModal = modal;
  }
}

if (!customElements.get("rogue-deltas-manager")) {
  customElements.define("rogue-deltas-manager", RogueDeltasManager);
}

class RogueDeltaModal extends HTMLElement {
  manager: RogueDeltasManager | null = null;
  target: HTMLElement | null = null;

  constructor() {
    super();
  }

  rewind(e: MouseEvent) {
    e.stopPropagation();
    if (!this.target) {
      console.error("no target");
      return;
    }

    this.manager?.rewind(this.target);
  }

  close() {
    if (this.manager) {
      this.manager.close();
    }
  }

  connectedCallback() {
    const rewindButton = document.createElement("button");
    rewindButton.textContent = "Undo";
    rewindButton.addEventListener("click", this.rewind.bind(this));
    this.append(rewindButton);
  }
}

if (!customElements.get("rogue-delta-modal")) {
  customElements.define("rogue-delta-modal", RogueDeltaModal);
}

function findAdjacentChanges(target: HTMLElement): HTMLElement[] {
  const adjacentChanges: HTMLElement[] = [target];
  let sibling = target.previousElementSibling;

  // Check previous siblings
  while (sibling && (sibling.tagName === "DEL" || sibling.tagName === "INS")) {
    adjacentChanges.push(sibling as HTMLElement);
    sibling = sibling.previousElementSibling;
  }

  sibling = target.nextElementSibling;

  // Check next siblings
  while (sibling && (sibling.tagName === "DEL" || sibling.tagName === "INS")) {
    adjacentChanges.push(sibling as HTMLElement);
    sibling = sibling.nextElementSibling;
  }

  return adjacentChanges;
}

function getDeltaIds(elements: HTMLElement[]): Id[] {
  const deltaIdStrs: string[] = [];
  elements.map((element) => {
    const start = element.getAttribute("data-delta-start") || "";
    if (start.length === 0) {
      console.error("no " + dataDeltaStart + " attribute for element", element);
    }
    deltaIdStrs.push(start);

    const end = element.getAttribute("data-delta-end") || "";
    if (end.length === 0) {
      console.error("no " + dataDeltaEnd + " attribute for element", element);
    }
    deltaIdStrs.push(end);
  });

  return deltaIdStrs.map(parseId);
}

function findLargestAndSmallestIndex(
  editor: RogueEditor,
  ids: Id[],
): { smallestId: Id | null; largestId: Id | null } {
  let smallestIndex = Number.MAX_SAFE_INTEGER;
  let largestIndex = Number.MIN_SAFE_INTEGER;
  let smallestId: Id | null = null;
  let largestId: Id | null = null;

  ids.forEach((id) => {
    const indexes = editor.getIndex(id);
    console.log("getIndex", id, indexes);
    if (indexes.error) {
      console.error("getIndex: " + id + ": " + indexes.error);
      return;
    }
    if (indexes.total < smallestIndex) {
      smallestIndex = indexes.total;
      smallestId = id;
    }
    if (indexes.total > largestIndex) {
      largestIndex = indexes.total;
      largestId = id;
    }
  });

  return { smallestId, largestId };
}
