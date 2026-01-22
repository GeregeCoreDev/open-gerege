# API Гэрээ (Монгол)

## Ерөнхий мэдээлэл

**Үндсэн URL:** `https://api.yourdomain.com`  
**Хувилбар:** v1  
**Протокол:** HTTPS  
**Контентийн төрөл:** `application/json`  
**Аутентификаци:** Session-based (Cookie)

---

## Аутентификаци

Хамгаалагдсан endpoint-уудад `sid` cookie шаардлагатай:

```http
Cookie: sid=session-id-here
```

---

## Стандарт Response

### Амжилттай
```json
{
  "code": "OK",
  "message": "",
  "request_id": "uuid",
  "data": {}
}
```

### Алдаатай
```json
{
  "code": "ERROR_CODE",
  "message": "Алдааны тайлбар",
  "request_id": "uuid",
  "details": {}
}
```

---

## API Endpoint-ууд

### 🔐 = Аутентификаци шаардлагатай

## 1. Health & Documentation

| Method | Endpoint | Тайлбар | Auth |
|--------|----------|---------|------|
| GET | `/health` | Серверийн төлөв | ❌ |
| GET | `/docs/*` | Swagger UI | ❌ |

---

## 2. Аутентификаци (`/auth`)

| Method | Endpoint | Тайлбар | Auth |
|--------|----------|---------|------|
| GET | `/auth/login` | SSO login redirect | ❌ |
| GET | `/auth/callback` | OAuth2 callback | ❌ |
| POST | `/auth/logout` | Гарах | ❌ |
| POST | `/auth/google/login` | Google OAuth | ❌ |
| GET | `/auth/verify` | Token шалгах | ❌ |
| POST | `/auth/org/change` | Байгууллага солих | 🔐 |

---

## 3. Хэрэглэгч (`/user`)

| Method | Endpoint | Тайлбар | Auth |
|--------|----------|---------|------|
| GET | `/user/me` | Миний мэдээлэл | 🔐 |
| GET | `/user` | Жагсаалт | 🔐 |
| POST | `/user` | Үүсгэх | 🔐 |
| PUT | `/user/:id` | Засварлах | 🔐 |
| DELETE | `/user/:id` | Устгах | 🔐 |
| POST | `/user/find-from-core` | Core-оос хайх | 🔐 |
| GET | `/user/profile` | Профайл | 🔐 |
| GET | `/user/profile/sso` | SSO профайл | 🔐 |
| GET | `/user/organizations` | Байгууллагууд | 🔐 |

### Жишээ: Хэрэглэгч үүсгэх
```bash
POST /user
Content-Type: application/json

{
  "id": 123,
  "reg_no": "УА12345678",
  "family_name": "Болд",
  "last_name": "Бат",
  "first_name": "Дорж",
  "gender": 1,
  "birth_date": "1990-01-01",
  "phone_no": "99119911",
  "email": "bat@example.com"
}
```

---

## 4. Хэрэглэгч-Эрх (`/user-role`)

| Method | Endpoint | Тайлбар | Auth |
|--------|----------|---------|------|
| GET | `/user-role/users?role_id=1` | Эрхийн хэрэглэгчид | 🔐 |
| GET | `/user-role/roles?user_id=1` | Хэрэглэгчийн эрхүүд | 🔐 |
| POST | `/user-role` | Эрх олгох | 🔐 |
| DELETE | `/user-role` | Эрх хасах | 🔐 |

---

## 5. Систем (`/system`)

| Method | Endpoint | Тайлбар | Auth |
|--------|----------|---------|------|
| GET | `/system` | Жагсаалт | 🔐 |
| GET | `/system/:id` | Дэлгэрэнгүй | 🔐 |
| POST | `/system` | Үүсгэх | 🔐 |
| PUT | `/system/:id` | Засварлах | 🔐 |
| DELETE | `/system/:id` | Устгах | 🔐 |
| GET | `/system/by-role?role_id=1` | Эрхийн системүүд | 🔐 |

---

## 6. Модуль (`/module`)

| Method | Endpoint | Тайлбар | Auth |
|--------|----------|---------|------|
| GET | `/module` | Жагсаалт | 🔐 |
| POST | `/module` | Үүсгэх | 🔐 |
| PUT | `/module/:id` | Засварлах | 🔐 |
| DELETE | `/module/:id` | Устгах | 🔐 |
| GET | `/module/by-role?role_id=1` | Эрхийн модулууд | 🔐 |
| GET | `/module/by-org-admin` | Админы модулууд | 🔐 |

---

## 7. Зөвшөөрөл (`/permission`)

| Method | Endpoint | Тайлбар | Auth |
|--------|----------|---------|------|
| GET | `/permission` | Жагсаалт | 🔐 |
| POST | `/permission` | Үүсгэх | 🔐 |
| PUT | `/permission/:id` | Засварлах | 🔐 |
| DELETE | `/permission/:id` | Устгах | 🔐 |

---

## 8. Эрх (`/role`)

| Method | Endpoint | Тайлбар | Auth |
|--------|----------|---------|------|
| GET | `/role` | Жагсаалт | 🔐 |
| POST | `/role` | Үүсгэх | 🔐 |
| PUT | `/role/:id` | Засварлах | 🔐 |
| DELETE | `/role/:id` | Устгах | 🔐 |
| GET | `/role/permissions?role_id=1` | Эрхийн зөвшөөрлүүд | 🔐 |
| POST | `/role/permissions` | Зөвшөөрөл олгох | 🔐 |

### Жишээ: Эрх үүсгэх
```bash
POST /role
Content-Type: application/json

{
  "system_id": 1,
  "code": "ADMIN",
  "name": "Администратор",
  "description": "Системийн админ",
  "is_active": true
}
```

### Жишээ: Эрхэд зөвшөөрөл олгох
```bash
POST /role/permissions
Content-Type: application/json

{
  "role_id": 5,
  "permission_ids": [1, 2, 3, 5, 8]
}
```

---

## 9. Байгууллага (`/organization`)

| Method | Endpoint | Тайлбар | Auth |
|--------|----------|---------|------|
| GET | `/organization/find?search_text=1234567` | Core-оос хайх | 🔐 |
| GET | `/organization` | Жагсаалт | 🔐 |
| POST | `/organization` | Үүсгэх | 🔐 |
| PUT | `/organization/:id` | Засварлах | 🔐 |
| DELETE | `/organization/:id` | Устгах | 🔐 |
| GET | `/organization/tree?org_id=1` | Модон бүтэц | 🔐 |

### Жишээ: Байгууллага үүсгэх
```bash
POST /organization
Content-Type: application/json

{
  "reg_no": "1234567",
  "name": "ХХК Компани",
  "short_name": "Компани",
  "type_id": 1,
  "phone_no": "75001122",
  "email": "info@company.mn",
  "is_active": true,
  "address_detail": "Баянзүрх дүүрэг"
}
```

---

## 10. Байгууллага-Хэрэглэгч (`/orguser`)

| Method | Endpoint | Тайлбар | Auth |
|--------|----------|---------|------|
| GET | `/orguser` | Жагсаалт | 🔐 |
| GET | `/orguser/users?org_id=1` | Байгууллагын хэрэглэгчид | 🔐 |
| GET | `/orguser/organizations?user_id=1` | Хэрэглэгчийн байгууллагууд | 🔐 |
| POST | `/orguser` | Хэрэглэгч нэмэх | 🔐 |
| DELETE | `/orguser` | Хэрэглэгч хасах | 🔐 |

---

## 11. Байгууллагын төрөл (`/orgtype`)

| Method | Endpoint | Тайлбар | Auth |
|--------|----------|---------|------|
| GET | `/orgtype` | Жагсаалт | 🔐 |
| POST | `/orgtype` | Үүсгэх | 🔐 |
| PUT | `/orgtype/:id` | Засварлах | 🔐 |
| DELETE | `/orgtype/:id` | Устгах | 🔐 |
| GET | `/orgtype/system?type_id=1` | Төрлийн системүүд | 🔐 |
| POST | `/orgtype/system` | Систем нэмэх | 🔐 |

---

## 12. Терминал (`/terminal`)

| Method | Endpoint | Тайлбар | Auth |
|--------|----------|---------|------|
| GET | `/terminal` | Жагсаалт | 🔐 |
| POST | `/terminal` | Үүсгэх | 🔐 |
| PUT | `/terminal/:id` | Засварлах | 🔐 |
| DELETE | `/terminal/:id` | Устгах | 🔐 |

---

## 13. OAuth Client (`/client`)

| Method | Endpoint | Тайлбар | Auth |
|--------|----------|---------|------|
| GET | `/client` | Жагсаалт | 🔐 |
| GET | `/client/scope` | Scope жагсаалт | 🔐 |
| POST | `/client/scope` | Scope үүсгэх | 🔐 |
| DELETE | `/client/scope` | Scope устгах | 🔐 |

---

## 14-17. App Icon Management

| Group | Endpoint | Тайлбар |
|-------|----------|---------|
| Service Icon | `/app-service-icon` | CRUD |
| Service Group | `/app-service-group` | CRUD |
| Desktop Icon | `/app-desktop-icon` | CRUD |
| Business Icon | `/app-business-icon` | CRUD |

Бүх endpoint-үүд `GET`, `POST`, `PUT/:id`, `DELETE/:id` дэмжинэ.

---

## 18. Файл (`/file`)

| Method | Endpoint | Тайлбар | Auth |
|--------|----------|---------|------|
| GET | `/file/list` | Жагсаалт | 🔐 |
| POST | `/file/upload` | Upload | 🔐 |
| DELETE | `/file` | Устгах | 🔐 |
| GET | `/file/:uuid` | Татаж авах | ❌ |

### Жишээ: Файл upload
```bash
POST /file/upload
Content-Type: multipart/form-data
Cookie: sid=session-id

file=@document.pdf
```

**Response:**
```json
{
  "code": "OK",
  "data": {
    "uuid": "file-uuid-123",
    "filename": "document.pdf",
    "size": 1024000,
    "url": "https://api.example.com/file/file-uuid-123"
  }
}
```

---

## 19. Мэдэгдэл (`/notification`)

| Method | Endpoint | Тайлбар | Auth |
|--------|----------|---------|------|
| GET | `/notification` | Жагсаалт | 🔐 |
| GET | `/notification/groups` | Бүлгүүд | 🔐 |
| POST | `/notification` | Илгээх | 🔐 |
| POST | `/notification/read` | Уншсан тэмдэглэх | 🔐 |
| POST | `/notification/read-all` | Бүгдийг уншсан | 🔐 |

### Жишээ: Мэдэгдэл илгээх
```bash
POST /notification
Content-Type: application/json

{
  "title": "Шинэ мэдэгдэл",
  "message": "Танд шинэ мэдэгдэл ирлээ",
  "user_ids": [1, 2, 3, 5],
  "type": "info"
}
```

---

## 20. Мэдээ (`/news`)

| Method | Endpoint | Тайлбар | Auth |
|--------|----------|---------|------|
| GET | `/news` | Жагсаалт | ❌ |
| GET | `/news/get/:id` | Дэлгэрэнгүй | ❌ |
| POST | `/news` | Үүсгэх | 🔐 |
| PUT | `/news/:id` | Засварлах | 🔐 |
| DELETE | `/news/:id` | Устгах | 🔐 |

---

## 21. Баталгаажуулалт (`/verify`)

| Method | Endpoint | Тайлбар | Auth |
|--------|----------|---------|------|
| GET | `/verify/dan` | DAN баталгаажуулалт | 🔐 |
| POST | `/verify/email` | Email код илгээх | 🔐 |
| POST | `/verify/email/confirm` | Email баталгаажуулах | 🔐 |
| POST | `/verify/phone` | Утасны код илгээх | 🔐 |
| POST | `/verify/phone/confirm` | Утас баталгаажуулах | 🔐 |

### Жишээ: Email баталгаажуулах
```bash
# 1. Код илгээх
POST /verify/email
{
  "email": "user@example.com"
}

# 2. Код баталгаажуулах
POST /verify/email/confirm
{
  "email": "user@example.com",
  "code": "123456"
}
```

---

## 22. Видео хурал (`/room`)

| Method | Endpoint | Тайлбар | Auth |
|--------|----------|---------|------|
| GET | `/room` | Жагсаалт | 🔐 |
| GET | `/room/token?room_id=1` | Token үүсгэх | 🔐 |
| POST | `/room` | Өрөө үүсгэх | 🔐 |
| POST | `/room/join` | Өрөөнд орох | 🔐 |
| POST | `/room/:id/users` | Хэрэглэгч нэмэх | 🔐 |
| DELETE | `/room/:id` | Өрөө устгах | 🔐 |
| DELETE | `/room/:id/users/:user_id` | Хэрэглэгч хасах | 🔐 |

### Жишээ: Өрөө үүсгэх
```bash
POST /room
Content-Type: application/json

{
  "name": "Долоо хоногийн хурал",
  "description": "Ажлын хурал",
  "max_participants": 10,
  "scheduled_at": "2025-12-10T14:00:00Z"
}
```

---

## 23. Терминал төлбөр (`/tpay`)

### Данс

| Method | Endpoint | Тайлбар | Auth |
|--------|----------|---------|------|
| GET | `/tpay/accounts/me` | Миний дансууд | 🔐 |
| PUT | `/tpay/accounts/set-default` | Үндсэн данс | 🔐 |
| GET | `/tpay/accounts/statement` | Дансны хуулга | 🔐 |
| POST | `/tpay/accounts/:id/qr` | QR үүсгэх | 🔐 |

### Гүйлгээ

| Method | Endpoint | Тайлбар | Auth |
|--------|----------|---------|------|
| POST | `/tpay/transaction/qr_pay` | QR төлбөр | 🔐 |
| POST | `/tpay/p2p` | Шилжүүлэг | 🔐 |

### Карт

| Method | Endpoint | Тайлбар | Auth |
|--------|----------|---------|------|
| GET | `/tpay/card/list` | Картын жагсаалт | 🔐 |
| POST | `/tpay/card/create` | Карт нэмэх | 🔐 |
| POST | `/tpay/card/confirm` | Карт баталгаажуулах | 🔐 |
| GET | `/tpay/card/send_otp` | OTP илгээх | 🔐 |
| POST | `/tpay/verify_card` | Карт шалгах | 🔐 |

### Жишээ: QR төлбөр
```bash
# 1. QR үүсгэх
POST /tpay/accounts/acc-123/qr
{
  "amount": 50000,
  "description": "Бараа худалдан авалт"
}

# Response:
{
  "code": "OK",
  "data": {
    "qr_code": "data:image/png;base64,...",
    "qr_string": "qr-text-value",
    "expires_at": "2025-12-08T15:00:00Z"
  }
}

# 2. QR-ээр төлөх
POST /tpay/transaction/qr_pay
{
  "qr_string": "qr-text-value",
  "account_id": "acc-456"
}
```

### Жишээ: Данс хооронд шилжүүлэг
```bash
POST /tpay/p2p
{
  "from_account_id": "acc-123",
  "to_account_id": "acc-456",
  "amount": 10000,
  "description": "Зээл төлөлт"
}
```

---

## 24. Чат (`/chat`)

| Method | Endpoint | Тайлбар | Auth |
|--------|----------|---------|------|
| GET | `/chat` | Жагсаалт | 🔐 |
| POST | `/chat` | Үүсгэх | 🔐 |
| PUT | `/chat/:id` | Засварлах | 🔐 |
| DELETE | `/chat/:id` | Устгах | 🔐 |
| POST | `/chat/key` | Key-ээр хайх | 🔐 |

---

## Pagination параметрүүд

Бүх жагсаалт endpoint-үүд pagination дэмжинэ:

| Параметр | Төрөл | Тайлбар | Default |
|----------|-------|---------|---------|
| `page` | int | Хуудасны дугаар | 1 |
| `size` | int | Нэг хуудсанд харуулах тоо (max: 500) | 20 |
| `q` | string | Хайлтын текст | - |
| `sort` | string | Эрэмбэлэх (`field:asc` эсвэл `field:desc`) | - |
| `created_from` | date | Эхлэх огноо (YYYY-MM-DD) | - |
| `created_to` | date | Дуусах огноо (YYYY-MM-DD) | - |

### Жишээ:
```bash
GET /user?page=2&size=50&q=бат&sort=created_at:desc&created_from=2025-01-01
```

---

## Алдааны кодууд

| Код | HTTP | Тайлбар |
|-----|------|---------|
| `OK` | 200 | Амжилттай |
| `CREATED` | 201 | Үүсгэгдсэн |
| `BAD_REQUEST` | 400 | Буруу хүсэлт |
| `UNAUTHORIZED` | 401 | Нэвтрээгүй |
| `FORBIDDEN` | 403 | Эрх хүрэхгүй |
| `NOT_FOUND` | 404 | Олдсонгүй |
| `CONFLICT` | 409 | Давхардсан |
| `VALIDATION_ERROR` | 422 | Баталгаажуулалтын алдаа |
| `INTERNAL_ERROR` | 500 | Серверийн алдаа |
| `SERVICE_UNAVAILABLE` | 503 | Үйлчилгээ боломжгүй |

---

## Request хугацаа

- **Үндсэн:** 5 секунд
- **Файл upload:** 30 секунд
- **Урт үргэлжлэх үйлдэл:** 60 секунд

---

## Rate Limiting

- **IP-ээр:** 100 хүсэлт/минут
- **Хэрэглэгчээр:** 1000 хүсэлт/минут
- **Файл upload:** 10 хүсэлт/минут

---

## Бүрдэл өгөгдлийн загварууд

### Хэрэглэгч
```json
{
  "id": 123,
  "reg_no": "УА12345678",
  "first_name": "Бат",
  "last_name": "Болд",
  "email": "bat@example.com",
  "phone_no": "99119911",
  "gender": 1,
  "birth_date": "1990-01-01"
}
```

### Байгууллага
```json
{
  "id": 10,
  "reg_no": "1234567",
  "name": "ХХК Компани",
  "short_name": "Компани",
  "type_id": 1,
  "is_active": true
}
```

### Эрх
```json
{
  "id": 5,
  "system_id": 1,
  "code": "ADMIN",
  "name": "Администратор",
  "is_active": true
}
```

---

## Дэмжлэг

- **Техникийн дэмжлэг:** support@example.com
- **API асуудал:** api@example.com

---

**Шинэчлэгдсэн огноо:** 2025-12-08  
**Баримтын хувилбар:** 1.0.0  
**Хариуцсан баг:** Gerege Core Team

