# SIAku WhatsApp Service (Baileys)

WhatsApp service menggunakan JavaScript dan Baileys untuk fitur-fitur yang tidak bisa dilakukan di Go/whatsmeow.

## ✅ Fitur

### Yang BISA via Baileys:
- ✅ **Ganti Nama Profil** - Update nama WhatsApp profil
- ✅ **Ganti Profile Picture** - Update PP bot
- ✅ **Auto-connect** - Session tersimpan, auto-reconnect
- ✅ **Auto-read** - Pesan otomatis dibaca
- ✅ **Owner System** - Verifikasi kode untuk jadi owner
- ✅ **Block/Unblock** - Manage blocked users
- ✅ **Kirim Pesan** - Send message via API
- ✅ **Broadcast** - Kirim ke multiple numbers

## 🚀 Cara Pakai

### 1. Install Dependencies
```bash
cd whatsapp
npm install
```

### 2. Jalankan Service
```bash
npm start
```

### 3. Scan QR Code
Saat pertama kali jalan, akan muncul QR code di terminal. Scan pakai WhatsApp.

### 4. Session Tersimpan
Session tersimpan di folder `auth_info/`. Restart tidak perlu scan QR lagi.

## 📡 API Endpoints

Service berjalan di `http://localhost:3000`

### Status Connection
```
GET /api/wa/status
```

### Get QR Code
```
GET /api/wa/qr-code
```

### Send Message
```
POST /api/wa/send
Body:
{
  "phone_number": "628123456789",
  "message": "Hello from SIAku!"
}
```

### Broadcast
```
POST /api/wa/broadcast
Body:
{
  "phone_numbers": ["628111", "628222"],
  "message": "Broadcast message"
}
```

## 🤖 Bot Commands

### Untuk Semua User:
- `/menu` - Tampilkan menu
- `/help` - Bantuan
- `/jadiowner` - Jadi owner (kode di terminal)
- `/cekowner` - Cek status owner

### Untuk Owner Only:
- `/gantinama [nama]` - Ganti nama profil WhatsApp ✅
- `/gantipp` - Cara ganti PP (kirim foto dengan caption `setpp`) ✅
- `/block [nomor]` - Block user
- `/unblock [nomor]` - Unblock user
- `/listblock` - List blocked users
- `/infobot` - Info bot
- `/keluarowner` - Keluar dari ownership

## 🔗 Integrasi dengan Go Backend

Go backend bisa hit API endpoints di atas untuk:
- Cek status connection
- Kirim notifikasi WhatsApp
- Broadcast pesan

Contoh dari Go:
```go
// Send WhatsApp message
resp, err := http.Post("http://localhost:3000/api/wa/send", 
    "application/json",
    bytes.NewBuffer(jsonData))
```

## 📂 File Structure

```
whatsapp/
├── index.js          # Express server
├── whatsapp.js       # Baileys WhatsApp logic
├── routes.js         # API routes
├── package.json      # Dependencies
├── auth_info/        # Session data (auto-created)
├── data.json         # Owners & blocked users
└── README.md         # This file
```

## 🔐 Security

- Kode verifikasi owner **hanya muncul di terminal** (secure!)
- Session tersimpan lokal (tidak ke-upload)
- Owners & blocked users tersimpan di `data.json`

## 🎯 Keunggulan vs whatsmeow

| Fitur | whatsmeow (Go) | Baileys (JS) |
|-------|----------------|--------------|
| Ganti Nama | ❌ | ✅ |
| Ganti PP | ❌ | ✅ |
| Auto-connect | ✅ | ✅ |
| Kirim Pesan | ✅ | ✅ |
| Performance | ⚡ Cepat | 🐢 Lebih lambat |
| Memory Usage | 💚 Rendah | 🟡 Lebih tinggi |

## 💡 Tips

1. **Pertama kali:** Scan QR code, lalu session tersimpan permanent
2. **Jadi Owner:** Kirim `/jadiowner`, cek terminal untuk kode
3. **Ganti Nama:** `/gantinama SIAku Bot` - langsung berubah!
4. **Ganti PP:** Kirim foto dengan caption `setpp` - langsung berubah!

