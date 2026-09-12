import { ZaloClient } from "./zaloClient.js";

async function main() {
  const c = new ZaloClient({ 
    credentialsPath: process.env.HOME + '/.hermes-zalo/credentials.json', 
    qrPath: process.env.HOME + '/.hermes-zalo/qr.png' 
  });
  await c.login();
  
  const groupId = "8813682027038154228";
  const info = await c.api.getGroupInfo(groupId);
  const gridInfo = info.gridInfoMap[groupId];
  console.log("Group name:", gridInfo.name);
  console.log("creatorId:", gridInfo.creatorId);
  
  // Try getGroupMembersInfo
  try {
    const members = await c.api.getGroupMembersInfo(groupId);
    console.log("Members info:", JSON.stringify(members, null, 2).substring(0, 1500));
  } catch(e) {
    console.error("getGroupMembersInfo error:", e.message);
  }
  
  // Also check currentMems / memberIds
  if (gridInfo.memberIds) {
    console.log("memberIds count:", gridInfo.memberIds.length);
    console.log("First 20:", gridInfo.memberIds.slice(0, 20));
  }
  process.exit(0);
}

main().catch(e => { console.error(e); process.exit(1); });
