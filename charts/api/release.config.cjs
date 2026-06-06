const { buildReleaseRules, buildReleaseNotesConfig } = require('../../release.utils.cjs');

module.exports = {
  branches: ['main'],
  tagFormat: 'chart-api-v${version}',
  plugins: [
    ['@semantic-release/commit-analyzer', {
      preset: 'conventionalcommits',
      releaseRules: buildReleaseRules('chart-api'),
    }],
    ['@semantic-release/release-notes-generator', buildReleaseNotesConfig('chart-api')],
    ['@semantic-release/github', {
      successComment: false,
      labels: false,
    }],
    ['@semantic-release/exec', {
      publishCmd: 'echo "${nextRelease.version}" > nextversion && echo "true" > released',
    }],
  ],
};
