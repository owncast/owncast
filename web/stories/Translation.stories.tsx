import { Meta, StoryObj } from '@storybook/react';
import { Translation } from '../components/ui/Translation/Translation';

const meta: Meta<typeof Translation> = {
  title: 'owncast/Components/Translation',
  component: Translation,
  parameters: {
    chromatic: { diffThreshold: 0.8 },
  },
  argTypes: {
    translationKey: {
      control: 'text',
      description: 'The translation key to use for the text',
    },
    vars: {
      control: 'object',
      description: 'Variables to interpolate into the translation',
    },
    className: {
      control: 'text',
      description: 'CSS class name to apply to the component',
    },
  },
};

export default meta;
type Story = StoryObj<typeof Translation>;

export const SimpleTranslation: Story = {
  args: {
    translationKey: 'Chat is offline',
  },
};

export const TranslationWithVariable: Story = {
  args: {
    translationKey: 'Last live ago',
    vars: {
      timeAgo: '2 hours',
    },
  },
};

export const ComplexHTMLTranslation: Story = {
  args: {
    translationKey: 'hello_world',
    vars: {
      name: 'Gabe',
    },
  },
};

export const NotificationMessage: Story = {
  args: {
    translationKey: 'notification_message',
    vars: {
      streamer: 'MyAwesomeStream',
    },
  },
};

export const ComplexMessage: Story = {
  args: {
    translationKey: 'complex_message',
    vars: {
      count: 42,
      status: 'live',
    },
  },
};

export const WithCustomClass: Story = {
  args: {
    translationKey: 'hello_world',
    vars: {
      name: 'Styled User',
    },
    className: 'custom-translation-style',
  },
};

export const TestDifferentLanguages: Story = {
  args: {
    translationKey: 'hello_world',
    vars: {
      name: 'Test User',
    },
  },
  parameters: {
    docs: {
      description: {
        story: 'Test different languages by adding ?lang=de or ?lang=fr to the URL',
      },
    },
  },
};