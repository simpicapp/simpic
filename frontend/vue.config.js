module.exports = {
  outputDir: "dist",
  lintOnSave: false,
  devServer: {
    proxy: {
      "^/data": {target: "http://localhost:8080"},
      "^/user": {target: "http://localhost:8080"},
      "^/login": {target: "http://localhost:8080"},
      "^/logout": {target: "http://localhost:8080"},
      "^/timeline": {target: "http://localhost:8080"},
      "^/albums": {target: "http://localhost:8080"},
      "^/photos": {target: "http://localhost:8080"},
    },
  },
};
