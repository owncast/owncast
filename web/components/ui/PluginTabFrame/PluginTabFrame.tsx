import { FC, useEffect, useRef } from 'react';
import { useRecoilValue } from 'recoil';
import { ClientConfig } from '../../../interfaces/client-config.model';
import { clientConfigStateAtom } from '../../stores/ClientConfigStore';

export type PluginTabFrameProps = {
  content: string;
};

// The generated element-baseline + token-defaults stylesheet (built from the
// design tokens by web/style-definitions, served from public/). Root-relative,
// which a srcdoc document resolves against the parent page's base URL.
const PLUGIN_STYLESHEET = '/styles/plugin.css';

// The cascade we want inside the frame, bottom to top:
//   1. plugin.css          — element baseline + token defaults.
//   2. author's own styles — in their HTML, so they override the baseline.
//   3. pluginStyles        — other loaded plugins' page CSS (e.g. a theme),
//                            so the tab inherits the same baseline theming
//                            as the main page.
//   4. appearanceVariables — the admin's :root variable overrides.
//   5. customStyles        — the admin's custom CSS.
//
// All of it is injected by the host (the author writes none of it). We inject
// into the DOM rather than string-templating the srcdoc because the author's
// HTML may be a full document with its own <head>: the parser always
// synthesizes exactly one <head>, so inserting at head.firstChild reliably
// places plugin.css ahead of the author's own head styles (lowest layer)
// regardless of how their markup is shaped. The runtime layers (3, 4, 5) are
// appended last so they win, in priority order, and set via textContent so a
// stray </style> in any of them can't break out of the tag. Admin layers
// (4, 5) come after plugin styles (3) so the admin wins on overlap, matching
// the main-page cascade.
const injectStyles = (
  doc: Document,
  pluginStyles: string,
  appearanceVars: string,
  customStyles: string,
) => {
  // Appended after the author's content, low to high priority.
  if (pluginStyles) {
    const plugins = doc.createElement('style');
    plugins.textContent = pluginStyles;
    doc.head.appendChild(plugins);
  }
  if (appearanceVars) {
    const vars = doc.createElement('style');
    vars.textContent = `:root { ${appearanceVars} }`;
    doc.head.appendChild(vars);
  }
  if (customStyles) {
    const custom = doc.createElement('style');
    custom.textContent = customStyles;
    doc.head.appendChild(custom);
  }
  // Prepended before the author's head styles → lowest priority.
  const baseline = doc.createElement('link');
  baseline.setAttribute('rel', 'stylesheet');
  baseline.setAttribute('href', PLUGIN_STYLESHEET);
  doc.head.insertBefore(baseline, doc.head.firstChild);
};

// A frame loads link clicks and form submits inside itself by default, so an
// <a href> would navigate the little embedded frame instead of the page. A
// <base target="_top"> makes links/forms target the top-level window like
// normal page links; its href gives relative URLs a real base to resolve
// against (srcdoc documents otherwise resolve against about:srcdoc).
const applyBase = (doc: Document) => {
  const base = doc.createElement('base');
  base.setAttribute('target', '_top');
  base.setAttribute('href', window.location.href);
  doc.head.insertBefore(base, doc.head.firstChild);
};

// PluginTabFrame renders a plugin tab's HTML as a real (same-origin) document
// via srcdoc. Scripts run in a normal document context — document.getElementById,
// document.currentScript, inline on* handlers and persistent globals all work —
// so plugin authors write ordinary interactive HTML with no special contract.
// The host injects the theme (baseline + appearance overrides + custom CSS)
// and auto-sizes the frame to its content. The frame starts hidden and is
// revealed once styles are in and height is set, so there's no flash of
// unstyled content and no resize jump.
export const PluginTabFrame: FC<PluginTabFrameProps> = ({ content }) => {
  const frameRef = useRef<HTMLIFrameElement>(null);
  const observerRef = useRef<ResizeObserver | null>(null);

  const clientConfig = useRecoilValue<ClientConfig>(clientConfigStateAtom);
  const { appearanceVariables, customStyles, pluginStyles } = clientConfig;
  const appearanceVars = Object.keys(appearanceVariables || {})
    .filter(variable => !!appearanceVariables[variable])
    .map(variable => `--${variable}: ${appearanceVariables[variable]}`)
    .join(';');

  const handleLoad = () => {
    const iframe = frameRef.current;
    const doc = iframe?.contentDocument;
    // Cast to include globals: the frame's own ResizeObserver (used below so
    // it observes nodes in its own document) isn't on the DOM Window type.
    const win = iframe?.contentWindow as (Window & typeof globalThis) | null | undefined;
    if (!iframe || !doc || !win) return;

    applyBase(doc);
    injectStyles(doc, pluginStyles, appearanceVars, customStyles);
    // The container behind the frame paints the panel background; let it
    // show through.
    doc.documentElement.style.background = 'transparent';
    doc.body.style.background = 'transparent';

    // Match the frame's height to its content so it reads as part of the
    // page rather than a scrollable box. Use the frame's own ResizeObserver
    // so it observes nodes in its own document.
    observerRef.current?.disconnect();
    const resize = () => {
      iframe.style.height = `${doc.documentElement.scrollHeight}px`;
    };
    resize();
    const observer = new win.ResizeObserver(resize);
    observer.observe(doc.documentElement);
    observerRef.current = observer;

    // Reveal now that it's styled and sized.
    iframe.style.visibility = 'visible';
  };

  useEffect(() => () => observerRef.current?.disconnect(), []);

  return (
    <iframe
      ref={frameRef}
      srcDoc={content}
      onLoad={handleLoad}
      title="Plugin tab"
      sandbox="allow-scripts allow-same-origin allow-forms allow-popups"
      style={{ border: 'none', width: '100%', display: 'block', visibility: 'hidden' }}
    />
  );
};
