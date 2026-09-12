import { ZaloClient } from "./zaloClient.js";

async function main() {
  const c = new ZaloClient({ 
    credentialsPath: process.env.HOME + '/.hermes-zalo/credentials.json', 
    qrPath: process.env.HOME + '/.hermes-zalo/qr.png' 
  });
  await c.login();

  console.log("=== Listening for message events (30s) ===");
  c.api.listener.on("message", (message) => {
    console.log("=== NEW MESSAGE ===");
    console.log("Thread ID:", message.threadId);
    console.log("Group ID:", message.groupId || "N/A");
    console.log("Sender ID:", message.uidFrom);
    console.log("Content:", message.data || message.content);
    console.log("Timestamp:", new Date(message.cliMsgId || Date.now()).toISOString());
  });

  // Timeout
  setTimeout(() => {
    console.log("\n=== Monitoring ended — no messages received ===");
    process.exit(0);
  }, 30000);
  
  console.log("Send a message in Zalo group 'Hermes' within 30s...");
}

main().catch(e => { console.error(e); process.exit(1); });
