const process = require('process');
const path = require('path');
const concurrently = require('concurrently');

const { result } = concurrently(
  [
    {
      command: 'npm exec start-hexlet-url-shortener-frontend',
      name: 'frontend',
    },
    {
      command: 'go run cmd/api/main.go',
      name: 'backend',
      cwd: path.resolve(__dirname, '..'),
    },
  ],
  {
    killOthersOn: ['failure', 'success'],
    cwd: path.resolve(__dirname),
  },
);

result
  .then(() => process.exit(0))
  .catch(() => process.exit(1));
