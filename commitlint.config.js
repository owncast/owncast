// Commit message linting for the commit-msg hook (see lefthook.yml). Kept
// self-contained rather than `extends: ['@commitlint/config-conventional']`
// so it works with just the commitlint CLI, without requiring the
// config-conventional package to be resolvable from the repo root. The allowed
// types mirror the self-contained fallback check in lefthook.yml's commit-msg
// hook, so both paths enforce the same Conventional Commits format.
module.exports = {
	rules: {
		'type-enum': [
			2,
			'always',
			[
				'build',
				'chore',
				'ci',
				'docs',
				'feat',
				'fix',
				'perf',
				'refactor',
				'revert',
				'style',
				'test',
			],
		],
		'type-empty': [2, 'never'],
		'subject-empty': [2, 'never'],
	},
};
