/* eslint-disable react/no-danger */
import React, { FC } from 'react';
import { useTranslation } from 'next-export-i18n';
import { LocalizationKey } from '../../../types/localization';

export interface TranslationProps {
  translationKey: LocalizationKey;
  vars?: Record<string, any>;
  className?: string;
  defaultText?: string;
  count?: number;
}

export const Translation: FC<TranslationProps> = ({
  translationKey,
  vars,
  className,
  defaultText,
  count,
}) => {
  const { t } = useTranslation();

  // Include count in vars for interpolation
  const allVars = count !== undefined ? { ...vars, count } : vars;

  let translatedText;

  if (count !== undefined) {
    const pluralKey = count === 1 ? `${translationKey}_one` : `${translationKey}_other`;
    translatedText = t(pluralKey, allVars);

    // Fall back to singular translation if plural translation is missing
    if (translatedText === pluralKey) {
      translatedText = t(translationKey, allVars);
    }
  } else {
    translatedText = t(translationKey, allVars);
  }

  // Use fallback if translation is missing (returns the key itself)
  if (translatedText === translationKey && defaultText) {
    translatedText = defaultText;

    // Interpolate variables manually into defaultText
    // eslint-disable-next-line no-restricted-syntax
    for (const [k, v] of Object.entries(allVars || {})) {
      const regex = new RegExp(`{{\\s*${k}\\s*}}`, 'g');
      translatedText = translatedText.replace(regex, String(v));
    }
  }

  return <span className={className} dangerouslySetInnerHTML={{ __html: translatedText }} />;
};

export default Translation;
