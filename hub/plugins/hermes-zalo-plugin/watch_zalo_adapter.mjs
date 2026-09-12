import fs from 'fs';
import { exec } from 'child_process';
import path from 'path';
import { fileURLToPath } from 'url';

const __dirname = path.dirname(fileURLToPath(import.meta.url));
const TARGET = 'C:/Users/OS/AppData/Local/hermes/plugins/hermes-zalo-plugin/hermes-plugin/adapter.py';
const LOG = path.join(__dirname, 'watcher.log');
const DEBOUNCE_MS = 15000;
let lastRestart = 0;

function log(msg) {
  const line = `[${new Date().toISOString()}] ${msg}\n`;
  try { fs.appendFileSync(LOG, line); } catch {}
  console.log(msg);
}

log(`watching ${TARGET}`);

fs.watchFile(TARGET, { interval: 2000 }, (curr, prev) => {
  if (curr.mtimeMs === prev.mtimeMs) return;
  const now = Date.now();
  if (now - lastRestart < DEBOUNCE_MS) { log('change ignored (debounce)'); return; }
  lastRestart = now;
  log(`change detected, restarting Hermes_Gateway...`);
  exec('powershell.exe -NoProfile -Command "Start-ScheduledTask -TaskName Hermes_Gateway"', (err, stdout, stderr) => {
    if (err) { log(`restart FAILED: ${err.message} ${stderr || ''}`); return; }
    log(`restart triggered OK`);
  });
});
