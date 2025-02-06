local state = { autocmd = {} }

local capabilities = vim.lsp.protocol.make_client_capabilities()
capabilities.textDocument.foldingRange = { dynamicRegistration = false, lineFoldingOnly = true }

local server = {
  name = 'KamaiZen',
  cmd = { vim.fn.stdpath 'data' .. '/lazy/KamaiZen/KamaiZen' },

  cmd_cwd = vim.fn.getcwd(),
  filetypes = { 'cfg', 'inc', 'kamailio' },
  root_dir = vim.fn.getcwd(),
  capabilities = capabilities,
  autostart = true,
  settings = {
    kamaizen = {
      enableDeprecatedCommentHint = false, -- to enable hints for '#' comments
      KamailioSourcePath = vim.fn.getcwd(),
      loglevel = 3,
    },
  },
  on_init = function(client, results)
    if results.offsetEncoding then
      client.offset_encoding = results.offsetEncoding
    end

    -- if client.config.settings then
    --   client.notify('workspace/didChangeConfiguration', {
    --     settings = client.config.settings,
    --   })
    -- end
  end,
  --on_exit = function(code, signal, client_id)
  on_exit = function(_, _, client_id)
    vim.schedule(function()
      vim.api.nvim_del_autocmd(state.autocmd[client_id])
    end)
  end,
}

return {
  server_opts = server,
}
