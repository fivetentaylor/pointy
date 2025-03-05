class RogueIdElement extends HTMLElement {
  private tooltip: HTMLSpanElement;

  constructor() {
    super();
    this.tooltip = document.createElement("span");
    this.setupTooltip();
  }

  private setupTooltip(): void {
    this.tooltip.style.position = "absolute";
    this.tooltip.style.display = "none";
    this.tooltip.style.backgroundColor = "black";
    this.tooltip.style.color = "white";
    this.tooltip.style.padding = "4px 8px";
    this.tooltip.style.borderRadius = "4px";
    this.tooltip.style.fontSize = "12px";
    this.tooltip.style.zIndex = "1000";
  }

  connectedCallback(): void {
    document.body.appendChild(this.tooltip);

    this.addEventListener("mouseover", this.handleMouseOver);
    this.addEventListener("mouseout", this.handleMouseOut);
    this.addEventListener("click", this.handleClick);
  }

  private handleMouseOver = (event: MouseEvent): void => {
    const rogueId = this.getAttribute("data-rogue-id");
    const isDel = this.getAttribute("data-is-del");
    if (rogueId) {
      const elements = document.querySelectorAll(
        `[data-rogue-id="${rogueId}"]`,
      );
      elements.forEach((el) => el.classList.add("userhovered"));

      // Show tooltip
      if (isDel === "true") {
        this.tooltip.style.backgroundColor = "#FFEBEB";
        this.tooltip.style.color = "#B30000";
      }

      this.tooltip.textContent = rogueId;
      this.tooltip.style.top = `${event.clientY - 40}px`;

      this.tooltip.style.display = "block";
      this.tooltip.style.left = `${event.clientX - 10}px`;
    }
  };

  private handleMouseOut = (): void => {
    const rogueId = this.getAttribute("data-rogue-id");
    if (rogueId) {
      const elements = document.querySelectorAll(
        `[data-rogue-id="${rogueId}"]`,
      );
      elements.forEach((el) => el.classList.remove("userhovered"));

      // Hide tooltip
      this.tooltip.style.display = "none";
    }
  };

  private handleClick = (event: MouseEvent): void => {
    const rogueId = this.getAttribute("data-rogue-id");
    if (!rogueId) return;
    navigator.clipboard.writeText(rogueId);
    this.style.borderColor = "green";
    this.tooltip.innerHTML = "ID copied to clipboard";
  };

  disconnectedCallback(): void {
    if (this.tooltip && this.tooltip.parentElement) {
      this.tooltip.parentElement.removeChild(this.tooltip);
    }
  }
}

customElements.define("rogue-id", RogueIdElement);
