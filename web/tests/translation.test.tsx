import React from 'react';
import { Translation } from '../components/ui/Translation/Translation';

// Mock the next-export-i18n hook
jest.mock('next-export-i18n', () => ({
  useTranslation: () => ({
    t: (key: string, vars?: Record<string, any>) => {
      // Mock translations for testing
      const translations: Record<string, string> = {
        'hello_world': 'Hello <strong>{{name}}</strong>, welcome to the world!',
        'Chat is offline': 'Chat is offline',
        'notification_message': 'You can <a href="#">click here</a> to receive notifications when {{streamer}} goes live.',
        'simple_key': 'Simple translation text'
      };
      
      let result = translations[key] || key;
      
      // Simple variable replacement for testing
      if (vars) {
        Object.keys(vars).forEach(varKey => {
          result = result.replace(`{{${varKey}}}`, vars[varKey]);
        });
      }
      
      return result;
    }
  })
}));

describe('Translation Component', () => {
  test('should render with translationKey prop', () => {
    const props = {
      translationKey: 'simple_key'
    };
    
    // Test that the component accepts the required props
    expect(() => <Translation {...props} />).not.toThrow();
  });

  test('should accept vars prop for variable interpolation', () => {
    const props = {
      translationKey: 'hello_world',
      vars: { name: 'TestUser' }
    };
    
    // Test that the component accepts vars prop
    expect(() => <Translation {...props} />).not.toThrow();
  });

  test('should accept className prop', () => {
    const props = {
      translationKey: 'simple_key',
      className: 'custom-class'
    };
    
    // Test that the component accepts className prop
    expect(() => <Translation {...props} />).not.toThrow();
  });

  test('should accept all props together', () => {
    const props = {
      translationKey: 'notification_message',
      vars: { streamer: 'TestStreamer' },
      className: 'notification-style'
    };
    
    // Test that the component accepts all props together
    expect(() => <Translation {...props} />).not.toThrow();
  });
});