local state = { autocmd = {} }

local capabilities = vim.lsp.protocol.make_client_capabilities()
capabilities.textDocument.foldingRange = { dynamicRegistration = false, lineFoldingOnly = true }

local server = {
  name = 'KamaiZen',
  cmd = { vim.fn.stdpath 'data' .. '/lazy/KamaiZen/KamaiZen' },

  cmd_cwd = vim.fn.getcwd(),
  filetypes = { 'kamailio' },
  root_dir = vim.fn.getcwd(),
  capabilities = capabilities,
  autostart = true,
  settings = {
    kamaizen = {
      enableDeprecatedCommentHint = true,
      KamailioSourcePath = vim.fn.getcwd(),
      enableDiagnostics = true,
      loglevel = 3,
    },
  },
  on_init = function(client, results)
    if results.offsetEncoding then
      client.offset_encoding = results.offsetEncoding
    end
  end,
  on_exit = function(_, _, client_id)
    vim.schedule(function()
      vim.api.nvim_del_autocmd(state.autocmd[client_id])
    end)
  end,
  on_attach = function(client, bufnr)
    if client then
      print('Attaching to: ' .. client.name .. ' attached to buffer ' .. bufnr)
      ------------------------------------------------------------------
      local bufkeymap = function(mode, keys, func, desc)
        vim.keymap.set(mode, keys, func, { buffer = bufnr, noremap = true, silent = true, desc = 'LSP: ' .. desc })
      end
      -- Quickfix list
      bufkeymap('n', '[q', vim.cmd.cprev, 'Previous quickfix item')
      bufkeymap('n', ']q', vim.cmd.cnext, 'Next quickfix item')
      -- Diagnostic keymaps
      bufkeymap('n', '[d', vim.diagnostic.goto_prev, 'Go to previous [D]iagnostic message')
      bufkeymap('n', ']d', vim.diagnostic.goto_next, 'Go to next [D]iagnostic message')
      bufkeymap('n', '<leader>e', vim.diagnostic.open_float, 'Show diagnostic [E]rror messages')
      bufkeymap('n', '<leader>q', vim.diagnostic.setloclist, 'Open diagnostic [Q]uickfix list')
      --
      if client.server_capabilities.hoverProvider then
        bufkeymap('n', 'K', vim.lsp.buf.hover, 'Hover Documentation')
      end
      if client.server_capabilities.signatureHelpProvider then
        bufkeymap({ 'i', 'n' }, '<C-k>', vim.lsp.buf.signature_help, 'Show signature')
      end
      if client.server_capabilities.declarationProvider then
        bufkeymap('n', 'gD', vim.lsp.buf.declaration, '[G]oto [D]eclaration')
      end
      if client.server_capabilities.definitionProvider then
        -- bufkeymap('n', 'gd', vim.lsp.buf.definition, 'Go to definition')
        bufkeymap('n', 'gd', require('telescope.builtin').lsp_definitions, '[G]oto [D]efinition')
      end
      if client.server_capabilities.typeDefinitionProvider then
        bufkeymap('n', 'gt', vim.lsp.buf.type_definition, '[G]oto [T]ype definition')
      end
      if client.server_capabilities.implementationProvider then
        bufkeymap('n', 'gI', vim.lsp.buf.implementation, '[G]oto [I]mplementation')
      end
      if client.server_capabilities.referencesProvider then
        -- bufkeymap('n', 'gr', vim.lsp.buf.references, 'List references')
        bufkeymap('n', 'gr', require('telescope.builtin').lsp_references, '[G]oto [R]eferences')
      end

      if client.server_capabilities.documentSymbolProvider then
        -- bufkeymap('n', '<leader>ds', vim.lsp.buf.document_symbol, '[D]ocument [S]ymbols' )
        bufkeymap('n', '<leader>ds', require('telescope.builtin').lsp_document_symbols, '[D]ocument [S]ymbols')
      end
      --
      if client.server_capabilities.documentFormattingProvider then
        bufkeymap({ 'n', 'x' }, '<F3>', function()
          require('conform').format({ bufnr = bufnr, async = true })
        end, 'format buffer')
      end
      --
      if client.server_capabilities.inlayHintProvider and vim.lsp.inlay_hint then
        bufkeymap('n', '<leader>lh', function()
          vim.lsp.inlay_hint.enable(not vim.lsp.inlay_hint.is_enabled({ bufnr = bufnr }), { bufnr = bufnr })
        end, '[l]sp [h]ints toggle')
      end
      --
      if client.server_capabilities.documentHighlightProvider then
        local highlight_augroup = vim.api.nvim_create_augroup('kickstart-lsp-highlight', { clear = false })
        vim.api.nvim_create_autocmd({ 'CursorHold', 'CursorHoldI' }, {
          buffer = bufnr,
          group = highlight_augroup,
          callback = vim.lsp.buf.document_highlight,
        })
        --
        vim.api.nvim_create_autocmd({ 'CursorMoved', 'CursorMovedI' }, {
          buffer = bufnr,
          group = highlight_augroup,
          callback = vim.lsp.buf.clear_references,
        })
        --
      end
    end
  
    -- --
    vim.lsp.handlers['textDocument/hover'] = vim.lsp.with(vim.lsp.handlers.hover, {
      border = 'rounded',
    })
    --
    vim.lsp.handlers['textDocument/signatureHelp'] = vim.lsp.with(vim.lsp.handlers.signature_help, {
      border = 'rounded',
    })
    --
    vim.lsp.handlers['textDocument/publishDiagnostics'] = vim.lsp.with(vim.lsp.diagnostic.on_publish_diagnostics, {
      signs = true,
      underline = true,
      virtual_text = {
        spacing = 5,
        min = vim.diagnostic.severity.HINT,
      },
      update_in_insert = true,
    })
    -- --
    -- vim.cmd([[autocmd FileType * set formatoptions-=ro]])
  end,

}

return {
  server_opts = server,
}
