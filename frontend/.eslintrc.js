module.exports = {
  root: true,
  env: {
    node: true,
  },
  extends: ["plugin:vue/essential", "@vue/typescript/recommended"],
  parserOptions: {
    ecmaVersion: 2020,
  },
  rules: {
    "no-console": process.env.NODE_ENV === "production" ? "error" : "off",
    "no-debugger": process.env.NODE_ENV === "production" ? "error" : "off",
    "vue/script-indent": ["error", 2, {baseIndent: 1}],
    "@typescript-eslint/camelcase": ["error", {properties: "never"}],
  },
  overrides: [
    {
      files: ["*.vue"],
      rules: {indent: "off"},
    },
  ],
};
