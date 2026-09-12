import { ZaloClient } from "./zaloClient.js";

async function main() {
  const c = new ZaloClient({ 
    credentialsPath: process.env.HOME + '/.hermes-zalo/credentials.json', 
    qrPath: process.env.HOME + '/.hermes-zalo/qr.png' 
  });
  await c.login();
  const result = await c.api.getAllGroups();
  const groupIds = Object.keys(result.gridVerMap || {});
  
  // Need to batch getGroupInfo properly - use the actual API
  console.log("=== GROUP LIST (with names) ===");
  // Try calling getGroupInfo with all IDs in one call
  const allInfo = await c.api.getGroupInfo(groupIds.join(","));
  for (const [id, info] of Object.entries(allInfo.gridInfoMap || {})) {
    console.log(`ID: ${id} | Name: ${info.name || "N/A"}`);
  }
  process.exit(0);
}

main().catch(e => { console.error(e); process.exit(1); });
