module.exports = {
  extends: [
    'stylelint-config-standard',
    'stylelint-config-prettier',
  ],
  plugins: [],
  rules: {
    // Custom rules can be added here
    'selector-class-pattern': '^[a-z][a-zA-Z0-9]*$',
    'declaration-no-important': true,
    'max-nesting-depth': 3,
    'color-no-hex': null,
    'length-zero-no-unit': true,
    'font-family-no-missing-generic-family-keyword': true,
    'no-descending-specificity': null,
    'declaration-block-trailing-semicolon': 'always',
    'declaration-block-single-line-max-declarations': 1,
  },
  ignoreFiles: [
    'node_modules/**/*',
    'dist/**/*',
    'build/**/*',
    '**/*.js',
    '**/*.jsx',
    '**/*.ts',
    '**/*.tsx',
  ],
};