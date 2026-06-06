const ALL_SCOPES = ['api', 'ui', 'cli', 'chart-api', 'chart-ui', 'repo'];

function buildReleaseRules(scope) {
  return [
    ...ALL_SCOPES.filter(s => s !== scope).map(s => ({ scope: s, release: false })),
    // Change these later, patch for everythng right now while it's very early 
    // build.
    { type: 'feat', scope: scope, release: 'patch' },
    { type: 'fix', scope: scope, release: 'patch' },
    { type: 'perf', scope: scope, release: 'patch' },
    { type: 'revert', scope: scope, release: 'patch' },
    { type: 'build', scope: scope, release: 'patch' },
    { type: 'chore', scope: scope, release: 'patch' },
    { type: 'refactor', scope: scope, release: 'patch' },
    { type: 'test', scope: scope, release: 'patch' },
  ];
}

function buildReleaseNotesConfig(scope) {
  const types = [
    { type: 'feat',     section: 'Features' },
    { type: 'fix',      section: 'Bug Fixes' },
    { type: 'perf',     section: 'Performance' },
    { type: 'revert',   section: 'Reverts' },
    { type: 'build',    section: 'Build' },
    { type: 'chore',    section: 'Chores' },
    { type: 'refactor', section: 'Refactors' },
    { type: 'test',     section: 'Tests' },
  ];

  const typeToSection = Object.fromEntries(types.map(t => [t.type, t.section]));

  return {
    preset: 'conventionalcommits',
    presetConfig: { types },
    writerOpts: {
      transform: (commit) => {
        if (commit.scope !== scope) return false;

        const out = { ...commit };

        if (typeToSection[out.type]) {
          out.type = typeToSection[out.type];
        }

        if (out.hash) {
          out.shortHash = out.hash.substring(0, 7);
        }

        return out;
      },
    },
  };
}

module.exports = { buildReleaseRules, buildReleaseNotesConfig };
