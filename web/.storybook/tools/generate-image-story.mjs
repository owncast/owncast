import fs from 'fs';
import path, { resolve } from 'path';

import { readdirSync, lstatSync } from 'fs';
import handlebars from 'handlebars';

handlebars.registerHelper('capitalize', function (str) {
  return str.charAt(0).toUpperCase() + str.slice(1);
});

const args = process.argv;
const dir = args[2];
const title = args[3];
const category = args[4];
const publicPath = args[5];

if (args.length < 6) {
  console.error('Usage: generate-image-story.mjs <dir> <title> <category> <webpublicpath>');
  process.exit(1);
}

const images = readdirSync(dir)
  .map(img => {
    const resolvedPath = path.resolve(dir, img);
    if (lstatSync(resolvedPath).isDirectory()) {
      return;
    }

    return { name: img, src: `${publicPath}/${img}` };
  })
  .filter(Boolean);

const template = fs.readFileSync('./Images.stories.mdx', 'utf8');
let t = handlebars.compile(template);
let output = t({ images, title, category });
console.log(output);
