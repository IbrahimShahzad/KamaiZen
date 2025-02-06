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
    - [x] Show documentation

> Note: This is a work in progress, and not all features are available yet.

## Installation

with [lazy.nvim](https://github.com/folke/lazy.nvim):
```lua
    {
      'batoaqaa/KamaiZen',
      dependencies = { 'batoaqaa/KamaiZen', build = 'go build' },
      opts = {
        settings = {
          kamaizen = {
            enableDeprecatedCommentHint = false, -- to enable hints for '#' comments
            KamailioSourcePath = vim.fn.getcwd(),
            loglevel = 3,
          },
        },
      },
    }
```

## Integration

### Neovim

- [x] [kamaizen.nvim](https://github.com/IbrahimShahzad/kamaizen.nvim)

### Vscode

> Not yet available

- [ ] [vscode-kamaizen](github.com/IbrahimShahzad/vscode-kamaizen)


