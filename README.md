# KamaiZen
A language server for kamailio configuration files

For syntax highlighting, use [tree-sitter-kamailio-cfg](https://github.com/IbrahimShahzad/tree-sitter-kamailio-cfg)


## Features
- [ ] Code completion
    - [ ] Variables
      - [x] Global variables (avps)
      - [ ] Local variables (vars)
    - [x] exported functions
    - [x] Modules
    - [x] Keywords
    - [ ] Parameters
- [ ] Code navigation
  - [ ] Go to definition for routes - In progress
  - [ ] Find references for routes - In progress
- [ ] Code Actions
  - [ ] Add missing modules
- [ ] Snippets
  - [ ] Route snippets
  - [ ] Module snippets
  - [ ] Ifblock snippets
  - [ ] loop snippets
  - [ ] switch snippets
- [ ] Code formatting
- [ ] Code folding
- [ ] Diagnostics
    - [x] Syntax Errors -- Buggy (requires re-work on the parser)
    - [x] Invalid statements
    - [x] Unreachable code
    - [x] Assignment Errors
    - [ ] Function calls from non-loaded modules
    - [ ] Unused variables
    - [ ] Unused modules
    - [ ] Unused parameters
- [ ] Hover
    - [x] Show documentation for functions
    - [ ] Show documentation for routes
    - [ ] Show documentation for avp vars
    - [ ] Show documentation for local vars

> [!Caution]
> This is a work in progress, and not all features are available yet and might contain bugs.

## Installation

### From source

```bash
git clone https://github.com/IbrahimShahzad/KamaiZen.git 
cd KamaiZen
go build
cp KamaiZen /usr/local/bin
```

### To use with Neovim

update your `init.lua` with the following

```lua
local client = vim.lsp.start_client {
  cmd = { '/usr/local/bin/KamaiZen' }, -- Path to KamaiZen executable
  name = 'KamaiZen',
  settings = {
    kamaizen = {
      logLevel = 1,
      kamailioSourcePath = '/path/to/kamailio-source', -- Path to kamailio source
      enableDeprecatedCommentHint = false, -- to enable hints for '#' comments
    },
  },
}

if not client then
  vim.notify('Failed to start LSP client', vim.log.levels.ERROR)
  return
end

vim.api.nvim_create_autocmd('FileType', {
  pattern = 'kamailio',
  callback = function()
    vim.lsp.buf_attach_client(0, client)
  end,
})
```

## Integration

### Neovim

- [x] [kamaizen.nvim](https://github.com/IbrahimShahzad/kamaizen.nvim)

### Vscode

> Not yet available

- [ ] [vscode-kamaizen](github.com/IbrahimShahzad/vscode-kamaizen)


