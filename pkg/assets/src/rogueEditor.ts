import * as Sentry from "@sentry/browser";
import { OperationManager } from "./operationManager";
import {
  RogueHighlights,
  HighlightOptions,
  RogueHighlight,
} from "./rogueHighlights";
import { RogueDeltasManager } from "./rogueDeltas";

const RogueIdAttiribute = "data-rid";
const RogueHTMLCopyMarker = 'data-editor="rogue"';
const RogueHTMLDraftIdRegex = /draft-id="([^"]*)"/;

const HIGHLIGHT_STYLES = {
  active: {
    backgroundColor: "hsla(48, 96%, 53%, 0.25)",
    borderBottom: "1px solid hsla(45, 93%, 47%, 1)",
  },
  inactive: {
    backgroundColor: "hsla(48, 96%, 53%, 0.05)",
    borderBottom: "1px solid hsla(45, 93%, 47%, 0.8)",
  },
};

// WASM types
export interface Go {
  importObject: any;
  run: (instance: any) => void;
}
declare global {
  var Go: new () => any;
  var NewRogue: (authorId: string) => any;
  var RegisterPanicCallback: (callback: (msg: string) => void) => void;
  var re: RogueEditor;
  var rogueVersion: string;
  var process: {
    env: {
      APP_HOST: string;
      WS_HOST: string;
      NODE_ENV: string;
      WEB_HOST: string;
      SEGMENT_KEY: string;
      PUBLIC_POSTHOG_KEY: string;
      PUBLIC_POSTHOG_HOST: string;
      IMAGE_TAG: string;
    };
  };
}

export type CallbackFunction<T> = (newValue: T) => void;

export type OpStats = {
  inserts: number[];
  deletes: number[];
  insertsByPrefix: { [key: string]: number[] };
  deletesByPrefix: { [key: string]: number[] };
  currentCharsByPrefix: { [key: string]: number };
  segments: number;
};

export type DocStats = {
  wordCount: number;
  paragraphCount: number;
};

export type SubscribableAttr =
  | "activeComments"
  | "address"
  | "editorMode"
  | "addressDescription"
  | "baseAddress"
  | "canUndo"
  | "canRedo"
  | "enabled"
  | "editing"
  | "loaded"
  | "syncing"
  | "connected"
  | "applyingDiff"
  | "showDiffHighlights"
  | "curLineFormat"
  | "curSpanFormat"
  | "selectedHtml"
  | "lastEdit"
  | "cursors";

export type EditorMode =
  | "diff"
  | "history"
  | "paste"
  | "edit"
  | "xray"
  | "scrub";

const supportedTextMimeTypes = [
  "text/plain",
  "text/html",
  "text/rtf",
  "text/markdown",

  // custom formats for now just used to distinguish the type of content
  "text/_notion",
  "application/x-vnd.google-docs-document-slice-clip+wrapped",

  // "text/csv", // once we support tables
];

const supportedImageMimeTypes = [
  "image/gif",
  "image/png",
  "image/jpeg",
  "image/tiff",

  /*
  "image/webp",
  "image/svg+xml",
  "image/bmp",
  "image/x-icon",
  "image/vnd.microsoft.icon",
  "image/vnd.wap.wbmp",
  "image/x-xbitmap",
  "image/x-xbm",
  "image/x-portable-bitmap",
  "image/x-portable-graymap",
  "image/x-portable-pixmap
  */
];

// Rogue types
export type Id = [string, number];
type RogueRange = [Id, Id];
type Op = [number, Id, ...any[]];
export type Event = {
  event: string;
  op?: Op;
};

export type AuthorInfo = {
  userID: string;
  authorID: string;
  color: string;
  name: string;
  editing: boolean;
};

const defaultAuthorInfo: AuthorInfo = {
  userID: "",
  authorID: "",
  color: "#312e81",
  name: "Anonymous",
  editing: false,
};

interface PasteItem {
  kind: string;
  mime: string;
  data: string;
}

export class RogueEditor extends HTMLElement {
  private docID: string | null = null;
  private contentDiv: HTMLElement | null = null;
  private rogue: any | null = null;
  private reconnectInterval: any | null = null;
  private operationManager: OperationManager | null = null;
  private isWasmLoaded: boolean = false;
  private apiHost: string = process.env.APP_HOST || "";
  private wsHost: string = process.env.WS_HOST || "";
  private highlights: RogueHighlights;
  private deltaManager: RogueDeltasManager;
  private subscribers: { [key: string]: Array<CallbackFunction<any>> } = {};

  private _activeComments: string[] = [];
  private _curRogueRange: RogueRange | null = null;
  private _canUndo: boolean = false;
  private _canRedo: boolean = false;
  private _active: boolean = false; // if the instance of RogueEditor is in the DOM
  private _syncing: boolean = false;
  private _loaded: boolean = false;
  private _connected: boolean = false;
  private _curSpanFormat: Record<string, any> = {};
  private _curLineFormat: Record<string, any> = {};
  private _applyingDiff: boolean = false;
  private _selectedHtml: string | null = null;
  private _lastEdit: Date | null = null;
  private _lastTargetAnchor: HTMLAnchorElement | null = null;
  private _baseAddress: string | null = null;
  private _address: string | null = null;
  private _editorMode: EditorMode = "edit";
  private _addressDescription: string | null = null;
  private _showDiffHighlights: boolean = false;
  private _enabled: boolean = false;
  private _cursors: Record<string, AuthorInfo> = {}; // authorID => authorInfo
  private _editing: boolean = false; // Has the user send ops since the subscription
  private _opStats: OpStats | null = null;
  private _opStatsAddress: string | null = null;
  private _docStats: DocStats | null = null;
  private _docStatsAddress: string | null = null;

  public scrubMode: "full" | "partial" = "full";
  public instanceID: string;
  public imageTag: string = process.env.IMAGE_TAG || "unknown";
  public uploadImage: (
    file: File,
    docId: string,
  ) => Promise<Image | undefined> = async () => undefined;

  public getImage: (
    docId: string,
    imageId: string,
  ) => Promise<Image | undefined> = async () => undefined;

  public listDocumentImages: (docId: string) => Promise<Image[] | undefined> =
    async () => [];

  public getImageSignedUrl: (
    docId: string,
    imageId: string,
  ) => Promise<string | undefined> = async () => undefined;

  recvEventBuffer: Event[] = [];

  network: WebSocket | null = null;

  config = { childList: true, subtree: true, characterData: true };
  debug = false;

  constructor() {
    super();

    this.instanceID = Math.random().toString(36).substring(2, 15);
    this.syncing = true;
    this.loaded = false;

    this.deltaManager = new RogueDeltasManager();
    this.deltaManager.editor = this;
    this.prepend(this.deltaManager);

    this.highlights = new RogueHighlights();
    this.prepend(this.highlights);

    this.docID = this.getAttribute("docid");
    if (!this.docID) {
      console.error("docid not found");
      return;
    }

    const apiHostOverride = this.getAttribute("apihost");
    if (apiHostOverride) {
      this.apiHost = apiHostOverride;
    }

    const wsHostOverride = this.getAttribute("wsHost");
    if (wsHostOverride) {
      this.wsHost = wsHostOverride;
    }

    this.contentDiv = this.querySelector(".content");
    if (!this.contentDiv) {
      console.trace("contentDiv not found");
      return;
    }
    this.connect();
    this.operationManager = new OperationManager(
      this.docID,
      this.updateAuthorId.bind(this),
    );

    window.re = this;
  }

  connectedCallback() {
    this._active = true;
    if (this.debug) {
      console.log("connectedCallback");
    }

    if (!this.contentDiv) {
      console.error("contentDiv not found");
      return;
    }

    this.contentDiv.addEventListener("copy", this.onCopy);
    this.contentDiv.addEventListener("paste", this.onPaste);
    this.contentDiv.addEventListener("keydown", this.onKeydown, true);
    this.contentDiv.addEventListener("beforeinput", this.onBeforeInput);
    this.contentDiv.addEventListener("input", this.onInput);
    this.contentDiv.addEventListener("click", this.onClick);

    document.addEventListener("selectionchange", this.onSelectionChange);
    document.addEventListener("click", this.onDocumentClick);
  }

  disconnectedCallback() {
    this._active = false;
    if (this.debug) {
      console.log("disconnectedCallback");
    }

    if (this.network && this.network.readyState === WebSocket.OPEN) {
      this.network.close();
    }

    this.contentDiv?.addEventListener("copy", this.onCopy);
    this.contentDiv?.removeEventListener("paste", this.onPaste);
    this.contentDiv?.removeEventListener("keydown", this.onKeydown);
    this.contentDiv?.removeEventListener("beforeinput", this.onBeforeInput);
    this.contentDiv?.removeEventListener("input", this.onInput);
    this.contentDiv?.removeEventListener("click", this.onClick);

    document.removeEventListener("selectionchange", this.onSelectionChange);
    document.removeEventListener("click", this.onDocumentClick);
  }

  componentWillUnmount() {
    this.disconnect();
  }

  clearSelection() {
    window.getSelection()?.removeAllRanges();
    this.clearCurrentRogueRange();
  }

  clearCurrentRogueRange() {
    this.curRogueRange = null;
    this.removeHighlight("selection-");
  }

  get curRogueRange(): RogueRange | null {
    return this._curRogueRange;
  }

  getCurRogueRangeContent(): string {
    if (this.debug) {
      console.log("getCurRogueRangeContent");
    }

    if (!this.curRogueRange) {
      return "";
    }

    const [startID, afterID] = this.curRogueRange;
    if (arraysEqual(startID, afterID)) {
      return "";
    }

    const endID = this.rogue.TotLeftOf(afterID);
    if (endID.error) {
      throw new Error(endID.error);
    }

    const getMarkdown = this.rogue.GetMarkdown(startID, endID);

    if (getMarkdown.error) {
      console.error(getMarkdown.error);
      return "";
    }

    return getMarkdown.text;
  }

  set editing(value: boolean) {
    if (this.debug) {
      console.log("set editing", value);
    }

    if (value === this._editing) {
      return;
    }
    this._editing = value;
    this.notifySubscribers("editing", value);

    // update other editors that this editor is no longer editing
    if (this._editing === false && this.curRogueRange) {
      this.sendCursorUpdate(this.curRogueRange);
    }
  }

  get editing(): boolean {
    return this._editing;
  }

  set curRogueRange(value: RogueRange | null) {
    if (this.debug) {
      console.log("set curRogueRange", JSON.stringify(value));
    }

    if (value === this._curRogueRange) {
      return;
    }
    this._curRogueRange = value;

    if (!this._curRogueRange) {
      this.selectedHtml = null;
      return;
    }

    if (arraysEqual(this._curRogueRange[0], this._curRogueRange[1])) {
      this.selectedHtml = null;
      return;
    }

    const [startID, afterID] = this._curRogueRange;

    const html = this.getCurHtml(startID, afterID, false);
    if (html.error) {
      console.error(html.error);
      return;
    }

    this.selectedHtml = html;
  }

  get baseAddress(): string | null {
    return this._baseAddress;
  }

  get address(): string | null {
    return this._address;
  }

  get editorMode(): EditorMode {
    return this._editorMode;
  }

  get addressDescription(): string | null {
    return this._addressDescription;
  }

  async getDocStats(): Promise<DocStats | null> {
    const currentContentAddress = this.currentContentAddress();
    if (
      currentContentAddress &&
      currentContentAddress == this._docStatsAddress
    ) {
      return this._docStats;
    }

    if (!this.rogue) {
      return null;
    }

    const stats = this.rogue.DocStats();
    this._docStats = stats;
    this._docStatsAddress = currentContentAddress;
    return stats;
  }

  async getOpStats(): Promise<OpStats | null> {
    const currentContentAddress = this.currentContentAddress();
    if (
      currentContentAddress &&
      currentContentAddress == this._opStatsAddress
    ) {
      return this._opStats;
    }

    if (!this.rogue) {
      return null;
    }

    const stats = this.rogue.OpStats();
    this._opStats = stats;
    this._opStatsAddress = currentContentAddress;
    return stats;
  }

  setHistoryDiff(startAddr: string, endAddr: string) {
    if (this.debug) {
      console.log("setHistoryDiff", startAddr, endAddr);
    }

    this._baseAddress = startAddr;
    this._address = endAddr;
    this._showDiffHighlights = true;
    this._editorMode = "history";

    this.notifySubscribers("baseAddress", this._baseAddress);
    this.notifySubscribers("address", this._address);
    this.notifySubscribers("showDiffHighlights", this._showDiffHighlights);

    this.renderRogue();
  }

  swapHistoryDiff() {
    if (this.debug) {
      console.log("swapHistoryDiff");
    }

    if (!this._baseAddress) {
      return;
    }

    this._baseAddress = this._address;
    this._address = this._baseAddress;
    this._showDiffHighlights = true;
    this._editorMode = "history";

    this.notifySubscribers("baseAddress", this._baseAddress);
    this.notifySubscribers("address", this._address);
    this.notifySubscribers("editorMode", this._editorMode);
    this.notifySubscribers("showDiffHighlights", this._showDiffHighlights);

    this.renderRogue();
  }

  setAddress(value: string, mode: EditorMode, showDiffHighlights = true) {
    if (this.debug) {
      console.log("setAddress", value);
    }

    if (value === this._address && mode === this._editorMode) {
      return;
    }

    if (this._address !== value) {
      this._baseAddress = null;
      this._address = value;
      this._showDiffHighlights = showDiffHighlights;
      this._editorMode = mode;

      this.notifySubscribers("baseAddress", this._baseAddress);
      this.notifySubscribers("address", this._address);
      this.notifySubscribers("editorMode", this._editorMode);
      this.notifySubscribers("showDiffHighlights", this._showDiffHighlights);
    }

    this.renderRogue();
  }

  toggleXRayMode() {
    this._editorMode = this._editorMode === "xray" ? "edit" : "xray";
    this._address = null;
    this.notifySubscribers("editorMode", this._editorMode);
    this.renderRogue();
  }

  setAddressDescription(value: string) {
    if (this.debug) {
      console.log("setAddressDescription", value);
    }

    this._addressDescription = value;
    this.notifySubscribers("addressDescription", value);
  }

  resetAddress() {
    if (this.debug) {
      console.log("resetAddress");
    }

    this._baseAddress = null;
    this._address = null;
    this._addressDescription = null;
    this._editorMode = "edit";
    this.notifySubscribers("baseAddress", null);
    this.notifySubscribers("address", null);
    this.notifySubscribers("editorMode", null);
    this.notifySubscribers("addressDescription", null);
    this.renderRogue();
  }

  set showDiffHighlights(value: boolean) {
    if (this.debug) {
      console.log("set showDiffHighlights", value);
    }

    this._showDiffHighlights = value;
    this.notifySubscribers("showDiffHighlights", value);
  }

  get showDiffHighlights(): boolean {
    return this._showDiffHighlights;
  }

  set canUndo(value: boolean) {
    this._canUndo = value;
    this.notifySubscribers("canUndo", value);
  }
  get canUndo(): boolean {
    return this._canUndo;
  }

  set canRedo(value: boolean) {
    this._canRedo = value;
    this.notifySubscribers("canRedo", value);
  }
  get canRedo(): boolean {
    return this._canRedo;
  }

  set enabled(value: boolean) {
    if (this.debug) {
      console.log("set enabled", value);
    }

    if (value === this._enabled) {
      return;
    }
    this._enabled = value;

    this.notifySubscribers("enabled", value);
    this.renderRogue();
  }

  get enabled(): boolean {
    return this._enabled;
  }

  subscribe<T>(
    attribute: SubscribableAttr,
    callback: CallbackFunction<T>,
  ): void {
    if (!this.subscribers[attribute]) {
      this.subscribers[attribute] = [];
    }
    this.subscribers[attribute].push(callback);
  }

  unsubscribe<T>(
    attribute: SubscribableAttr,
    callback: CallbackFunction<T>,
  ): void {
    if (this.subscribers[attribute]) {
      this.subscribers[attribute] = this.subscribers[attribute].filter(
        (cb) => cb !== callback,
      );
    }
  }

  private notifySubscribers<T>(attribute: SubscribableAttr, newValue: T): void {
    if (this.debug) {
      console.log("notifySubscribers:", attribute, newValue);
    }

    if (this.subscribers[attribute]) {
      this.subscribers[attribute].forEach((callback) => callback(newValue));
    }
  }

  set activeComments(value: string[]) {
    this._activeComments = value;
    this.notifySubscribers("activeComments", value);
  }

  get activeComments(): string[] {
    return this._activeComments;
  }

  set syncing(value: boolean) {
    if (this.debug) {
      console.log("set syncing", value);
    }

    if (value === this._syncing) {
      return;
    }
    this._syncing = value;
    this.notifySubscribers("syncing", value);
  }

  get syncing(): boolean {
    return this._syncing;
  }

  set loaded(value: boolean) {
    if (this.debug) {
      console.log("set loaded", value);
    }

    if (value === this._loaded) {
      return;
    }
    this._loaded = value;
    this.notifySubscribers("loaded", value);
  }

  get loaded(): boolean {
    return this._loaded;
  }

  set applyingDiff(value: boolean) {
    if (this.debug) {
      console.log("set applyingDiff", value);
    }

    if (value === this._applyingDiff) {
      return;
    }
    this._applyingDiff = value;
    this.notifySubscribers("applyingDiff", value);
  }

  get applyingDiff(): boolean {
    return this._applyingDiff;
  }

  set connected(value: boolean) {
    if (this.debug) {
      console.log("set connected", value);
    }

    if (value === this._connected) {
      return;
    }
    this._connected = value;
    this.notifySubscribers("connected", value);
  }

  get connected(): boolean {
    return this._connected;
  }

  set curSpanFormat(value: Record<string, any>) {
    if (this.debug) {
      console.log("set curSpanFormat", value);
    }

    if (shallowEqual(value, this._curSpanFormat)) {
      return;
    }
    this._curSpanFormat = value;
    this.notifySubscribers("curSpanFormat", value);
  }

  get curSpanFormat(): Record<string, any> {
    return this._curSpanFormat;
  }

  set curLineFormat(value: Record<string, any>) {
    if (this.debug) {
      console.log("set curLineFormat", value);
    }

    if (shallowEqual(value, this._curLineFormat)) {
      return;
    }
    this._curLineFormat = value;
    this.notifySubscribers("curLineFormat", value);
  }

  get curLineFormat(): Record<string, any> {
    return this._curLineFormat;
  }

  get selectedHtml(): string | null {
    return this._selectedHtml;
  }

  set selectedHtml(value: string | null) {
    if (this.debug) {
      console.log("set selectedHtml", value);
    }

    if (value === this._selectedHtml) {
      return;
    }
    this._selectedHtml = value;
    this.notifySubscribers("selectedHtml", value);
  }

  get lastEdit(): Date | null {
    return this._lastEdit;
  }

  set lastEdit(value: Date) {
    if (this.debug) {
      console.log("set lastEdit", value);
    }

    if (value === this._lastEdit) {
      return;
    }
    this._lastEdit = value;
    this.notifySubscribers("lastEdit", value);
  }

  currentContentAddress(): string | null {
    if (!this.contentDiv) {
      console.error("contentDiv not found");
      return null;
    }

    if (!this.rogue) {
      return null;
    }

    const resp = this.rogue.GetAddress();
    if (resp && resp.error) {
      console.error(resp.error);
      return null;
    }

    return resp;
  }

  disable() {
    if (this.debug) {
      console.log("disable");
    }

    if (!this.enabled) {
      return;
    }

    if (!this.contentDiv) {
      console.error("contentDiv not found");
      return;
    }
    this.enabled = false;
    this.contentDiv.contentEditable = "false";
    this.contentDiv.spellcheck = false;
    this.contentDiv.translate = false;
  }

  enable() {
    if (this.debug) {
      console.log("enable");
    }

    if (this.enabled) {
      return;
    }

    if (!this.contentDiv) {
      console.error("contentDiv not found");
      return;
    }
    this.enabled = true;
    this.contentDiv.contentEditable = "true";
    this.contentDiv.spellcheck = true;
    this.contentDiv.translate = false;
  }

  formatImage(
    src: string,
    alt: string,
    width: string,
    height: string,
  ): [string, number] | undefined {
    if (!this.curRogueRange) {
      return;
    }
    const [startId, afterId] = this.curRogueRange;

    const startIx = this.rogue.GetIndex(startId);
    if (startIx.error) {
      console.error(startIx.error);
      return;
    }
    const afterIx = this.rogue.GetIndex(afterId);
    if (afterIx.error) {
      console.error(afterIx.error);
      return;
    }

    const format = {
      img: src,
      alt: alt,
      width: width,
      height: height,
    };

    const mop = this.rogue.Format(
      startIx.visible,
      afterIx.visible - startIx.visible,
      format,
    );

    if (mop.error) {
      console.error(mop.error);
      return;
    }

    // make sure op is a format op
    if (mop.length == 0 || (mop[0] != 2 && mop[0] != 6)) {
      return;
    }

    this.sendOp(mop);
    this.renderRogue();
    this.setCurSpanFormat();

    if (mop[0] == 2) {
      return mop[2]; // return the startID of the image format
    } else if (mop[0] == 6) {
      for (const op in mop[2]) {
        // get the startID from the first formatOp
        if (op[0] == 2) {
          return op[2];
        }
      }
    }
  }

  format(style: string, value: string) {
    if (this.debug) {
      console.log("format", style, value);
    }

    if (!this.curRogueRange) {
      return;
    }
    const [startId, afterId] = this.curRogueRange;

    // Special handling for switching between ul and ol
    if (style === "ul" || style === "ol") {
      const curFormat = this.rogue.GetCurLineFormat(
        this.curRogueRange[0],
        this.curRogueRange[1],
      );

      if (curFormat.error) {
        console.error(curFormat.error);
        return;
      }

      if (style === "ol" && curFormat["ul"]) {
        value = curFormat["ul"];
      } else if (style === "ul" && curFormat["ol"]) {
        value = curFormat["ol"];
      }
    }

    const startIx = this.rogue.GetIndex(startId);
    if (startIx.error) {
      console.error(startIx.error);
      return;
    }
    const afterIx = this.rogue.GetIndex(afterId);
    if (afterIx.error) {
      console.error(afterIx.error);
      return;
    }

    let format = { [style]: value };
    if (style === "text") {
      format = {};
    }
    const mop = this.rogue.Format(
      startIx.visible,
      afterIx.visible - startIx.visible,
      format,
    );

    if (mop.error) {
      console.error(mop.error);
      return;
    }

    this.sendOp(mop);
    this.renderRogue(mop);
    this.setCurSpanFormat();
  }

  undo() {
    if (this.debug) {
      console.log("undo");
    }

    let [startID, endID]: [any, any] = this.curRogueRange || [null, null];
    if (!this.curRogueRange) {
      startID = this.rogue.GetFirstTotID();
      if (startID.error) {
        console.error(startID.error);
        return;
      }

      endID = this.rogue.GetLastTotID();
      if (endID.error) {
        console.error(endID.error);
        return;
      }
    }

    const op = this.rogue.Undo(startID, endID);

    if (op.error) {
      console.error(op.error);
      return;
    }

    this.sendOp(op);
    this.renderRogue();
  }

  rewind(startid: Id, endid: Id) {
    if (this.debug) {
      console.log("rewind", startid, endid);
    }

    if (!this.address || !startid || !endid) {
      console.error("No address", startid, endid);
      return;
    }

    const mop = this.rogue.Rewind(startid, endid, this.address);
    if (mop.error) {
      console.error(mop.error);
      return;
    }

    this.sendOp(mop);
    this.renderRogue();
  }

  rewindAll() {
    if (this.debug) {
      console.log("rewind all");
    }

    const firstID = this.rogue.GetFirstTotID();
    if (firstID.error) {
      console.error(firstID.error);
      return;
    }

    const lastID = this.rogue.GetLastTotID();
    if (lastID.error) {
      console.error(lastID.error);
      return;
    }

    this.rewind(firstID, lastID);
  }

  redo() {
    if (this.debug) {
      console.log("redo");
    }

    const mop = this.rogue.Redo();

    if (mop.error) {
      console.error(mop.error);
      return;
    }

    this.sendOp(mop);
    this.renderRogue();
  }

  async copy() {
    if (this.debug) {
      console.log("copy");
    }

    if (!this.curRogueRange) {
      await this.copyDoc();
      return;
    }

    await this.copyBetween(this.curRogueRange[0], this.curRogueRange[1]);
  }

  async cut() {
    if (this.debug) {
      console.log("cut");
    }

    await this.copy();
    this.deleteText();
    this.renderRogue();
  }

  async copyDoc() {
    if (this.debug) {
      console.log("copyDoc");
    }

    const firstID = this.rogue.GetFirstTotID();
    if (firstID.error) {
      console.error(firstID.error);
      return;
    }

    const lastID = this.rogue.GetLastTotID();
    if (lastID.error) {
      console.error(lastID.error);
      return;
    }

    await this.copyBetween(firstID, lastID);
  }

  async copyBetween(startID: Id, endID: Id) {
    if (this.debug) {
      console.log("copyBetween", startID, endID);
    }

    const innerHtml = this.getCurHtml(startID, endID, false);
    if (innerHtml.error) {
      console.error(innerHtml.error);
      return;
    }

    const html = `<div ${RogueHTMLCopyMarker} draft-id="${this.docID}">${innerHtml}</div>`;

    const plaintext = this.getCurPlaintext(startID, endID);
    if (plaintext.error) {
      console.error(plaintext.error);
      return;
    }

    console.log("📋 copy", { plaintext: plaintext, html: html });

    await navigator.clipboard.write([
      new ClipboardItem({
        "text/plain": new Blob([plaintext], { type: "text/plain" }),
        "text/html": new Blob([html], { type: "text/html" }),
      }),
    ]);
  }

  toggleSpanFormat(style: string) {
    if (this.debug) {
      console.log("toggleSpanFormat", style);
    }

    const f = this.curSpanFormat[style];
    let v = "true";
    if (f === "true") {
      v = "";
    }

    if (
      this.curRogueRange &&
      arraysEqual(this.curRogueRange[0], this.curRogueRange[1])
    ) {
      // it's just a cursor right now, so we need to stash the format changes
      this.curSpanFormat = {
        ...this.curSpanFormat,
        [style]: v,
      };
    } else {
      this.format(style, v);
    }
  }

  get container(): HTMLElement | null {
    return this.contentDiv;
  }

  get authorId(): string {
    return this.operationManager?.authorId || "";
  }

  withAuthorId(value: string, fn: () => void): void {
    const originalAuthorId = this.authorId;
    this.authorId = value;

    try {
      fn();
    } finally {
      this.authorId = originalAuthorId;
    }
  }

  withAuthorPrefix(prefix: string, fn: () => void): void {
    this.withAuthorId(prefix + this.authorId, fn);
  }

  set authorId(value: string) {
    if (this.operationManager) {
      this.operationManager.authorId = value;
    } else {
      throw new Error("No operation manager");
    }
  }

  getIndex(id: Id): { error?: string; visible: number; total: number } {
    return this.rogue.GetIndex(id) as {
      error?: string;
      visible: number;
      total: number;
    };
  }

  createHighlight(key: string, ids: [Id, Id], options: HighlightOptions) {
    if (this.debug) {
      console.log("createHighlight", key, ids, options);
    }

    if (!this.highlights) {
      console.error("highlights not found");
      return;
    }

    this.highlights.highlightRange(key, ids, null, options);
  }

  createHighlightWithAddress(
    key: string,
    ids: [Id, Id],
    address: string,
    options: HighlightOptions,
  ) {
    if (this.debug) {
      console.log("createHighlightWithAddress", key, ids, options);
    }

    if (!this.highlights) {
      console.error("highlights not found");
      return;
    }

    this.highlights.highlightRange(key, ids, address, options);
  }

  removeHighlight(ident: string) {
    if (this.debug) {
      console.log("removeHighlight", ident);
    }

    if (!this.highlights) {
      console.error("highlights not found");
      return;
    }
    this.highlights.removeHighlight(ident);
  }

  setCurSpanFormat() {
    if (this.debug) {
      console.log("setCurSpanFormat");
    }

    if (this.editorMode !== "edit") {
      return;
    }

    if (!this.curRogueRange) {
      this.curSpanFormat = {};
      return;
    }

    const format = this.rogue.GetCurSpanFormat(
      this.curRogueRange[0],
      this.curRogueRange[1],
    );

    if (format.error) {
      console.error(format.error);
      return;
    }

    this.curSpanFormat = format;
  }

  setCurLineFormat() {
    if (this.debug) {
      console.log("setCurLineFormat");
    }

    if (this.editorMode !== "edit") {
      return;
    }

    if (!this.curRogueRange) {
      this.curLineFormat = {};
      return;
    }

    const format = this.rogue.GetCurLineFormat(
      this.curRogueRange[0],
      this.curRogueRange[1],
    );

    if (format.error) {
      console.error(format.error);
      return;
    }

    this.curLineFormat = format;
  }

  showCommentHighlights(
    comments: [{ eventId: string; startId: Id; endId: Id; address: string }],
  ) {
    for (const comment of comments) {
      try {
        this.createHighlightWithAddress(
          "comment-" + comment.eventId,
          [comment.startId, comment.endId],
          comment.address,
          {
            caret: false,
            dataAttrs: {
              commentId: comment.eventId,
            },
            styles: {
              borderBottom: HIGHLIGHT_STYLES.inactive.borderBottom,
              backgroundColor: HIGHLIGHT_STYLES.inactive.backgroundColor,
            },
            eventListeners: {
              click: () => {
                // if the comment is already in the active comments array, do nothing
                if (this.activeComments.includes(comment.eventId)) {
                  return;
                }
                this.showActiveCommentHighlight(comment);
                this.activeComments = [...this.activeComments, comment.eventId];
              },
              mouseover: function () {
                const highlight = this as RogueHighlight;
                highlight.updateStyles({
                  backgroundColor: HIGHLIGHT_STYLES.active.backgroundColor,
                  borderBottom: HIGHLIGHT_STYLES.active.borderBottom,
                });
              },
              mouseout: function () {
                const highlight = this as RogueHighlight;
                highlight.updateStyles({
                  backgroundColor: HIGHLIGHT_STYLES.inactive.backgroundColor,
                  borderBottom: HIGHLIGHT_STYLES.inactive.borderBottom,
                });
              },
            },
          },
        );
      } catch (e) {
        console.error("showCommentHighlights", e);
      }
    }
  }

  hideCommentHighlights(eventIds?: string[]) {
    if (!eventIds || eventIds.length === 0) {
      this.highlights.removeHighlightsWithPrefix("comment-");
    } else {
      for (const eventId of eventIds) {
        this.removeHighlight("comment-" + eventId);
      }
    }
  }

  showActiveCommentHighlight(comment: {
    eventId: string;
    startId: Id;
    endId: Id;
    address: string;
  }) {
    this.createHighlightWithAddress(
      "comment-active-" + comment.eventId,
      [comment.startId, comment.endId],
      comment.address,
      {
        styles: {
          backgroundColor: HIGHLIGHT_STYLES.active.backgroundColor,
          borderBottom: HIGHLIGHT_STYLES.active.borderBottom,
        },
        caret: false,
        dataAttrs: {
          commentId: comment.eventId,
        },
      },
    );
  }

  activateComment({ eventId }: { eventId: string }) {
    // this.activeComments = [eventId];
    const highlight = document.getElementById(`comment-active-${eventId}`);
    if (highlight) {
      highlight.firstElementChild?.scrollIntoView({
        behavior: "smooth",
        block: "center",
      });
    }
  }

  hideActiveCommentHighlight(eventId: string) {
    this.removeHighlight("comment-active-" + eventId);
  }

  hideAllActiveCommentHighlights() {
    this.highlights.removeHighlightsWithPrefix("comment-active-");
  }

  onSelectionChange = () => {
    if (this.debug) {
      console.log("onSelectionChange");
    }

    const selection = document.getSelection();
    if (selection && selection.rangeCount > 0) {
      // console.log("selectionchange: selection.rangeCount > 0", selection);
      const range = selection.getRangeAt(0);
      if (!this.contentDiv) {
        console.error("contentDiv not found");
        return;
      }

      // Check if the selection is within contentDiv
      if (
        range.commonAncestorContainer === this.contentDiv ||
        this.contentDiv.contains(range.commonAncestorContainer)
      ) {
        this.curRogueRange = this._getCurRogueIds();
        this.setCurLineFormat();
        this.setCurSpanFormat();

        if (!this.curRogueRange) {
          return;
        }

        this.sendCursorUpdate(this.curRogueRange);

        if (this.debug) {
          console.log("rogue editor address", this.address);
        }

        if (!this.address) {
          this.createHighlight("selection-", this.curRogueRange, {
            styles: {
              backgroundColor: "hsla(var(--reviso))",
            },
            caret: false,
          });
        } else {
          this.removeHighlight("selection-");
        }
      }
    }
  };

  onCopy = (e: ClipboardEvent) => {
    if (this.debug) {
      console.log("onCopy");
    }

    e.preventDefault();
    this.copy();
  };

  getPreferredImageType(e: ClipboardEvent): File | null {
    const items = e.clipboardData?.items;
    const files = e.clipboardData?.files;
    let index = supportedImageMimeTypes.length;
    let result: File | null = null;

    if (items) {
      for (const item of items) {
        const ix = supportedImageMimeTypes.indexOf(item.type);
        if (ix < index) {
          index = ix;
          result = item.getAsFile();
        }
      }
    }

    if (files) {
      for (const file of files) {
        const ix = supportedImageMimeTypes.indexOf(file.type);
        if (ix < index) {
          index = ix;
          result = file;
        }
      }
    }

    return result;
  }

  checkImageUploadStatus(
    targetID: [string, number],
    url: string,
    docId: string,
    imageId: string,
  ) {
    console.log("checkImageUploadStatus", docId, imageId);

    this.getImage(docId, imageId)
      .then((image) => {
        if (!image) {
          console.error("Image not found", imageId);
          return;
        }

        console.log("IMAGE", image);

        if (image.status === "LOADING") {
          setTimeout(() => {
            this.checkImageUploadStatus(targetID, url, docId, imageId);
          }, 1000);
        } else if (image.status === "SUCCESS") {
          const mop = this.rogue.FormatLineByID(targetID, { img: image.url });
          if (mop.error) {
            console.error(mop.error);
            return;
          }

          this.sendOp(mop);
          this.renderRogue();
          this.setCurSpanFormat();
        } else {
          console.error("Image upload failed", image);
        }
      })
      .catch((error) => {
        console.error("Error checking image upload status", error);
      });
  }

  pasteImage(file: File) {
    // Print out all items and files
    console.log(`File ${file.name}:`, file);

    if (!this.docID) {
      console.error("No docID available. Cannot upload image.");
      return;
    }

    if (!this.uploadImage) {
      console.error("uploadImage function is not defined.");
      return;
    }

    this.uploadImage(file, this.docID)
      .then((image) => {
        console.log(`Successfully uploaded image ${file.name}:`, image);
        // Here you might want to do something with the result,
        // like inserting the image into your editor
        if (image) {
          const targetID = this.formatImage(image.url, "", "", "");

          if (targetID) {
            this.checkImageUploadStatus(
              targetID,
              image.url,
              image.docId,
              image.id,
            );
          }
        }
      })
      .catch((error) => {
        console.error(`Error uploading image ${file.name}:`, error);
        if (error.graphQLErrors) {
          error.graphQLErrors.forEach(
            ({
              message,
              locations,
              path,
            }: {
              message: string;
              locations?: ReadonlyArray<{ line: number; column: number }>;
              path?: ReadonlyArray<string | number>;
            }) => {
              console.log(
                `[GraphQL error]: Message: ${message}, Location: ${JSON.stringify(locations)}, Path: ${path}`,
              );
            },
          );
        }
        if (error.networkError) {
          console.log(`[Network error]: ${error.networkError}`);
        }
        // Here you might want to show an error message to the user
      });
  }

  findContentParent(node: Element): Element {
    let current = node;
    while (current.parentElement) {
      if (current.parentElement.classList.contains("content")) {
        return current;
      }
      current = current.parentElement;
    }
    return node;
  }

  scrubInit(isWholeDoc: boolean): number {
    if (this.debug) {
      console.log("scrubInit");
    }

    if (!this.rogue) {
      return -1;
    }

    if (this._editorMode === "scrub") {
      return this.rogue.ScrubMax();
    }

    this._editorMode = "scrub";
    this.disable();

    if (this.curRogueRange && !isWholeDoc) {
      this.scrubMode = "partial";

      const x = this.rogue.ScrubInit(
        this.curRogueRange[0],
        this.curRogueRange[1],
      );

      if (x.error) {
        console.error(x.error);
        return -1;
      }

      return x;
    } else {
      this.scrubMode = "full";

      const x = this.rogue.ScrubInit();

      if (x.error) {
        console.error(x.error);
        return -1;
      }

      return x;
    }
  }

  scrubTo(n: number) {
    if (this.debug) {
      console.log("scrubTo", n);
    }

    if (!this.rogue || this._editorMode !== "scrub") {
      return;
    }

    const resp = this.rogue.ScrubTo(n);

    if (resp.error) {
      console.error(resp.error);
      return;
    }

    if (!resp) {
      return;
    }

    this.replaceBetweenNodes(resp.firstBlockID, resp.lastBlockID, resp.html);
    this.moveCursorTo([resp.cursorStartID, resp.cursorEndID]);
  }

  scrubExit() {
    if (this.debug) {
      console.log("scrubExit");
    }

    if (!this.rogue) {
      return;
    }

    this.rogue.ScrubExit();
    this._editorMode = "edit";
    this.renderRogue();
  }

  scrubRevert() {
    if (this.debug) {
      console.log("scrubRevert");
    }

    const resp = this.rogue.ScrubRevert();
    if (resp.error) {
      console.error(resp.error);
      return;
    }
    this.sendOp(resp.mop);
    this._editorMode = "edit";
    this.renderRogue();
    this.moveCursorTo([resp.cursorID, resp.cursorID]);
  }

  replaceBetweenNodes(
    startID: [string, string],
    endID: [string, string],
    newHtmlContent: string,
  ) {
    if (this.debug) {
      console.log("replaceBetweenNodes", startID, endID, newHtmlContent);
    }

    // Find nodes by data-rid
    let startNode = document.querySelector(
      `[data-rid="${startID[0]}_${startID[1]}"]`,
    );
    let endNode = document.querySelector(
      `[data-rid="${endID[0]}_${endID[1]}"]`,
    );

    if (!startNode && endNode) {
      startNode = endNode;
    } else if (startNode && !endNode) {
      endNode = startNode;
    }

    if (!startNode || !endNode) {
      throw new Error(
        `Could not find nodes with data-rid ${startID}, ${endID}`,
      );
    }

    startNode = this.findContentParent(startNode as Element);
    endNode = this.findContentParent(endNode as Element);

    // Validate nodes have same parent
    if (!startNode.parentNode || !endNode.parentNode) {
      throw new Error("Nodes must have parent elements");
    }

    if (startNode.parentNode !== endNode.parentNode) {
      throw new Error("Start and end nodes must share the same parent");
    }

    // Create a range
    const range = document.createRange();
    range.setStartBefore(startNode); // Changed from setStartAfter
    range.setEndAfter(endNode); // Changed from setEndBefore

    // Delete existing content (including the nodes)
    range.deleteContents();

    // Insert new content directly using Range.createContextualFragment()
    const fragment = range.createContextualFragment(newHtmlContent);
    range.insertNode(fragment);
  }

  onPaste = (e: ClipboardEvent) => {
    if (this.debug) {
      console.log("onPaste");
    }

    e.preventDefault();
    if (!this.curRogueRange) {
      return;
    }

    const file = this.getPreferredImageType(e);
    if (file) {
      this.pasteImage(file);
      return;
    }

    // Print out all items and files
    const items = e.clipboardData?.items;
    const files = e.clipboardData?.files;

    const pasteItems: PasteItem[] = [];
    let pasteFromRogue = false;

    if (items) {
      for (let i = 0; i < items.length; i++) {
        const item = items[i];
        const data = e.clipboardData?.getData(item.type);

        pasteItems.push({
          kind: item.kind,
          mime: item.type,
          data: data,
        });

        if (item.type === "text/html") {
          if (data && data.includes(RogueHTMLCopyMarker)) {
            const match = data.match(RogueHTMLDraftIdRegex);
            if (match && match[1]) {
              if (this.debug) {
                console.log("paste from rogue", match[1]);
              }
              if (match[1] === this.docID) {
                pasteFromRogue = true;
                break;
              }
            }
          }
        }
      }
    }

    // TODO: handle files
    if (files) {
      for (let i = 0; i < files.length; i++) {
        console.log(`File ${i}:`, files[i]);
      }
    }

    const prefex = pasteFromRogue ? "" : "#";

    this.withAuthorPrefix(prefex, () => {
      let [startIx, endIx] = this.getCurRogueIndexes();
      if (startIx < 0) {
        if (this.rogue.Size() > 1) {
          console.error("startIx not found for newline");
          this.renderRogue();
          return;
        }
        startIx = 0;
        endIx = 0;
      }

      // remove the selection highlight
      this.removeHighlight("selection-");

      const contentAddressBefore = this.currentContentAddress();

      const result = this.rogue.Paste(
        startIx,
        endIx - startIx,
        this.curSpanFormat,
        pasteItems,
      );

      if (result.error) {
        console.error(result.error);
        this.renderRogue();
        return;
      }

      for (const op of result.ops) {
        this.sendOp(op);
      }
      this.curRogueRange = [result.cursorID, result.cursorID];

      const contentAddressAfter = this.currentContentAddress();

      console.log("Pasted", pasteItems);

      if (!pasteFromRogue) {
        this.sendEvent("paste", {
          contentAddressBefore,
          contentAddressAfter,
        });
      }

      this.renderRogue();
    });
  };

  onKeydown = (ev: KeyboardEvent): void => {
    if (this.debug) {
      console.log(
        "onKeydown",
        ev.key,
        ev.ctrlKey ? "ctrl" : "",
        ev.altKey ? "alt" : "",
        ev.metaKey ? "meta" : "",
      );
    }

    if (ev.altKey || ev.ctrlKey || ev.metaKey) {
      if (ev.key === "z") {
        ev.preventDefault();
        if (ev.shiftKey) {
          this.redo();
        } else {
          this.undo();
        }
        return;
      }
      if (ev.key === "c") {
        ev.preventDefault();
        this.copy();
        return;
      }

      if (ev.key === "b") {
        ev.preventDefault();
        this.toggleSpanFormat("b");
        return;
      }

      if (ev.key === "i") {
        ev.preventDefault();
        this.toggleSpanFormat("i");
        return;
      }

      if (ev.key === "u") {
        ev.preventDefault();
        this.toggleSpanFormat("u");
        return;
      }

      return;
    }

    if (ev.shiftKey) {
      if (ev.key === "Tab") {
        ev.preventDefault();
        const shiftTabIdentifier = "(1+4cT5lP9";
        this.insertText(shiftTabIdentifier);
      } else if (ev.key === "Enter") {
        ev.preventDefault();
        const shiftEnterIdentifier = "oLEcI0yPY9";
        this.insertText(shiftEnterIdentifier);
      }
      return;
    }

    if (ev.key === "Shift" || ev.key === "CapsLock") {
      return;
    }

    if (ev.key.startsWith("Arrow")) {
      return;
    }

    ev.preventDefault();

    if (ev.key === "Enter") {
      this.insertText("\n");
    } else if (ev.key === "Escape") {
      this.contentDiv?.blur();
    } else if (ev.key === "Tab") {
      this.insertText("\t");
    } else if (ev.key === "Backspace") {
      this.deleteText();
    } else if (ev.key === "Delete") {
      this.deleteText(true);
    } else if (/^[\x20-\x7E]$/.test(ev.key)) {
      // Check for visible ASCII characters
      this.insertText(ev.key);
    } else {
      console.log(`Unhandled key: ${ev.key}`);
    }
  };

  onClick = (event: MouseEvent) => {
    // open URLs if the user control clicks them
    if (event.ctrlKey || (event.metaKey && this.curSpanFormat["a"])) {
      const url = this.curSpanFormat["a"];
      window.open(url, "_blank");
    }

    const target = event.target as HTMLElement;
    if (
      target.tagName.toLowerCase() === "del" ||
      target.tagName.toLowerCase() === "ins"
    ) {
      event.preventDefault();
      event.stopPropagation();
      this.deltaManager.show(target);
      return;
    }

    if (target.tagName.toLowerCase() === "a") {
      event.preventDefault(); // Prevent default link behavior
      const selection = window.getSelection();
      if (selection) {
        if (event.target === this._lastTargetAnchor) {
          this._lastTargetAnchor = null;
        } else {
          this._lastTargetAnchor = event.target as HTMLAnchorElement;

          // If the selection doesn't match, proceed with selecting the link content
          const range = document.createRange();
          let firstNode: Node = target,
            lastNode: Node = target;

          // Find the first and last text nodes
          while (firstNode.firstChild) firstNode = firstNode.firstChild;
          while (lastNode.lastChild) lastNode = lastNode.lastChild;

          // Ensure we're dealing with text nodes
          if (
            firstNode.nodeType === Node.TEXT_NODE &&
            lastNode.nodeType === Node.TEXT_NODE
          ) {
            range.setStart(firstNode, 0);
            range.setEnd(lastNode, (lastNode as Text).length);
            selection.removeAllRanges();
            selection.addRange(range);
          }
        }
      }
    } else {
      this._lastTargetAnchor = null;
    }

    this.hideAllActiveCommentHighlights();
    this.highlights.triggerClicksForCoordinates(event.clientX, event.clientY, {
      onFirstClick: () => (this.activeComments = []),
    });
  };

  onDocumentClick = (_e: MouseEvent) => {
    this.deltaManager.close();
  };

  // Called when the contenteditable div is updated
  onBeforeInput = (e: any) => {
    if (this.debug) {
      console.log("onBeforeInput");
    }

    if (e.inputType === "insertReplacementText") {
      //we will rely on onInput for this
      console.log("irt", e);
      return;
    }

    e.preventDefault();
    // console.log(e.inputType, e.data, e);
    switch (e.inputType) {
      case "insertText":
      case "insertCompositionText":
        // Different input types that show up with the emoji picker
        // or other input devices
        this.insertText(e.data);
        break;
      case "historyUndo":
        console.error("historyUndo not implemented yet, srry");
        break;
      case "deleteByCut":
        this.cut();
        break;
      case "deleteContentBackward":
        this.deleteText(false);
        break;
      case "deleteContentForward":
        this.deleteText(true);
        break;
      case "deleteWordBackward":
        this.deleteWord(false);
        break;
      case "deleteSoftLineBackward":
        this.deleteLine();
        break;
      case "deleteHardLineBackward":
        this.deleteLine();
        break;
      default:
        console.error("Unknown input type:", e.inputType);
    }
  };

  onInput = (e: any) => {
    if (this.debug) {
      console.log("onInput");
    }

    e.preventDefault();
    if (e.inputType === "insertReplacementText") {
      //we will rely on onInput for this
      this.insertText(e.data);
    }
  };

  insertBlankState() {
    if (this.debug) {
      console.log("insertBlankState");
    }

    const firstID = this.rogue.GetFirstID();
    if (firstID.error) {
      console.error(firstID.error);
      return;
    }

    this.moveCursorTo([firstID, firstID]);
    this.insertText("Untitled");
    this.format("h", "1");
  }

  insertText(text: string) {
    if (this.debug) {
      console.log("insertText", text);
    }

    if (!this.curRogueRange) {
      console.error("curRogueRange not found");
      return;
    }

    let [startIx, endIx] = this.getCurRogueIndexes();
    if (startIx < 0) {
      if (this.rogue.Size() > 1) {
        console.error("startIx not found for newline");
        this.renderRogue();
        return;
      }
      startIx = 0;
      endIx = 0;
    }

    const result = this.rogue.RichInsert(
      startIx,
      endIx - startIx,
      this.curSpanFormat,
      text,
    );
    if (result.error) {
      console.error(result.error);
      this.renderRogue();
      return;
    }

    for (const op of result.ops) {
      this.sendOp(op);
      this.renderRogue(op);
    }

    this.curRogueRange = [result.cursorID, result.cursorID];
    this.moveCursorTo(this.curRogueRange);
  }

  insertTextAtEnd(text: string) {
    if (this.debug) {
      console.log("insertTextAtEnd", text);
    }

    const idx = Math.max(this.rogue.Size() - 1, 0);

    const id = this.rogue.GetID(idx);
    if (id.error) {
      console.error(id.error);
      return;
    }

    this.curRogueRange = [id, id];

    this.insertText(text);
  }

  deleteLine() {
    if (this.debug) {
      console.log("deleteLine");
    }
    //
    let [startIx] = this.getCurRogueIndexes();
    if (startIx < 0) {
      console.error("startIx not found");
      return;
    }

    const delop = this.rogue.RichDeleteLine(startIx);
    if (delop.error) {
      console.error(delop.error);
      return;
    }

    if (!delop.op) {
      return;
    }

    this.sendOp(delop.op);

    const id = this.rogue.GetID(delop.startIx);
    if (id.error) {
      console.error(id.error);
      return;
    }

    this.curRogueRange = [id, id];
    this.renderRogue(delop.op);
  }

  deleteWord(forward?: boolean) {
    if (this.debug) {
      console.log("deleteWord", forward);
    }

    let [startIx] = this.getCurRogueIndexes();
    if (startIx < 0) {
      console.error("startIx not found");
      return;
    }

    const delop = this.rogue.RichDeleteWord(startIx, forward);
    if (delop.error) {
      console.error(delop.error);
      return;
    }

    if (!delop.op) {
      return;
    }

    this.sendOp(delop.op);

    const id = this.rogue.GetID(delop.startIx);
    if (id.error) {
      console.error(id.error);
      return;
    }

    this.curRogueRange = [id, id];
    this.renderRogue(delop.op);
  }

  deleteText(forward?: boolean) {
    if (this.debug) {
      console.log("deleteText");
    }

    let [startIx, endIx] = this.getCurRogueIndexes();
    if (startIx < 0) {
      console.error("startIx not found");
      return;
    }

    if (startIx === endIx) {
      if (forward) {
        endIx += 1;
      } else {
        startIx -= 1;
      }
    }

    this.richDelete(startIx, endIx);
  }

  richDelete(startIx: number, endIx: number) {
    const delop = this.rogue.RichDelete(startIx, endIx - startIx);
    if (delop.error) {
      console.error(delop.error);
      return;
    }

    if (!delop.op) {
      return;
    }

    this.sendOp(delop.op);

    const id = this.rogue.GetID(delop.startIx);
    if (id.error) {
      console.error(id.error);
      return;
    }

    this.curRogueRange = [id, id];
    this.renderRogue(delop.op);
  }

  sendCursorUpdate(range: RogueRange) {
    if (this.debug) {
      console.log("sendCursorUpdate", range);
    }

    const msg = JSON.stringify({
      type: "cursor",
      range: range,
      authorID: this.authorId,
      editing: this.editing,
    });

    this.send(msg);
  }

  moveCursorTo(ids: [Id, Id]) {
    if (this.debug) {
      console.log("moveCursorTo", ids);
    }

    const sel = window.getSelection();
    if (!sel) {
      console.error("selection not found, can not restore cursor position!!!");
      return null;
    }

    const range = this.rangeFor(ids);
    if (!range) {
      console.error("range not found, can not restore cursor position!!!");
      return;
    }

    sel.removeAllRanges();
    sel.addRange(range);

    this.curRogueRange = ids;
  }

  rangeFor(ids: [Id, Id]): Range | null {
    if (this.debug) {
      console.log("rangeFor", ids);
    }

    const isScrub = this.editorMode === "scrub";
    const selection = this.rogue.GetSelection(ids[0], ids[1], isScrub);
    if (selection.error) {
      console.error("rangeFor", selection.error);
      return null;
    }

    if (this.debug) {
      console.log("rangeFor selection", selection);
    }

    const startNode = this.getSpanTextNode(selection.startSpanID);
    if (!startNode) {
      console.error("startNode not found");
      return null;
    }
    const endNode = this.getSpanTextNode(selection.endSpanID);
    if (!endNode) {
      console.error("endNode not found");
      return null;
    }

    const range = document.createRange();
    const startOffset = Math.min(
      selection.startOffset,
      startNode.textContent?.length || 0,
    );
    const endOffset = Math.min(
      selection.endOffset,
      endNode.textContent?.length || 0,
    );
    range.setStart(startNode, startOffset);
    range.setEnd(endNode, endOffset);

    return range;
  }

  rangeForAddress(ids: [Id, Id], address: string): Range | null {
    if (this.debug) {
      console.log("rangeForAddress", ids, address);
    }

    const selection = this.rogue.GetSelectionAt(ids[0], ids[1], address);
    if (selection.error) {
      console.error("rangeForAddress", selection.error);
      return null;
    }

    if (this.debug) {
      console.log("rangeForAddress selection", ids, address, selection);
    }

    const startNode = this.getSpanTextNode(selection.startSpanID);
    if (!startNode) {
      console.error("startNode not found");
      return null;
    }
    const endNode = this.getSpanTextNode(selection.endSpanID);
    if (!endNode) {
      console.error("endNode not found");
      return null;
    }

    const range = document.createRange();
    const startOffset = Math.min(
      selection.startOffset,
      startNode.textContent?.length || 0,
    );
    const endOffset = Math.min(
      selection.endOffset + 1,
      endNode.textContent?.length || 0,
    );
    range.setStart(startNode, startOffset);
    range.setEnd(endNode, endOffset);

    return range;
  }

  restoreAnchorsAndOffsets(anchors: {
    start: [string, number];
    end: [string, number];
  }) {
    if (this.debug) {
      console.log("restoreAnchorsAndOffsets", anchors);
    }

    if (!anchors) {
      return;
    }
    const sel = window.getSelection();
    if (!sel) {
      console.error("selection not found, can not restore cursor position!!!");
      return;
    }

    const startEl = document.querySelector(`[data-rid="${anchors.start[0]}"]`);
    const endEl = document.querySelector(`[data-rid="${anchors.end[0]}"]`);
    if (!startEl || !endEl) {
      console.error(
        "can not restore cursor position!!! start or end not found",
        startEl,
        endEl,
        anchors,
      );
      return;
    }

    if (!startEl.childNodes[0] || !endEl.childNodes[0]) {
      console.error(
        "start or end not found, can not restore cursor position!!!",
        startEl.childNodes,
        endEl.childNodes,
      );
      return;
    }

    if (startEl.childNodes.length !== 1 || endEl.childNodes.length !== 1) {
      console.error(
        "too many children, can not restore cursor position!!!",
        startEl.childNodes,
        endEl.childNodes,
      );
      return;
    }

    const startText = startEl.childNodes[0] as Text;
    const endText = endEl.childNodes[0] as Text;

    const range = document.createRange();
    range.setStart(startText, anchors.start[1]);
    range.setEnd(endText, anchors.end[1]);

    sel.removeAllRanges();
    sel.addRange(range);
  }

  getCurRogueIndexes(): [number, number] {
    if (this.debug) {
      console.log("getCurRogueIndexes");
    }

    const ids = this.curRogueRange;
    if (!ids) {
      return [-1, -1];
    }

    const [startId, endId] = ids;
    const startIxs = this.rogue.GetIndex(startId);
    if (startIxs.error) {
      console.error(startIxs.error);
      return [-1, -1];
    }
    const startIx = startIxs.visible;
    if (arraysEqual(startId, endId)) {
      return [startIx, startIx];
    }

    const endIxs = this.rogue.GetIndex(endId);
    if (endIxs.error) {
      console.error(endIxs.error);
      return [-1, -1];
    }

    const endIx = endIxs.visible;

    return [startIx, endIx];
  }

  _getCurRogueIds(): RogueRange | null {
    if (this.debug) {
      console.log("getCurRogueIds");
    }

    const range = this.getCurRange();
    if (!range) {
      return null;
    }

    if (!range.startContainer.parentNode || !range.endContainer.parentNode) {
      console.error(
        "start or end not found",
        range.startContainer,
        range.startOffset,
        range.endContainer,
        range.endOffset,
      );
      return null;
    }

    const startId = this.containerAndOffsetToId(
      range.startContainer as HTMLElement,
      range.startOffset,
    );
    if (!startId) {
      console.error(
        "start not found",
        range,
        range.startContainer.parentNode,
        range.startOffset,
      );
      return null;
    }

    if (
      range.startContainer == range.endContainer &&
      range.startOffset == range.endOffset
    ) {
      return [startId, startId];
    }

    const afterId = this.containerAndOffsetToId(
      range.endContainer as HTMLElement,
      range.endOffset, // The browser use a non inclusive end offset
    );
    if (!afterId) {
      console.error("end not found", range.endContainer);
      return null;
    }

    return [startId, afterId];
  }

  containerAndOffsetToId(container: HTMLElement, offset: number): Id | null {
    if (this.debug) {
      console.log("containerAndOffsetToId", container, offset);
    }

    if (container == this.contentDiv) {
      // if the container is the content div, then we're at the end of the document

      const lastID = this.rogue.GetLastID();
      if (lastID.error) {
        console.error(lastID.error);
        return null;
      }

      return lastID;
    }

    if (container instanceof Text) {
      return this.containerAndOffsetToId(
        container.parentNode as HTMLElement,
        offset,
      );
    }

    const mostNested = findMostNestedElementWithRid(container);
    if (mostNested) {
      container = mostNested;
    }

    const rid = container.getAttribute(RogueIdAttiribute);
    if (!rid) {
      return this.containerAndOffsetToId(
        container.parentNode as HTMLElement,
        offset,
      );
    }

    const containerId = RidToRogueID(rid);
    let returnId: any = null;
    if (this.address) {
      returnId = this.rogue.IDFromIDAndOffset(
        containerId,
        offset,
        this.address,
      );
    } else {
      returnId = this.rogue.IDFromIDAndOffset(containerId, offset);
    }

    if (returnId.error) {
      console.error(returnId.error);
      return null;
    }

    return returnId;
  }

  connect() {
    // If there's an existing connection, close it before creating a new one
    if (this.network && this.network.readyState === WebSocket.OPEN) {
      this.network.close();
    }

    console.log(`Connecting: ${this.docID}`);
    this.resetEditorState();

    this.network = new WebSocket(
      `${this.wsHost}/api/v1/documents/${this.docID}/rogue/ws`,
    );

    this.network.addEventListener("open", this.onOpen.bind(this));
    this.network.addEventListener("message", this.onMessage.bind(this));

    // Clean up: Remove event listeners when the connection is closed
    this.network.addEventListener("close", this.onClose);
  }

  resetEditorState() {
    this.editing = false;
    if (Object.keys(this._cursors).length > 0) {
      this._cursors = {};
      this.notifySubscribers("cursors", this._cursors);
    }
  }

  onClose = (e: any) => {
    console.log("Connection closed", e);
    this.connected = false;
    this.loaded = false;
    this.network?.removeEventListener("open", this.onOpen);
    this.network?.removeEventListener("message", this.onMessage);
    this.reconnectInterval = setInterval(this.reconnect.bind(this), 2000);
  };

  reconnect = () => {
    if (!this._active) {
      clearInterval(this.reconnectInterval);
      return;
    }
    if (!this.network || this.network.readyState !== WebSocket.OPEN) {
      console.log("Attempting to reconnect...", this.instanceID, this.docID);
      this.connect();
      if (this.reconnectInterval) {
        clearInterval(this.reconnectInterval);
      }
    }
  };

  disconnect = () => {
    console.log("Disconnecting from WS Server");
    if (this.network) {
      this.network.removeEventListener("open", this.onOpen);
      this.network.removeEventListener("message", this.onMessage);
      this.network.removeEventListener("close", this.onClose);
      if (this.reconnectInterval) {
        clearInterval(this.reconnectInterval);
      }
      this.network.close();
    }
    this.connected = false;
    this.loaded = false;
  };

  getCurHtml(firstID: Id, lastID: Id, includeIDs: boolean): any {
    let html: any;

    if (this.address && this.editorMode) {
      switch (this.editorMode) {
        case "diff":
          if (this.showDiffHighlights) {
            html = this.rogue.GetHtmlDiff(
              firstID,
              lastID,
              this.address,
              includeIDs,
            );
            console.log("🎨 getHtmlDiff", this.address);
          } else {
            html = this.rogue.GetHtml(firstID, lastID, includeIDs);
            console.log("🎨 getHtml (no diff highlights)", this.address);
          }
          break;
        case "history":
          if (this.showDiffHighlights) {
            if (this.baseAddress) {
              html = this.rogue.GetHtmlDiffBetween(
                firstID,
                lastID,
                this.baseAddress,
                this.address,
                includeIDs,
              );
              console.log(
                "🎨 getHtmlDiffBetween",
                this.baseAddress,
                this.address,
              );
            } else {
              html = this.rogue.GetHtmlDiff(
                firstID,
                lastID,
                this.address,
                includeIDs,
              );
              console.log("🎨 getHtmlDiff", this.address);
            }
          } else {
            html = this.rogue.GetHtmlAtAddress(
              firstID,
              lastID,
              this.address,
              includeIDs,
            );
            console.log("🎨 getHtmlAtAddress", this.address);
          }
          break;
        default:
          console.error("Unrecognized address mode", this.editorMode);
          return "";
      }
    } else {
      switch (this.editorMode) {
        case "xray":
          html = this.rogue.GetHtmlXRay(firstID, lastID, includeIDs);
          console.log("🎨 getHtmlXRay", this.address);
          break;
        default:
          html = this.rogue.GetHtml(firstID, lastID, includeIDs);
          if (this.debug) {
            console.log("🎨 getHtml");
          }
      }
    }

    return html;
  }

  getCurPlaintext(firstID: Id, lastID: Id): any {
    let plaintext: any;

    if (this.address) {
      plaintext = this.rogue.GetPlaintext(firstID, lastID, this.address);
    } else {
      plaintext = this.rogue.GetPlaintext(firstID, lastID);
    }

    return plaintext;
  }

  renderRogue(op?: any) {
    if (!this.contentDiv) {
      console.error("contentDiv not found");
      return;
    }

    if (!this.rogue) {
      return;
    }

    // parital render
    if (op) {
      const span = this.rogue.RenderOp(JSON.stringify(op));
      if (span.error) {
        console.error(span.error);
        return;
      }

      this.replaceBetweenNodes(span.firstBlockID, span.lastBlockID, span.html);
    } else {
      const firstID = this.rogue.GetFirstTotID();
      if (firstID.error) {
        console.error(firstID.error);
        return;
      }

      const lastID = this.rogue.GetLastTotID();
      if (lastID.error) {
        console.error(lastID.error);
        return;
      }

      const html = this.getCurHtml(firstID, lastID, true);
      if (html.error) {
        console.error(html.error);
        return;
      }

      this.contentDiv.innerHTML = html;
    }

    if (
      this.curRogueRange &&
      this.contentDiv &&
      this.contentDiv.contains(document.activeElement)
    ) {
      this.moveCursorTo(this.curRogueRange);
    }
    this.canUndo = this.rogue.CanUndo();
    this.canRedo = this.rogue.CanRedo();

    switch (this.editorMode) {
      case "diff":
        this.disable();
        break;
      case "history":
        this.disable();
        break;
      case "xray":
        this.disable();
        break;
      case "scrub":
        this.disable();
        break;
      default:
        this.enable();
        break;
    }

    this.highlights.recalculate();
  }

  async onOpen() {
    if (!this.contentDiv) {
      console.error("contentDiv not found");
      return;
    }

    if (this.debug) {
      console.log("Connected to WS Server");
    }

    this.network?.send(
      JSON.stringify({
        type: "subscribe",
        docID: this.docID,
        authorID: this.operationManager?.authorId,
      }),
    );

    this.resetAddress();
  }

  async onLoaded() {
    if (!this.contentDiv) {
      console.error("contentDiv not found");
      return;
    }
    if (!this.operationManager) {
      console.error("operationManager not found");
      return;
    }

    if (this.operationManager.hasOperations()) {
      console.log("onLoaded: Merge and send all operations while offline");
      // Merge and send all operations while offline
      const localOps = this.operationManager.getAllOperationsOrderedByIndex();
      for (const op of localOps) {
        const err = this.mergeAndApply(op);
        if (err) {
          console.error(err);
          continue;
        }
        this.sendOp(op);
      }
    } else {
      this.syncing = false;
    }

    this.loaded = true;
    if (!this.address) {
      this.enable();
    } else {
      this.enabled = false;
    }

    if (this.rogue.Size() === 1) {
      this.insertBlankDocContent();
    }
    this.connected = true;

    this.showCursors();
  }

  onMessage(event: any) {
    if (this.debug) {
      console.log("onMessage", event.data);
    }

    const op = JSON.parse(event.data);
    if (op.type === "auth") {
      if (!this.operationManager) {
        throw new Error("No operation manager");
      }
      this.operationManager.authorId = op.authorID;
      // possibly not the best way to do this, but
      // currently the broken websocket connection sends
      // a new auth event and we need to rebuild the rogue
      // a better fix could be to modify subscribe and say
      // "hey I already have a rogue ready to go"
      if (this.rogue) {
        this.rogue.Exit();
      }

      // Reset the rogue editor
      this.rogue = null;
      this.buildRogue(this.operationManager.authorId);

      return;
    }

    if (op.type) {
      switch (op.type) {
        case "newCursor":
        case "cursor":
          if (op.authorID !== this.authorId) {
            this._cursors[op.authorID] = op;
            this.notifySubscribers("cursors", this._cursors);
            this.upsertOtherAuthorCursor(op);
          }
          break;
        case "deleteCursor":
          delete this._cursors[op.authorID];
          this.notifySubscribers("cursors", this._cursors);
          this.removeHighlight("cursor-" + op.authorID);
          break;
        case "event":
          if (op.event === "loaded") {
            this.recvEventBuffer.push({ event: "loaded" });
            setTimeout(this.drainBuffer.bind(this), 0);
          }

          if (op.event === "ping") {
            this.send(
              JSON.stringify({
                type: "event",
                event: "pong",
              }),
            );
          }
          break;
      }

      return;
    }

    this.syncing = true;

    this.recvEventBuffer.push({ event: "op", op });
    setTimeout(this.drainBuffer.bind(this), 0);
  }

  upsertOtherAuthorCursor(msg: any) {
    if (!this.rogue) {
      return;
    }

    if (this.debug) {
      console.log("upsertOtherAuthorCursor", msg);
    }

    if (!msg.range) {
      return;
    }

    const authorInfo = this._cursors[msg.authorID] || defaultAuthorInfo;

    this.createHighlight("cursor-" + msg.authorID, msg.range, {
      styles: {
        backgroundColor: authorInfo.color,
      },
      caret: true,
      caretFlagValue: authorInfo.name,
    });
  }

  // showCursors is called on load to show the cursors that are already loaded in _cursors
  showCursors() {
    for (const authorID in this._cursors) {
      const cursor = this._cursors[authorID];
      this.upsertOtherAuthorCursor(cursor);
    }
  }

  updateCursor(cursor: AuthorInfo) {
    if (this.debug) {
      console.log("updateCursor", cursor);
    }
    this._cursors[cursor.authorID] = cursor;
    this.notifySubscribers("cursors", this._cursors);
    this.upsertOtherAuthorCursor(cursor);
  }

  updateAuthorId() {
    if (this.debug) {
      console.log("updateAuthorId");
    }

    if (!this.rogue) {
      setTimeout(this.updateAuthorId.bind(this), 100);
      return;
    }

    this.rogue.SetAuthor(this.authorId);
  }

  drainBuffer() {
    if (this.debug) {
      console.log("drainBuffer");
    }

    if (!this.rogue) {
      setTimeout(this.drainBuffer.bind(this), 100);
      return;
    }

    let needsRender = false;

    while (this.recvEventBuffer.length > 0) {
      const next = this.recvEventBuffer.shift();
      if (!next) {
        continue;
      }

      if (next.event === "op") {
        if (!next.op) {
          continue;
        }
        const op = next.op;
        if (this.debug) {
          console.log("[wire] recv", JSON.stringify(op));
        }
        if (next.op[0] !== 3) {
          // ignore snapshot ops from because they
          // originate from the server and are not
          // in local storage
          const exists = this.operationManager?.removeOperation(op[1]);
          if (!exists) {
            // was someone else's operation
            needsRender = true;
          }
        } else {
          needsRender = true;
        }

        const err = this.mergeAndApply(op as Op);
        if (err) {
          console.error(err);
          continue;
        }
      } else if (next.event === "loaded") {
        this.onLoaded();
      } else {
        console.warn("Unknown event", next);
      }
    }

    if (!this.hasPendingOperations()) {
      this.syncing = false;
    }

    if (needsRender) {
      this.renderRogue();
    }
  }

  mergeAndApply(op: Op): null | Error {
    if (this.debug) {
      console.log("mergeAndApply", op);
    }

    const r = this.rogue.MergeOp(JSON.stringify(op));
    if (r && r.error) {
      return new Error(r.error);
    }

    return null;
  }

  getRange(): Range | null {
    if (this.debug) {
      console.log("getRange");
    }

    const sel = window.getSelection();
    if (sel && sel.rangeCount > 0) {
      return sel.getRangeAt(0).cloneRange();
    }
    return null;
  }

  restoreRange(range: Range | null) {
    if (this.debug) {
      console.log("restoreRange", range);
    }

    if (!range) {
      return;
    }
    const sel = window.getSelection();
    if (sel) {
      sel.removeAllRanges();
      sel.addRange(range);
    }
  }

  saveCursorPosition(): { start: number; end: number } {
    if (this.debug) {
      console.log("saveCursorPosition");
    }

    if (!this.contentDiv) {
      console.error("contentDiv not found");
      return { start: 0, end: 0 };
    }

    const sel = window.getSelection();
    if (sel && sel.rangeCount > 0) {
      const range = sel.getRangeAt(0);
      const preSelectionRange = range.cloneRange();
      preSelectionRange.selectNodeContents(this.contentDiv);
      preSelectionRange.setEnd(range.startContainer, range.startOffset);
      const start = preSelectionRange.toString().length;
      const end = start + range.toString().length;
      return { start, end };
    }

    return { start: 0, end: 0 };
  }

  sendOp(op: Op) {
    if (this.debug) {
      console.log("[wire] send", JSON.stringify(op));
    }

    // don't send empty multiops
    if (op[0] === 6 && op[2].length === 0) {
      return;
    }

    if (!this.operationManager || !this.operationManager.authed) {
      console.error("Not connected!");
      return;
    }

    this.lastEdit = new Date();

    const opJson = JSON.stringify(op);
    // console.log("opJson", opJson);
    this.operationManager.storeOperation(op);
    const msg = JSON.stringify({ type: "op", op: opJson });
    this.send(msg);
    this.editing = true;
  }

  sendEvent(event: string, data?: any) {
    if (this.debug) {
      console.log("[wire] send event", JSON.stringify(event));
    }

    const payload = { type: "event", event: event, data: {} };
    if (data) {
      payload.data = data;
    }

    const msg = JSON.stringify(payload);
    this.send(msg);
  }

  send(msg: string) {
    if (this.debug) {
      console.log("send", msg);
    }

    if (this.network && this.network.readyState === WebSocket.OPEN) {
      this.network.send(msg);
      return;
    }
    console.error("Could not send message: " + msg);
  }

  hasPendingOperations(): boolean {
    return !!this.operationManager?.hasOperations();
  }

  buildRogue = async (authorID: string) => {
    const loadWasmExecScript = () => {
      return new Promise<void>((resolve, reject) => {
        const script = document.createElement("script");
        script.src = `${this.apiHost}/static/wasm_exec.js?${process.env.IMAGE_TAG}`;
        script.onload = () => resolve();
        script.onerror = () => reject(new Error("Failed to load wasm_exec.js"));
        document.head.appendChild(script);
        this.isWasmLoaded = true;
      });
    };

    try {
      if (!this.isWasmLoaded) {
        await loadWasmExecScript();
      }

      if (this.debug) {
        console.log("Starting WebAssembly");
      }

      const go = new Go();
      const result = await WebAssembly.instantiateStreaming(
        fetch(`${this.apiHost}/static/rogueV3.wasm?${process.env.IMAGE_TAG}`),
        go.importObject,
      );
      go.run(result.instance);

      RegisterPanicCallback((msg: string) => {
        // TODO add sentry alert here
        console.error("Panic:", msg);
        alert(
          "We've hit a unrecoverable error. The reviso team has been notified. Reloading page...",
        );
        window.location.reload();
      });

      this.rogue = NewRogue(authorID);
      console.log(`Rogue ${this.rogue.Version()}-${this.rogue.ImageTag()}`);

      if (process.env.IMAGE_TAG !== this.rogue.ImageTag()) {
        const errMsg = `Rogue image tag mismatch: ${process.env.IMAGE_TAG} vs ${this.rogue.ImageTag()}. Please clear cache and reload the page.`;
        console.error(errMsg);
        Sentry.captureMessage(errMsg);
      }
    } catch (err) {
      alert("Failed to load the doc. Please reload the page.");
      console.error("WebAssembly failed!", err);
    }
  };

  getCurRange(): Range | null {
    if (this.debug) {
      console.log("getCurRange");
    }

    const selection = document.getSelection();
    if (!selection || selection.rangeCount === 0) {
      return null;
    }

    if (selection.rangeCount > 0) {
      // console.log("getCurRange: selection.rangeCount > 0", selection);
      const range = selection.getRangeAt(0);
      if (!range) {
        return null;
      }
      return range;
    }

    return null;
  }

  getSpanTextNode(id: Id): Node | null {
    if (this.debug) {
      console.log("getNodeWithOffset", id);
    }

    const node = document.querySelector(`[data-rid="${id[0]}_${id[1]}"]`);

    if (!node) {
      console.warn("getNodeWithOffset: node not found for id", id);
      return null;
    }

    const textNode = this.getTextNodeFor(node);
    if (!textNode) {
      return node;
    }

    return textNode;
  }

  getAllSpanIds(): Id[] {
    if (this.debug) {
      console.log("getAllSpanIds");
    }

    const elementsWithDataRid = document.querySelectorAll("[data-rid]");
    const uniqueDataRidValues = new Set<string>();
    elementsWithDataRid.forEach((element) => {
      const rid = element.getAttribute("data-rid");
      if (rid) {
        uniqueDataRidValues.add(rid);
      }
    });

    const uniqueDataRidArr = Array.from(uniqueDataRidValues);

    const ids: Id[] = [];
    for (let i = 0; i < uniqueDataRidArr.length; i++) {
      const rid = uniqueDataRidArr[i];
      ids.push(RidToRogueID(rid));
    }

    return ids;
  }

  getTextNodeFor(element: Element): Text | null {
    if (this.debug) {
      console.log("getTextNodeFor", element);
    }

    const subNode = element.querySelector("[data-rid]");
    if (subNode) {
      return this.getTextNodeFor(subNode);
    }

    // Todo: actually check if it's a text node
    return element.childNodes[0] as Text | null;
  }

  getDistanceBetween(a: Id, b: Id): number {
    if (this.debug) {
      console.log("getDistanceBetween", a, b);
    }

    const aIndexes = this.rogue.GetIndex(a);
    if (aIndexes.error) {
      console.error(aIndexes.error);
      return -1;
    }
    const bIndexes = this.rogue.GetIndex(b);
    if (bIndexes.error) {
      console.error(bIndexes.error);
      return -1;
    }
    return aIndexes.visible - bIndexes.visible;
  }

  insertBlankDocContent() {
    if (this.debug) {
      console.log("insertBlankDocContent");
    }

    if (!this.contentDiv) {
      console.error("contentDiv not found");
      return;
    }

    this.contentDiv.focus();
    this.insertBlankState();
    setTimeout(() => {
      const spanElement = this.container?.querySelector("h1 span");
      if (spanElement && window.getSelection && document.createRange) {
        const selection = window.getSelection();
        if (!selection) {
          return;
        }

        const range = document.createRange();
        let firstNode: Node = spanElement,
          lastNode: Node = spanElement;
        while (firstNode.firstChild) firstNode = firstNode.firstChild;
        while (lastNode.lastChild) lastNode = lastNode.lastChild;

        // Ensure lastNode is a text node
        if (lastNode.nodeType !== Node.TEXT_NODE) {
          console.error("Last node is not a text node");
          return;
        }

        // Set the range to encompass all text
        range.setStart(firstNode, 0);
        range.setEnd(lastNode, (lastNode as Text).length);

        // Apply the selection
        selection.removeAllRanges();
        selection.addRange(range);
      }
    }, 0);
  }

  aiMessageSelection(): any {
    if (!this.curRogueRange) {
      return null;
    }

    const [startID, afterID] = this.curRogueRange;
    if (arraysEqual(startID, afterID)) {
      return null;
    }

    const beforeID = this.rogue.TotLeftOf(startID);
    if (beforeID.error) {
      console.error(beforeID.error);
      return null;
    }

    return {
      start: beforeID[0] + "_" + beforeID[1],
      end: afterID[0] + "_" + afterID[1],
      content: this.getCurRogueRangeContent(),
    };
  }
}

export function RidToRogueID(rid: string): [string, number] {
  const parts = rid.split("_");
  const id = parseInt(parts[1]);
  return [parts[0], id];
}

if (!customElements.get("rogue-editor")) {
  customElements.define("rogue-editor", RogueEditor);
}

function shallowEqual(obj1: any, obj2: any): boolean {
  if (obj1 === obj2) {
    return true;
  }

  if (
    typeof obj1 !== "object" ||
    obj1 === null ||
    typeof obj2 !== "object" ||
    obj2 === null
  ) {
    return false;
  }

  const keys1 = Object.keys(obj1);
  const keys2 = Object.keys(obj2);

  if (keys1.length !== keys2.length) {
    return false;
  }

  for (const key of keys1) {
    if (!obj2.hasOwnProperty(key) || obj1[key] !== obj2[key]) {
      return false;
    }
  }

  return true;
}

function arraysEqual<T>(a: T[], b: T[]): boolean {
  if (a.length !== b.length) {
    return false;
  }

  for (let i = 0; i < a.length; i++) {
    if (a[i] !== b[i]) {
      return false;
    }
  }

  return true;
}

function findMostNestedElementWithRid(
  element: HTMLElement | null,
): HTMLElement | null {
  if (!element) return null;

  let current: HTMLElement | null = element;
  let result = null;

  while (current) {
    if (current.hasAttribute("data-rid")) {
      result = current;
    }

    current = current.firstElementChild as HTMLElement | null;
  }

  return result;
}
