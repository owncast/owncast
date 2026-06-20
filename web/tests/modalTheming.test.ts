import fs from 'fs';
import path from 'path';

// These tests guard against a regression where modal colors were hardcoded to
// fixed palette entries (e.g. --theme-color-palette-3) instead of deriving from
// the admin-configurable theme variables. Because they were not indirected
// through the themeable variables, modals always rendered as a "dark header +
// light body" card and ignored any custom theme.
//
// The generated variables.css is the source of truth the browser actually
// loads, so we parse its :root block and follow the var() chains the same way
// the browser would.

const variablesCssPath = path.join(__dirname, '..', 'styles', 'variables.css');
const css = fs.readFileSync(variablesCssPath, 'utf8');

function parseRootVariables(source: string): Record<string, string> {
  const rootMatch = source.match(/:root\s*\{([\s\S]*?)\}/);
  if (!rootMatch) {
    throw new Error('No :root block found in variables.css');
  }
  const vars: Record<string, string> = {};
  const declRegex = /(--[\w-]+)\s*:\s*([^;]+);/g;
  let match: RegExpExecArray | null;
  // eslint-disable-next-line no-cond-assign
  while ((match = declRegex.exec(rootMatch[1])) !== null) {
    // Values may be wrapped across multiple lines by the generator, e.g.
    // `var(\n  --foo\n)`. Collapse internal whitespace so var() references
    // resolve regardless of formatting.
    vars[match[1].trim()] = match[2].replace(/\s+/g, ' ').trim();
  }
  return vars;
}

const vars = parseRootVariables(css);

// Resolve a variable's final value by following single var(--x) references.
// `overrides` simulates the runtime appearance variables that the admin theme
// editor (and theme plugins) inject into :root, which take precedence.
function resolve(
  name: string,
  overrides: Record<string, string> = {},
  seen: Set<string> = new Set(),
): string {
  if (overrides[name] !== undefined) {
    return overrides[name];
  }
  if (seen.has(name)) {
    throw new Error(`Cyclic variable reference at ${name}`);
  }
  seen.add(name);

  const value = vars[name];
  if (value === undefined) {
    return '';
  }

  const varMatch = value.match(/^var\(\s*(--[\w-]+)\s*\)$/);
  if (varMatch) {
    return resolve(varMatch[1], overrides, seen);
  }
  return value;
}

describe('modal theming', () => {
  // The four modal color slots and the themeable variable each must follow.
  const followsTheme: Array<[string, string]> = [
    [
      '--theme-color-components-modal-content-background',
      '--theme-color-components-content-background',
    ],
    ['--theme-color-components-modal-content-text', '--theme-color-components-text-on-light'],
    ['--theme-color-components-modal-header-background', '--theme-color-background-header'],
    ['--theme-color-components-modal-header-text', '--theme-color-components-text-on-dark'],
  ];

  test.each(followsTheme)('%s follows %s when a custom theme is applied', (modalVar, themeVar) => {
    const customColor = '#abcdef';
    expect(resolve(modalVar, { [themeVar]: customColor })).toBe(customColor);
  });

  test('the Ant modal-content-bg override follows the themed modal content background', () => {
    const customColor = '#123456';
    expect(
      resolve('--modal-content-bg', {
        '--theme-color-components-content-background': customColor,
      }),
    ).toBe(customColor);
  });

  test('the modal close icon color follows the themed modal header text', () => {
    const customColor = '#fedcba';
    expect(
      resolve('--modal-close-color', {
        '--theme-color-components-text-on-dark': customColor,
      }),
    ).toBe(customColor);
  });

  test('modal colors are not hardcoded directly to a fixed palette entry', () => {
    // Each modal variable must point at a themeable component variable, never
    // straight at a --theme-color-palette-* entry (the original bug).
    followsTheme.forEach(([modalVar]) => {
      expect(vars[modalVar]).toBeDefined();
      expect(vars[modalVar]).not.toMatch(/var\(\s*--theme-color-palette-\d+\s*\)/);
    });
  });
});
