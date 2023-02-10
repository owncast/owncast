import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import { RecoilRoot } from 'recoil';
import { ChatUserMessage } from './ChatUserMessage';
import { ChatMessage } from '../../../interfaces/chat-message.model';
import Mock from '../../../stories/assets/mocks/chatmessage-user.png';

export default {
  title: 'owncast/Chat/Messages/Standard user',
  component: ChatUserMessage,
  parameters: {
    design: {
      type: 'image',
      url: Mock,
      scale: 0.5,
    },
    docs: {
      description: {
        component: `This is the standard text message design that is used when a user sends a message in Owncast chat.`,
      },
    },
  },
} as ComponentMeta<typeof ChatUserMessage>;

const Template: ComponentStory<typeof ChatUserMessage> = args => (
  <RecoilRoot>
    <ChatUserMessage {...args} />
  </RecoilRoot>
);

const standardMessage: ChatMessage = JSON.parse(`{
  "type": "CHAT",
  "id": "wY-MEXwnR",
  "timestamp": "2022-04-28T20:30:27.001762726Z",
  "user": {
    "id": "h_5GQ6E7R",
    "displayName": "EliteMooseTaskForce",
    "displayColor": 3,
    "createdAt": "2022-03-24T03:52:37.966584694Z",
    "previousNames": ["gifted-nobel", "EliteMooseTaskForce"],
    "nameChangedAt": "2022-04-26T23:56:05.531287897Z",
    "scopes": []
  },
  "body": "Test message from a regular user."}`);

const messageWithLinkAndCustomEmoji: ChatMessage = JSON.parse(`{
		"type": "CHAT",
		"id": "wY-MEXwnR",
		"timestamp": "2022-04-28T20:30:27.001762726Z",
		"user": {
			"id": "h_5GQ6E7R",
			"displayName": "EliteMooseTaskForce",
			"displayColor": 3,
			"createdAt": "2022-03-24T03:52:37.966584694Z",
			"previousNames": ["gifted-nobel", "EliteMooseTaskForce"],
			"nameChangedAt": "2022-04-26T23:56:05.531287897Z",
			"scopes": []
		},
		"body": "Test message with a link https://owncast.online and a custom emoji <img src='/img/emoji/mutant/skull.svg' width='30px'/> ."}`);

const moderatorMessage: ChatMessage = JSON.parse(`{
    "type": "CHAT",
    "id": "wY-MEXwnR",
    "timestamp": "2022-04-28T20:30:27.001762726Z",
    "user": {
      "id": "h_5GQ6E7R",
      "displayName": "EliteMooseTaskForce",
      "displayColor": 2,
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
        "displayColor": 7,
        "createdAt": "2022-03-24T03:52:37.966584694Z",
        "previousNames": ["gifted-nobel", "EliteMooseTaskForce"],
        "nameChangedAt": "2022-04-26T23:56:05.531287897Z",
        "authenticated": true,
        "scopes": []
      },
      "body": "I am an authenticated user."}`);

const botUserMessage: ChatMessage = JSON.parse(`{
				"type": "CHAT",
				"id": "wY-MEXwnR",
				"timestamp": "2022-04-28T20:30:27.001762726Z",
				"user": {
					"id": "h_5GQ6E7R",
					"displayName": "EliteMooseTaskForce",
					"displayColor": 7,
					"createdAt": "2022-03-24T03:52:37.966584694Z",
					"previousNames": ["gifted-nobel", "EliteMooseTaskForce"],
					"nameChangedAt": "2022-04-26T23:56:05.531287897Z",
					"authenticated": true,
					"scopes": ["bot"]
				},
				"body": "I am a bot."}`);

export const WithoutModeratorMenu = Template.bind({});
WithoutModeratorMenu.args = {
  message: standardMessage,
  showModeratorMenu: false,
};

export const WithLinkAndCustomEmoji = Template.bind({});
WithLinkAndCustomEmoji.args = {
  message: messageWithLinkAndCustomEmoji,
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
  isAuthorModerator: true,
};

export const FromAuthenticatedUser = Template.bind({});
FromAuthenticatedUser.args = {
  message: authenticatedUserMessage,
  showModeratorMenu: false,
  isAuthorAuthenticated: true,
};

export const FromBotUser = Template.bind({});
FromBotUser.args = {
  message: botUserMessage,
  showModeratorMenu: false,
  isAuthorBot: true,
};

export const WithStringHighlighted = Template.bind({});
WithStringHighlighted.args = {
  message: standardMessage,
  showModeratorMenu: false,
  highlightString: 'message',
};
