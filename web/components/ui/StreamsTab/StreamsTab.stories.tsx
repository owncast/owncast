import React from 'react';
import { Meta, StoryFn } from '@storybook/react';
import { StreamsTab, StreamsTabProps, FederatedServer } from './StreamsTab';

export default {
  title: 'owncast/Components/StreamsTab',
  component: StreamsTab,
  parameters: {
    docs: {
      description: {
        component:
          'A tab component that displays a grid of featured Owncast streams with their streaming status.',
      },
    },
  },
} as Meta<typeof StreamsTab>;

const Template: StoryFn<StreamsTabProps> = args => <StreamsTab {...args} />;

const mockServers: FederatedServer[] = [
  {
    id: '1',
    name: 'Tech Talk Central',
    url: 'https://techtalk.example.com',
    logo: 'https://picsum.photos/seed/10/64/64',
    isOnline: true,
    streamTitle: 'Building Microservices with Go',
    streamDescription: 'Learn how to build scalable microservices using Go and Docker.',
    tags: ['technology', 'golang', 'microservices'],
    thumbnail: 'https://picsum.photos/seed/11/640/360',
  },
  {
    id: '2',
    name: 'Gaming Paradise',
    url: 'https://gaming.example.com',
    logo: 'https://picsum.photos/seed/12/64/64',
    isOnline: true,
    streamTitle: 'Speedrun Marathon Day 3',
    streamDescription: 'Day 3 of our annual speedrun marathon featuring classic games.',
    tags: ['gaming', 'speedrun', 'retro'],
    thumbnail: 'https://picsum.photos/seed/13/640/360',
  },
  {
    id: '3',
    name: 'Music Studio Live',
    url: 'https://music.example.com',
    logo: 'https://picsum.photos/seed/14/64/64',
    isOnline: false,
    tags: ['music', 'live', 'jazz'],
  },
  {
    id: '4',
    name: 'Art & Design Hub',
    url: 'https://art.example.com',
    logo: 'https://picsum.photos/seed/15/64/64',
    isOnline: true,
    streamTitle: 'Digital Painting Workshop',
    streamDescription: 'Learn digital painting techniques using Procreate and Photoshop.',
    tags: ['art', 'design', 'tutorial'],
    thumbnail: 'https://picsum.photos/seed/16/640/360',
  },
  {
    id: '5',
    name: 'Cooking Corner',
    url: 'https://cooking.example.com',
    logo: 'https://picsum.photos/seed/17/64/64',
    isOnline: false,
    tags: ['cooking', 'recipes'],
  },
  {
    id: '6',
    name: 'Fitness & Wellness',
    url: 'https://fitness.example.com',
    logo: 'https://picsum.photos/seed/18/64/64',
    isOnline: true,
    streamTitle: 'Morning Yoga Session',
    streamDescription: 'Start your day with a relaxing 30-minute yoga flow.',
    tags: ['fitness', 'yoga', 'wellness'],
    thumbnail: 'https://picsum.photos/seed/19/640/360',
  },
];

export const Default = Template.bind({});
Default.args = {
  servers: mockServers,
};

export const AllOnline = Template.bind({});
AllOnline.args = {
  servers: mockServers.map(server => ({ ...server, isOnline: true })),
};

export const AllOffline = Template.bind({});
AllOffline.args = {
  servers: mockServers.map(server => ({
    ...server,
    isOnline: false,
    streamTitle: undefined,
    streamDescription: undefined,
    thumbnail: undefined,
  })),
};

export const Loading = Template.bind({});
Loading.args = {
  loading: true,
};

export const Empty = Template.bind({});
Empty.args = {
  servers: [],
};

export const Error = Template.bind({});
Error.args = {
  error: 'Failed to load featured streams. Please try again later.',
};

export const SingleServer = Template.bind({});
SingleServer.args = {
  servers: [mockServers[0]],
};

export const ManyServers = Template.bind({});
ManyServers.args = {
  servers: [
    ...mockServers,
    ...mockServers.map((server, index) => ({
      ...server,
      id: `duplicate-${index}`,
      name: `${server.name} (2)`,
    })),
  ],
};
