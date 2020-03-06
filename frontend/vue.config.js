module.exports = {
  outputDir: "dist",
  lintOnSave: false,
  devServer: {
    proxy: {
      "^/data": {target: "http://localhost:8080"},
      "^/api": {target: "http://localhost:8080"},
    },
  },
};
