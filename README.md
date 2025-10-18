# 🎓 SIAku - Sistem Informasi Akademik

<div align="center">

![SIAku Banner](https://img.shields.io/badge/SIAku-Academic%20Information%20System-blue?style=for-the-badge)
![Status](https://img.shields.io/badge/Status-Active%20Development-green?style=for-the-badge)
![License](https://img.shields.io/badge/License-MIT-yellow?style=for-the-badge)

**Sistem Informasi Akademik Modern dengan Integrasi WhatsApp Bot**

[Overview](#-overview) • [Tech Stack](#%EF%B8%8F-tech-stack) • [Setup](#-setup) • [Contributing](#-contributing) • [License](#-license) • [Support](#-support)

</div>

---

## 📋 Overview

SIAku adalah sistem informasi akademik berbasis REST API yang terintegrasi dengan WhatsApp Bot, memungkinkan mahasiswa, dosen, kajur, dan rektor untuk mengelola aktivitas akademik secara efisien melalui web maupun WhatsApp.

### ✨ Key Features

- 🔐 **Multi-Role Authentication** - Mahasiswa, Dosen, Kajur, Rektor
- 📱 **WhatsApp Integration** - Bot untuk notifikasi dan query akademik
- 📊 **KRS Management** - Sistem persetujuan KRS dengan approval flow
- 🎯 **Penilaian Lengkap** - Tugas, UTS, UAS dengan kalkulasi IPK otomatis
- ✅ **Absensi Digital** - Tracking kehadiran per pertemuan
- 📚 **Materi Kuliah** - Upload dan manage materi per pertemuan
- 📅 **Jadwal Kuliah** - Management jadwal perkuliahan
- 🔔 **Security Alerts** - Notifikasi keamanan untuk login tidak sah
- 📈 **Dashboard Analytics** - Dashboard untuk Kajur dan Rektor

---

## 🛠️ Tech Stack

### **Backend**
![Go](https://img.shields.io/badge/Go-1.24-00ADD8?style=flat&logo=go&logoColor=white)
![Gin](https://img.shields.io/badge/Gin-Web%20Framework-00ADD8?style=flat&logo=go)
![GORM](https://img.shields.io/badge/GORM-ORM-00ADD8?style=flat)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-Database-336791?style=flat&logo=postgresql&logoColor=white)
![JWT](https://img.shields.io/badge/JWT-Authentication-000000?style=flat&logo=jsonwebtokens)

```
📦 Backend (Go)
├── 🌐 Gin Framework - Web server & routing
├── 🗄️ GORM - ORM untuk database
├── 🐘 PostgreSQL - Database utama
├── 🔐 JWT - Token-based authentication
├── 🔒 Bcrypt - Password hashing
└── ✅ Validator - Input validation
```

### **WhatsApp Bot**
![Node.js](https://img.shields.io/badge/Node.js-20.x-339933?style=flat&logo=node.js&logoColor=white)
![Express](https://img.shields.io/badge/Express-Server-000000?style=flat&logo=express)
![WhatsApp](https://img.shields.io/badge/WhatsApp-Web.js-25D366?style=flat&logo=whatsapp&logoColor=white)

```
📦 WhatsApp Service (Node.js)
├── 💬 whatsapp-web.js - WhatsApp Web automation
├── 🚀 Express.js - REST API server
├── 🖼️ Jimp - Image processing
├── 📡 Axios - HTTP client untuk backend
└── 🔔 Real-time notifications
```

### **Architecture**

```
┌─────────────────────────────────────────────────────────────┐
│                        CLIENT LAYER                          │
├──────────────────────┬──────────────────────────────────────┤
│   Web/Mobile App     │        WhatsApp User                  │
│   (REST API Client)  │   (Chat dengan Bot)                   │
└──────────────────────┴──────────────────────────────────────┘
            │                           │
            ▼                           ▼
┌─────────────────────┐      ┌──────────────────────┐
│   Backend (Go)      │◄────►│  WhatsApp Bot (JS)   │
│   Port: 8080        │      │  Port: 3000          │
├─────────────────────┤      ├──────────────────────┤
│ • REST API          │      │ • Message Handler    │
│ • JWT Auth          │      │ • Command Parser     │
│ • Business Logic    │      │ • Notification       │
│ • Data Validation   │      │ • Security Alert     │
└─────────────────────┘      └──────────────────────┘
            │
            ▼
┌─────────────────────────────────────┐
│      PostgreSQL Database            │
├─────────────────────────────────────┤
│ • Users, Mahasiswa, Dosen           │
│ • Courses, KRS, Nilai               │
│ • Jadwal, Absensi, Materi           │
│ • Kajur, Rektor                     │
└─────────────────────────────────────┘
```

---

## 🚀 Setup

### Prerequisites

Pastikan sudah terinstall:
- 🐹 **Go** 1.24 atau lebih baru ([Download](https://golang.org/dl/))
- 🟢 **Node.js** 20.x atau lebih baru ([Download](https://nodejs.org/))
- 🐘 **PostgreSQL** 14+ ([Download](https://www.postgresql.org/download/))
- 📱 **WhatsApp** account untuk bot

### 1️⃣ Clone Repository

```bash
git clone https://github.com/RasyaGtps/Full-Stack-SIAku.git
cd Full-Stack-SIAku
```

### 2️⃣ Setup Backend (Go)

#### Install Dependencies
```bash
cd backend
go mod download
```

#### Configure Environment
Buat file `.env` di folder `backend/`:

```env
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=siaku_db

# JWT Configuration
JWT_SECRET=your_super_secret_jwt_key_change_this_in_production

# Server Configuration
SERVER_PORT=8080
GIN_MODE=release

# WhatsApp Service URL
WHATSAPP_SERVICE_URL=http://localhost:3000
```

#### Run Backend
```bash
go run main.go
```

Backend akan berjalan di `http://localhost:8080`

**Output yang diharapkan:**
```
✅ Connected to DB
🚀 Server running on :8080
📋 Environment: release
🔍 Checking WhatsApp Bot Service...
✅ WhatsApp Bot Service: ONLINE & CONNECTED
```

### 3️⃣ Setup WhatsApp Bot

#### Install Dependencies
```bash
cd whatsapp
npm install
```

#### Configure Environment
Buat file `.env` di folder `whatsapp/`:

```env
PORT=3000
BACKEND_URL=http://localhost:8080
NODE_ENV=production
```

#### Run WhatsApp Bot
```bash
npm start
```

#### Scan QR Code
Saat pertama kali dijalankan, QR code akan muncul di terminal:
1. Buka WhatsApp di HP
2. Pilih **WhatsApp Web**
3. Scan QR code dari terminal
4. Bot akan auto-connect

**Output yang diharapkan:**
```
🚀 WhatsApp Service running on port 3000
📱 Initializing WhatsApp connection...
╔═══════════════════════════════════════════════════════╗
║            📱 SCAN QR CODE DENGAN WHATSAPP           ║
╚═══════════════════════════════════════════════════════╝
✅ WhatsApp Client Ready!
```

Session akan tersimpan, tidak perlu scan QR lagi di run berikutnya.

---

## 🤝 Contributing

Kami sangat welcome untuk kontribusi! Berikut cara berkontribusi:

### 1️⃣ Fork & Clone
```bash
# Fork repository ini di GitHub
# Lalu clone fork kamu
git clone https://github.com/YOUR_USERNAME/Full-Stack-SIAku.git
cd Full-Stack-SIAku
```

### 2️⃣ Create Branch
```bash
git checkout -b feature/amazing-feature
```

### 3️⃣ Make Changes
- Ikuti coding style yang ada
- Tambahkan comments untuk kode yang kompleks
- Test perubahan kamu secara menyeluruh

### 4️⃣ Commit Changes
```bash
git add .
git commit -m "feat: add amazing feature"
```

**Commit Message Convention:**
- `feat:` - Fitur baru
- `fix:` - Bug fix
- `docs:` - Update dokumentasi
- `style:` - Format code (tidak mengubah logic)
- `refactor:` - Refactor code
- `test:` - Tambah/update tests
- `chore:` - Update dependencies, dll

### 5️⃣ Push & Pull Request
```bash
git push origin feature/amazing-feature
```

Lalu buat Pull Request di GitHub dengan deskripsi yang jelas.

### 📋 Contribution Guidelines

- ✅ Pastikan kode tidak ada error
- ✅ Test fitur baru sebelum PR
- ✅ Update README jika perlu
- ✅ Ikuti struktur project yang ada
- ✅ Gunakan bahasa Indonesia untuk comments
- ✅ Jangan commit `.env` atau credentials

---

## 📄 License

This project is licensed under the **MIT License** - see the [LICENSE](LICENSE) file for details.

---

## 💬 Support

Butuh bantuan? Ada beberapa cara untuk mendapatkan support:

### 📧 Email Support
Kirim email ke: **rasyarayhandev@gmail.com**

### 🐛 Report Issues
Temukan bug? [Create an Issue](https://github.com/RasyaGtps/Full-Stack-SIAku/issues)

### 💡 Feature Request
Punya ide fitur baru? [Submit Feature Request](https://github.com/RasyaGtps/Full-Stack-SIAku/issues/new?labels=enhancement)

### 📖 Documentation
Dokumentasi lengkap: [Wiki](https://github.com/RasyaGtps/Full-Stack-SIAku/wiki)

### ❓ FAQ

**Q: Apakah gratis?**  
A: Ya, SIAku adalah open source dan gratis untuk digunakan.

**Q: Bisa deploy ke production?**  
A: Ya, tapi pastikan ganti semua secret keys dan konfigurasi keamanan.

**Q: Apakah aman?**  
A: Ya, menggunakan JWT authentication, bcrypt password hashing, dan input validation.

**Q: Bisa custom untuk universitas saya?**  
A: Tentu! Fork repository ini dan sesuaikan dengan kebutuhan.

**Q: WhatsApp bot butuh WhatsApp Business?**  
A: Tidak, cukup WhatsApp biasa.

---

## 🙏 Acknowledgments

Terima kasih kepada:
- 🎯 **Go Community** - Untuk Gin & GORM framework yang luar biasa
- 💬 **WhatsApp-Web.js** - Library WhatsApp automation yang powerful
- 🐘 **PostgreSQL Team** - Database yang reliable dan scalable
- 🌟 **Open Source Community** - Untuk semua library yang digunakan

---

## 📊 Project Stats

![GitHub stars](https://img.shields.io/github/stars/RasyaGtps/Full-Stack-SIAku?style=social)
![GitHub forks](https://img.shields.io/github/forks/RasyaGtps/Full-Stack-SIAku?style=social)
![GitHub issues](https://img.shields.io/github/issues/RasyaGtps/Full-Stack-SIAku)
![GitHub pull requests](https://img.shields.io/github/issues-pr/RasyaGtps/Full-Stack-SIAku)
![GitHub last commit](https://img.shields.io/github/last-commit/RasyaGtps/Full-Stack-SIAku)

---

## 👨‍💻 Author

<div align="center">

**Developed with ❤️ by**

### Rasya

[![GitHub](https://img.shields.io/badge/GitHub-Follow-181717?style=for-the-badge&logo=github)](https://github.com/RasyaGtps)
[![Email](https://img.shields.io/badge/Email-Contact-D14836?style=for-the-badge&logo=gmail&logoColor=white)](mailto:rasyarayhandev@gmail.com)

---

### ⭐ Star Project Ini Jika Bermanfaat!

**Made with ❤️ in Indonesia** 🇮🇩

[⬆ Back to Top](#-siaku---sistem-informasi-akademik)

</div>
