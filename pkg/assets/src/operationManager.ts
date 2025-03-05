type Id = [string, number];

export type StorableOperation = [number, Id, ...any[]];

export class OperationManager {
  private documentId: string;
  private authorIdChangeChannel: BroadcastChannel;
  private onAuthorChange: () => void;

  /**
   * Constructor for OperationManager.
   * @param documentId A unique identifier for the document.
   */
  constructor(documentId: string, onAuthorChange: () => void) {
    this.documentId = documentId;
    this.onAuthorChange = onAuthorChange;

    this.authorIdChangeChannel = new BroadcastChannel(
      this.authorIdChannelKey(),
    );
    this.authorIdChangeChannel.onmessage = () => {
      this.onAuthorChange();
    };
  }

  get authed(): boolean {
    return this.authorId !== "";
  }

  get authorId(): string {
    const storedID = localStorage.getItem(this.authorIdKey());
    if (storedID) {
      return storedID;
    }
    return "";
  }

  set authorId(value: string) {
    localStorage.setItem(this.authorIdKey(), value);
    this.onAuthorChange();
    this.authorIdChangeChannel.postMessage({});
  }

  /**
   * Returns if their are any operations stored in localStorage.
   * @returns True if there are operations stored, false otherwise.
   */
  hasOperations(): boolean {
    for (let i = 0; i < localStorage.length; i++) {
      const key = localStorage.key(i);
      if (key && key.startsWith(`op-${this.documentId}-${this.authorId}-`)) {
        return true;
      }
    }

    return false;
  }

  remainingOperations(): number {
    let count = 0;
    for (let i = 0; i < localStorage.length; i++) {
      const key = localStorage.key(i);
      if (key && key.startsWith(`op-${this.documentId}-${this.authorId}-`)) {
        count++;
      }
    }

    return count;
  }

  /**
   * Stores an operation in localStorage.
   * @param operation The operation to be stored.
   */
  storeOperation(operation: StorableOperation): void {
    const key = this.createStorageKey(operation[1]);
    localStorage.setItem(key, JSON.stringify(operation));
  }

  /**
   * Removes a specific operation from localStorage.
   * @param id The unique identifier for the operation.
   * @returns True if the operation was removed, false otherwise.
   */
  removeOperation(id: Id): boolean {
    const key = this.createStorageKey(id);
    const exists = localStorage.getItem(key) !== null;
    if (exists) {
      localStorage.removeItem(key);
    }
    return exists;
  }

  /**
   * Retrieves all stored operations ordered by their index.
   * @returns An array of operation strings ordered by index.
   */
  getAllOperationsOrderedByIndex(): StorableOperation[] {
    const operations: StorableOperation[] = [];
    for (let i = 0; i < localStorage.length; i++) {
      const key = localStorage.key(i);
      if (key && key.startsWith(`op-${this.documentId}-${this.authorId}-`)) {
        const value = localStorage.getItem(key);
        if (!value) {
          continue;
        }
        operations.push(JSON.parse(value) as StorableOperation);
      }
    }

    // Sorting based on the id
    operations.sort((a, b) => {
      // Compare by the seq first
      if (a[1][1] !== b[1][1]) {
        return a[1][1] - b[1][1];
      }

      // If numbers are equal, compare by the author
      return a[1][0].localeCompare(b[1][0]);
    });

    return operations;
  }

  /**
   * Creates a unique key for localStorage based on id and operation.
   * @param id The unique identifier.
   * @param operation The operation string.
   * @returns A string representing the unique localStorage key.
   */
  private createStorageKey(id: Id): string {
    // remove the * from the authorId which indicates
    // reviso ai generated this operation
    const authorId = id[0].replace(/\*$/, "");
    return `op-${this.documentId}-${authorId}-${id[1]}`;
  }

  private authorIdKey(): string {
    return `doc-${this.documentId}-author`;
  }

  private authorIdChannelKey(): string {
    return `doc-${this.documentId}-author-chan`;
  }
}
