import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import UserChatMessage from '../components/chat/ChatUserMessage';
import { ChatMessage } from '../interfaces/chat-message.model';

export default {
  title: 'owncast/Chat/Messages/Standard user',
  component: UserChatMessage,
  parameters: {},
} as ComponentMeta<typeof UserChatMessage>;

const Template: ComponentStory<typeof UserChatMessage> = args => <UserChatMessage {...args} />;

const standardMessage: ChatMessage = JSON.parse(`{
  "type": "CHAT",
  "id": "wY-MEXwnR",
  "timestamp": "2022-04-28T20:30:27.001762726Z",
  "user": {
    "id": "h_5GQ6E7R",
    "displayName": "EliteMooseTaskForce",
    "displayColor": 329,
    "createdAt": "2022-03-24T03:52:37.966584694Z",
    "previousNames": ["gifted-nobel", "EliteMooseTaskForce"],
    "nameChangedAt": "2022-04-26T23:56:05.531287897Z",
    "scopes": []
  },
  "body": "Test message from a regular user."}`);

const moderatorMessage: ChatMessage = JSON.parse(`{
    "type": "CHAT",
    "id": "wY-MEXwnR",
    "timestamp": "2022-04-28T20:30:27.001762726Z",
    "user": {
      "id": "h_5GQ6E7R",
      "displayName": "EliteMooseTaskForce",
      "displayColor": 329,
      "createdAt": "2022-03-24T03:52:37.966584694Z",
      "previousNames": ["gifted-nobel", "EliteMooseTaskForce"],
      "nameChangedAt": "2022-04-26T23:56:05.531287897Z",
      "scopes": ["moderator"]
    },
    "body": "I am a moderator user."}`);

const authenticatedUserMessage: ChatMessage = JSON.parse(`{
      "type": "CHAT",
      "id": "wY-MEXwnR",
      "timestamp": "2022-04-28T20:30:27.001762726Z",
      "user": {
        "id": "h_5GQ6E7R",
        "displayName": "EliteMooseTaskForce",
        "displayColor": 329,
        "createdAt": "2022-03-24T03:52:37.966584694Z",
        "previousNames": ["gifted-nobel", "EliteMooseTaskForce"],
        "nameChangedAt": "2022-04-26T23:56:05.531287897Z",
        "authenticated": true,
        "scopes": []
      },
      "body": "I am an authenticated user."}`);

export const WithoutModeratorMenu = Template.bind({});
WithoutModeratorMenu.args = {
  message: standardMessage,
  showModeratorMenu: false,
};

export const WithModeratorMenu = Template.bind({});
WithModeratorMenu.args = {
  message: standardMessage,
  showModeratorMenu: true,
};

export const FromModeratorUser = Template.bind({});
FromModeratorUser.args = {
  message: moderatorMessage,
  showModeratorMenu: false,
};

export const FromAuthenticatedUser = Template.bind({});
FromAuthenticatedUser.args = {
  message: authenticatedUserMessage,
  showModeratorMenu: false,
};
