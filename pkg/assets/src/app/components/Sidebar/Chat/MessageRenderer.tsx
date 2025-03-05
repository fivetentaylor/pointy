import React, { ReactNode, Suspense, Children } from "react";
import ReactMarkdown from "react-markdown";
import remarkbreaks from "remark-breaks";
import { Skeleton } from "@/components/ui/skeleton";
import { cn } from "@/lib/utils";

const replaceNewlines = (content: string) => {
  return content
    .replace(/```[\s\S]*?```/g, (m) => m.replace(/\n/g, "\n "))
    .replace(/(?<=\n)(?![*-])\n/g, "&nbsp;\n ");
};

// @:user:([a-zA-Z0-9-]+):  Match @:user: followed by alphanumeric ID
// [^@]+ Match one or more non-@ characters
// (?:@[^@]+)? Optionally match @ followed by non-@ characters
const mentionRegex = /@:user:([a-zA-Z0-9-]+):([^@]+(?:@[^@]+)?)@/g;

const UserMention = ({ userName }: { userName: string }) => {
  return <span className="bg-reviso-highlight text-primary">@{userName}</span>;
};

const processText = (text: string) => {
  const output = [];
  let lastEnd = 0;
  let match;

  while ((match = mentionRegex.exec(text)) !== null) {
    const start = match.index;
    const end = mentionRegex.lastIndex;

    // Add text before the match
    if (start > lastEnd) {
      output.push(text.substring(lastEnd, start));
    }

    // Add the UserMention component for the match
    const userId = match[1];
    const userName = match[2];
    output.push(<UserMention key={userId} userName={userName} />);

    lastEnd = end;
  }

  // Add any remaining text after the last match
  if (lastEnd < text.length) {
    output.push(text.substring(lastEnd));
  }

  return output;
};

const MessageRenderer = function ({
  fallback,
  content,
  variant = "default",
}: {
  fallback?: ReactNode;
  content: string;
  variant?: "default" | "timeline" | "timeline-update";
}) {
  if (!fallback) {
    fallback = (
      <>
        <Skeleton className="bg-reviso w-full h-4 mb-2" />
        <Skeleton className="bg-reviso w-full h-4" />
      </>
    );
  }
  return (
    <Suspense fallback={fallback}>
      <ReactMarkdown
        remarkPlugins={[remarkbreaks]}
        components={{
          p: (
            props: React.DetailedHTMLProps<
              React.HTMLAttributes<HTMLParagraphElement>,
              HTMLParagraphElement
            >,
          ) => {
            return (
              <p
                className={cn(
                  "mb-3 last:mb-0",
                  (variant === "timeline" || variant === "timeline-update") &&
                    "text-sm mb-2",
                  variant === "timeline-update" && "inline",
                )}
              >
                {Children.map(props.children, (child) => {
                  if (typeof child === "string") {
                    return processText(child);
                  } else {
                    // Non-string children (such as other React elements) are passed through unchanged
                    return child;
                  }
                })}
              </p>
            );
          },
          h1: ({ children }) => (
            <h1 className="mb-3 font-semibold text-[1.25rem] leading-[2]">
              {children}
            </h1>
          ),
          h2: ({ children }) => (
            <h2 className="mb-3 font-semibold text-[1.125rem] leading-[1.67]">
              {children}
            </h2>
          ),
          h3: ({ children }) => (
            <h3 className="mb-3 font-semibold text-[1rem] leading-[1.5]">
              {children}
            </h3>
          ),
          ul: ({ children }) => (
            <ul
              className={cn(
                "list-disc ml-4 mb-3 text-[0.875rem] leading-[1.25] font-normal",
                variant === "timeline" && "text-sm",
              )}
            >
              {children}
            </ul>
          ),
          ol: ({ children }) => (
            <ol className="list-decimal ml-4 mb-3 text-[0.875rem] leading-[1.3125] font-normal">
              {children}
            </ol>
          ),
          li: ({ children }) => (
            <li className="mt-1 leading-[1.3125rem]">{children}</li>
          ),
          strong: ({ children }) => (
            <strong className="font-semibold">{children}</strong>
          ),
          blockquote: ({ children }) => (
            <blockquote className="border-l-2 border-border pl-3 mb-3">
              {children}
            </blockquote>
          ),
          code: ({ children }) => (
            <code className="whitespace-pre-wrap">{children}</code>
          ),
        }}
      >
        {content}
      </ReactMarkdown>
    </Suspense>
  );
};

export default MessageRenderer;
