module.exports = {
  apps: [
    {
      name: 'template',
      script: 'node_modules/next/dist/bin/next',
      args: 'start -p 3000',
      cwd: __dirname,
      exec_mode: 'fork',
      instances: 1,
      watch: false,
      max_memory_restart: '512M',
      error_file: 'pm2-logs/err.log',
      out_file: 'pm2-logs/out.log',
      kill_timeout: 5000,
      listen_timeout: 8000,
    },
  ],
}
