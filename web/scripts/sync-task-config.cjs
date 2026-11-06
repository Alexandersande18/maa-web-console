// Keep the backend's standalone page on the same task schema as the Vue app.
const fs = require('node:fs');
const path = require('node:path');
const destination = path.resolve(__dirname, '../../resource/public/html/task-config');
fs.mkdirSync(destination, { recursive: true });
for (const name of ['taskDefinitions.js', 'taskConfigUtils.js']) {
  fs.copyFileSync(path.resolve(__dirname, '../src/components/task-config', name), path.join(destination, name));
}
