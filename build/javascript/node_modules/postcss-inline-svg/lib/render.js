const { parseDOM } = require("htmlparser2");
const serialize = require("dom-serializer");

module.exports = function render(code, ...processors) {
  const dom = parseDOM(code, { xmlMode: true });

  processors.forEach(processor => processor(dom));

  return serialize(dom);
};
