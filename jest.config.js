// force timezone to UTC to allow tests to work regardless of local timezone
// generally used by snapshots, but can affect specific tests
process.env.TZ = 'UTC';

const path = require('path');
const { grafanaESModules, nodeModulesToTransform } = require('./.config/jest/utils');

const baseConfig = require('./.config/jest.config');

module.exports = {
  // Jest configuration provided by Grafana scaffolding
  ...baseConfig,

  moduleNameMapper: {
    ...baseConfig.moduleNameMapper,
    // Override the scaffolded react-inlinesvg mock with one that forwards props
    // (e.g. role/aria-label) so accessible queries against Icon-based controls
    // (like the Select "Clear value" button in @grafana/ui v13) keep working.
    'react-inlinesvg': path.resolve(__dirname, 'tests', 'mocks', 'react-inlinesvg.tsx'),
  },

  transformIgnorePatterns: [nodeModulesToTransform([...grafanaESModules, '@marcbachmann/cel-js'])],
};
