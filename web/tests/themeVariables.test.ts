import fs from 'fs';
import path from 'path';

// Guards against regressions in the generated theme variables that ship in
// styles/variables.css. The token source (style-definitions) is the source of
// truth, but variables.css is the artifact the app actually loads, so we assert
// against it directly.

const variablesCssPath = path.join(__dirname, '..', 'styles', 'variables.css');
const css = fs.readFileSync(variablesCssPath, 'utf8');

describe('generated theme variables', () => {
  // A YAML indentation bug in the token source once collapsed the eight
  // per-user chat colors (--theme-color-users-0..7) into a single broken
  // --theme-color-users variable. That breaks chat username colors and the
  // admin "Chat User Colors" pickers (which reference theme-color-users-0..7).
  test('declares all eight distinct chat user colors', () => {
    for (let i = 0; i < 8; i += 1) {
      const declaration = `--theme-color-users-${i}: var(--color-owncast-user-${i});`;
      expect(css).toContain(declaration);
    }
  });

  test('does not declare a collapsed --theme-color-users variable', () => {
    // The trailing colon distinguishes the broken collapsed variable from the
    // valid indexed ones (--theme-color-users-0, etc.).
    expect(css).not.toMatch(/--theme-color-users:/);
  });
});
