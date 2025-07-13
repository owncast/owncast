import { Meta, StoryObj } from '@storybook/react';
import { Translation } from '../components/ui/Translation/Translation';
import { Localization } from '../types/localization';

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
    translationKey: Localization.Frontend.chatOffline,
  },
};

export const TranslationWithVariable: Story = {
  args: {
    translationKey: Localization.Frontend.lastLiveAgo,
    vars: {
      timeAgo: '2 hours',
    },
  },
};

export const ComplexHTMLTranslation: Story = {
  args: {
    translationKey: Localization.Frontend.helloWorld,
    vars: {
      name: 'Gabe',
    },
  },
};

export const NotificationMessage: Story = {
  args: {
    translationKey: Localization.Frontend.notificationMessage,
    vars: {
      streamer: 'MyAwesomeStream',
    },
  },
};

export const ComplexMessage: Story = {
  args: {
    translationKey: Localization.Frontend.complexMessage,
    vars: {
      count: 42,
      status: 'live',
    },
  },
};

export const WithCustomClass: Story = {
  args: {
    translationKey: Localization.Frontend.helloWorld,
    vars: {
      name: 'Styled User',
    },
    className: 'custom-translation-style',
  },
};

export const TestDifferentLanguages: Story = {
  args: {
    translationKey: Localization.Frontend.helloWorld,
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
