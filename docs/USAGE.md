# Usage of OpenTofu Language Server

This guide assumes you have installed the server by following instructions
in the [README.md](../README.md) if that is applicable to your client
(i.e. if the client doesn't download the server itself).

The following filetypes are supported by the OpenTofu Language Server:

- `terraform` - standard `*.tf` config files
- `terraform-vars` - variable files (`*.tfvars`)

_NOTE_ Clients should be configured to follow the above language ID conventions
and do **not** send `*.tf.json`, `*.tfvars.json` nor Packer HCL config
nor any other HCL config files as the server is not
equipped to handle these file types.

In most clients with a dedicated OpenTofu extension/plugin this is
already the default configuration, so you should not need to worry about it.

Instructions for popular IDEs are below and pull requests
for updates or addition of more IDEs are welcomed.

See also [settings](./SETTINGS.md) to understand
how you may configure the settings.

## Workspaces / Folders / Files

Most editors support opening folders. Such a root folder is commonly referred to
as "workspace". Opening folders is always preferred over individual files
as it allows the language server to index the whole folder and keep track
of changes more easily. We do however support "single-file mode" which provides
limited IntelliSense.

Indexing enables IntelliSense related to `module` blocks,
such as go-to-definition, completion of `module.*` references,
or workspace-wide symbol lookup.

The server will _not_ index any folders or files above the workspace root
initially opened in the editor.

## Emacs

If you are using `use-package`, you can put this in the [init.el](https://www.gnu.org/software/emacs/manual/html_node/emacs/Init-File.html)
file to install `lsp-mode`:

```emacs-lisp
(use-package lsp-mode
  :ensure t
  :hook ((terraform-mode . lsp-deferred)))
```

There are various other ways to install `lsp-mode` and they are
documented [here.](https://emacs-lsp.github.io/lsp-mode/page/installation/#installation)

<!--
TODO: We don't have an Emacs LSP yet
The `lsp-mode` language client for OpenTofu supports various features
like semantic tokens, code lens for references etc. There is more
detailed documentation [here](https://emacs-lsp.github.io/lsp-mode/page/lsp-terraform-ls/). -->

## IntelliJ IDE

- Install [LSP Support plugin](https://plugins.jetbrains.com/plugin/10209-lsp-support)
- Open Settings
- Go to `Languages & Frameworks → Language Server Protocol → Server Definitions`
  - Pick `Executable`
  - set `Extension` to `tf`
  - set `Path` to `tofu-ls`
  - set `Args` to `serve`
- Confirm by clicking `Apply`

Please note that the [Terraform plugin](https://plugins.jetbrains.com/plugin/7808-hashicorp-terraform--hcl-language-support)
provides overlapping functionality (and more features at the time of writing).
As a result having both enabled at the same time may result in suboptimal UX,
such as duplicate completion candidates.

## Sublime Text

- Install the [LSP package](https://github.com/sublimelsp/LSP#installation)
- Install the [LSP-terraform package](https://github.com/sublimelsp/LSP-terraform#installation)

## Vim / NeoVim

### coc.nvim

- Install the [coc.nvim plugin](https://github.com/neoclide/coc.nvim)
- Add the following snippet to the `coc-setting.json` file (editable via `:CocConfig` in NeoVim)

```json
{
  "languageserver": {
    "terraform": {
      "command": "tofu-ls",
      "args": ["serve"],
      "filetypes": ["terraform", "tf"],
      "initializationOptions": {},
      "settings": {}
    }
  }
}
```

Make sure to read through the [example vim configuration](https://github.com/neoclide/coc.nvim#example-vim-configuration) of the plugin, especially key remapping, which is required for completion to work correctly:

```vim
" Use <c-space> to trigger completion.
inoremap <silent><expr> <c-space> coc#refresh()
```

### vim-lsp

- [Install](https://opensource.com/article/20/2/how-install-vim-plugins) the following plugins:
  - [async.vim plugin](https://github.com/prabirshrestha/async.vim)
  - [vim-lsp plugin](https://github.com/prabirshrestha/vim-lsp)
  - [asyncomplete.vim plugin](https://github.com/prabirshrestha/asyncomplete.vim)
  - [asyncomplete-lsp.vim plugin](https://github.com/prabirshrestha/asyncomplete-lsp.vim)
- Add the following to your `.vimrc`:

```vim
if executable('tofu-ls')
    au User lsp_setup call lsp#register_server({
        \ 'name': 'tofu-ls',
        \ 'cmd': {server_info->['tofu-ls', 'serve']},
        \ 'whitelist': ['tofu'],
        \ })
endif
```

### YouCompleteMe

- [Install](https://opensource.com/article/20/2/how-install-vim-plugins) the following plugins:
  - [YouCompleteMe plugin](https://github.com/ycm-core/YouCompleteMe)
- Add the following to your `.vimrc`:

```vim
" Remove this line if additional custom language servers are set elsewhere
let g:ycm_language_server = []

if executable('tofu-ls')
    let g:ycm_language_server += [
        \   {
        \     'name': 'tofu',
        \     'cmdline': [ 'tofu-ls', 'serve' ],
        \     'filetypes': [ 'tofu', 'terraform' ],
        \     'project_root_files': [ '*.tf', '*.tfvars' ],
        \   },
        \ ]
endif
```

### LanguageClient-neovim

- Install the [LanguageClient-neovim plugin](https://github.com/autozimu/LanguageClient-neovim)
- Add the following to your `.vimrc`:

```vim
let g:LanguageClient_serverCommands = {
    \ 'terraform': ['tofu-ls', 'serve'],
    \ }
```

### Neovim v0.5.0+

- Install the [nvim-lspconfig plugin](https://github.com/neovim/nvim-lspconfig)
- Add the following to your `.vimrc` or `init.vim`:

```vim
lua <<EOF
  require'lspconfig'.terraformls.setup{}
EOF
autocmd BufWritePre *.tfvars lua vim.lsp.buf.formatting_sync()
autocmd BufWritePre *.tf lua vim.lsp.buf.formatting_sync()
```

- If you are using `init.lua`:

```lua
require'lspconfig'.terraformls.setup{}
vim.api.nvim_create_autocmd({"BufWritePre"}, {
  pattern = {"*.tf", "*.tfvars"},
  callback = vim.lsp.buf.formatting_sync(),
})
```

### Neovim v0.8.0+

- Install the [nvim-lspconfig plugin](https://github.com/neovim/nvim-lspconfig)
- Add the following to your `.vimrc` or `init.vim`:

```vim
lua <<EOF
  require'lspconfig'.terraformls.setup{}
EOF
autocmd BufWritePre *.tfvars lua vim.lsp.buf.format()
autocmd BufWritePre *.tf lua vim.lsp.buf.format()
```

- If you are using `init.lua`:

```lua
require'lspconfig'.terraformls.setup{}
vim.api.nvim_create_autocmd({"BufWritePre"}, {
  pattern = {"*.tf", "*.tfvars"},
  callback = function()
    vim.lsp.buf.format()
  end,
})
```

Make sure to read through to [server_configurations.md#terraformls](https://github.com/neovim/nvim-lspconfig/blob/master/doc/server_configurations.md#terraformls) if you need more detailed settings.

## VS Code

- Install [OpenTofu VS Code Extension](https://marketplace.visualstudio.com/items?itemName=opentofu.vscode-opentofu) `>=2.24.0`
- Latest compatible version of the language server is bundled with the extension
- See [Configuration](https://github.com/hashicorp/vscode-terraform/blob/main/README.md#configuration) in case you need to tweak anything. Default settings should work for majority of users though.

## Zed

- Install the [OpenTofu Extension](https://zed.dev/extensions?query=OpenTofu) or add the following lines to your zed settings

  ```json
  {
    "auto_install_extensions": {
      "opentofu": true
    }
  }
  ```

- Latest compatible version of the language server will be installed with this extension, if the binary is not already installed
- Configuration see extension [GitHub repo](https://github.com/ashpool37/zed-extension-opentofu?tab=readme-ov-file#file-association-conflicts)

## BBEdit

_BBEdit 14 [added support](https://www.barebones.com/support/bbedit/lsp-notes.html) for the Language Server Protocol so you'll need to upgrade to version 14 to use; this won't work for older versions of BBEdit_.

- Open Preferences > Languages
- In _Language-specific settings_ section, add an entry for OpenTofu
- In the Server tab, Set _Command_ to `tofu-ls` and _Arguments_ to `serve`
- Once you've correctly installed `tofu-ls` and configured BBEdit, the status indicator on this settings panel will flip to green
- If you'd like to pass any [settings](./SETTINGS.md) to the server you can do so via the _Arguments_ field.

## Kate

KDE [Kate editor](https://kate-editor.org/) supports LSP and is user configurable.

- Install the `tofu-ls` package (or the equivalent package name appropriate to your distro)
- Open Kate configuration (Settings Menu -> `Configure` Kate or Kate -> `Preferences` on macOS)
- Select _LSP Client_ in the left pane
- Select _User Server Settings_ tab
- Paste the following JSON and _Save_:

```json
{
  "servers": {
    "opentofu": {
      "command": ["tofu-ls", "serve"],
      "url": "https://github.com/opentofu/tofu-ls",
      "highlightingModeRegex": "^Terraform$",
      "rootIndicationFileNames": ["*.tf", "*.tfvars"]
    }
  }
}
```

- Restart of the editor should _not_ be necessary.

## Other text editors

> [!WARNING]
> Be careful when installing on extensions outside the OpenTofu organization, always read the source code to make sure what you're installing is safe.

There are two ways of finding extensions implementing `tofu-ls` for other text editors:

1. There's a topic on Github called `tofu-ls`. You can find it at [here](https://github.com/topics/tofu-ls). The expectation is if you create an extension, you're going to add a topic on your project to be easily discoverable by other people on Github.
1. There's a curated list at https://awesome-opentofu.com/#helpers.
