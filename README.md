# KamaiZen

<table style="width: 100%;">
  <tr>
    <td><img src="docs/logo.png" alt="KamaiZen" width="200"></td>
    <td>
      <h3>KamaiZen</h3>
      A language server for Kamailio configuration files.<br><br>
    </td>
  </tr>
</table>

> [!TIP]
> It uses grammar based on [tree-sitter-kamailio-cfg](https://github.com/IbrahimShahzad/tree-sitter-kamailio-cfg)
> You can use that if you just want to have syntax highlighting.

## Features

### Demo Video

Watch a demo of KamaiZen in action:

https://github.com/user-attachments/assets/99cedc2d-92c2-4e15-b5c8-aa0fc9f49fc6


### Code completion

- [x] Variables
  - [x] AVP 
  - [x] Local variables (vars)
  - [x] Dialog variables
- [x] Core Cookbook items
- [x] exported functions
- [x] Modules
- [x] SIP Keywords
- [x] Parameters

### Diagnostics

- [x] Syntax Errors -- Buggy (requires re-work on the parser, you can disable it by setting `enableDiagnostics` to false)
- [x] Invalid statements
- [x] Unreachable code
- [x] Assignment Errors

### Hover

- [x] Show documentation for functions
- [x] Show documentation for variables
- [x] Core Cookbook items
- [x] Variables

### Code navigation

- [x] Go to definition for routes: This currently works with the routes defined in current file.

---

> [!Note]
> This is a work in progress, and not all features are available yet and might contain bugs.

## Installation

### Neovim

with `lazy.nvim`

> [!Note]
> Make sure to have [go 1.24](https://go.dev/doc/install) installed.

```lua
{
  'IbrahimShahzad/KamaiZen',
  tag = 'v0.1.1', -- or branch = 'master'
  build = 'go build',
  opts = {
    settings = {
      kamaizen = {
        enableDeprecatedCommentHint = true, -- to enable hints for '#' comments
        enableDiagnostics = true, -- to enable/disable diagnostics
        KamailioSourcePath = "/path/to/kamailio", -- or use current dir vim.fn.getcwd()
        loglevel = 3,
      },
    },
  },
}
```


## How To Contrribute

Contributions are welcome! To help improve KamaiZen, please follow these guidelines:


1. Fork the Repository:

Click the "Fork" button on the GitHub repository page to create your own copy.

2. Clone Your Fork:

```sh
git clone https://github.com/<your-username>/KamaiZen.git
cd KamaiZen
```

3. Create a New Branch:
Create a branch for your feature or bugfix:

```sh
git checkout -b feature/your-feature-name
```

4. Make Your Changes:

Ensure your changes follow the project's code style. Write tests if applicable.

5. Commit and Push:

Commit your changes with a descriptive message:

```sh
git add .
git commit -m "Description of your change"
git push origin feature/your-feature-name
```

6. Submit a Pull Request:

Open a pull request against the main branch of the original repository. Please provide a clear description of your changes and reference any relevant issues.

7. Code Reviews and Feedback:

Your pull request will be reviewed. Be prepared to make adjustments based on feedback.

For any questions or suggestions, please open an issue on GitHub.

## Future Plans

These are the features that are planned to be implemented in the future:

- [ ] scratch-parser implementation ?
- [ ] LSP for Workspace Folder instead of open file
- [ ] Code navigation
  - [ ] Find references for routes
- [ ] Code Actions
  - [ ] Add missing modules
  - [ ] string evaluations
  - [ ] regex check
- [ ] Snippets
  - [ ] Route snippets
  - [ ] Module snippets
  - [ ] Ifblock snippets
  - [ ] loop snippets
  - [ ] switch snippets
- [ ] Code formatting
- [ ] Code folding
- [ ] Diagnostics
  - [ ] Function calls from non-loaded modules
  - [ ] Unused variables
  - [ ] Unused modules
  - [ ] Unused parameters
- [ ] Other Editors
  - [x] Neovim
  - [ ] VSCode
  - [ ] CLion

> [!Note]
> These are not in any particular order and might change in the future.

---

Feel free to contribute or open issues if you have suggestions or encounter problems!

---
