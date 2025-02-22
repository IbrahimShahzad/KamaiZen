## Manual

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
      enableDiagnostics = true, -- to enable/disable diagnostics
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

