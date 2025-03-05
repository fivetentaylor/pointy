import React, { ReactNode, Suspense, Children } from "react";
import ReactMarkdown from "react-markdown";
import remarkbreaks from "remark-breaks";
import { Skeleton } from "@/components/ui/skeleton";

const replaceNewlines = (content: string) => {
  return content
    .replace(/```[\s\S]*?```/g, (m) => m.replace(/\n/g, "\n "))
    .replace(/(?<=\n)(?![*-])\n/g, "&nbsp;\n ");
};

// @:user:([a-zA-Z0-9-]+):  Match @:user: followed by alphanumeric ID
// ((?:(?!@:user:).)*?)  Capture any characters except the start of another user mention, non-greedy
// (?=@(?::user:|\s|$))  Lookahead for @ followed by either :user:, whitespace, or end of string
const oldMentionRegex =
  /@:user:([a-zA-Z0-9-]+):((?:(?!@:user:).)*?)(?=@(?::user:|\s|$))/g;

const urlRegex = /(https?:\/\/[^\s/$.?#].[^\s]*)/gi;

const UserMention = ({
  userId,
  userName,
}: {
  userId: string;
  userName: string;
}) => {
  return <span className="bg-reviso-highlight text-primary">@{userName}</span>;
};

const convertTextToLinks = (text: string) => {
  const parts = text.split(urlRegex);
  return parts.map((part, index) => {
    if (index % 2 === 1) {
      // This is a URL match
      return (
        <a
          className="break-all text-[#0ea5e9] underline"
          key={index}
          href={part.trim()}
          target="_blank"
          rel="noopener noreferrer"
        >
          {part}
        </a>
      );
    }
    // This is regular text
    return part;
  });
};

const processText = (text: string) => {
  const base64EmbedRegex = /@@([a-zA-Z0-9+/=]+)@@/g;
  const oldMentionRegex =
    /@:user:([a-zA-Z0-9-]+):((?:(?!@:user:).)*?)(?=@(?::user:|\s|$))/g;
  const combinedRegex = new RegExp(
    base64EmbedRegex.source + "|" + oldMentionRegex.source,
    "g",
  );
  const output = [];
  let lastIndex = 0;
  let match;

  while ((match = combinedRegex.exec(text)) !== null) {
    const [fullMatch, base64Match, userId, userName] = match;
    const matchIndex = match.index;

    // Add text before the match
    if (matchIndex > lastIndex) {
      output.push(convertTextToLinks(text.substring(lastIndex, matchIndex)));
    }

    if (base64Match) {
      //parse out the decoded base64 embed. Format is: :user:userId:userName
      const decodedEmbed = atob(match[1]).substring(1, match[1].length - 1);
      const [embedType, embedUserId, embedUserName] = decodedEmbed.split(":");

      if (embedType === "user" && embedUserId && embedUserName) {
        // Add the UserMention component for the match
        output.push(
          <UserMention
            key={embedUserId}
            userId={embedUserId}
            userName={embedUserName}
          />,
        );
      }
      lastIndex = matchIndex + fullMatch.length;
    } else if (userId && userName) {
      output.push(
        <UserMention key={userId} userId={userId} userName={userName} />,
      );
      lastIndex = matchIndex + fullMatch.length + 1;
    }
  }

  // Add any remaining text after the last match
  if (lastIndex < text.length) {
    output.push(convertTextToLinks(text.substring(lastIndex)));
  }

  return output;
};

const renderWithMentions = (children: ReactNode) => {
  return Children.map(children, (child) => {
    if (typeof child === "string") {
      return processText(child);
    } else {
      // Non-string children (such as other React elements) are passed through unchanged
      return child;
    }
  });
};

export const Comment = function ({
  fallback,
  content,
}: {
  fallback?: ReactNode;
  content: string;
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
          p: ({ children }) => (
            <p className="text-sm break-words">
              {renderWithMentions(children)}
            </p>
          ),
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
            <ul className="list-disc ml-4 mb-3 text-[0.875rem] leading-[1.25] font-normal">
              {children}
            </ul>
          ),
          ol: ({ children }) => (
            <ol className="list-decimal ml-4 mb-3 text-[0.875rem] leading-[1.25] font-normal">
              {children}
            </ol>
          ),
          li: ({ children }) => (
            <li className="mt-1">{renderWithMentions(children)}</li>
          ),
          strong: ({ children }) => (
            <strong className="font-semibold">{children}</strong>
          ),
          blockquote: ({ children }) => (
            <blockquote className="border-l-2 border-border pl-3 mb-3">
              {children}
            </blockquote>
          ),
        }}
      >
        {replaceNewlines(content)}
      </ReactMarkdown>
    </Suspense>
  );
};
