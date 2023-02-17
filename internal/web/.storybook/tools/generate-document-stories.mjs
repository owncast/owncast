import fs from 'fs';
import handlebars from 'handlebars';

const template = fs.readFileSync('./Document.stories.mdx', 'utf8');
let t = handlebars.compile(template, { noEscape: true });

const documents = [
  {
    title: 'Product Definition',
    name: 'ProductDefinition',
    path: '../../../docs/product-definition.md',
  },
  { title: 'Design', name: 'Design', path: '../../../.design/DESIGN.md' },
  {
    title: 'Building Frontend Components',
    name: 'WebComponents',
    path: '../../../web/components/_COMPONENT_HOW_TO.md',
  },
  {
    title: 'Get Started with Owncast Development',
    name: 'Development',
    path: '/tmp/development.md',
  },
];

documents.forEach(doc => {
  if (!fs.existsSync(doc.path)) {
    return;
  }

  const document = fs.readFileSync(doc.path, 'utf8');
  const output = t({ name: doc.name, title: doc.title, content: document });
  fs.writeFileSync(`../stories-category-doc-pages/${doc.name}.stories.mdx`, output);
});
