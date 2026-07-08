/*
Generates web/styles/antd.css: the static component CSS for Ant Design.

The app runs antd in zeroRuntime mode (see AntdProvider): components do
not generate their style rules at runtime, they expect this pre-extracted
stylesheet to be present. The rules reference CSS variables whose VALUES
are still provided at runtime by ConfigProvider from the Owncast theme, so
dynamic theming keeps working; this file only carries the rules plus the
default token values as a pre-hydration fallback.

Run `npm run antd:extract` after upgrading antd or changing
antd-default-theme.json, and commit the regenerated CSS (same convention as
the committed style-dictionary output in styles/variables.css).

The render-every-component approach is adapted from MIT-licensed
@ant-design/static-style-extract, kept local so extraction renders only the
components in the manifest below. This script uses the hoisted
@ant-design/cssinjs, which is the same module instance antd itself uses, so
the style cache is shared.
*/

/* eslint-disable no-underscore-dangle -- antd's own extraction renders the
	 _InternalPanelDoNotUseOrYouWillBeFired pure panels; there is no public way
	 to render Modal/message/notification chrome statically. */

process.env.NODE_ENV = 'production';

const React = require('react');
const { renderToString } = require('react-dom/server');
// Deliberately the hoisted copy that antd itself resolves: the style cache
// is only shared when both use the same module instance, so declaring our
// own dependency (which could drift or nest) would silently break extraction.
// eslint-disable-next-line import/no-extraneous-dependencies
const { createCache, extractStyle, StyleProvider } = require('@ant-design/cssinjs');
const antd = require('antd');
const fs = require('node:fs');
const path = require('node:path');
// eslint-disable-next-line import/no-extraneous-dependencies -- dev-only build script
const prettier = require('prettier');

const defaultTheme = require('../components/theme/antd-default-theme.json');

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
// new antd components it uses here and re-run `npm run antd:extract`. A
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
  // Admin main layout (components/admin/MainLayout)
  'Layout',
  'Badge',
  'Tooltip',
  // Remaining admin pages and components (bulk batch: cards, stats, lists,
  // uploads-adjacent chrome, static message/notification holders)
  'Card',
  'Statistic',
  'Avatar',
  'Skeleton',
  'Divider',
  'Result',
  'Progress',
  'Empty',
  'Checkbox',
  'message',
  'notification',
  // Admin table/upload batch (log/viewer/client/banned-IP/variant/webhook/
  // access-token/user/follower/emoji/plugin tables, logo/favicon uploads)
  'Table',
  'Pagination',
  'Upload',
  // Viewer batch (chat text field emoji picker, moderation menus, stream
  // cards, follower collection, offline banner, viewer layout)
  'Popover',
];

const unknown = INCLUDED_COMPONENTS.filter(name => !antd[name]);
if (unknown.length) {
  throw new Error(`Unknown antd export(s) in INCLUDED_COMPONENTS: ${unknown.join(', ')}`);
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

// The theme here must mirror the runtime AntdProvider exactly (hashed:false)
// so the extracted selectors match the classes components render with.
// zeroRuntime must NOT be set during extraction: it is the flag that
// suppresses exactly the rule generation being captured here.
const cache = createCache();
renderToString(
  h(
    StyleProvider,
    { cache },
    h(
      antd.ConfigProvider,
      {
        theme: { ...defaultTheme, hashed: false },
      },
      allComponents,
    ),
  ),
);

const css = extractStyle(cache, true);
const outFile = path.join(__dirname, '..', 'styles', 'antd.css');
const header =
  '/* Generated by build-scripts/extract-antd-styles.js - do not edit.\n' +
  '   Regenerate with `npm run antd:extract` after upgrading antd or\n' +
  '   changing components/theme/antd-default-theme.json. */\n';

// extractStyle emits minified CSS; format it so the committed file stays
// reviewable and regeneration diffs are purely additive. antd.css is in
// .prettierignore/.stylelintignore, so this is the only formatting pass.
prettier
  .resolveConfig(outFile, { useCache: false })
  .then(config => prettier.format(header + css, { ...config, parser: 'css' }))
  .then(formatted => {
    fs.writeFileSync(outFile, formatted);
    // eslint-disable-next-line no-console
    console.log(`Wrote ${outFile} (${(formatted.length / 1024).toFixed(1)} kB)`);
  });
