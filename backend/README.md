# SIAku API Documentation

## 🚀 Fitur yang Sudah Diimplementasi

### ✅ **Architecture Improvements**
- **Clean Architecture**: Controllers, Middleware, Routes terpisah
- **JWT Authentication**: Secure token-based authentication
- **Input Validation**: Comprehensive validation dengan error messages
- **CORS Support**: Cross-origin requests support
- **Pagination**: Efficient data pagination
- **Error Handling**: Consistent error responses
- **Password Hashing**: Secure bcrypt password hashing

### ✅ **Security Features**
- JWT token authentication
- Password hashing dengan bcrypt
- Request validation
- CORS middleware
- Rate limiting ready (middleware tersedia)

## 📋 **API Endpoints**

### **Health Check**
```
GET / - API health check
```

### **Authentication** (No auth required)
```
POST /api/auth/register - Register mahasiswa baru
POST /api/auth/login    - Login mahasiswa
```

### **Profile** (Auth required)
```
GET /api/profile - Get user profile
```

### **Mahasiswa** (Auth required)
```
GET    /api/mahasiswa     - Get all mahasiswa (paginated)
GET    /api/mahasiswa/:id - Get mahasiswa by ID
PUT    /api/mahasiswa/:id - Update mahasiswa (own data only)
DELETE /api/mahasiswa/:id - Delete mahasiswa (own data only)
```

### **Courses** (Auth required)
```
GET    /api/courses           - Get all courses (paginated)
GET    /api/courses/:id       - Get course by ID
POST   /api/courses           - Create new course
PUT    /api/courses/:id       - Update course
DELETE   /api/courses/:id       - Delete course
POST   /api/courses/:id/enroll   - Enroll in course
DELETE /api/courses/:id/enroll   - Unenroll from course
```

## 📝 **Request/Response Examples**

### **1. Register**
```bash
POST /api/auth/register
Content-Type: application/json

{
  "nim": "12345678",
  "nama": "John Doe",
  "jurusan": "Teknik Informatika",
  "password": "password123"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Registration successful",
  "data": {
    "id": 1,
    "nim": "12345678",
    "nama": "John Doe",
    "jurusan": "Teknik Informatika",
    "created_at": "2025-01-01T00:00:00Z",
    "updated_at": "2025-01-01T00:00:00Z"
  },
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

### **2. Login**
```bash
POST /api/auth/login
Content-Type: application/json

{
  "nim": "12345678",
  "password": "password123"
}
```

### **3. Get Courses (with pagination)**
```bash
GET /api/courses?page=1&limit=10
Authorization: Bearer <your_jwt_token>
```

**Response:**
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "code": "TI101",
      "name": "Pemrograman Dasar",
      "credits": 3,
      "created_at": "2025-01-01T00:00:00Z",
      "updated_at": "2025-01-01T00:00:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 10,
    "total": 15
  }
}
```

### **4. Create Course**
```bash
POST /api/courses
Authorization: Bearer <your_jwt_token>
Content-Type: application/json

{
  "code": "TI102",
  "name": "Struktur Data",
  "credits": 3
}
```

### **5. Enroll in Course**
```bash
POST /api/courses/1/enroll
Authorization: Bearer <your_jwt_token>
```

## 🔐 **Authentication**

Semua protected endpoints memerlukan JWT token di header:
```
Authorization: Bearer <your_jwt_token>
```

Token didapat dari response login/register dan valid selama 24 jam.

## ⚠️ **Validation Rules**

### **Mahasiswa**
- `nim`: required, min 8 chars, max 20 chars
- `nama`: required, min 2 chars, max 100 chars  
- `jurusan`: required, min 2 chars, max 100 chars
- `password`: required, min 6 chars

### **Course**
- `code`: required, min 3 chars, max 10 chars, unique
- `name`: required, min 3 chars, max 100 chars
- `credits`: required, min 1, max 6

## 🗃️ **Environment Variables**

Create `.env` file (gunakan `.env.example` sebagai template):
```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=siaku_db
JWT_SECRET=your_super_secret_jwt_key_here
SERVER_PORT=8080
GIN_MODE=development
```

## 🚀 **How to Run**

1. **Setup Database**: Pastikan PostgreSQL running
2. **Setup Environment**: Copy `.env.example` ke `.env` dan isi values
3. **Install Dependencies**: `go mod tidy`
4. **Run Server**: `go run main.go`

Server akan running di `http://localhost:8080`

## 📊 **What's Changed from Original**

### **Before (main.go only):**
- ❌ No authentication
- ❌ No validation
- ❌ No middleware
- ❌ No pagination
- ❌ Basic CRUD only
- ❌ No proper error handling

### **After (Clean Architecture):**
- ✅ JWT Authentication
- ✅ Input validation with proper error messages
- ✅ CORS, logging, recovery middleware
- ✅ Pagination support
- ✅ Complete CRUD + enrollment system
- ✅ Consistent error handling
- ✅ Secure password hashing
- ✅ Clean separation of concerns
- ✅ Simple API structure (`/api/` instead of `/api/v1/`)
- ✅ Clean code (no unnecessary comments)

## 🎯 **Next Possible Improvements**
- Rate limiting middleware
- Email verification
- Password reset functionality
- Role-based access control (Admin/Student)
- File upload for profile pictures
- API documentation dengan Swagger
- Unit tests
- Docker containerization