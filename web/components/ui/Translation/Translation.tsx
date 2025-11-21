/* eslint-disable react/no-danger */
import React, { FC } from 'react';
import { useTranslation } from 'next-export-i18n';
import { LocalizationKey } from '../../../types/localization';

export interface TranslationProps {
  translationKey: LocalizationKey;
  vars?: Record<string, any>;
  className?: string;
  id?: string;
  defaultText?: string;
  count?: number;
}

export const Translation: FC<TranslationProps> = ({
  translationKey,
  vars,
  className,
  id,
  defaultText,
  count,
}) => {
  const { t } = useTranslation();

  // Include count in vars for interpolation
  const allVars = count !== undefined ? { ...vars, count } : vars;

  let translatedText;

  if (count !== undefined) {
    if (count === 1) {
      // For singular, try _one key first, fall back to original key
      const singularKey = `${translationKey}_one`;
      translatedText = t(singularKey, allVars);

      // Fall back to original key if _one key doesn't exist
      if (translatedText === singularKey) {
        translatedText = t(translationKey, allVars);
      }
    } else {
      // For plural, always use the original key (no _other suffix needed)
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

  return (
    <span className={className} id={id} dangerouslySetInnerHTML={{ __html: translatedText }} />
  );
};

export default Translation;
