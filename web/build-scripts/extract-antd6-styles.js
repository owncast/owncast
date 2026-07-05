/*
Generates web/styles/antd6.css: the static component CSS for Ant Design v6
(installed under the "antd6" npm alias).

The app runs antd v6 in zeroRuntime mode (see Antd6Provider): components do
not generate their style rules at runtime, they expect this pre-extracted
stylesheet to be present. The rules reference CSS variables (--ant6-*) whose
VALUES are still provided at runtime by ConfigProvider from the Owncast
theme, so dynamic theming keeps working; this file only carries the rules
plus the default token values as a pre-hydration fallback.

Run `npm run antd6:extract` after upgrading the antd6 package or changing
antd6-default-theme.json, and commit the regenerated CSS (same convention as
the committed style-dictionary output in styles/variables.css).

The render-every-component approach is adapted from MIT-licensed
@ant-design/static-style-extract, which cannot be used directly: it resolves
the bare "antd" specifier instead of the "antd6" alias, and its peer
dependency on antd>=6 conflicts with the v4 package that remains installed
during the incremental migration. This script uses the hoisted
@ant-design/cssinjs, which is the same module instance antd6 itself uses, so
the style cache is shared.
*/

/* eslint-disable no-underscore-dangle -- antd's own extraction renders the
	 _InternalPanelDoNotUseOrYouWillBeFired pure panels; there is no public way
	 to render Modal/message/notification chrome statically. */

process.env.NODE_ENV = 'production';

const React = require('react');
const { renderToString } = require('react-dom/server');
// Deliberately the hoisted copy that antd6 itself resolves: the style cache
// is only shared when both use the same module instance, so declaring our
// own dependency (which could drift or nest) would silently break extraction.
// eslint-disable-next-line import/no-extraneous-dependencies
const { createCache, extractStyle, StyleProvider } = require('@ant-design/cssinjs');
const antd = require('antd6');
const fs = require('node:fs');
const path = require('node:path');

const defaultTheme = require('../components/theme/antd6-default-theme.json');

const h = React.createElement;

// Components that cannot render bare and need minimal props or subcomponent
// coverage (adapted from @ant-design/static-style-extract).
const customRender = {
  Affix: Comp => h(Comp, null, h('div')),
  BackTop: () => h(antd.FloatButton.BackTop),
  Cascader: Comp => h(React.Fragment, null, h(Comp), h(Comp.Panel)),
  Dropdown: Comp => h(Comp, { menu: { items: [] } }, h('div')),
  Menu: Comp => h(Comp, { items: [] }),
  QRCode: Comp => h(Comp, { value: 'https://owncast.online' }),
  Tree: Comp => h(Comp, { treeData: [] }),
  Tag: Comp =>
    h(
      React.Fragment,
      null,
      h(Comp, { color: 'blue' }, 'Tag'),
      h(Comp, { color: 'success' }, 'Tag'),
    ),
  Badge: Comp => h(React.Fragment, null, h(Comp), h(Comp.Ribbon)),
  Space: Comp =>
    h(
      React.Fragment,
      null,
      h(Comp),
      h(Comp.Compact, null, h(antd.Button), h(Comp.Addon, null, '1')),
    ),
  Input: Comp =>
    h(
      React.Fragment,
      null,
      h(Comp),
      h(Comp.Group, null, h(Comp), h(Comp)),
      h(Comp.Search),
      h(Comp.TextArea),
      h(Comp.Password),
      h(Comp.OTP),
    ),
  Modal: Comp =>
    h(
      React.Fragment,
      null,
      h(Comp),
      h(Comp._InternalPanelDoNotUseOrYouWillBeFired),
      h(Comp._InternalPanelDoNotUseOrYouWillBeFired, { type: 'confirm' }),
    ),
  message: m => h(m._InternalPanelDoNotUseOrYouWillBeFired),
  notification: n => h(n._InternalPanelDoNotUseOrYouWillBeFired),
  Layout: Comp =>
    h(
      Comp,
      null,
      h(Comp.Header, null, 'Header'),
      h(Comp.Sider, null, 'Sider'),
      h(Comp.Content, null, 'Content'),
      h(Comp.Footer, null, 'Footer'),
    ),
};

// Only the components actually used by migrated surfaces are extracted: the
// full-component extract weighs ~111 kB gzipped, which is not worth shipping
// while most of the app is still on v4. WHEN MIGRATING A SURFACE, add any
// new antd6 components it uses here and re-run `npm run antd6:extract`. A
// missing entry is loud, not subtle: the component renders visibly unstyled.
const INCLUDED_COMPONENTS = [
  // Shared primitives
  'Button',
  'Typography',
  'Space',
  'Row',
  'Col',
  // Follow modal (components/modals/FollowModal)
  'Modal',
  'Input',
  'Alert',
  'Spin',
  // Browser notifications modal (components/modals/BrowserNotifyModal)
  // (uses the same set as above)
  // Auth modal (components/modals/AuthModal + FediAuthModal + IndieAuthModal)
  'Tabs',
  'Collapse',
  // User dropdown (components/common/UserDropdown)
  'Dropdown',
  'Menu',
  // Admin form wrappers (components/admin/TextField, ToggleSwitch,
  // TextFieldWithSubmit)
  'Form',
  'InputNumber',
  'Switch',
  // Admin select/slider/tag wrappers (CodecSelector, AutoplaySelector,
  // VideoLatency, EditValueArray)
  'Select',
  'Slider',
  'Tag',
  'Popconfirm',
];

const unknown = INCLUDED_COMPONENTS.filter(name => !antd[name]);
if (unknown.length) {
  throw new Error(`Unknown antd6 export(s) in INCLUDED_COMPONENTS: ${unknown.join(', ')}`);
}

const allComponents = h(
  React.Fragment,
  null,
  INCLUDED_COMPONENTS.map(name => {
    const Comp = antd[name];
    const render = customRender[name];
    return h(React.Fragment, { key: name }, render ? render(Comp) : h(Comp));
  }),
);

// The theme here must mirror the runtime Antd6Provider exactly (prefix and
// hashed:false) so the extracted selectors match the classes components
// render with. zeroRuntime must NOT be set during extraction: it is the flag
// that suppresses exactly the rule generation being captured here.
const cache = createCache();
renderToString(
  h(
    StyleProvider,
    { cache },
    h(
      antd.ConfigProvider,
      {
        prefixCls: 'ant6',
        iconPrefixCls: 'ant6icon',
        theme: { ...defaultTheme, hashed: false },
      },
      allComponents,
    ),
  ),
);

const css = extractStyle(cache, true);
const outFile = path.join(__dirname, '..', 'styles', 'antd6.css');
const header =
  '/* Generated by build-scripts/extract-antd6-styles.js - do not edit.\n' +
  '   Regenerate with `npm run antd6:extract` after upgrading antd6 or\n' +
  '   changing components/theme/antd6-default-theme.json. */\n';
fs.writeFileSync(outFile, header + css);

// eslint-disable-next-line no-console
console.log(`Wrote ${outFile} (${(css.length / 1024).toFixed(1)} kB raw)`);
