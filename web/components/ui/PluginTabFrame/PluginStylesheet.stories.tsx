import { Meta, StoryFn } from '@storybook/nextjs';
import { FC, useEffect, useRef } from 'react';

// The served baseline stylesheet (web/public/styles/plugin.css), which
// Storybook serves at /styles/plugin.css via its staticDirs. Editing that file
// and reloading the story reflects here immediately, so this story is a live
// sandbox for iterating on the plugin stylesheet.
const PLUGIN_STYLESHEET = '/styles/plugin.css';

// A self-contained mirror of PluginTabFrame's rendering: a srcdoc iframe with
// the baseline linked into its head, a transparent background so the host
// panel shows through, and auto-sizing to content. It deliberately omits
// PluginTabFrame's Recoil reads (pluginStyles / appearanceVariables /
// customStyles) so the
// story has no store dependency and stays a pure stylesheet preview.
const PluginFramePreview: FC<{ html: string }> = ({ html }) => {
  const frameRef = useRef<HTMLIFrameElement>(null);
  const observerRef = useRef<ResizeObserver | null>(null);

  const handleLoad = () => {
    const iframe = frameRef.current;
    const doc = iframe?.contentDocument;
    const win = iframe?.contentWindow as (Window & typeof globalThis) | null | undefined;
    if (!iframe || !doc || !win) return;

    // <base target="_top"> so links/forms behave; href gives relative URLs a
    // real base (srcdoc otherwise resolves against about:srcdoc).
    const base = doc.createElement('base');
    base.setAttribute('target', '_top');
    base.setAttribute('href', window.location.href);
    doc.head.insertBefore(base, doc.head.firstChild);

    // The baseline, prepended so the author's own head styles would win.
    const link = doc.createElement('link');
    link.setAttribute('rel', 'stylesheet');
    link.setAttribute('href', PLUGIN_STYLESHEET);
    doc.head.insertBefore(link, doc.head.firstChild);

    // Let the panel behind the frame paint the background.
    doc.documentElement.style.background = 'transparent';
    doc.body.style.background = 'transparent';

    // Match the frame height to its content so it reads as part of the page.
    observerRef.current?.disconnect();
    const resize = () => {
      iframe.style.height = `${doc.documentElement.scrollHeight}px`;
    };
    resize();
    const observer = new win.ResizeObserver(resize);
    observer.observe(doc.documentElement);
    observerRef.current = observer;

    iframe.style.visibility = 'visible';
  };

  useEffect(() => () => observerRef.current?.disconnect(), []);

  return (
    <iframe
      ref={frameRef}
      srcDoc={html}
      onLoad={handleLoad}
      title="Plugin stylesheet preview"
      sandbox="allow-scripts allow-same-origin allow-forms allow-popups"
      style={{ border: 'none', width: '100%', display: 'block', visibility: 'hidden' }}
    />
  );
};

// Sample content exercising every element the baseline styles plus all the
// opt-in helper classes (card, card-grid, tag, stack, row, muted). Plain
// semantic HTML, no inline styles, so what you see is entirely the stylesheet.
const SAMPLE_HTML = `
<h1>Plugin stylesheet preview</h1>
<p>
  This is how plain semantic HTML renders inside a plugin tab. Everything below
  is styled only by the injected baseline. Edit
  <code>web/public/styles/plugin.css</code> and reload to iterate.
</p>

<h2>Typography</h2>
<p>
  Body copy with a <a href="#">link</a>, some <code>inline code</code>, and
  <span class="muted">muted secondary text</span> for captions.
</p>
<h3>Heading level 3</h3>
<h4>Heading level 4</h4>
<ul>
  <li>Unordered list item one</li>
  <li>Unordered list item two</li>
</ul>
<pre><code>// preformatted block
const answer = 42;</code></pre>
<hr />

<h2>Helper: cards &amp; grid</h2>
<div class="card-grid">
  <article class="card interactive">
    <h3>Album A</h3>
    <p class="muted">Artist A</p>
    <div class="row"><span class="tag">jazz</span><span class="tag">2024</span></div>
  </article>
  <article class="card interactive">
    <h3>Album B</h3>
    <p class="muted">Artist B</p>
    <div class="row"><span class="tag">ambient</span></div>
  </article>
  <article class="card interactive">
    <h3>Album C</h3>
    <p class="muted">Artist C</p>
  </article>
</div>

<h2>Helper: stack</h2>
<div class="card">
  <div class="stack">
    <strong>Vertical stack</strong>
    <span>First item</span>
    <span>Second item</span>
    <span>Third item</span>
  </div>
</div>

<h2>Forms</h2>
<form>
  <label>
    Text field
    <input type="text" placeholder="Type here" />
  </label>
  <label>
    Select
    <select>
      <option>Option A</option>
      <option>Option B</option>
    </select>
  </label>
  <label>
    Textarea
    <textarea placeholder="Notes"></textarea>
  </label>
  <label class="row"><input type="checkbox" /> Checkbox option</label>
  <div class="row">
    <button type="button">Primary button</button>
    <button type="button" class="secondary">Secondary button</button>
    <button type="button" disabled>Disabled</button>
  </div>
</form>

<h2>Table</h2>
<table>
  <thead>
    <tr><th>Name</th><th>Role</th></tr>
  </thead>
  <tbody>
    <tr><td>Ada</td><td>Engineer</td></tr>
    <tr><td>Lin</td><td>Designer</td></tr>
    <tr><td>Sam</td><td>Streamer</td></tr>
  </tbody>
</table>

<h2>Fieldset</h2>
<fieldset>
  <legend>Group</legend>
  <p>A fieldset groups related controls with a light bordered box.</p>
</fieldset>
`;

const meta = {
  title: 'owncast/Plugins/Plugin Stylesheet',
  component: PluginFramePreview,
  parameters: {
    docs: {
      description: {
        component: `Live sandbox for the plugin baseline stylesheet (\`web/public/styles/plugin.css\`),
the CSS Owncast injects into plugin tab and admin-page iframes. The frame here mirrors
how PluginTabFrame renders real plugin content: the baseline is linked in, the body is
transparent so the host panel shows through, and the frame auto-sizes to its content.

The sample HTML uses only plain semantic elements and the opt-in helper classes
(\`card\`, \`card-grid\`, \`tag\`, \`stack\`, \`row\`, \`muted\`), so everything you see comes from
the stylesheet. Edit \`plugin.css\` and reload to iterate.`,
      },
    },
  },
} satisfies Meta<typeof PluginFramePreview>;

export default meta;

// Renders the frame on a backdrop that mimics the viewer content panel
// (--theme-color-components-content-background) sitting inside the page, so the
// flush-on-panel effect the baseline relies on is visible.
const Template: StoryFn<typeof PluginFramePreview> = args => (
  <div
    style={{
      background: 'var(--theme-color-background-main)',
      padding: 24,
      minHeight: '100vh',
    }}
  >
    <div
      style={{
        background: 'var(--theme-color-components-content-background)',
        borderRadius: 'var(--theme-rounded-corners)',
        padding: 32,
        maxWidth: 1100,
        margin: '0 auto',
      }}
    >
      <PluginFramePreview {...args} />
    </div>
  </div>
);

// The default preview: full sample content on the content panel.
export const Default = {
  render: Template,
  args: {
    html: SAMPLE_HTML,
  },
};

// Just the helper classes, for iterating on cards / grid / tags in isolation.
export const HelperClasses = {
  render: Template,
  args: {
    html: `
<h2>Cards &amp; grid</h2>
<div class="card-grid">
  <article class="card interactive"><h3>Card one</h3><p class="muted">Hover for the lift</p><div class="row"><span class="tag">tag</span><span class="tag">tag</span></div></article>
  <article class="card interactive"><h3>Card two</h3><p class="muted">Secondary detail</p></article>
  <article class="card"><h3>Static card</h3><p class="muted">No hover</p></article>
</div>
<h2>Stack &amp; row</h2>
<div class="card"><div class="stack"><strong>Stack</strong><span>One</span><span>Two</span></div></div>
<div class="row"><span class="tag">one</span><span class="tag">two</span><span class="tag">three</span></div>
`,
  },
};
