const { buildReleaseRules, buildReleaseNotesConfig } = require('../release.utils.cjs');

module.exports = {
  branches: ['main'],
  tagFormat: 'cli-v${version}',
  plugins: [
    ['@semantic-release/commit-analyzer', {
      preset: 'conventionalcommits',
      releaseRules: buildReleaseRules('cli'),
    }],
    ['@semantic-release/release-notes-generator', buildReleaseNotesConfig('cli')],
    ['@semantic-release/github', {
      successComment: false,
      labels: false,
    }],
    ['@semantic-release/exec', {
      publishCmd: 'echo "${nextRelease.version}" > nextversion && echo "true" > released',
    }],
  ],
};
