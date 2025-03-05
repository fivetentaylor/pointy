import React from "react";
import { MsgLlm } from "@/__generated__/graphql";

import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { AI_SWITCH_LLM } from "@/lib/events";
import { analytics } from "@/lib/segment";

export default function LLMSelector({
  onChange,
}: {
  onChange: (llm: MsgLlm) => void;
}) {
  return (
    <Select
      defaultValue="CLAUDE"
      onValueChange={(value: string) => {
        analytics.track(AI_SWITCH_LLM, { llm: value });
        onChange(value as MsgLlm);
      }}
    >
      <SelectTrigger className="w-28 border-none focus:ring-0 focus:ring-offset-0 bg-transparent ml-[-0.25rem]">
        <SelectValue />
      </SelectTrigger>
      <SelectContent>
        <SelectGroup>
          <SelectItem value="CLAUDE">Claude 3.5</SelectItem>
          <SelectItem value="GPT4O">GPT-4o</SelectItem>
        </SelectGroup>
      </SelectContent>
    </Select>
  );
}
