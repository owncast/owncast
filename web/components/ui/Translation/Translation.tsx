/* eslint-disable react/no-danger */
import React, { FC } from 'react';
import { useTranslation } from 'next-export-i18n';
import { LocalizationKey } from '../../../types/localization';

export interface TranslationProps {
  translationKey: LocalizationKey;
  vars?: Record<string, any>;
  className?: string;
  defaultText?: string;
}

export const Translation: FC<TranslationProps> = ({
  translationKey,
  vars,
  className,
  defaultText,
}) => {
  const { t } = useTranslation();

  let translatedText = t(translationKey, vars);

  // Use fallback if translation is missing (returns the key itself)
  if (translatedText === translationKey && defaultText) {
    translatedText = defaultText;

    // Interpolate variables manually into defaultText
    // eslint-disable-next-line no-restricted-syntax
    for (const [k, v] of Object.entries(vars || {})) {
      const regex = new RegExp(`{{\\s*${k}\\s*}}`, 'g');
      translatedText = translatedText.replace(regex, String(v));
    }
  }

  return <span className={className} dangerouslySetInnerHTML={{ __html: translatedText }} />;
};

export default Translation;
