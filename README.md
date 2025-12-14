# SIMPUS - Sistem Informasi Manajemen Perpustakaan

Sistem manajemen perpustakaan modern dengan Go, MySQL, HTMX, dan JWT authentication.

## Fitur

### Fungsional Utama
- ✅ **Manajemen Data Buku** - CRUD buku dengan judul, kategori, stok, penulis
- ✅ **Manajemen Anggota** - Mahasiswa, guru, karyawan
- ✅ **Peminjaman & Pengembalian** - Tracking lengkap dengan perhitungan denda
- ✅ **Notifikasi Keterlambatan** - Alert untuk buku terlambat dikembalikan
- ✅ **Riwayat & Laporan** - History transaksi per periode

### Fitur Teknis
- JWT Authentication untuk admin dan anggota
- HTMX untuk interaksi tanpa reload halaman
- Responsive design modern
- Search dan filter dengan pagination

## Tech Stack

- **Backend**: Go dengan Chi Router
- **Database**: MySQL
- **Frontend**: Go Templates + HTMX
- **Authentication**: JWT (JSON Web Token)
- **Styling**: Vanilla CSS dengan design system modern

## Instalasi

### Prasyarat
- Go 1.21+
- MySQL 8.0+

### Setup Database

1. Buat database MySQL:
```sql
CREATE DATABASE simpus;
```

2. Import schema:
```bash
mysql -u root -p simpus < database/migrations/001_init.sql
```

### Konfigurasi

1. Edit file `.env` sesuai konfigurasi MySQL Anda:
```env
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your_password
DB_NAME=simpus
```

### Menjalankan Aplikasi

```bash
# Build
go build -o simpus.exe ./cmd/main.go

# Run
./simpus.exe
```

Atau langsung:
```bash
go run ./cmd/main.go
```

Aplikasi akan berjalan di `http://localhost:8080`

## Default Login

### Admin
- **Username**: admin
- **Password**: admin123

### Sample Member
- **Email**: budi@student.ac.id
- **Password**: admin123

## Struktur Project

```
SIMPUS/
├── cmd/
│   └── main.go              # Entry point
├── config/
│   └── config.go            # Configuration
├── database/
│   ├── connection.go        # DB connection
│   └── migrations/          # SQL schema
├── internal/
│   ├── handlers/            # HTTP handlers
│   ├── middleware/          # Auth middleware
│   ├── models/              # Data models
│   ├── repository/          # Database operations
│   └── services/            # Business logic
├── static/
│   ├── css/style.css        # Styling
│   └── js/htmx.min.js       # HTMX library
├── templates/
│   ├── layouts/             # Base templates
│   ├── components/          # Reusable components
│   ├── auth/                # Login pages
│   ├── admin/               # Admin pages
│   └── member/              # Member pages
├── .env                     # Environment config
├── go.mod
└── go.sum
```

## API Routes

### Public
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/login` | Admin login page |
| POST | `/login` | Process admin login |
| GET | `/login/member` | Member login page |
| POST | `/login/member` | Process member login |
| GET | `/logout` | Logout |

### Admin (Protected)
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/admin/dashboard` | Dashboard |
| GET/POST | `/admin/books` | Manage books |
| GET/POST | `/admin/categories` | Manage categories |
| GET/POST | `/admin/authors` | Manage authors |
| GET/POST | `/admin/members` | Manage members |
| GET/POST | `/admin/borrowings` | Manage borrowings |
| POST | `/admin/borrowings/{id}/return` | Return book |
| GET | `/admin/reports` | Reports |

### Member (Protected)
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/member/dashboard` | Member dashboard |

## Perhitungan Denda

- Denda keterlambatan: **Rp 1.000 per hari**
- Denda otomatis dihitung saat pengembalian

## License

MIT License
