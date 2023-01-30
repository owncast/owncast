import { ComponentMeta, ComponentStory } from '@storybook/react';

const Template = ({
  origin,
  query,
  title,
  width,
  height,
}: {
  origin: string;
  query: string;
  title: string;
  width: number;
  height: number;
}) => (
  <iframe
    src={`${origin}/embed/video?${query}`}
    title={title}
    height={`${height}px`}
    width={`${width}px`}
    referrerPolicy="origin"
    scrolling="no"
    allowFullScreen
  />
);

const origins = {
  DemoServer: `https://watch.owncast.online`,
  RetroStrangeTV: `https://live.retrostrange.com`,
  localhost: `http://localhost:3000`,
};

export default {
  title: 'owncast/Player/Embeds',
  component: Template,
  argTypes: {
    origin: {
      options: Object.keys(origins),
      mapping: origins,
      control: {
        type: 'select',
      },
      defaultValue: origins.DemoServer,
    },
    query: {
      type: 'string',
    },
    title: {
      defaultValue: 'My Title',
      type: 'string',
    },
    height: {
      defaultValue: 350,
      type: 'number',
    },
    width: {
      defaultValue: 550,
      type: 'number',
    },
  },
} satisfies ComponentMeta<typeof Template>;

export const Default: ComponentStory<typeof Template> = Template.bind({});
Default.args = {};

export const InitiallyMuted: ComponentStory<typeof Template> = Template.bind({});
InitiallyMuted.args = {
  query: 'initiallyMuted=true',
};
