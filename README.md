# KamaiZen

<table>
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

[![KamaiZen Demo](https://img.youtube.com/vi/IbnZwrY13IY/hqdefault.jpg)](https://www.youtube.com/watch?v=IbnZwrY13IY)

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
- [x] Syntax Errors -- Buggy (requires re-work on the parser)
- [x] Invalid statements
- [x] Unreachable code
- [x] Assignment Errors
- [x] Function calls from non-loaded modules
- [x] Unused variables
- [x] Unused modules
- [x] Unused parameters

### Hover
- [x] Show documentation for functions
- [x] Show documentation for variables
- [x] Core Cookbook items
- [x] Variables

> [!Caution]
> This is a work in progress, and not all features are available yet and might contain bugs.

## Installation

with `lazy.nvim`

```lua
    {
      'IbrahimShahzad/KamaiZen',
      branch = 'master', -- or tag = 'v0.0.5'
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


## Integration

- [x] Neovim [kamaizen.nvim](https://github.com/IbrahimShahzad/kamaizen.nvim)

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

- [ ] scratch-parser implementation ?
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
    - [ ] Function calls from non-loaded modules
    - [ ] Unused variables
    - [ ] Unused modules
    - [ ] Unused parameters
---

Feel free to contribute or open issues if you have suggestions or encounter problems!

---
