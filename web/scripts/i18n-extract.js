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
      return objectPath ? `${objectPath}.${prop}` : prop; // skip base if empty
    }
  } else if (node.type === 'Identifier' && node.name === 'Localization') {
    return ''; // treat as base, skip
  }
  return null;
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

        if (key && defaultText && !results[key]) {
          results[key] = defaultText;
        }
      },
    });
  }

  return results;
}

function updateTranslationFile(newTranslations) {
  let existing = {};

  if (fs.existsSync(TRANSLATIONS_PATH)) {
    existing = JSON.parse(fs.readFileSync(TRANSLATIONS_PATH, 'utf8'));
  }

  let changed = false;

  for (const [key, value] of Object.entries(newTranslations)) {
    if (!(key in existing)) {
      existing[key] = value;
      changed = true;
      console.log(`[i18n] Added: ${key}`);
    }
  }

  if (changed) {
    const sorted = Object.fromEntries(
      Object.entries(existing).sort(([a], [b]) => a.localeCompare(b)),
    );
    fs.writeFileSync(TRANSLATIONS_PATH, JSON.stringify(sorted, null, 2));
    console.log(`[i18n] Updated ${TRANSLATIONS_PATH}`);
  } else {
    console.log('[i18n] No new keys to add.');
  }
}

const extracted = scanTranslationKeys();
updateTranslationFile(extracted);
