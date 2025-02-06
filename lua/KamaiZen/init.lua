local config = require 'KamaiZen.config'
local M = {}

M.setup = function(opts)
  -- opts = opts or {}
  opts = vim.tbl_deep_extend('force', config.server_opts or {}, opts)

  vim.filetype.add {
    extension = {
      cfg = 'kamailio',
    },
    filename = {
      ['kamctlrc'] = 'kamailio',
    },
    pattern = {
      ['.*'] = {
        --set files that start with '#!KAMAILIO' as kamailio file type
        --function(path, bufnr)
        function(_, bufnr)
          local content = vim.api.nvim_buf_get_lines(bufnr, 0, 1, false)[1] or ''
          if vim.regex([[^#!.*KAMAILIO]]):match_str(content) ~= nil then
            return 'kamailio'
          end
        end,
      },
    },
  }

  vim.api.nvim_create_autocmd('FileType', {
    pattern = 'kamailio',
    --callback = function(args)
    callback = function()
      ----------------------------------------
      local parser_config = require('nvim-treesitter.parsers').get_parser_configs()
      local parsers = require 'nvim-treesitter.parsers'

      -- if not parser_config["kamailio"] then
      if not parsers.has_parser 'kamailio' then
        parser_config['kamailio'] = {
          install_info = {
            -- url = 'https://github.com/IbrahimShahzad/tree-sitter-kamailio-cfg',
            url = 'https://github.com/batoaqaa/tree-sitter-kamailio', --'/home/batoaqaa/tree-sitter-kamailio', --
            files = { 'src/parser.c' }, -- note that some parsers also require src/scanner.c or src/scanner.cc
            -- optional entries:
            -- branch = 'v0.1.2', -- default branch in case of git repo if different from master
            branch = 'main', -- default branch in case of git repo if different from master
            generate_requires_npm = false, -- if stand-alone parser without npm dependencies
            requires_generate_from_grammar = false, -- if folder contains pre-generated src/parser.c
          },
          filetype = 'kamailio', -- if filetype does not match the parser name
        }

        vim.schedule_wrap(function()
          vim.cmd 'TSInstall kamailio'
        end)()
      end
      ----------------------------------------
      local id = vim.lsp.start(opts)
      if not id then
        vim.notify('Failed to start LSP client', vim.log.levels.ERROR)
        return
      end
    end,
  })
end

return M
