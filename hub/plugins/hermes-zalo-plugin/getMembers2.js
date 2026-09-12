import { ZaloClient } from "./zaloClient.js";

async function main() {
  const c = new ZaloClient({ 
    credentialsPath: process.env.HOME + '/.hermes-zalo/credentials.json', 
    qrPath: process.env.HOME + '/.hermes-zalo/qr.png' 
  });
  await c.login();
  
  const groupId = "8813682027038154228";
  
  // Try getGroupChatHistory to see recent members
  try {
    const hist = await c.api.getGroupChatHistory(groupId, 30);
    console.log("History keys:", Object.keys(hist));
    const msgs = hist.messages || hist.gridMessage || hist.data || [];
    console.log("msgs type:", Array.isArray(msgs) ? `array ${msgs.length}` : typeof msgs);
    if (Array.isArray(msgs)) {
      const seen = new Set();
      for (const m of msgs.slice(0, 30)) {
        const uid = m.uidFrom || m.uid || m.senderId;
        const name = m.dName || m.senderName || m.displayName;
        if (uid && !seen.has(uid)) {
          seen.add(uid);
          console.log("UID:", uid, "| Name:", name);
        }
      }
    } else {
      console.log("Raw:", JSON.stringify(hist).substring(0, 1200));
    }
  } catch(e) {
    console.error("History error:", e.message);
  }
  process.exit(0);
}

main().catch(e => { console.error(e); process.exit(1); });
