/**
 * DragContext provides state management for drag and drop operations in the document sidebar.
 * It tracks the currently dragged item and provides this information to all components within
 * its provider, enabling features like:
 * - Showing drop indicators when dragging documents between folders
 * - Maintaining drag state across different components in the tree
 * - Automatically cleaning up drag state when drag operations end
 */

import React, { createContext, useContext, useState, useEffect } from "react";

/**
 * Represents a document being dragged in the sidebar.
 * - id: The unique identifier of the document
 * - folderID: The ID of the folder containing the document, or null if it's a top-level document
 */
type DraggedItem = {
  id: string;
  folderID: string | null;
} | null;

type DragContextType = {
  /** The currently dragged document, or null if no drag operation is in progress */
  draggedItem: DraggedItem;
  /** Updates the currently dragged document state */
  setDraggedItem: (item: DraggedItem) => void;
};

const DragContext = createContext<DragContextType | null>(null);

/**
 * Provider component that makes drag state available to its children.
 * Should be placed at a common ancestor of all components that need drag and drop functionality.
 *
 * Features:
 * - Maintains draggedItem state
 * - Automatically cleans up drag state when drag operations end
 * - Provides drag state and setter through context
 *
 * Example:
 * ```tsx
 * <DragProvider>
 *   <DocumentList />
 * </DragProvider>
 * ```
 */
export function DragProvider({ children }: { children: React.ReactNode }) {
  const [draggedItem, setDraggedItem] = useState<DraggedItem>(null);

  useEffect(() => {
    // Clean up drag state when drag operations end
    // This ensures we don't have stale drag state if the dragend event
    // happens outside our components
    const handleDragEnd = () => {
      setDraggedItem(null);
    };
    document.addEventListener("dragend", handleDragEnd);
    return () => document.removeEventListener("dragend", handleDragEnd);
  }, []);

  return (
    <DragContext.Provider value={{ draggedItem, setDraggedItem }}>
      {children}
    </DragContext.Provider>
  );
}

/**
 * Hook to access drag state within components.
 * Must be used within a DragProvider component.
 *
 * Returns:
 * - draggedItem: The currently dragged document or null
 * - setDraggedItem: Function to update the dragged document state
 *
 * Example:
 * ```tsx
 * function DocumentItem() {
 *   const { draggedItem, setDraggedItem } = useDragContext();
 *
 *   const handleDragStart = (e) => {
 *     setDraggedItem({ id: doc.id, folderID: doc.folderID });
 *   };
 *
 *   // ... rest of component
 * }
 * ```
 */
export function useDragContext() {
  const context = useContext(DragContext);
  if (!context) {
    throw new Error("useDragContext must be used within a DragProvider");
  }
  return context;
}
