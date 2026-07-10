import { execFileSync } from 'node:child_process';
import { createHash } from 'node:crypto';
import { readdirSync, readFileSync, renameSync, writeFileSync } from 'node:fs';
import { dirname, join, relative, resolve } from 'node:path';
import { fileURLToPath } from 'node:url';

const root = resolve(dirname(fileURLToPath(import.meta.url)), '..');
const webRoot = join(root, 'web');
const outputPath = join(root, 'THIRD-PARTY-LICENSES.txt');
const targets = [
  ['darwin', 'amd64'],
  ['darwin', 'arm64'],
  ['linux', '386'],
  ['linux', 'amd64'],
  ['linux', 'arm', '7'],
  ['linux', 'arm64'],
];
const knownMissing = new Map([
  [
    'Go module: github.com/nareix/joy5@v0.0.0-20210317075623-2c912ca30590',
    'No license is published with this version or in the upstream repository. Redistribution permission is unresolved.',
  ],
  [
    'Go module: github.com/schollz/sqlite3dump@v1.3.1',
    'The pinned archive predates the MIT license now published by the upstream repository.',
  ],
]);
const legalName = /^(licen[cs]e|notice|copying|copyright|third[-_]party[-_]notices?)([._-].*)?$/i;

function run(command, args, options = {}) {
  return execFileSync(command, args, {
    encoding: 'utf8',
    maxBuffer: 16 * 1024 * 1024,
    ...options,
  }).trim();
}

function legalDocuments(directory) {
  return readdirSync(directory, { withFileTypes: true })
    .filter(entry => entry.isFile() && legalName.test(entry.name))
    .map(entry => {
      const text = `${readFileSync(join(directory, entry.name), 'utf8').replaceAll('\r\n', '\n').trim()}\n`;
      return {
        name: entry.name,
        text,
        hash: createHash('sha256').update(text).digest('hex'),
      };
    })
    .sort((a, b) => a.name.localeCompare(b.name));
}

function goComponents() {
  const template =
    '{{if not .Standard}}{{with .Module}}{{if not .Main}}{{.Path}}\t{{.Version}}\t{{if .Replace}}{{.Replace.Dir}}{{else}}{{.Dir}}{{end}}{{end}}{{end}}{{end}}';
  const modules = new Map();

  for (const [goos, goarch, goarm = ''] of targets) {
    const output = run(
      'go',
      ['list', '-deps', '-tags=sqlite_omit_load_extension', '-f', template, '.'],
      {
        cwd: root,
        env: {
          ...process.env,
          CGO_ENABLED: '1',
          GOOS: goos,
          GOARCH: goarch,
          GOARM: goarm,
        },
      },
    );
    for (const line of output.split('\n').filter(Boolean)) {
      const [name, version, directory] = line.split('\t');
      modules.set(`${name}@${version}`, { name, version, directory });
    }
  }

  const components = [...modules.values()].map(module => ({
    id: `Go module: ${module.name}@${module.version}`,
    source: `https://pkg.go.dev/${module.name}@${module.version}`,
    documents: legalDocuments(module.directory),
  }));

  const goRoot = run('go', ['env', 'GOROOT']);
  components.push({
    id: `Go runtime: ${run('go', ['env', 'GOVERSION'])}`,
    source: 'https://go.dev/LICENSE',
    documents: legalDocuments(goRoot),
  });

  return components;
}

function webComponents() {
  const packages = JSON.parse(run('npm', ['query', '.prod', '--json'], { cwd: webRoot }));
  const components = new Map();

  for (const pkg of packages) {
    if (!pkg.name || !pkg.version || resolve(pkg.path) === webRoot) continue;
    const id = `Web package: ${pkg.name}@${pkg.version}`;
    const item = {
      id,
      source: pkg.resolved || `https://www.npmjs.com/package/${pkg.name}/v/${pkg.version}`,
      declared: typeof pkg.license === 'string' ? pkg.license : JSON.stringify(pkg.license || ''),
      documents: legalDocuments(pkg.path),
    };
    const existing = components.get(id);
    if (existing && JSON.stringify(existing.documents) !== JSON.stringify(item.documents)) {
      throw new Error(`${id} has inconsistent license files across installed copies`);
    }
    components.set(id, existing || item);
  }

  return [...components.values()];
}

function bundledComponents() {
  const emojiRoot = join(root, 'static', 'img', 'emoji');
  return readdirSync(emojiRoot, { withFileTypes: true })
    .filter(entry => entry.isDirectory())
    .flatMap(entry => {
      const directory = join(emojiRoot, entry.name);
      const documents = legalDocuments(directory);
      return documents.length
        ? [
            {
              id: `Bundled asset: ${relative(root, directory)}`,
              source: relative(root, directory),
              documents,
            },
          ]
        : [];
    });
}

function render(components) {
  components.sort((a, b) => a.id.localeCompare(b.id));
  const documents = new Map();

  const inventory = components
    .map(component => {
      let missing = '';
      if (!component.documents.length) {
        missing =
          knownMissing.get(component.id) ||
          (component.declared
            ? `The package contains no root license text. package.json declares: ${component.declared}.`
            : 'No license file or declaration was found.');
        if (!knownMissing.has(component.id) && !component.declared) {
          throw new Error(`${component.id} has no license file or declaration`);
        }
      }
      for (const document of component.documents) {
        const current = documents.get(document.hash) || {
          ...document,
          uses: [],
        };
        current.uses.push(`${component.id} (${document.name})`);
        documents.set(document.hash, current);
      }
      return [
        component.id,
        `  Source: ${component.source}`,
        component.declared ? `  Declared license: ${component.declared}` : '',
        missing ? `  LICENSE TEXT NOT FOUND: ${missing}` : '',
        ...component.documents.map(document => `  ${document.name}: SHA256 ${document.hash}`),
      ]
        .filter(Boolean)
        .join('\n');
    })
    .join('\n\n');

  const texts = [...documents.values()]
    .sort((a, b) => a.hash.localeCompare(b.hash))
    .map(document =>
      [
        '='.repeat(80),
        `SHA256: ${document.hash}`,
        ...document.uses.sort().map(use => `Used by: ${use}`),
        '-'.repeat(80),
        document.text,
      ].join('\n'),
    )
    .join('\n');

  return [
    'THIRD-PARTY LICENSES',
    '',
    'Generated by: npm run licenses',
    'Inputs: release-target Go dependencies, npm production dependencies, and bundled emoji licenses.',
    'Do not edit this file by hand.',
    '',
    'COMPONENT INVENTORY',
    '',
    inventory,
    '',
    'LICENSE AND NOTICE TEXTS',
    '',
    texts,
  ].join('\n');
}

const contents = render([...goComponents(), ...webComponents(), ...bundledComponents()]);
const temporaryPath = `${outputPath}.tmp`;
writeFileSync(temporaryPath, contents);
renameSync(temporaryPath, outputPath);
console.log(`Wrote ${relative(root, outputPath)}`);
