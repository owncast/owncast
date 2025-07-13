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

  let translatedText = t(translationKey, allVars);

  // Handle pluralization if count is provided
  if (count !== undefined) {
    const pluralKey = count === 1 ? `${translationKey}_one` : `${translationKey}_other`;
    const pluralTranslation = t(pluralKey, allVars);

    // Use plural translation if it exists (not returning the key itself)
    if (pluralTranslation !== pluralKey) {
      translatedText = pluralTranslation;
    }
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
