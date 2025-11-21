import { Meta, StoryObj } from '@storybook/nextjs';
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
    count: {
      control: 'number',
      description:
        'Count for pluralization support (1 = singular with _one key, others = plural with original key)',
    },
    defaultText: {
      control: 'text',
      description: 'Default text to use when translation is missing',
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

export const ComplexHTMLMessage: Story = {
  args: {
    translationKey: Localization.Frontend.offlineNotifyOnly,
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

export const PluralizationSingular: Story = {
  args: {
    translationKey: Localization.Testing.itemCount,
    count: 1,
  },
  parameters: {
    docs: {
      description: {
        story:
          'Pluralization example with singular form (count = 1) - looks for _one key first, falls back to original key',
      },
    },
  },
};

export const PluralizationPlural: Story = {
  args: {
    translationKey: Localization.Testing.itemCount,
    count: 5,
  },
  parameters: {
    docs: {
      description: {
        story: 'Pluralization example with plural form (count > 1) - uses original key directly',
      },
    },
  },
};

export const PluralizationWithVariables: Story = {
  args: {
    translationKey: Localization.Testing.messageCount,
    count: 3,
    vars: {
      sender: 'Alice',
    },
  },
  parameters: {
    docs: {
      description: {
        story: 'Pluralization with additional variables',
      },
    },
  },
};

export const PluralizationFallback: Story = {
  args: {
    translationKey: Localization.Testing.noPluralKey,
    count: 2,
  },
  parameters: {
    docs: {
      description: {
        story:
          'Pluralization fallback when no _one key exists (uses original key for both singular and plural)',
      },
    },
  },
};

export const PluralizationWithDefaultText: Story = {
  args: {
    translationKey: 'non_existent_plural_key' as any,
    count: 4,
    defaultText: 'This is test number {{count}} for {{name}}',
    vars: {
      name: 'Bob',
    },
  },
  parameters: {
    docs: {
      description: {
        story: 'Pluralization with default text fallback',
      },
    },
  },
};
