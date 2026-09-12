# Dashboard UI & Design System Guidelines

## Layout & Architecture (OpenClaw-Inspired)

OpenClaw dashboard layout:
1. **Layer 1: Hero / Node Status**
   - Agent node identity + status badge (live pulse indicator)
   - Active gateway port + current running model
   - Quick node-level actions (Làm mới / New Session)
   - Live channel status pills (Telegram, Zalo, Discord)
2. **Layer 2: Metric Tiles (4-column grid on desktop, 2-col on mobile)**
   - High-contrast mono numbers for counts
   - Subtext showing context (e.g. `6/6 skills active`)
   - Distinct icon + color per metric: Sky (chat), Amber (skills), Indigo (MCP), Emerald (cron)
3. **Layer 3: Split Workboard (2:1 ratio)**
   - Left (2 cols): Recent activity / conversation history with relative timestamps
   - Right (1 col): Quick action toolbox with direct navigation links

## Color Palette (Linear / Slate Dark - Eye-Comfort Optimized)

Avoid pure black `#000000` or aggressive cyber-red `#ff5c5c` as primary backgrounds — causes eye fatigue.

```css
/* Backgrounds */
--bg-main: #0b0f19 (or tailwind slate-950)
--card-bg: rgba(15, 23, 42, 0.6) (tailwind slate-900/60)
--card-hover: rgba(30, 41, 59, 0.5) (tailwind slate-800/50)
--border-subtle: rgba(51, 65, 85, 0.8) (tailwind slate-800/80)
--border-highlight: rgba(71, 85, 105, 0.6) (tailwind slate-700/60)

/* Typography */
--text-primary: #f1f5f9 (tailwind slate-100)
--text-secondary: #cbd5e1 (tailwind slate-300)
--text-muted: #94a3b8 (tailwind slate-400)
--text-dim: #64748b (tailwind slate-500)

/* Semantic Accents */
--accent-primary: #4f46e5 (tailwind indigo-600)
--accent-success: #10b981 (tailwind emerald-500)
--accent-warning: #f59e0b (tailwind amber-500)
--accent-info: #0ea5e9 (tailwind sky-500)
```

## Anti-Caching for SPAs

When gateway serves static assets via wildcard route:
```ts
c.header('Cache-Control', 'no-cache, no-store, must-revalidate');
c.header('Pragma', 'no-cache');
c.header('Expires', '0');
```
Ensures browser loads new bundle immediately after `pnpm --filter @ox/web build`.

## Chat & Canvas Scroll Container Contract (Double Scroll Bug Prevention)

When building nested Flexbox chat lists or split canvases:
- **Rule**: Flex items default to `min-height: auto`, which prevents children with `overflow-y-auto` from shrinking or triggering scrollbars when contents expand.
- **Outer Container (`Layout.tsx`)**: MUST use `overflow-hidden`, NEVER `overflow-y-auto`. Having an outer scroll container steals wheel events from the chat list.
- **Flex Parent Chain**: Every flex container wrapping the scrollable list MUST include `flex-1 min-h-0 overflow-hidden` (or `h-full min-h-0`).
- **Scrollable Target (`ChatWindow.tsx`)**: Use `className="flex-1 min-h-0 overflow-y-auto overscroll-contain"` and explicitly set `scrollbarWidth: 'thin'` with `::-webkit-scrollbar` styling so the scrollbar thumb is visible and draggable.
- **Selection**: Avoid blanket `select-none` on parent containers; ensure chat bubbles retain `select-text` for copying.

## Go Native Gateway Agent Extension Patterns

### 1. Live Task Abort / Cancel
- Store `context.CancelFunc` per active `conversationId` in a thread-safe `sync.Map`.
- Endpoint `POST /api/chat/stop` looks up the ID, triggers `cancel()`, and deletes the entry.
- All downstream calls (`Run(ctx, ...)`, `exec_shell`, HTTP requests) must respect `ctx.Done()`.

### 2. Vision / Multimodal Payload Formatting
- Define `ChatMessage.Content` as `interface{}` rather than `string`.
- When an image (base64 data URL) is attached, format `Content` as OpenAI-compatible array:
  `[{"type": "text", "text": prompt}, {"type": "image_url", "image_url": {"url": dataUrl}}]`.
- Implement `GetContentString() string` helper to extract text safely when decoding assistant responses.

### 3. Multi-Agent Delegation (`delegate_task`)
- Expose `delegate_task` function tool to the ReAct agent accepting `tasks: [{goal, context}]`.
- Spawn each subtask into a dedicated worker goroutine with `sync.WaitGroup` and buffered result channel.
- Aggregate all subagent outputs into a structured summary for the parent ReAct loop.

### 4. Deep Research & Web Grounding Pipeline
- **Root-cause of "naive" research**: Agents calling only `web_search` (snippet-only, 50 chars), capped by artificial "stop after 1-2 searches" rules, lacking real-time anchors (searching old years).
- **Tool Pairing**:
  - `web_search`: URL discovery via DuckDuckGo / SearXNG / Tavily.
  - `web_extract`: URL inspection fetching deep content, stripping boilerplate (scripts, nav, ads) into markdown (capped at 15k chars).
- **Prompt Architecture**:
  - Inject live anchor timestamp: `time.Now().Format("Monday, January 02, 2006, 15:04:05 MST")` into ReAct system prompt every turn.
  - Require multi-hop investigation: search -> inspect links -> cross-verify >= 2 sources -> cite direct URLs and exact data. Remove early-termination directives.

### 5. Animated Owl Expressive Thinking Indicator (Long Task Engagement)
- **Problem**: Long-running autonomous workflows (10-30s+ for deep research, compiling, multi-agent subtasks) make users think the UI is frozen if only a generic static spinner or 3 pulsing dots are shown.
- **Dynamic Time-Stage States**: Break execution time into humorous, progressive personality stages:
  - `0 - 3s (Curious)`: Owl tilts head 15° side-to-side, wide eyes blinking: *"Cú đang lắng nghe & ngẫm nghĩ..."*
  - `4 - 9s (Searching)`: Head rotates 180°, pupils dart scanning: *"Đang xoay đầu 180° lục lọi tri thức..."*
  - `10 - 17s (Overclock)`: Owl vibrates rapidly (`animate-vibrate`), sweat drop 💦, wings typing fast: *"Não cú đang ép xung 120%!"*
  - `18 - 29s (Tea Sip)`: Eyes half-closed, sipping hot tea/coffee ☕, blowing steam: *"Task này khoai ghê... Hớp ngụm cà phê tính kế"*
  - `30s+ (Snooze & Panic)`: Head nods asleep with `Zzz`, suddenly startles wide-eyed: *"Khò... Ơ giật mình! Đợi tí sub-agent sắp về tới rồi! 🚀"*
- **Implementation**: Pure SVG element manipulation (swapping `<circle>` and `<ellipse>` pupil coordinates) + pure CSS keyframes (`owlTilt`, `owlHeadTurn`, `owlVibrate`, `owlNod`). Zero external animation libraries or heavy GIFs.



