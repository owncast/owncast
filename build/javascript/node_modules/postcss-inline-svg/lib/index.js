const postcss = require("postcss");
const { parseRuleDefinition, getRuleParams } = require("./parseRule.js");
const parseDeclValue = require("./parseDeclValue.js");
const resolveId = require("./resolveId.js");
const load = require("./load.js");

function removeLoader(loader) {
  if (!loader.error && loader.node.type === "atrule") {
    loader.node.remove();
  }
}

function applyInliner(inliner) {
  if (!inliner.loader.error) {
    inliner.valueNode.value = "url";
    inliner.valueNode.nodes = [
      {
        type: "word",
        value: inliner.loader.svg
      }
    ];
  }
}

function stringifyInliner(inliner) {
  if (!inliner.loader.error) {
    inliner.node.value = String(inliner.parsedValue);
  }
}

module.exports = postcss.plugin(
  "postcss-inline-svg",
  (opts = {}) => (css, result) => {
    const loadersMap = {};
    const loaders = [];
    const inliners = [];

    css.walk(node => {
      if (node.type === "atrule") {
        if (node.name === "svg-load") {
          try {
            const file =
              node.source && node.source.input && node.source.input.file;
            const { name, url } = parseRuleDefinition(node.params);
            const { params, selectors } = getRuleParams(node);
            const loader = {
              id: resolveId(file, url, opts),
              parent: file,
              params,
              selectors,
              node
            };
            loaders.push(loader);
            loadersMap[name] = loader;
          } catch (e) {
            node.warn(result, e.message);
          }
        }
      } else if (node.type === "decl") {
        if (
          node.value.indexOf("svg-load(") !== -1 ||
          node.value.indexOf("svg-inline(") !== -1
        ) {
          try {
            const file =
              node.source && node.source.input && node.source.input.file;
            const statements = parseDeclValue(node.value);
            statements.loaders.forEach(
              ({ url, params, valueNode, parsedValue }) => {
                const loader = {
                  id: resolveId(file, url, opts),
                  parent: file,
                  params,
                  selectors: {},
                  node
                };
                loaders.push(loader);
                inliners.push({
                  loader,
                  node,
                  valueNode,
                  parsedValue
                });
              }
            );
            statements.inliners.forEach(({ name, valueNode, parsedValue }) => {
              const loader = loadersMap[name];
              if (loader) {
                inliners.push({
                  loader,
                  node,
                  valueNode,
                  parsedValue
                });
              } else {
                node.warn(result, `"${name}" svg is not defined`);
              }
            });
          } catch (e) {
            node.warn(result, e.message);
          }
        }
      }
    });

    const promises = loaders.map(loader => {
      return load(loader.id, loader.params, loader.selectors, opts)
        .then(code => {
          loader.svg = code;
          result.messages.push({
            type: "dependency",
            file: loader.id,
            parent: loader.parent
          });
        })
        .catch(err => {
          loader.error = true;
          loader.node.warn(result, err.message);
        });
    });

    return Promise.all(promises).then(() => {
      loaders.forEach(removeLoader);
      inliners.forEach(applyInliner);
      inliners.forEach(stringifyInliner);
    });
  }
);
