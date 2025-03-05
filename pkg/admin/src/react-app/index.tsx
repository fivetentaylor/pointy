import React from "react";
import { createRoot } from "react-dom/client";
import { EditorInterface } from "./components/EditorInterface";

const container = document.getElementById("react-root");
if (!container) {
  throw new Error("react-root not found");
}

const root = createRoot(container); // createRoot(container!) if you use TypeScript
const currentPath = window.location.pathname;

function extractIdFromPath(path: string) {
  const parts = path.split("/");
  return parts[3]; // Index 3 corresponds to the ID segment in the given path format
}

const documentId = extractIdFromPath(currentPath);

// SIGNAL STUFF
import { signal } from "@preact/signals-react";

const hideBar = signal(false);
const spanFormat = signal({});
const lineFormat = signal({});

document.addEventListener("DOMContentLoaded", () => {
  const editorRef = document.getElementById(documentId) as any;

  const toggleSpanFormat = (format: string) => {
    editorRef.toggleSpanFormat(format);
    spanFormat.value = editorRef.curSpanFormat;
  };

  const setLineFormat = (format: string, value: string) => {
    editorRef.format(format, value);
    lineFormat.value = editorRef.curLineFormat;
  };

  const undo = () => {
    editorRef.undo();
  };

  const redo = () => {
    editorRef.redo();
  };

  editorRef.subscribe("curSpanFormat", (value: any) => {
    spanFormat.value = value;
  });

  editorRef.subscribe("curLineFormat", (value: any) => {
    lineFormat.value = value;
  });

  root.render(
    <EditorInterface
      hideBar={hideBar}
      spanFormat={spanFormat}
      lineFormat={lineFormat}
      toggleSpanFormat={toggleSpanFormat}
      setLineFormat={setLineFormat}
      undo={undo}
      redo={redo}
    />,
  );
});
// END SIGNAL STUFF
