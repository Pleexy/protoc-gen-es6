module.exports = api => {
  api.cache(true);
  return {
    babelrcRoots: [
      './**',
    ],
    exclude: "/node_modules/",
    presets: [
      '@babel/preset-flow'
    ],
    plugins: [
      '@babel/plugin-transform-modules-commonjs',
//      '@babel/plugin-proposal-private-methods',
    [ '@babel/plugin-proposal-class-properties', {"loose": true}]
    ]

  };
};