import React from "react";
import { StoryObj, Meta } from "@storybook/react";
import Title, { DocumentState, TitleProps } from "../Title";

export default {
  title: "Components/Title",
  component: Title,
  argTypes: {
    deleteDocument: { action: "deleteDocument" },
    handleCopy: { action: "handleCopy" },
    toggleMaximize: { action: "toggleMaximize" },
    updateTitle: { action: "updateTitle" },
  },
} as Meta<typeof Title>;

const Template: StoryObj<typeof Title> = (args: TitleProps) => (
  <Title {...args} />
);

export const Default = Template.bind({});
Default.args = {
  deleteLoading: false,
  document: {
    id: "1",
    title: "Sample Document Title",
    updatedAt: new Date(),
  },
  documentState: DocumentState.Saved,
  lastEdit: new Date(),
  loading: false,
  maximized: false,
  me: {
    isAdmin: true,
  },
};

export const Loading = Template.bind({});
Loading.args = {
  deleteLoading: false,
  document: null,
  documentState: DocumentState.Loading,
  lastEdit: null,
  loading: true,
  maximized: false,
  me: null,
};

export const Maximized = Template.bind({});
Maximized.args = {
  deleteLoading: false,
  document: {
    id: "1",
    title: "Sample Document Title",
    updatedAt: new Date(),
  },
  documentState: DocumentState.Disconnected,
  lastEdit: new Date(),
  loading: false,
  maximized: true,
  me: {
    isAdmin: true,
  },
};
