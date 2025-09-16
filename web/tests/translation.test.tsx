import React from 'react';
import { render, screen } from '@testing-library/react';
import { Translation } from '../components/ui/Translation/Translation';
import { Localization } from '../types/localization';

// Mock the next-export-i18n hook
jest.mock('next-export-i18n', () => ({
  useTranslation: () => ({
    t: (key: string, vars?: Record<string, any>) => {
      // Simulate the actual translation structure from the JSON files
      const translations: Record<string, any> = {
        // Frontend translations
        'Frontend.helloWorld': 'Hello <strong>{{name}}</strong>, welcome to the world!',
        'Frontend.componentError': 'Error: {{message}}',
        'Frontend.offlineBasic': 'This stream is offline. Check back soon!',
        'Frontend.offlineNotifyOnly':
          "This stream is offline. <span class='notify-link'>Be notified</span> the next time {{streamer}} goes live.",

        // Testing translations
        'Testing.simpleKey': 'Simple translation text',
        'Testing.itemCount': 'You have {{count}} items',
        'Testing.itemCount_one': 'You have {{count}} item',
        'Testing.messageCount': 'You have {{count}} messages from {{sender}}',
        'Testing.messageCount_one': 'You have {{count}} message from {{sender}}',
        'Testing.noPluralKey': 'This key has no plural variants - {{count}} things',

        // Legacy flat keys for backwards compatibility
        hello_world: 'Hello <strong>{{name}}</strong>, welcome to the world!',
        chat_offline: 'Chat is offline',
        notification_message:
          'You can <a href="#">click here</a> to receive notifications when {{streamer}} goes live.',
        component_error: 'Error: {{message}}',
        offline_basic: 'This stream is offline. Check back soon!',
      };

      let result = translations[key];

      // If not found, return the key itself (as real i18n would do)
      if (!result) {
        result = key;
      }

      // Simple variable replacement for testing
      if (vars && typeof result === 'string') {
        Object.keys(vars).forEach(varKey => {
          result = result.replace(new RegExp(`{{${varKey}}}`, 'g'), vars[varKey]);
        });
      }

      return result;
    },
  }),
}));

describe('Translation Component', () => {
  test('should render simple translation text', () => {
    render(<Translation translationKey={Localization.Testing.simpleKey} />);

    expect(screen.getByText('Simple translation text')).toBeInTheDocument();
  });

  test('should render translation with variable interpolation', () => {
    render(
      <Translation translationKey={Localization.Frontend.helloWorld} vars={{ name: 'TestUser' }} />,
    );

    // Check that the text contains the interpolated variable
    // Use a function matcher to handle text across multiple elements, targeting the span
    const element = screen.getByText((_, e) => {
      const hasText = e?.textContent === 'Hello TestUser, welcome to the world!';
      const isSpan = e?.tagName === 'SPAN';
      return hasText && isSpan;
    });
    expect(element).toBeInTheDocument();
  });

  test('should render HTML content correctly', () => {
    render(
      <Translation translationKey={Localization.Frontend.helloWorld} vars={{ name: 'TestUser' }} />,
    );

    // Check that HTML tags are rendered (strong tag in this case)
    const strongElement = screen.getByText('TestUser');
    expect(strongElement.tagName).toBe('STRONG');
  });

  test('should apply className prop', () => {
    render(
      <Translation translationKey={Localization.Testing.simpleKey} className="custom-class" />,
    );

    const element = screen.getByText('Simple translation text');
    expect(element).toHaveClass('custom-class');
  });

  test('should render notification message with HTML content', () => {
    render(
      <Translation
        translationKey={Localization.Frontend.offlineNotifyOnly}
        vars={{ streamer: 'TestStreamer' }}
      />,
    );

    // Check that the HTML content is rendered
    const linkElement = screen.getByText('Be notified');
    expect(linkElement.tagName).toBe('SPAN');
    expect(linkElement).toHaveClass('notify-link');

    // Check that the variable is interpolated
    expect(screen.getByText(/TestStreamer/)).toBeInTheDocument();
  });

  test('should render with all props combined', () => {
    render(
      <Translation
        translationKey={Localization.Frontend.offlineNotifyOnly}
        vars={{ streamer: 'TestStreamer' }}
        className="notification-style"
      />,
    );

    // Check that the content is rendered correctly
    const element = screen.getByText((_, e) => {
      const hasText =
        e?.textContent ===
        'This stream is offline. Be notified the next time TestStreamer goes live.';
      const isSpan = e?.tagName === 'SPAN';
      return hasText && isSpan;
    });
    expect(element).toBeInTheDocument();
    expect(element).toHaveClass('notification-style');
  });

  test('should handle translation without variables', () => {
    render(<Translation translationKey={Localization.Frontend.chatOffline} />);

    expect(screen.getByText('Chat is offline')).toBeInTheDocument();
  });

  test('should render defaultText when translation key is missing', () => {
    // Use a key that doesn't exist in our mock translations
    render(
      <Translation
        translationKey={'non_existent_key' as any}
        defaultText="This is the default text"
      />,
    );

    expect(screen.getByText('This is the default text')).toBeInTheDocument();
  });

  test('should render defaultText with variable interpolation when translation key is missing', () => {
    // Use a key that doesn't exist in our mock translations
    render(
      <Translation
        translationKey={'non_existent_key' as any}
        defaultText="Hello {{name}}, this is default text with {{count}} items"
        vars={{ name: 'John', count: 5 }}
      />,
    );

    expect(screen.getByText('Hello John, this is default text with 5 items')).toBeInTheDocument();
  });

  test('should render defaultText with HTML content when translation key is missing', () => {
    // Use a key that doesn't exist in our mock translations
    render(
      <Translation
        translationKey={'non_existent_key' as any}
        defaultText="This is <strong>bold</strong> default text with <em>emphasis</em>"
      />,
    );

    // Check that HTML tags are rendered correctly
    const strongElement = screen.getByText('bold');
    expect(strongElement.tagName).toBe('STRONG');

    const emElement = screen.getByText('emphasis');
    expect(emElement.tagName).toBe('EM');
  });

  test('should use actual translation when key exists, ignoring defaultText', () => {
    render(
      <Translation
        translationKey={Localization.Testing.simpleKey}
        defaultText="This default text should be ignored"
      />,
    );

    // Should render the actual translation, not the default text
    expect(screen.getByText('Simple translation text')).toBeInTheDocument();
    expect(screen.queryByText('This default text should be ignored')).not.toBeInTheDocument();
  });

  test('should render translation key as fallback when no defaultText is provided and key is missing', () => {
    // Use a key that doesn't exist in our mock translations
    render(<Translation translationKey={'missing_key' as any} />);

    // Should render the key itself as fallback
    expect(screen.getByText('missing_key')).toBeInTheDocument();
  });

  // Pluralization tests
  describe('Pluralization Support', () => {
    test('should use singular form when count is 1', () => {
      render(<Translation translationKey={Localization.Testing.itemCount} count={1} />);

      expect(screen.getByText('You have 1 item')).toBeInTheDocument();
    });

    test('should use plural form when count is greater than 1', () => {
      render(<Translation translationKey={Localization.Testing.itemCount} count={5} />);

      expect(screen.getByText('You have 5 items')).toBeInTheDocument();
    });

    test('should use plural form when count is negative', () => {
      render(<Translation translationKey={Localization.Testing.itemCount} count={-3} />);

      expect(screen.getByText('You have -3 items')).toBeInTheDocument();
    });

    test('should interpolate count and other variables in pluralized text', () => {
      render(
        <Translation
          translationKey={Localization.Testing.messageCount}
          count={3}
          vars={{ sender: 'Alice' }}
        />,
      );

      expect(screen.getByText('You have 3 messages from Alice')).toBeInTheDocument();
    });

    test('should interpolate count and other variables in singular text', () => {
      render(
        <Translation
          translationKey={Localization.Testing.messageCount}
          count={1}
          vars={{ sender: 'Bob' }}
        />,
      );

      expect(screen.getByText('You have 1 message from Bob')).toBeInTheDocument();
    });

    test('should fallback to original key when pluralization keys do not exist', () => {
      render(<Translation translationKey={Localization.Testing.noPluralKey} count={2} />);

      expect(screen.getByText('This key has no plural variants - 2 things')).toBeInTheDocument();
    });

    test('should work without count prop (backward compatibility)', () => {
      render(<Translation translationKey={Localization.Testing.simpleKey} />);

      expect(screen.getByText('Simple translation text')).toBeInTheDocument();
    });

    test('should use defaultText with count interpolation when translation is missing', () => {
      render(
        <Translation
          translationKey={'missing_plural_key' as any}
          count={3}
          defaultText="Default text with {{count}} items"
        />,
      );

      expect(screen.getByText('Default text with 3 items')).toBeInTheDocument();
    });

    test('should use defaultText with count and other vars when translation is missing', () => {
      render(
        <Translation
          translationKey={'missing_plural_key' as any}
          count={1}
          vars={{ name: 'John' }}
          defaultText="{{name}} has {{count}} item"
        />,
      );

      expect(screen.getByText('John has 1 item')).toBeInTheDocument();
    });

    test('should handle count of 0 as plural form', () => {
      render(<Translation translationKey={Localization.Testing.itemCount} count={0} />);

      expect(screen.getByText('You have 0 items')).toBeInTheDocument();
    });

    test('should handle fractional count as plural form', () => {
      render(<Translation translationKey={Localization.Testing.itemCount} count={1.5} />);

      expect(screen.getByText('You have 1.5 items')).toBeInTheDocument();
    });

    test('should work with className and pluralization', () => {
      render(
        <Translation
          translationKey={Localization.Testing.itemCount}
          count={2}
          className="count-display"
        />,
      );

      const element = screen.getByText('You have 2 items');
      expect(element).toHaveClass('count-display');
    });
  });
});
