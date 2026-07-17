#!/usr/bin/env node
/**
 * Capture before/after screenshots of Storybook stories affected by a
 * change set, for the PR comment posted by
 * .github/workflows/storybook-screenshots.yml.
 *
 * "Before" comes from the published develop Storybook (BASELINE_URL);
 * "after" from a locally served build of this branch (STATIC_URL). A story
 * is affected when a changed file lives under the directory of the story's
 * importPath in the build's index.json -- stories are co-located with
 * their components in this repo, so a directory match is the whole mapping.
 *
 * Environment:
 *   CHANGED_FILES    whitespace-separated paths relative to web/ (required)
 *   STATIC_URL       URL serving this branch's storybook-static (required)
 *   BASELINE_URL     published Storybook to diff against
 *                    (default https://ui.owncast.online)
 *   PUBLIC_BASE_URL  public URL prefix where the PNGs will be hosted
 *   OUT_DIR          output directory (default screenshot-diff)
 *   CHROME           chrome binary (default google-chrome)
 *
 * Writes OUT_DIR/img/<story-id>-{before,after}.png and
 * OUT_DIR/comment-body.md, and appends stories=<n> to $GITHUB_OUTPUT.
 */
import { execFileSync } from 'node:child_process';
import fs from 'node:fs';
import path from 'node:path';
import puppeteer from 'puppeteer-core';

function required(name) {
  const value = process.env[name];
  if (!value) {
    console.error(`${name} is required`);
    process.exit(1);
  }
  return value;
}

const STATIC_URL = required('STATIC_URL');
const BASELINE_URL = process.env.BASELINE_URL || 'https://ui.owncast.online';
const PUBLIC_BASE_URL = process.env.PUBLIC_BASE_URL || '';
const OUT_DIR = process.env.OUT_DIR || 'screenshot-diff';
// puppeteer-core stats executablePath literally, so resolve bare command
// names (e.g. the runner's google-chrome) through PATH first.
const chromeEnv = process.env.CHROME || 'google-chrome';
const CHROME = chromeEnv.includes('/')
  ? chromeEnv
  : execFileSync('/bin/sh', ['-c', `command -v ${chromeEnv}`], { encoding: 'utf8' }).trim();

// Keep the comment reviewable on sweeping changes; Chromatic covers the rest.
const MAX_STORIES = 12;
const WIDTH = 1000;
const HEIGHT = 800;

const changed = (process.env.CHANGED_FILES || '').split(/\s+/).filter(Boolean);

async function loadIndex(base) {
  const res = await fetch(`${base}/index.json`);
  if (!res.ok) throw new Error(`${base}/index.json: HTTP ${res.status}`);
  return res.json();
}

const [head, baseline] = await Promise.all([
  loadIndex(STATIC_URL),
  loadIndex(BASELINE_URL).catch(err => {
    console.warn(`baseline index unavailable (${err.message}); treating all stories as new`);
    return { entries: {} };
  }),
]);

const affected = Object.values(head.entries)
  .filter(entry => entry.type === 'story')
  .filter(entry => {
    const dir = path.posix.dirname(entry.importPath.replace(/^\.\//, ''));
    return changed.some(file => file.startsWith(`${dir}/`));
  })
  .sort((a, b) => a.id.localeCompare(b.id));

const shown = affected.slice(0, MAX_STORIES);
const omitted = affected.slice(MAX_STORIES);

fs.rmSync(OUT_DIR, { recursive: true, force: true });
fs.mkdirSync(path.join(OUT_DIR, 'img'), { recursive: true });

const browser = await puppeteer.launch({
  executablePath: CHROME,
  args: ['--no-sandbox', '--disable-gpu', '--hide-scrollbars', '--force-device-scale-factor=1'],
});

async function capture(base, id, file) {
  const url = `${base}/iframe.html?id=${id}&viewMode=story`;
  const page = await browser.newPage();
  try {
    await page.setViewport({ width: WIDTH, height: HEIGHT });
    for (let attempt = 1; ; attempt += 1) {
      await page.goto(url, { waitUntil: 'networkidle0', timeout: 60_000 });
      try {
        // Storybook flips the body to sb-show-main and populates
        // #storybook-root once the story (including msw-gated fetches) has
        // rendered; an error lands on sb-show-errordisplay instead.
        await page.waitForFunction(
          () =>
            document.body.classList.contains('sb-show-main') &&
            document.getElementById('storybook-root')?.childElementCount > 0,
          { timeout: 30_000 }
        );
        break;
      } catch {
        // Chunk loads from the hosted baseline occasionally fail
        // transiently; reload once, then capture whatever is on screen
        // (a Storybook error page is useful evidence too).
        if (attempt >= 2) {
          console.warn(`story never rendered for ${id} at ${base}; capturing current state`);
          break;
        }
      }
    }
    await page.evaluate(() => document.fonts.ready);
    // Brief settle for image decoding/paint after render is signaled.
    await new Promise(resolve => setTimeout(resolve, 500));
    await page.screenshot({ path: file, fullPage: true });
    return true;
  } catch (err) {
    console.warn(`capture failed for ${id} at ${base}: ${err.message}`);
    return false;
  } finally {
    await page.close();
  }
}

const lines = [];
for (const story of shown) {
  const isNew = !baseline.entries?.[story.id];
  const beforeName = `${story.id}-before.png`;
  const afterName = `${story.id}-after.png`;
  const hasBefore =
    !isNew && (await capture(BASELINE_URL, story.id, path.join(OUT_DIR, 'img', beforeName)));
  const hasAfter = await capture(STATIC_URL, story.id, path.join(OUT_DIR, 'img', afterName));

  const img = (name, alt) => `<img src="${PUBLIC_BASE_URL}/${name}" alt="${alt}" width="420">`;
  const beforeCell = isNew
    ? '_new story_'
    : hasBefore
      ? img(beforeName, `${story.title} — ${story.name} on develop`)
      : '_capture failed_';
  const afterCell = hasAfter
    ? img(afterName, `${story.title} — ${story.name} on this PR`)
    : '_capture failed_';

  lines.push(
    `#### [${story.title} — ${story.name}](${BASELINE_URL}/?path=/story/${story.id})`,
    '',
    '| develop | this PR |',
    '| --- | --- |',
    `| ${beforeCell} | ${afterCell} |`,
    ''
  );
}
if (omitted.length) {
  lines.push(`_+${omitted.length} more affected ${omitted.length === 1 ? 'story' : 'stories'} not shown._`, '');
}
fs.writeFileSync(path.join(OUT_DIR, 'comment-body.md'), lines.join('\n'));
await browser.close();

if (process.env.GITHUB_OUTPUT) {
  fs.appendFileSync(process.env.GITHUB_OUTPUT, `stories=${affected.length}\n`);
}
console.log(`${affected.length} affected stories; captured ${shown.length}.`);
