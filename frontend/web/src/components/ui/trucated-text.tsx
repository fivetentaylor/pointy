import React from "react";

interface TruncateTextProps {
  text: string;
  length: number;
}

const TruncateText: React.FC<TruncateTextProps> = ({ text, length }) => {
  const truncatedText =
    text.length > length ? text.substring(0, length) + "..." : text;

  return <span className="truncate">{truncatedText}</span>;
};

export default TruncateText;
