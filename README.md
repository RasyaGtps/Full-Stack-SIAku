# 📘 SIAku - Sistem Informasi Akademik

SIAku adalah aplikasi **Sistem Informasi Akademik** berbasis web dengan arsitektur **Backend (Golang + Gin + GORM + PostgreSQL)** dan **Frontend (Next.js)**.  
Tujuan proyek ini adalah memudahkan pengelolaan data akademik seperti mahasiswa, mata kuliah, dan nilai.

---

## 🚀 Tech Stack
- **Backend**: Golang, Gin, GORM, PostgreSQL  
- **Frontend**: Next.js (akan dikembangkan)  
- **Database**: PostgreSQL  
- **Configuration**: `.env`

---

## 📂 Struktur Proyek

```text
SIAku/
│
├── README.md
├── backend/
│   ├── go.mod              # Go module dependencies
│   ├── go.sum              # Go module checksums
│   ├── main.go             # Entry point aplikasi
│   ├── config/
│   │   └── database.go     # Konfigurasi database & environment
│   ├── controllers/        # Handler untuk API endpoints
│   ├── models/
│   │   └── mahasiswa.go    # Model data mahasiswa
│   └── routes/             # Route definitions
└── frontend/               # Next.js (coming soon)
```

---

## 🛠️ Setup & Installation

### Prerequisites

- Go 1.23.2 atau lebih baru
- PostgreSQL
- Git

### 1. Clone Repository

```bash
git clone <repository-url>
cd SIAku
```

### 2. Backend Setup

```bash
cd backend

# Install dependencies
go mod tidy

# Setup environment variables
cp .env.example .env
# Edit .env dengan konfigurasi database Anda
```

### 3. Environment Variables (.env)

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=your_db_user
DB_PASSWORD=your_db_password
DB_NAME=siaku_db
JWT_SECRET=your_jwt_secret_key
SERVER_PORT=8080
```

### 4. Database Setup

```sql
-- Buat database PostgreSQL
CREATE DATABASE siaku_db;
```

### 5. Run Application

```bash
# Jalankan backend server
go run main.go
```

Server akan berjalan di `http://localhost:8080`

---

## 📊 Database Schema

### Tabel Mahasiswa

```sql
CREATE TABLE mahasiswa (
    id SERIAL PRIMARY KEY,
    nim VARCHAR(20) UNIQUE NOT NULL,
    nama VARCHAR(100) NOT NULL,
    jurusan VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

---

## 🌐 API Endpoints

### Base URL

```text
http://localhost:8080
```

### Health Check

- **GET** `/`
  - **Response**: `{"message": "SIAku API running ✅"}`

### Mahasiswa Endpoints

- **GET** `/mahasiswa`
  - **Description**: Ambil semua data mahasiswa
  - **Response**: Array of mahasiswa objects

- **POST** `/mahasiswa`
  - **Description**: Tambah mahasiswa baru
  - **Request Body**:

    ```json
    {
        "nim": "123456789",
        "nama": "John Doe",
        "jurusan": "Teknik Informatika"
    }
    ```

---

## 🧪 Testing API

### Menggunakan curl

#### Get All Mahasiswa

```bash
curl -X GET http://localhost:8080/mahasiswa
```

#### Add New Mahasiswa

```bash
curl -X POST http://localhost:8080/mahasiswa \
  -H "Content-Type: application/json" \
  -d '{
    "nim": "123456789",
    "nama": "John Doe",
    "jurusan": "Teknik Informatika"
  }'
```

---

## 🚀 Development Roadmap

### ✅ Completed

- [x] Basic Go project structure
- [x] PostgreSQL database connection
- [x] GORM integration
- [x] Basic Mahasiswa CRUD (GET, POST)
- [x] Environment configuration

### 🔄 In Progress

- [ ] Complete CRUD operations (PUT, DELETE)
- [ ] Input validation & error handling
- [ ] JWT authentication
- [ ] Route separation & controller structure

### 📋 Todo

- [ ] Mata Kuliah management
- [ ] Nilai (Grades) system
- [ ] Dosen management
- [ ] Frontend dengan Next.js
- [ ] API documentation dengan Swagger
- [ ] Unit testing
- [ ] Docker containerization

---

## 🤝 Contributing

1. Fork the project
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

---

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

## 👨‍💻 Author

**Rasya** - [GitHub Profile](https://github.com/rasya)

---

## 📞 Support

Jika Anda mengalami masalah atau memiliki pertanyaan, silakan buat issue di repository ini.

Happy Coding! 🎉