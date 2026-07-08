import { StoryFn, Meta } from '@storybook/nextjs';
import { Provider } from 'jotai';
import { Footer } from './Footer';

const meta = {
  title: 'owncast/Layout/Footer',
  component: Footer,
} satisfies Meta<typeof Footer>;

export default meta;

const Template: StoryFn<typeof Footer> = args => (
  <Provider>
    <Footer {...args} />
  </Provider>
);

export const Example = {
  render: Template,

  args: {
    version: 'v1.2.3',
  },
};
