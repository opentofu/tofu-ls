# Usage of OpenTofu Language Server

This guide assumes you have installed the server by following instructions
in the [installation.md](./installation.md) if the client doesn't download the server itself.
And make sure `tofu-ls` command is accessible.

> [!NOTE]
> If you are unable to configure `tofu-ls` command to be directly accessible by the editor.
> You can enter the absolute path of the `tofu-ls` binary instead.

The following filetypes are supported by the OpenTofu Language Server:

- `opentofu` - standard `*.tf` and `*.tofu` config files
- `opentofu-vars` - variable files (`*.tfvars`)

We also accept `terraform` and `terraform-vars` as language IDs, to support wider range of editors.
For consistent behavior we encourage users to remap them to corresponding opentofu IDs.

> [!NOTE]
> Clients should be configured to follow the above language ID conventions
> and do **not** send `*.tf.json`, `*.tfvars.json` nor Packer HCL config
> nor any other HCL config files as the server is not equipped to handle these file types.

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

### Eglot

```elisp
;; if you have doom emacs:
(set-eglot-client! '(terraform-mode :language-id "opentofu") '("tofu-ls" "serve"))

;; or without it, after loading `eglot`:
(add-to-list 'eglot-server-programs '((terraform-mode :language-id "opentofu") . ("tofu-ls" "serve")))
```

## IntelliJ IDE

_We do not have an officially supported way to use `tofu-ls` with IntelliJ IDEs. If you must use `tofu-ls`, you can try to find different ways to add generic language servers._

We recommend using the official [Terraform and HCL plugin.](https://plugins.jetbrains.com/plugin/7808-terraform-and-hcl)
It provides overlapping functionality with this language server and in some way has more features.

## Vim / NeoVim

### coc.nvim

- Install the [coc.nvim plugin](https://github.com/neoclide/coc.nvim)
- Add the following snippet to the `coc-setting.json` file (editable via `:CocConfig` in NeoVim)

```json
{
  "languageserver": {
    "opentofu": {
      "command": "tofu-ls",
      "args": ["serve"],
      "filetypes": ["terraform", "tf", "tofu"],
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
        \ 'whitelist': ['terraform'],
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
        \     'name': 'opentofu',
        \     'cmdline': [ 'tofu-ls', 'serve' ],
        \     'filetypes': [ 'terraform' ],
        \     'project_root_files': [ '.terraform', .git' ]
        \   },
        \ ]
endif
```

### Neovim v0.11.0+

We can natively configure LSP in Neovim (0.11.0+).
The following is the lua configuration for `tofu-ls`:

```lua
-- tofu-ls lsp setup
vim.lsp.config['tofu_ls'] = {
  cmd = { 'tofu-ls', 'serve' },
  -- Base filetypes
  filetypes = { 'terraform', 'terraform-vars' },
  root_markers = {'.terraform', '.git'},
}

vim.lsp.enable('tofu_ls')
```

If you want to enable auto-formatting on save, add the following configuration

```lua
vim.api.nvim_create_autocmd('LspAttach', {
  callback = function(args)
    local client = assert(vim.lsp.get_client_by_id(args.data.client_id))
    -- Auto-format on save
    if client:supports_method('textDocument/formatting') then
      vim.api.nvim_create_autocmd('BufWritePre', {
        group = vim.api.nvim_create_augroup('tofu-ls', {clear=false}),
        buffer = args.buf,
        callback = function()
          vim.lsp.buf.format({ bufnr = args.buf, id = client.id, timeout_ms = 1000 })
        end,
      })
    end
  end,
})

```

In case you are using '.tofu' files, you also need to add the filetype during the initialization.
At the time of writing this, Neovim doesn't support opentofu file extensions by default.

```lua
-- Add OpenTofu filetype
vim.filetype.add({
  extension = {
    tofu = 'opentofu'
  },
})
```

Make sure to read through [Neovim LSP documentation](https://neovim.io/doc/user/lsp.html) if you need more detailed settings.

## VS Code

- Install [OpenTofu VS Code Extension](https://marketplace.visualstudio.com/items?itemName=opentofu.vscode-opentofu)
- Latest compatible version of the language server is bundled with the extension
- See [Configuration](https://github.com/opentofu/vscode-opentofu/blob/main/README.md#configuration) in case you need to tweak anything. Default settings should work for majority of users though.

## Zed

- Install the [OpenTofu Extension](https://zed.dev/extensions?query=OpenTofu) or add the following lines to your zed settings

  ```json
  {
    "auto_install_extensions": {
      "opentofu": true
    }
  }
  ```

- Latest compatible version of the language server will be installed with this extension, if the binary is not already installed.
- For configuration options, see the corresponding GitHub repository of the extension you installed.

## BBEdit (Might require update)

_BBEdit 14 [added support](https://www.barebones.com/support/bbedit/lsp-notes.html) for the Language Server Protocol so you'll need to upgrade to version 14 to use; this won't work for older versions of BBEdit_.

- Open Preferences > Languages
- In _Language-specific settings_ section, add an entry for OpenTofu
- In the Server tab, Set _Command_ to `tofu-ls` and _Arguments_ to `serve`
- Once you've correctly installed `tofu-ls` and configured BBEdit, the status indicator on this settings panel will flip to green
- If you'd like to pass any [settings](./SETTINGS.md) to the server you can do so via the _Arguments_ field.

## Kate

KDE [Kate editor](https://kate-editor.org/) supports LSP and is user configurable.

- Open Kate configuration (`Settings` -> `Configure Kate` or Kate -> `Preferences` on macOS)
- Select _LSP Client_ in the left pane
- Select _User Server Settings_ tab
- Paste the following JSON and _Save_:

```json
{
  "servers": {
    "opentofu": {
      "command": ["/path/to/tofu-ls", "serve"],
      "url": "https://github.com/opentofu/tofu-ls",
      "highlightingModeRegex": "^(OpenTofu|OpenTofu-Vars|Terraform)$",
      "rootIndicationFileNames": [".terraform", ".git"]
    }
  }
}
```

- Restart of the editor should _not_ be necessary.

At the time of writing this guide, Kate along with most other editors do not have a separate language mode for OpenTofu.
Hence, this configuration will work on .tf and .tfvars files, in case you are using .tofu files, you will need to add `Sources/OpenTofu` as a new filetype with all appropriate extensions.
New filetypes can be configured from `Settings` -> `Configure Kate` > `Open/Save` > `Modes & Filetypes`.

## Helix Editor

Add the following config to your defined `languages.toml`:

```toml
[language-server.tofu-ls]
command = "tofu-ls"
args = ["serve"]

[[language]]
name = "hcl"
language-id = "opentofu"
scope = "source.hcl"
file-types = ["tf", "tofu", "tfvars"]
auto-format = true
comment-token = "#"
block-comment-tokens = { start = "/*", end = "*/" }
indent = { tab-width = 2, unit = "  " }
language-servers = [ "tofu-ls" ]
```

Then, you need to rebuild your grammars with the following two commands:

- hx -g fetch
- hx -g build

Check the health of the language with:

- hx --health hcl

## Other text editors

> [!WARNING]
> Be careful when installing on extensions outside the OpenTofu organization, always read the source code to make sure what you're installing is safe.

There are two ways of finding extensions implementing `tofu-ls` for other text editors:

1. There's a topic on Github called `tofu-ls`. You can find it at [here](https://github.com/topics/tofu-ls). The expectation is if you create an extension, you're going to add a topic on your project to be easily discoverable by other people on Github.
1. There's a curated list at https://awesome-opentofu.com/#helpers.
