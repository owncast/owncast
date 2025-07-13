import React from 'react';
import { render, screen } from '@testing-library/react';
import { Translation } from '../components/ui/Translation/Translation';
import { Localization } from '../types/localization';

// Mock the next-export-i18n hook
jest.mock('next-export-i18n', () => ({
  useTranslation: () => ({
    t: (key: string, vars?: Record<string, any>) => {
      // Mock translations for testing
      const translations: Record<string, string> = {
        hello_world: 'Hello <strong>{{name}}</strong>, welcome to the world!',
        'Chat is offline': 'Chat is offline',
        notification_message:
          'You can <a href="#">click here</a> to receive notifications when {{streamer}} goes live.',
        simple_key: 'Simple translation text',
      };

      let result = translations[key] || key;

      // Simple variable replacement for testing
      if (vars) {
        Object.keys(vars).forEach(varKey => {
          result = result.replace(`{{${varKey}}}`, vars[varKey]);
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

  test('should render notification message with HTML link', () => {
    render(
      <Translation
        translationKey={Localization.Frontend.notificationMessage}
        vars={{ streamer: 'TestStreamer' }}
      />,
    );

    // Check that the link is rendered
    const linkElement = screen.getByText('click here');
    expect(linkElement.tagName).toBe('A');
    expect(linkElement).toHaveAttribute('href', '#');

    // Check that the variable is interpolated
    expect(screen.getByText(/TestStreamer/)).toBeInTheDocument();
  });

  test('should render with all props combined', () => {
    render(
      <Translation
        translationKey={Localization.Frontend.notificationMessage}
        vars={{ streamer: 'TestStreamer' }}
        className="notification-style"
      />,
    );

    // Check that the content is rendered correctly
    const element = screen.getByText((_, e) => {
      const hasText =
        e?.textContent ===
        'You can click here to receive notifications when TestStreamer goes live.';
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
});
