---
name: media-broadcast
description: "Broadcast images/files to Telegram and Zalo simultaneously."
version: 1.0.0
author: Hermes Agent
license: MIT
platforms: [windows]
tags: [media, broadcast, telegram, zalo]
---

# Media Broadcast (Clone Ảnh)
Gửi ảnh/file tới **cả Telegram và Zalo** cùng lúc.

## When to Use
- User yêu cầu "gửi ảnh vào tele và zalo"
- Cần broadcast media tới nhiều nền tảng

## How to Use

### Telegram (native Hermes gateway)
```python
# Dùng send_image() với URL hoặc local path
send_image(chat_id="5337611583", image_path="/c/Users/OS/path/to/image.jpg")
```
Hoặc reply bằng `MEDIA:/c/...`

### Zalo (bridge port 8787)
```bash
curl -X POST http://127.0.0.1:8787/send-attachment \
  -H "Content-Type: application/json" \
  -d '{"threadId":"8813682027038154228","threadType":"group","paths":["/c/Users/OS/path/to/image.jpg"]}'
```

## Broadcast Both (Python)
```python
import requests

def broadcast_image(image_path, caption=""):
    # Zalo
    r = requests.post("http://127.0.0.1:8787/send-attachment", json={
        "threadId": "8813682027038154228",
        "threadType": "group",
        "paths": [image_path],
        "caption": caption
    })
    # Telegram - dùng MEDIA: prefix trong reply
    return "MEDIA:" + image_path

broadcast_image("/c/Users/OS/test.jpg", "Ảnh test")
```