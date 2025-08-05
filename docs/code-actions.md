# Code Actions

The OpenTofu Language Server implements a set of Code Actions which perform different actions on the current document. These commands are typically code fixes to either refactor code, fix problems or to beautify/refactor code.

## Available Actions

### `source.formatAll.opentofu`

The server will format a given document according to OpenTofu formatting conventions.

## Usage

### VS Code

To enable the format code action globally, set `source.formatAll.opentofu` to _true_ for the `editor.codeActionsOnSave` setting and set `editor.formatOnSave` to _false_.

```json
"editor.formatOnSave": false,
"editor.codeActionsOnSave": {
  "source.formatAll.opentofu": true
},
"[opentofu]": {
  "editor.defaultFormatter": "opentofu.vscode-opentofu",
}
```

> _Important:_ Disable `editor.formatOnSave` if you are using `source.formatAll.opentofu` in `editor.codeActionsOnSave`. The `source.formatAll.opentofu` code action is meant to be used instead of `editor.formatOnSave`, as it provides a [guarantee of order of execution](https://github.com/microsoft/vscode-docs/blob/71643d75d942e2c32cfd781c2b5322521775fb4a/release-notes/v1_44.md#explicit-ordering-for-editorcodeactionsonsave) based on the list provided. If you have both settings enabled, then your document will be formatted twice.

If you would like `editor.formatOnSave` to be _true_ for other extensions but _false_ for the OpenTofu extension, you can configure your settings as follows:

```json
"editor.formatOnSave": true,
"editor.codeActionsOnSave": {
  "source.formatAll.opentofu": true
},
"[opentofu]": {
  "editor.defaultFormatter": "opentofu.vscode-opentofu",
  "editor.formatOnSave": false,
},
```

Alternatively, you can include all OpenTofu related Code Actions inside the language specific setting if you prefer:

```json
"editor.formatOnSave": true,
"[opentofu]": {
  "editor.defaultFormatter": "opentofu.vscode-opentofu",
  "editor.formatOnSave": false,
  "editor.codeActionsOnSave": {
    "source.formatAll.opentofu": true
  },
},
```
