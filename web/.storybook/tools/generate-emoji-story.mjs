import fs from 'fs';
import path from 'path';

import { readdirSync } from 'fs';
import handlebars from 'handlebars';

handlebars.registerHelper('capitalize', function (str) {
  return str.charAt(0).toUpperCase() + str.slice(1);
});

function getDirectories(path) {
  return fs.readdirSync(path).filter(function (file) {
    return fs.statSync(path + '/' + file).isDirectory();
  });
}

const emojiDir = path.resolve('../../../static/img/emoji');

const emojiCollectionDirs = getDirectories(emojiDir).map(dir => {
  return dir;
});

let emojiCollections = {};

emojiCollectionDirs.forEach(collection => {
  const emojiCollection = readdirSync(path.resolve(emojiDir, collection))
    .filter(f => f.toLowerCase() !== 'license.md')
    .map(emoji => {
      return { name: emoji, src: `img/emoji/${collection}/${emoji}` };
    });
  emojiCollections[collection] = { name: collection, images: emojiCollection };
});

const template = fs.readFileSync('./Emoji.stories.mdx', 'utf8');
let t = handlebars.compile(template);
let output = t({ emojiCollections });
console.log(output);
