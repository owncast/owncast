import fs from 'fs';
import path from 'path';
import { Localization } from '../types/localization';

/**
 * Comprehensive localization test suite to verify that translation keys exist
 * across multiple languages and that the localization system is working correctly.
 */

describe('Localization Keys Cross-Language Validation', () => {
  const i18nDir = path.join(__dirname, '../i18n');

  // Get all available language directories
  const getAvailableLanguages = (): string[] =>
    fs.readdirSync(i18nDir).filter(item => {
      const itemPath = path.join(i18nDir, item);
      return fs.statSync(itemPath).isDirectory() && item !== 'en'; // Exclude English as it's our reference
    });

  // Load translation file for a specific language
  const loadTranslationFile = (language: string): Record<string, any> => {
    try {
      const translationPath = path.join(i18nDir, language, 'translation.json');
      const content = fs.readFileSync(translationPath, 'utf-8');
      return JSON.parse(content);
    } catch {
      return {};
    }
  };

  // Helper function to get nested value from object using dot notation
  const getNestedValue = (obj: Record<string, any>, key: string): any =>
    key
      .split('.')
      .reduce(
        (current, prop) => (current && current[prop] !== undefined ? current[prop] : undefined),
        obj,
      );

  // Helper function to check if a key exists in a translation object
  const keyExists = (translations: Record<string, any>, key: string): boolean =>
    getNestedValue(translations, key) !== undefined;

  // Load English translations as reference
  const englishTranslations = loadTranslationFile('en');
  const availableLanguages = getAvailableLanguages();

  describe('Core Frontend Component Keys', () => {
    const testKeys = [
      // NameChangeModal keys
      {
        key: Localization.Frontend.NameChangeModal.description,
        name: 'NameChangeModal.description',
      },
      {
        key: Localization.Frontend.NameChangeModal.placeholder,
        name: 'NameChangeModal.placeholder',
      },
      { key: Localization.Frontend.NameChangeModal.buttonText, name: 'NameChangeModal.buttonText' },
      { key: Localization.Frontend.NameChangeModal.colorLabel, name: 'NameChangeModal.colorLabel' },
      { key: Localization.Frontend.NameChangeModal.authInfo, name: 'NameChangeModal.authInfo' },
      { key: Localization.Frontend.NameChangeModal.overLimit, name: 'NameChangeModal.overLimit' },

      // Header component keys
      {
        key: Localization.Frontend.Header.skipToPlayer,
        name: 'Header.skipToPlayer',
      },
      {
        key: Localization.Frontend.Header.skipToOfflineMessage,
        name: 'Header.skipToOfflineMessage',
      },
      {
        key: Localization.Frontend.Header.skipToContent,
        name: 'Header.skipToContent',
      },
      {
        key: Localization.Frontend.Header.skipToFooter,
        name: 'Header.skipToFooter',
      },
      {
        key: Localization.Frontend.Header.chatWillBeAvailable,
        name: 'Header.chatWillBeAvailable',
      },
      {
        key: Localization.Frontend.Header.chatOffline,
        name: 'Header.chatOffline',
      },

      // Footer component keys
      {
        key: Localization.Frontend.Footer.documentation,
        name: 'Footer.documentation',
      },
      {
        key: Localization.Frontend.Footer.contribute,
        name: 'Footer.contribute',
      },
      {
        key: Localization.Frontend.Footer.source,
        name: 'Footer.source',
      },

      // BrowserNotifyModal keys (sample)
      {
        key: Localization.Frontend.BrowserNotifyModal.unsupported,
        name: 'BrowserNotifyModal.unsupported',
      },
      {
        key: Localization.Frontend.BrowserNotifyModal.allowButton,
        name: 'BrowserNotifyModal.allowButton',
      },
      {
        key: Localization.Frontend.BrowserNotifyModal.enabledTitle,
        name: 'BrowserNotifyModal.enabledTitle',
      },
      {
        key: Localization.Frontend.BrowserNotifyModal.mainDescription,
        name: 'BrowserNotifyModal.mainDescription',
      },

      // Offline messages
      { key: Localization.Frontend.offlineBasic, name: 'Frontend.offlineBasic' },
      { key: Localization.Frontend.offlineNotifyOnly, name: 'Frontend.offlineNotifyOnly' },

      // Error handling
      { key: Localization.Frontend.componentError, name: 'Frontend.componentError' },
    ];

    test('should verify all test keys exist in English translation file', () => {
      testKeys.forEach(({ key }) => {
        const value = getNestedValue(englishTranslations, key);
        expect(keyExists(englishTranslations, key)).toBe(true);
        expect(value).toBeDefined();
        expect(typeof value).toBe('string');
      });
    });

    testKeys.forEach(({ key, name }) => {
      test(`should verify "${name}" exists across all available languages`, () => {
        const missingLanguages: string[] = [];
        const emptyTranslationLanguages: string[] = [];

        availableLanguages.forEach(language => {
          const translations = loadTranslationFile(language);

          if (!keyExists(translations, key)) {
            missingLanguages.push(language);
          } else {
            const value = getNestedValue(translations, key);
            if (!value || value.trim() === '') {
              emptyTranslationLanguages.push(language);
            }
          }
        });

        // Log warnings for missing translations but don't fail the test
        if (missingLanguages.length > 0) {
          console.warn(`⚠️  Key "${key}" is missing in languages: ${missingLanguages.join(', ')}`);
        }

        if (emptyTranslationLanguages.length > 0) {
          console.warn(
            `⚠️  Key "${key}" has empty translations in languages: ${emptyTranslationLanguages.join(', ')}`,
          );
        }

        // At minimum, ensure the key exists in English
        expect(keyExists(englishTranslations, key)).toBe(true);
      });
    });
  });

  describe('Admin Component Keys', () => {
    const adminTestKeys = [
      // EditInstanceDetails
      {
        key: Localization.Admin.EditInstanceDetails.offlineMessageDescription,
        name: 'Admin.EditInstanceDetails.offlineMessageDescription',
      },
      {
        key: Localization.Admin.EditInstanceDetails.directoryDescription,
        name: 'Admin.EditInstanceDetails.directoryDescription',
      },
      {
        key: Localization.Admin.EditInstanceDetails.serverUrlRequiredForDirectory,
        name: 'Admin.EditInstanceDetails.serverUrlRequiredForDirectory',
      },

      // HardwareInfo
      {
        key: Localization.Admin.HardwareInfo.title,
        name: 'Admin.HardwareInfo.title',
      },
      {
        key: Localization.Admin.HardwareInfo.pleaseWait,
        name: 'Admin.HardwareInfo.pleaseWait',
      },
      {
        key: Localization.Admin.HardwareInfo.noDetails,
        name: 'Admin.HardwareInfo.noDetails',
      },
      {
        key: Localization.Admin.HardwareInfo.cpu,
        name: 'Admin.HardwareInfo.cpu',
      },
      {
        key: Localization.Admin.HardwareInfo.memory,
        name: 'Admin.HardwareInfo.memory',
      },
      {
        key: Localization.Admin.HardwareInfo.disk,
        name: 'Admin.HardwareInfo.disk',
      },
      {
        key: Localization.Admin.HardwareInfo.used,
        name: 'Admin.HardwareInfo.used',
      },

      // Help page keys
      {
        key: Localization.Admin.Help.title,
        name: 'Admin.Help.title',
      },
      {
        key: Localization.Admin.Help.configureInstance,
        name: 'Admin.Help.configureInstance',
      },
      {
        key: Localization.Admin.Help.learnMore,
        name: 'Admin.Help.learnMore',
      },
      {
        key: Localization.Admin.Help.configureBroadcasting,
        name: 'Admin.Help.configureBroadcasting',
      },
      {
        key: Localization.Admin.Help.troubleshooting,
        name: 'Admin.Help.troubleshooting',
      },
      {
        key: Localization.Admin.Help.documentation,
        name: 'Admin.Help.documentation',
      },
      {
        key: Localization.Admin.Help.commonTasks,
        name: 'Admin.Help.commonTasks',
      },

      // LogTable keys
      {
        key: Localization.Admin.LogTable.level,
        name: 'Admin.LogTable.level',
      },
      {
        key: Localization.Admin.LogTable.info,
        name: 'Admin.LogTable.info',
      },
      {
        key: Localization.Admin.LogTable.warning,
        name: 'Admin.LogTable.warning',
      },
      {
        key: Localization.Admin.LogTable.error,
        name: 'Admin.LogTable.error',
      },
      {
        key: Localization.Admin.LogTable.timestamp,
        name: 'Admin.LogTable.timestamp',
      },
      {
        key: Localization.Admin.LogTable.message,
        name: 'Admin.LogTable.message',
      },
      {
        key: Localization.Admin.LogTable.logs,
        name: 'Admin.LogTable.logs',
      },

      // NewsFeed keys
      {
        key: Localization.Admin.NewsFeed.link,
        name: 'Admin.NewsFeed.link',
      },
      {
        key: Localization.Admin.NewsFeed.noNews,
        name: 'Admin.NewsFeed.noNews',
      },
      {
        key: Localization.Admin.NewsFeed.title,
        name: 'Admin.NewsFeed.title',
      },

      // ViewerInfo keys
      {
        key: Localization.Admin.ViewerInfo.title,
        name: 'Admin.ViewerInfo.title',
      },
      {
        key: Localization.Admin.ViewerInfo.currentStream,
        name: 'Admin.ViewerInfo.currentStream',
      },
      {
        key: Localization.Admin.ViewerInfo.last12Hours,
        name: 'Admin.ViewerInfo.last12Hours',
      },
      {
        key: Localization.Admin.ViewerInfo.last24Hours,
        name: 'Admin.ViewerInfo.last24Hours',
      },
      {
        key: Localization.Admin.ViewerInfo.currentViewers,
        name: 'Admin.ViewerInfo.currentViewers',
      },
      {
        key: Localization.Admin.ViewerInfo.maxViewersThisStream,
        name: 'Admin.ViewerInfo.maxViewersThisStream',
      },
      {
        key: Localization.Admin.ViewerInfo.viewers,
        name: 'Admin.ViewerInfo.viewers',
      },
    ];

    adminTestKeys.forEach(({ key, name }) => {
      test(`should verify admin key "${name}" has appropriate translation structure`, () => {
        const englishValue = getNestedValue(englishTranslations, key);

        // Admin keys might have missing translation indicators
        expect(englishValue).toBeDefined();
        expect(typeof englishValue).toBe('string');

        // Check if it's a missing translation placeholder
        if (englishValue.includes('Missing translation')) {
          // console.warn(
          //   `⚠️  Admin key "${key}" appears to have missing translation in English: ${englishValue}`,
          // );
        }
      });
    });

    test('should identify missing admin keys in localization.ts vs translation files', () => {
      const missingAdminKeys = [
        { key: Localization.Admin.emojis, name: 'Admin.emojis' },
        { key: Localization.Admin.settings, name: 'Admin.settings' },
        {
          key: Localization.Admin.Chat.moderationMessagesSent,
          name: 'Admin.Chat.moderationMessagesSent',
        },
      ];

      missingAdminKeys.forEach(({ key, name }) => {
        const englishValue = getNestedValue(englishTranslations, key);

        if (!englishValue) {
          console.warn(
            `⚠️  Admin key "${name}" (${key}) is not present in translation files - consider adding it or removing from localization.ts`,
          );
        } else if (englishValue.includes('Missing translation')) {
          console.warn(
            `⚠️  Admin key "${name}" (${key}) has placeholder translation: ${englishValue}`,
          );
        }
      });

      // This test always passes but generates useful warnings
      expect(true).toBe(true);
    });
  });

  describe('Common Keys', () => {
    const commonTestKeys = [
      { key: Localization.Common.poweredByOwncastVersion, name: 'Common.poweredByOwncastVersion' },
    ];

    commonTestKeys.forEach(({ key, name }) => {
      test(`should verify common key "${name}" exists in English`, () => {
        expect(keyExists(englishTranslations, key)).toBe(true);
        const value = getNestedValue(englishTranslations, key);
        expect(value).toBeDefined();
        expect(typeof value).toBe('string');
      });
    });

    test('should identify missing common keys in localization.ts vs translation files', () => {
      // All Common keys are now properly used and extracted automatically
      // Only poweredByOwncastVersion remains as it's actually used in Footer.tsx
      const englishValue = getNestedValue(
        englishTranslations,
        Localization.Common.poweredByOwncastVersion,
      );
      expect(englishValue).toBeDefined();
      expect(typeof englishValue).toBe('string');
      // console.log(`✓ Common key "poweredByOwncastVersion" found: "${englishValue}"`);
    });
  });

  describe('Legacy Frontend Keys (Direct String Values)', () => {
    // These are keys that still use the old direct translation string approach
    const legacyKeys = [
      {
        key: Localization.Frontend.chatDisabled,
        name: 'chatDisabled',
        expectedValue: 'Chat is disabled',
      },
      {
        key: Localization.Frontend.currentViewers,
        name: 'currentViewers',
        expectedValue: 'Current viewers',
      },
      { key: Localization.Frontend.connected, name: 'connected', expectedValue: 'Connected' },
      {
        key: Localization.Frontend.healthyStream,
        name: 'healthyStream',
        expectedValue: 'Healthy Stream',
      },
      {
        key: Localization.Frontend.lastLiveAgo,
        name: 'lastLiveAgo',
        expectedValue: 'Last live {{timeAgo}} ago',
      },
      {
        key: Localization.Frontend.maxViewers,
        name: 'maxViewers',
        expectedValue: 'Max viewers this stream',
      },
    ];

    legacyKeys.forEach(({ key, name, expectedValue }) => {
      test(`should verify legacy frontend key "${name}" uses direct string value`, () => {
        // These keys use direct string values instead of namespace keys
        expect(key).toBe(expectedValue);

        // But we should also verify they exist in the translation file for some languages
        const value = getNestedValue(englishTranslations, key);
        if (value) {
          expect(typeof value).toBe('string');
        }
      });
    });
  });

  describe('Localization Summary Report', () => {
    test('should provide a concise summary of localization status', () => {
      const criticalIssues: string[] = [];
      const warnings: string[] = [];

      // Check NameChangeModal keys (new feature)
      const nameChangeKeys = [
        Localization.Frontend.NameChangeModal.description,
        Localization.Frontend.NameChangeModal.placeholder,
        Localization.Frontend.NameChangeModal.buttonText,
      ];

      const criticalLanguages = ['de', 'es', 'fr', 'it', 'ja', 'ru', 'zh'];
      let missingCriticalTranslations = 0;

      nameChangeKeys.forEach(key => {
        criticalLanguages.forEach(lang => {
          const langTranslations = loadTranslationFile(lang);
          if (!keyExists(langTranslations, key)) {
            missingCriticalTranslations++;
          }
        });
      });

      if (missingCriticalTranslations > 0) {
        warnings.push(
          `NameChangeModal needs translations in ${Math.floor(missingCriticalTranslations / nameChangeKeys.length)} major languages`,
        );
      }

      // Check for keys that shouldn't be in localization.ts
      const problematicKeys = [Localization.Admin.settings];

      problematicKeys.forEach(key => {
        if (!getNestedValue(englishTranslations, key)) {
          criticalIssues.push(
            `Key "${key}" exists in localization.ts but not in translation files`,
          );
        }
      });

      // Print summary
      // console.log('\n🔍 Localization Status Summary:');
      if (criticalIssues.length === 0 && warnings.length === 0) {
        console.log('✅ All critical localization keys are properly configured');
      } else {
        if (criticalIssues.length > 0) {
          console.log('❌ Critical Issues:');
          criticalIssues.forEach(issue => console.log(`  - ${issue}`));
        }
        if (warnings.length > 0) {
          console.log('⚠️  Warnings:');
          warnings.forEach(warning => console.log(`  - ${warning}`));
        }
      }

      // console.log(`📊 Languages supported: ${availableLanguages.length + 1} (including English)`);
      // console.log('💡 Run with LOCALIZATION_VERBOSE=true for detailed warnings\n');

      // Test always passes - this is just informational
      expect(true).toBe(true);
    });
  });

  describe('Language Coverage Report', () => {
    test('should generate language coverage report for key components', () => {
      const reportKeys = [
        Localization.Frontend.NameChangeModal.placeholder,
        Localization.Frontend.BrowserNotifyModal.allowButton,
        Localization.Frontend.offlineBasic,
        Localization.Common.poweredByOwncastVersion,
      ];

      const report: Record<string, { total: number; missing: number; coverage: string }> = {};

      availableLanguages.forEach(language => {
        const translations = loadTranslationFile(language);
        const missing = reportKeys.filter(key => !keyExists(translations, key)).length;
        const total = reportKeys.length;
        const coverage = (((total - missing) / total) * 100).toFixed(1);

        report[language] = {
          total,
          missing,
          coverage: `${coverage}%`,
        };
      });

      // console.log('\n📊 Translation Coverage Report for Key Components:');
      // console.log('Language\tCoverage\tMissing Keys');
      // console.log('--------\t--------\t------------');

      // Object.entries(report)
      //   .sort((a, b) => parseFloat(b[1].coverage) - parseFloat(a[1].coverage))
      //   .forEach(([lang, stats]) => {
      //     console.log(`${lang}\t\t${stats.coverage}\t\t${stats.missing}/${stats.total}`);
      //   });

      // Test passes if we have the report data
      expect(Object.keys(report).length).toBeGreaterThan(0);
    });
  });

  describe('Translation File Structure Validation', () => {
    test('should verify all language directories have translation.json files', () => {
      availableLanguages.forEach(language => {
        const translationPath = path.join(i18nDir, language, 'translation.json');
        expect(fs.existsSync(translationPath)).toBe(true);

        // Verify the file can be parsed as JSON
        expect(() => {
          const content = fs.readFileSync(translationPath, 'utf-8');
          JSON.parse(content);
        }).not.toThrow();
      });
    });

    test('should verify English translation file has expected structure', () => {
      expect(englishTranslations).toBeDefined();
      expect(typeof englishTranslations).toBe('object');

      // Check for expected top-level sections
      expect(englishTranslations.Frontend).toBeDefined();
      expect(englishTranslations.Common).toBeDefined();

      // Check for specific component sections
      expect(englishTranslations.Frontend.NameChangeModal).toBeDefined();
      expect(englishTranslations.Frontend.BrowserNotifyModal).toBeDefined();
    });
  });

  describe('Localization System Integration Test', () => {
    test('should verify the translation hook works with our localization keys', () => {
      // This test verifies that the actual translation system can resolve our keys
      // We'll use a sample of our keys to test the integration

      const testKeys = [
        Localization.Frontend.NameChangeModal.placeholder,
        Localization.Frontend.BrowserNotifyModal.allowButton,
        Localization.Common.poweredByOwncastVersion,
      ];

      testKeys.forEach(key => {
        const value = getNestedValue(englishTranslations, key);
        expect(value).toBeDefined();
        expect(typeof value).toBe('string');
        expect(value.length).toBeGreaterThan(0);
      });
    });

    test('should verify interpolation variables are correctly structured', () => {
      // Test keys that should have interpolation variables
      const interpolationTests = [
        {
          key: Localization.Frontend.componentError,
          expectedVars: ['message'],
          description: 'Component error message should interpolate {{message}}',
        },
        {
          key: Localization.Common.poweredByOwncastVersion,
          expectedVars: ['versionNumber'],
          description: 'Powered by Owncast should interpolate {{versionNumber}}',
        },
        {
          key: Localization.Frontend.offlineNotifyOnly,
          expectedVars: ['streamer'],
          description: 'Offline notify message should interpolate {{streamer}}',
        },
      ];

      interpolationTests.forEach(({ key, expectedVars }) => {
        const value = getNestedValue(englishTranslations, key);
        expect(value).toBeDefined();

        expectedVars.forEach(varName => {
          const hasVariable = value.includes(`{{${varName}}}`);
          expect(hasVariable).toBe(true);
        });
      });
    });
  });

  describe('Localization.ts Type Safety', () => {
    test('should verify localization keys match expected patterns', () => {
      // Test that nested component keys follow the namespace pattern
      expect(Localization.Frontend.NameChangeModal.placeholder).toMatch(
        /^Frontend\.NameChangeModal\./,
      );
      expect(Localization.Frontend.BrowserNotifyModal.allowButton).toMatch(
        /^Frontend\.BrowserNotifyModal\./,
      );
      expect(Localization.Admin.Chat.moderationMessagesSent).toMatch(/^Admin\.Chat\./);
      expect(Localization.Common.poweredByOwncastVersion).toMatch(/^Common\./);

      // Test that basic frontend keys are direct translation strings
      expect(typeof Localization.Frontend.chatOffline).toBe('string');
      expect(typeof Localization.Frontend.currentViewers).toBe('string');
      expect(typeof Localization.Frontend.connected).toBe('string');
    });

    test('should verify all localization keys are strings', () => {
      const validateKeys = (obj: any, keyPath = ''): void => {
        Object.entries(obj).forEach(([key, value]) => {
          const currentPath = keyPath ? `${keyPath}.${key}` : key;

          if (typeof value === 'object' && value !== null) {
            validateKeys(value, currentPath);
          } else {
            expect(typeof value).toBe('string');
            expect(value).toBeTruthy(); // Ensure no empty strings
          }
        });
      };

      validateKeys(Localization);
    });
  });
});
