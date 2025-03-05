import type { Meta, StoryObj } from "@storybook/react";
import { CommentThread } from "./CommentThread";

const meta = {
  title: "Components/Sidebar/Timeline/CommentThread",
  component: CommentThread,
  tags: ["autodocs"],
  parameters: {
    layout: "padded",
  },
} satisfies Meta<typeof CommentThread>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {
  args: {
    comment: {
      id: "1",
      content:
        "This looks really good overall. I think I'd just simplify in a couple spots which I commented above. But I'm really excited about this post.",
      user: { id: "user1", name: "Colleen" },
      createdAt: "2023-04-15T12:00:00Z",
      replies: [
        {
          id: "2",
          content: "Reply 1",
          user: { id: "user2", name: "Jane Smith" },
          createdAt: new Date().toISOString(),
          replies: [],
        },
        {
          id: "3",
          content: "Reply 2",
          user: { id: "user3", name: "Bob Johnson" },
          createdAt: new Date().toISOString(),
          replies: [],
        },
      ],
    },
  },
};

export const NoReplies: Story = {
  args: {
    comment: {
      id: "1",
      content: "This is a comment without any replies.",
      user: { id: "user1", name: "John Doe" },
      createdAt: new Date().toISOString(),
      replies: [],
    },
  },
};

export const LongComment: Story = {
  args: {
    comment: {
      id: "1",
      content:
        "This is a very long comment that spans multiple lines. It should demonstrate how the component handles longer text content. We want to make sure that the layout remains consistent and readable even with extended comments.",
      user: { id: "user1", name: "Alice Johnson" },
      createdAt: "2023-04-14T10:00:00Z",
      replies: [
        {
          id: "2",
          content: "Short reply",
          user: { id: "user2", name: "Bob Smith" },
          createdAt: new Date().toISOString(),
          replies: [],
        },
      ],
    },
  },
};
