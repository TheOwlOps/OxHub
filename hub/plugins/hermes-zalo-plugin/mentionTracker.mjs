import fs from 'node:fs';
import path from 'node:path';
import { Zalo, ThreadType } from "zca-js";
import { credentialsPath, qrPath } from "./paths.js";

const STORAGE = path.join(path.dirname(credentialsPath()), "mentioned_uids.json");

function loadStorage() {
  try { return JSON.parse(fs.readFileSync(STORAGE, "utf-8")); } catch { return { uids: [], last_updated: "" }; }
}

function saveStorage(data) {
  data.last_updated = new Date().toISOString();
  fs.writeFileSync(STORAGE, JSON.stringify(data, null, 2), "utf-8");
}

async function main() {
  const zalo = new Zalo({ selfListen: true, logging: false });
  const creds = JSON.parse(fs.readFileSync(credentialsPath(), "utf-8"));
  const api = await zalo.login(creds);

  console.log("[tracker] Logged in, listening for mentions...");

  api.listener.on("message", (msg) => {
    const isGroup = msg.type === ThreadType.Group;
    const data = msg.data || {};
    const mentions = Array.isArray(data.mentions) ? data.mentions : [];

    if (isGroup && mentions.length > 0) {
      const storage = loadStorage();
      for (const mn of mentions) {
        const uid = String(mn.uid);
        if (!storage.uids.includes(uid)) {
          storage.uids.push(uid);
        }
      }
      saveStorage(storage);
      console.log("[tracker] Captured mentions:", JSON.stringify(mentions.map(m => m.uid)));
      console.log("[tracker] UIDs in storage:", storage.uids.length);
    }
  });

  api.listener.on("error", (err) => console.error("[tracker] error:", err));
  api.listener.on("close", (code) => console.log("[tracker] listener closed:", code));
  api.listener.on("reconnect", (info) => console.log("[tracker] reconnecting...", info));

  // Keep alive ping
  setInterval(() => {
    api.updateActiveStatus(true).catch(() => {});
  }, 30000);

  console.log("[tracker] Running...");
}

main().catch(e => {
  console.error("[tracker] fatal:", e);
  process.exit(1);
});
