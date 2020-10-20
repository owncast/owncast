"use strict";

Object.defineProperty(exports, "__esModule", {
  value: true
});
exports.default = void 0;

function createPlugin(plugin, config) {
  return {
    handler: plugin,
    config
  };
}

createPlugin.withOptions = function (pluginFunction, configFunction = () => ({})) {
  const optionsFunction = function (options) {
    return {
      handler: pluginFunction(options),
      config: configFunction(options)
    };
  };

  optionsFunction.__isOptionsFunction = true;
  return optionsFunction;
};

var _default = createPlugin;
exports.default = _default;