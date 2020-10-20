const valueParser = require("postcss-value-parser");

function parseRuleDefinition(params) {
  const { nodes } = valueParser(params);
  if (
    nodes.length !== 3 ||
    nodes[0].type !== "word" ||
    nodes[1].type !== "space" ||
    nodes[2].type !== "function" ||
    nodes[2].value !== "url" ||
    nodes[2].nodes.length === 0
  ) {
    throw Error('Invalid "@svg-load" definition');
  }
  return {
    name: nodes[0].value,
    url: nodes[2].nodes[0].value
  };
}

function getRuleParams(rule) {
  const params = {};
  const selectors = {};

  rule.each(node => {
    if (node.type === "decl") {
      params[node.prop] = node.value;
    } else if (node.type === "rule") {
      const selector = selectors[node.selectors] || {};
      node.each(child => {
        if (child.type === "decl") {
          selector[child.prop] = child.value;
        }
      });
      selectors[node.selectors] = selector;
    }
  });

  return {
    params,
    selectors
  };
}

module.exports = {
  parseRuleDefinition,
  getRuleParams
};
