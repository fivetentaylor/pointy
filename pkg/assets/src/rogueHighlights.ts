import * as RangeFix from "rangefix";
import { RogueEditor, Id } from "./rogueEditor";

export type HighlightOptions = {
  caret: boolean;
  caretColor?: string;
  caretFlagValue?: string;
  scrollTo?: boolean;
  styles: {
    backgroundColor?: string;
    borderBottom?: string;
    opacity?: number;
  };
  dataAttrs?: Record<string, string>;
  eventListeners?: {
    click?: (event: Event) => void;
    mouseover?: (event: Event) => void;
    mouseout?: (event: Event) => void;
  };
};

export class RogueHighlights extends HTMLElement {
  constructor() {
    super();
  }

  triggerClicksForCoordinates(
    x: number,
    y: number,
    options: { onFirstClick: () => void },
  ) {
    const highlightBlocks = this.querySelectorAll("rogue-highlight-block");
    let firstClick = true;

    highlightBlocks.forEach((block) => {
      const rect = block.getBoundingClientRect();
      if (
        x >= rect.left &&
        x <= rect.right &&
        y >= rect.top &&
        y <= rect.bottom
      ) {
        if (firstClick) {
          options.onFirstClick();
          firstClick = false;
        }
        (block as RogueHighlightBlock).click();
      }
    });
  }

  recalculate() {
    this.querySelectorAll("rogue-highlight").forEach((highlight) => {
      (highlight as RogueHighlight).recalulate();
    });
  }

  removeHighlightsWithPrefix(identifier: string) {
    const highlights = this.querySelectorAll("rogue-highlight");

    if (highlights) {
      for (const highlight of highlights) {
        if (highlight.id.startsWith(identifier)) {
          highlight.remove();
        }
      }
    }
  }

  highlightRange(
    identifier: string,
    ids: [Id, Id],
    address: string | null,
    options: HighlightOptions,
  ) {
    let highlight = this.querySelector(
      "rogue-highlight#" + identifier,
    ) as RogueHighlight;
    if (!highlight) {
      highlight = document.createElement("rogue-highlight") as RogueHighlight;
      highlight.setAttribute("id", identifier);
      this.append(highlight);
      if (options.scrollTo) {
        highlight.scrollIntoView({ behavior: "smooth" });
      }
    }
    const rogueEditor: RogueEditor | null = this.closest("rogue-editor");

    if (!rogueEditor) {
      console.error("rogueEditor not found");
      return;
    }

    highlight.configure(ids, address, rogueEditor, options);
  }

  removeHighlight(identifier: string) {
    const highlight = this.querySelector("rogue-highlight#" + identifier);
    if (highlight) {
      highlight.remove();
    }
  }
}

if (!customElements.get("rogue-highlights")) {
  customElements.define("rogue-highlights", RogueHighlights);
}

export class RogueHighlight extends HTMLElement {
  private ids: [Id, Id] | null = null;
  private address: string | null = null;
  private range: Range | null = null;
  private editor: RogueEditor | null = null;

  private caretColor: string | undefined;
  private eventListeners: Record<string, (event: Event) => void> | undefined;
  private largestRectangle: DOMRect | null = null;
  private caret: RogueCaret | null = null;
  private styles: HighlightOptions["styles"] | undefined;

  constructor() {
    super();
    this.redraw = this.redraw.bind(this);
    this.mouseMove = this.mouseMove.bind(this);
    window.addEventListener("resize", this.redraw);
    window.addEventListener("scroll", this.redraw);
    window.addEventListener("mousemove", this.mouseMove);
  }

  configure(
    ids: [Id, Id],
    address: string | null,
    editor: RogueEditor,
    options: HighlightOptions,
  ) {
    this.ids = ids;
    this.address = address;
    this.editor = editor;
    this.styles = options.styles;
    this.caretColor = options.caretColor;

    if (options.dataAttrs) {
      for (const [key, value] of Object.entries(options.dataAttrs)) {
        this.setAttribute(key, value);
      }
    }

    if (options.eventListeners) {
      for (const [key, value] of Object.entries(options.eventListeners)) {
        this.addEventListener(key, value);
      }
      this.eventListeners = options.eventListeners;
    }

    if (options.caret) {
      if (!this.caret) {
        this.caret = document.createElement("rogue-caret") as RogueCaret;
        this.append(this.caret);
      }
    } else {
      if (this.caret) {
        this.caret.remove();
        this.caret = null;
      }
    }

    if (options.caretFlagValue) {
      if (!this.caret) {
        console.error("Caret flag value given but no caret");
      } else {
        this.caret.configure(
          options.caretFlagValue,
          this.caretColor || this.styles?.backgroundColor || "black",
        );
      }
    }

    let range;
    if (address) {
      range = editor.rangeForAddress(ids, address);
    } else {
      range = editor.rangeFor(ids);
    }
    if (!range) {
      this.remove();
      return;
    }

    if (
      range.startContainer === range.endContainer &&
      range.startOffset === range.endOffset
    ) {
      // @jmreidy 9-19 - we're going to display carets even when range is empty
      //this.remove();
      //return;
    }

    this.range = range;

    this.redraw();
  }

  updateStyles(options: HighlightOptions["styles"]) {
    this.styles = options;
    this.redraw();
  }

  disconnectedCallback() {
    if (this.eventListeners && this.eventListeners.click) {
      this.removeEventListener("click", this.eventListeners.click);
    }

    window.removeEventListener("resize", this.redraw);
    window.removeEventListener("scroll", this.redraw);
    window.removeEventListener("mousemove", this.mouseMove);
  }

  mouseMove(e: MouseEvent) {
    // check for mouseover
    if (this.eventListeners?.mouseover || this.eventListeners?.mouseout) {
      const highlightBlocks = this.querySelectorAll("rogue-highlight-block");
      const x = e.clientX;
      const y = e.clientY;

      highlightBlocks.forEach((block) => {
        const rect = block.getBoundingClientRect();
        if (
          x >= rect.left &&
          x <= rect.right &&
          y >= rect.top &&
          y <= rect.bottom
        ) {
          if (this.eventListeners?.mouseover) {
            block.setAttribute("mouseover", "true");
            this.eventListeners?.mouseover.apply(this, [e]);
          }
        } else if (
          this.eventListeners?.mouseout &&
          block.hasAttribute("mouseover")
        ) {
          block.removeAttribute("mouseover");
          this.eventListeners?.mouseout.apply(this, [e]);
        }
      });
    }

    // check for caret
    if (!this.largestRectangle || !this.caret) {
      return;
    }
    if (isMouseEventInDOMRect(e, this.largestRectangle, 30)) {
      this.caret.mouseOver();
      return;
    }

    this.caret.mouseOut();
  }

  getRange(): Range | null {
    return this.range;
  }

  recalulate() {
    if (!this.editor || !this.ids) {
      return;
    }

    let range;
    if (this.address) {
      range = this.editor.rangeForAddress(this.ids, this.address);
    } else {
      range = this.editor.rangeFor(this.ids);
    }
    if (!range) {
      this.remove();
      return;
    }

    if (
      range.startContainer === range.endContainer &&
      range.startOffset === range.endOffset
    ) {
      this.remove();
      return;
    }

    this.range = range;
    this.redraw();
  }

  redraw() {
    const range = this.range;
    const domRectList = RangeFix.getClientRects(range);
    const existingBlocks = this.querySelectorAll("rogue-highlight-block");

    const [nonOverlappingRectangles, largestRectangle] =
      findNonOverlappingRectangles(domRectList);
    this.largestRectangle = largestRectangle;

    let offsetX = 0,
      offsetY = 0;

    if (this.editor) {
      const { x, top } = this.editor.getBoundingClientRect();
      offsetX = x;
      offsetY = top;
    }

    let next = 0;
    for (let i = 0; i < nonOverlappingRectangles.length; i++) {
      const domRect = nonOverlappingRectangles[i];
      let block = existingBlocks[next] as RogueHighlightBlock;
      if (!block) {
        block = document.createElement(
          "rogue-highlight-block",
        ) as RogueHighlightBlock;
        this.append(block);
      }

      const lineHeight = this.getLineHeight(domRect);

      // Adjust the top position to account for italic text
      const topAdjustment = Math.floor(
        Math.max(0, (lineHeight - domRect.height) / 2),
      );

      // Adjust for Firefox if necessary
      const isFirefox = navigator.userAgent.toLowerCase().includes("firefox");
      const firefoxAdjustment = isFirefox ? -15 : 0; // Adjust this value as needed

      block.style.left = domRect.left - offsetX + "px";
      block.style.top =
        domRect.top - offsetY - topAdjustment - firefoxAdjustment + "px";
      block.style.width = domRect.width + "px";
      block.style.height = lineHeight + "px";

      //apply the this.styles
      if (this.styles) {
        for (const [key, value] of Object.entries(this.styles)) {
          if (key === "opacity" && typeof value === "number") {
            block.style.opacity = value.toString();
          } else {
            (block.style as any)[key] = value;
          }
        }
      }

      next++;
    }

    // remove extra
    for (let i = next; i < existingBlocks.length; i++) {
      existingBlocks[i].remove();
    }

    if (this.caret && nonOverlappingRectangles.length > 0) {
      const lastRect =
        nonOverlappingRectangles[nonOverlappingRectangles.length - 1];
      const caretColor = this.caretColor || this.styles?.backgroundColor;
      this.caret.style.border = `1px solid ${caretColor}`;
      this.caret.style.top = lastRect.top - offsetY + "px";
      this.caret.style.left = lastRect.left + lastRect.width - offsetX + "px";
      this.caret.style.width = "1px";
      this.caret.style.height = lastRect.height + "px";

      this.caret.boop();
    }
  }

  private getLineHeight(rect: DOMRect): number {
    if (!this.range) {
      return rect.height; // Fallback to the rectangle's height if we can't calculate
    }

    const elementsInRect = this.getElementsInRect(rect);
    if (elementsInRect.length === 0) {
      return rect.height; // Fallback if no elements found
    }

    // Calculate the maximum line height of elements in this rectangle
    const maxLineHeight = Math.max(
      ...elementsInRect.map((el) => {
        const style = window.getComputedStyle(el);
        const lineHeight = parseFloat(style.lineHeight);
        if (isNaN(lineHeight)) {
          const fontSize = parseFloat(style.fontSize);
          return fontSize * 1.2; // Approximate for 'normal' line-height
        }
        return lineHeight;
      }),
    );

    return maxLineHeight;
  }

  private getElementsInRect(rect: DOMRect): Element[] {
    if (!this.range) {
      return [];
    }

    const elements: Element[] = [];
    const container =
      this.range.commonAncestorContainer.nodeType === Node.TEXT_NODE
        ? this.range.commonAncestorContainer.parentElement
        : (this.range.commonAncestorContainer as Element);

    if (!container) {
      return [];
    }

    // Check the container itself
    if (this.elementOverlapsRect(container, rect)) {
      elements.push(container);
    }

    // Use a recursive function to check all child elements
    const checkChildren = (element: Element) => {
      for (const child of element.children) {
        if (this.elementOverlapsRect(child, rect)) {
          elements.push(child);
        }
        checkChildren(child);
      }
    };

    checkChildren(container);

    return elements;
  }

  private elementOverlapsRect(element: Element, rect: DOMRect): boolean {
    const elementRect = element.getBoundingClientRect();
    const overlaps = this.rectsOverlap(elementRect, rect);
    return overlaps;
  }

  private rectsOverlap(rect1: DOMRect, rect2: DOMRect): boolean {
    const overlap = !(
      rect2.left > rect1.right ||
      rect2.right < rect1.left ||
      rect2.top > rect1.bottom ||
      rect2.bottom < rect1.top
    );

    return overlap;
  }
}

if (!customElements.get("rogue-highlight")) {
  customElements.define("rogue-highlight", RogueHighlight);
}

export class RogueHighlightBlock extends HTMLElement {
  constructor() {
    super();
  }
}

if (!customElements.get("rogue-highlight-block")) {
  customElements.define("rogue-highlight-block", RogueHighlightBlock);
}

export class RogueCaret extends HTMLElement {
  flag: HTMLDivElement;
  caret: HTMLDivElement | null = null;
  boopTimeout: NodeJS.Timeout | null = null;

  constructor() {
    super();
    this.attachShadow({ mode: "open" });

    this.flag = document.createElement("div");
    this.flag.className = "rogue-cursor-flag";

    if (!this.shadowRoot) {
      console.error("No shadow root for:", this);
      return;
    }

    this.shadowRoot.innerHTML += `
            <style>
                :host {
                    position: relative;
                    display: inline-block;
                }
                .rogue-cursor-flag {
                    font-family: var(--font-sans);
                    z-index: 200;
                    color: white;
                    display: inline-flex;
                    padding: 2px 7px;
                    justify-content: center;
                    align-items: center;
                    gap: 10px;
                    position: absolute;
                    bottom: 1.2rem;
                    left: 0px;
                    font-size: 12px;
                    border-radius: 10px 10px 10px 0;
                    opacity: 0.0;  
                    transition: opacity 0.3s ease;
                }
            </style>
        `;
    this.shadowRoot?.appendChild(this.flag);
  }

  configure(flagValue: string, color: string) {
    this.flag.textContent = flagValue;
    this.flag.style.backgroundColor = color;
    this.flag.style.display = "block";
  }

  mouseOver() {
    this.flag.style.opacity = "1.0";
  }

  mouseOut() {
    this.flag.style.opacity = "0.0";
  }

  boop() {
    if (!this.flag || this.flag.style.opacity === "1.0") {
      return;
    }

    this.flag.style.opacity = "1.0";

    if (this.boopTimeout) {
      clearTimeout(this.boopTimeout);
    }

    this.boopTimeout = setTimeout(() => {
      this.flag.style.opacity = "0.0";
    }, 2000);
  }
}
if (!customElements.get("rogue-caret")) {
  customElements.define("rogue-caret", RogueCaret);
}

function findNonOverlappingRectangles(
  domRects: DOMRect[],
): [DOMRect[], DOMRect] {
  let minTopX = Infinity;
  let minTopY = Infinity;
  let maxBottomX = -Infinity;
  let maxBottomY = -Infinity;

  const rectangles = Array.from(domRects);
  // Sort rectangles by height in ascending order
  rectangles.sort((a, b) => a.height - b.height);

  const nonOverlapping: DOMRect[] = [];

  const tolerance = 10;

  function doesOverlap(rect1: DOMRect, rect2: DOMRect): boolean {
    return (
      rect1.x < rect2.x + rect2.width - tolerance &&
      rect1.x + rect1.width > rect2.x + tolerance &&
      rect1.y < rect2.y + rect2.height - tolerance &&
      rect1.y + rect1.height > rect2.y + tolerance
    );
  }

  for (const rect of rectangles) {
    let overlapFound = false;
    for (const selected of nonOverlapping) {
      if (doesOverlap(rect, selected)) {
        overlapFound = true;
        break;
      }
    }
    if (!overlapFound) {
      nonOverlapping.push(rect);
      if (rect.x < minTopX) {
        minTopX = rect.x;
      }
      if (rect.y < minTopY) {
        minTopY = rect.y;
      }
      if (rect.x + rect.width > maxBottomX) {
        maxBottomX = rect.x + rect.width;
      }
      if (rect.y + rect.height > maxBottomY) {
        maxBottomY = rect.y + rect.height;
      }
    }
  }

  return [
    nonOverlapping,
    new DOMRect(minTopX, minTopY, maxBottomX - minTopX, maxBottomY - minTopY),
  ];
}

function isMouseEventInDOMRect(
  event: MouseEvent,
  rect: DOMRect,
  buffer: number = 0,
): boolean {
  const mouseX = event.clientX;
  const mouseY = event.clientY;

  return (
    mouseX >= rect.x - buffer &&
    mouseX <= rect.x + rect.width + buffer &&
    mouseY >= rect.y - buffer &&
    mouseY <= rect.y + rect.height + buffer
  );
}
