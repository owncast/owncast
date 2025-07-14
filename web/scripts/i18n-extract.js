/* eslint-disable no-continue */
/* eslint-disable no-restricted-syntax */
const fs = require('fs');
const path = require('path');
const glob = require('glob');
const parser = require('@babel/parser');
const traverse = require('@babel/traverse').default;

const TRANSLATIONS_PATH = path.join(process.cwd(), 'i18n/en/translation.json');

function getDotPath(node) {
  if (node.type === 'MemberExpression') {
    const objectPath = getDotPath(node.object);
    const prop = node.property.name || node.property.value;

    if (objectPath !== null && prop) {
      return objectPath ? `${objectPath}.${prop}` : prop;
    }
  } else if (node.type === 'Identifier' && node.name === 'Localization') {
    return '';
  }
  return null;
}

function sortObjectKeys(obj) {
  if (Array.isArray(obj)) {
    return obj.map(sortObjectKeys);
  }

  if (obj !== null && typeof obj === 'object') {
    return Object.keys(obj)
      .sort()
      .reduce((acc, key) => ({ ...acc, [key]: sortObjectKeys(obj[key]) }), {});
  }

  return obj;
}

function scanTranslationKeys() {
  const files = glob.sync('**/*.{ts,tsx,js,jsx}', {
    ignore: ['node_modules/**', '.next/**', 'out/**'],
  });

  const results = {};

  for (const file of files) {
    const source = fs.readFileSync(file, 'utf8');

    let ast;
    try {
      ast = parser.parse(source, {
        sourceType: 'module',
        plugins: ['jsx', 'typescript'],
      });
    } catch (e) {
      console.warn(`[parse error] ${file}: ${e.message}`);
      continue;
    }

    traverse(ast, {
      JSXElement(p) {
        const opening = p.node.openingElement;
        const tagName = opening.name;

        if (tagName.type !== 'JSXIdentifier' || tagName.name !== 'Translation') return;

        let key = null;
        let defaultText = null;

        for (const attr of opening.attributes) {
          if (attr.type !== 'JSXAttribute') continue;

          const attrName = attr.name.name;
          const { value } = attr;

          if (!value) continue;

          if (attrName === 'translationKey') {
            if (value.expression) {
              const dotPath = getDotPath(value.expression);
              if (dotPath) {
                key = dotPath;
              }
            } else if (value.type === 'StringLiteral') {
              key = value.value;
            }
          }

          if (attrName === 'defaultText' && value.type === 'StringLiteral') {
            defaultText = value.value;
          }
        }

        if (key) {
          // Eventually enable this to not allow empty strings.
          // Then remove the 'missing translation' fallback below.
          // if (!defaultText || defaultText.trim() === '') {
          //   process.exit(1);
          // }
          results[key] =
            defaultText || `<strong><em>Missing translation ${key}: Please report</em></strong>`;
        }
      },
    });
  }

  return results;
}

// Recursively sets a nested value using a dot-notated key
function setNestedKey(obj, keyPath, value) {
  const keys = keyPath.split('.');
  let current = obj;

  keys.forEach((key, index) => {
    if (index === keys.length - 1) {
      current[key] = value;
    } else {
      if (!current[key] || typeof current[key] !== 'object') {
        current[key] = {};
      }
      current = current[key];
    }
  });
}

// Deep merge of two objects
function mergeDeep(target, source) {
  const output = { ...target };
  for (const key of Object.keys(source)) {
    if (source[key] && typeof source[key] === 'object' && !Array.isArray(source[key])) {
      output[key] = mergeDeep(output[key] || {}, source[key]);
    } else {
      output[key] = source[key];
    }
  }
  return output;
}

function updateTranslationFile(flatTranslations) {
  let existing = {};

  if (fs.existsSync(TRANSLATIONS_PATH)) {
    existing = JSON.parse(fs.readFileSync(TRANSLATIONS_PATH, 'utf8'));
  }

  let changed = false;
  let newNestedTranslations = {};

  for (const [flatKey, value] of Object.entries(flatTranslations)) {
    const tempObj = {};
    setNestedKey(tempObj, flatKey, value);

    // Detect if this key is already present
    const flatKeyParts = flatKey.split('.');
    const alreadyExists = flatKeyParts.reduce((acc, part) => acc && acc[part], existing);

    if (!alreadyExists) {
      newNestedTranslations = mergeDeep(newNestedTranslations, tempObj);
      changed = true;
      console.log(`[i18n] Added: ${flatKey}`);
    }
  }

  if (changed) {
    const merged = sortObjectKeys(mergeDeep(existing, newNestedTranslations));
    fs.writeFileSync(TRANSLATIONS_PATH, JSON.stringify(merged, null, 2));
    console.log(`[i18n] Updated ${TRANSLATIONS_PATH}`);
  } else {
    console.log('[i18n] No new keys to add.');
  }
}

const extracted = scanTranslationKeys();
updateTranslationFile(extracted);
