import React from "react";
import { FileTextIcon, LinkIcon } from "lucide-react";

interface ContentTypeDisplayInfo {
  label: string;
  color: string;
  icon: JSX.Element;
  iconFG: JSX.Element;
  lilIcon: JSX.Element;
}

export const ContentTypeDisplayInfoMap: Record<string, ContentTypeDisplayInfo> =
  {
    "text/plain": {
      label: "TXT",
      color: "bg-zinc-900",
      icon: <FileTextIcon className="w-4 h-4 text-background" />,
      iconFG: <FileTextIcon className="w-4 h-4 text-foreground" />,
      lilIcon: <FileTextIcon className="w-2 h-2 text-background" />,
    },
    ".md": {
      label: "MD",
      color: "bg-zinc-900",
      icon: <FileTextIcon className="w-4 h-4 text-background" />,
      iconFG: <FileTextIcon className="w-4 h-4 text-foreground" />,
      lilIcon: <FileTextIcon className="w-2 h-2 text-background" />,
    },
    "text/markdown": {
      label: "MD",
      color: "bg-zinc-900",
      icon: <FileTextIcon className="w-4 h-4 text-background" />,
      iconFG: <FileTextIcon className="w-4 h-4 text-foreground" />,
      lilIcon: <FileTextIcon className="w-2 h-2 text-background" />,
    },
    ".csv": {
      label: "CSV",
      color: "bg-zinc-900",
      icon: <FileTextIcon className="w-4 h-4 text-background" />,
      iconFG: <FileTextIcon className="w-4 h-4 text-foreground" />,
      lilIcon: <FileTextIcon className="w-2 h-2 text-background" />,
    },
    "text/csv": {
      label: "CSV",
      color: "bg-zinc-900",
      icon: <FileTextIcon className="w-4 h-4 text-background" />,
      iconFG: <FileTextIcon className="w-4 h-4 text-foreground" />,
      lilIcon: <FileTextIcon className="w-2 h-2 text-background" />,
    },
    "application/pdf": {
      label: "PDF",
      color: "bg-rose-500",
      icon: <FileTextIcon className="w-4 h-4 text-background" />,
      iconFG: <FileTextIcon className="w-4 h-4 text-foreground" />,
      lilIcon: <FileTextIcon className="w-2 h-2 text-background" />,
    },
    "image/png": {
      label: "PNG",
      color: "text-emerald-500",
      icon: <FileTextIcon className="w-4 h-4 text-background" />,
      iconFG: <FileTextIcon className="w-4 h-4 text-foreground" />,
      lilIcon: <FileTextIcon className="w-2 h-2 text-background" />,
    },
    "image/jpeg": {
      label: "JPEG",
      color: "text-emerald-500",
      icon: <FileTextIcon className="w-4 h-4 text-background" />,
      iconFG: <FileTextIcon className="w-4 h-4 text-foreground" />,
      lilIcon: <FileTextIcon className="w-2 h-2 text-background" />,
    },
    "image/jpg": {
      label: "JPG",
      color: "text-emerald-500",
      icon: <FileTextIcon className="w-4 h-4 text-background" />,
      iconFG: <FileTextIcon className="w-4 h-4 text-foreground" />,
      lilIcon: <FileTextIcon className="w-2 h-2 text-background" />,
    },
    "image/gif": {
      label: "GIF",
      color: "text-emerald-500",
      icon: <FileTextIcon className="w-4 h-4 text-background" />,
      iconFG: <FileTextIcon className="w-4 h-4 text-foreground" />,
      lilIcon: <FileTextIcon className="w-2 h-2 text-background" />,
    },
    "image/svg+xml": {
      label: "SVG",
      color: "text-emerald-500",
      icon: <FileTextIcon className="w-4 h-4 text-background" />,
      iconFG: <FileTextIcon className="w-4 h-4 text-foreground" />,
      lilIcon: <FileTextIcon className="w-2 h-2 text-background" />,
    },
    "application/vnd.openxmlformats-officedocument.wordprocessingml.document": {
      label: "DOCX",
      color: "bg-sky-500",
      icon: <FileTextIcon className="w-4 h-4 text-background" />,
      iconFG: <FileTextIcon className="w-4 h-4 text-foreground" />,
      lilIcon: <FileTextIcon className="w-2 h-2 text-background" />,
    },
    "application/msword": {
      label: "DOC",
      color: "bg-sky-500",
      icon: <FileTextIcon className="w-4 h-4 text-background" />,
      iconFG: <FileTextIcon className="w-4 h-4 text-foreground" />,
      lilIcon: <FileTextIcon className="w-2 h-2 text-background" />,
    },
    "application/vnd.ms-word": {
      label: "DOC",
      color: "bg-sky-500",
      icon: <FileTextIcon className="w-4 h-4 text-background" />,
      iconFG: <FileTextIcon className="w-4 h-4 text-foreground" />,
      lilIcon: <FileTextIcon className="w-2 h-2 text-background" />,
    },
    "application/vnd.openxmlformats-officedocument.presentationml.presentation":
      {
        label: "PPTX",
        color: "bg-sky-500",
        icon: <FileTextIcon className="w-4 h-4 text-background" />,
        iconFG: <FileTextIcon className="w-4 h-4 text-foreground" />,
        lilIcon: <FileTextIcon className="w-2 h-2 text-background" />,
      },
    "application/vnd.oasis.opendocument.text": {
      label: "ODT",
      color: "bg-sky-500",
      icon: <FileTextIcon className="w-4 h-4 text-background" />,
      iconFG: <FileTextIcon className="w-4 h-4 text-foreground" />,
      lilIcon: <FileTextIcon className="w-2 h-2 text-background" />,
    },
    "text/url": {
      label: "URL",
      color: "bg-indigo-500",
      icon: <LinkIcon className="w-4 h-4 text-background" />,
      iconFG: <LinkIcon className="w-4 h-4 text-foreground" />,
      lilIcon: <LinkIcon className="w-2 h-2 text-background" />,
    },
    "application/rtf": {
      label: "RTF",
      color: "bg-sky-500",
      icon: <FileTextIcon className="w-4 h-4 text-background" />,
      iconFG: <FileTextIcon className="w-4 h-4 text-foreground" />,
      lilIcon: <FileTextIcon className="w-2 h-2 text-background" />,
    },
    "application/x-rtf": {
      label: "RTF",
      color: "bg-sky-500",
      icon: <FileTextIcon className="w-4 h-4 text-background" />,
      iconFG: <FileTextIcon className="w-4 h-4 text-foreground" />,
      lilIcon: <FileTextIcon className="w-2 h-2 text-background" />,
    },
    "text/rtf": {
      label: "RTF",
      color: "bg-sky-500",
      icon: <FileTextIcon className="w-4 h-4 text-background" />,
      iconFG: <FileTextIcon className="w-4 h-4 text-foreground" />,
      lilIcon: <FileTextIcon className="w-2 h-2 text-background" />,
    },
    "text/richtext": {
      label: "RTF",
      color: "bg-sky-500",
      icon: <FileTextIcon className="w-4 h-4 text-background" />,
      iconFG: <FileTextIcon className="w-4 h-4 text-foreground" />,
      lilIcon: <FileTextIcon className="w-2 h-2 text-background" />,
    },
  };
