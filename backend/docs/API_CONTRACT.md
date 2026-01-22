# API Contract / API Гэрээ

## Ерөнхий мэдээлэл / Overview

**Base URL:** `https://api.yourdomain.com`  
**API Version:** v1  
**Protocol:** HTTPS  
**Content-Type:** `application/json`  
**Authentication:** Session-based (Cookie)

---

## Аутентификаци / Authentication

### Cookie-based Session
Хамгаалагдсан endpoint-уудад хандахдаа `sid` (session ID) cookie ашиглана.

```http
Cookie: sid=your-session-id
```

### Session үүсгэх
`/auth/login` эсвэл `/auth/callback` дуудаж session авна.

---

## Стандарт Response Format

### Амжилттай Response
```json
{
  "code": "OK",
  "message": "",
  "request_id": "uuid-v4",
  "data": {}
}
```

### Алдааны Response
```json
{
  "code": "ERROR_CODE",
  "message": "Алдааны тайлбар",
  "request_id": "uuid-v4",
  "details": {}
}
```

### Paginated Response
```json
{
  "code": "OK",
  "data": {
    "meta": {
      "total": 100,
      "page": 1,
      "size": 20,
      "pages": 5,
      "has_next": true,
      "has_prev": false,
      "start_idx": 0,
      "end_idx": 19
    },
    "links": {
      "self": "/user?page=1&size=20",
      "next": "/user?page=2&size=20",
      "prev": null
    },
    "items": []
  }
}
```

---

## API Endpoints

### 1. Health & Documentation

#### GET /health
**Тайлбар:** Server-ийн төлөв шалгах (database холболт)  
**Auth:** Шаардлагагүй  
**Response:**
```json
{
  "code": "OK",
  "data": {
    "status": "ok"
  }
}
```

#### GET /docs/*
**Тайлбар:** Swagger UI documentation  
**Auth:** Шаардлагагүй

---

### 2. Authentication Routes (`/auth`)

#### GET /auth/login
**Тайлбар:** SSO login хуудас руу redirect хийнэ  
**Auth:** Шаардлагагүй  
**Response:** 302 Redirect to SSO

#### GET /auth/callback
**Тайлбар:** OAuth2 callback (SSO-оос буцаж ирэх)  
**Auth:** Шаардлагагүй  
**Query Parameters:**
- `code` (required): Authorization code from SSO

**Response:** 302 Redirect + Set-Cookie

#### POST /auth/logout
**Тайлбар:** Session устгаж, logout хийх  
**Auth:** Шаардлагагүй  
**Response:**
```json
{
  "code": "OK",
  "data": {
    "redirect_url": "https://sso.example.com/logout"
  }
}
```

#### POST /auth/google/login
**Тайлбар:** Google OAuth login  
**Auth:** Шаардлагагүй  
**Request Body:**
```json
{
  "token": "google-id-token"
}
```

#### GET /auth/verify
**Тайлбар:** Session validation шалгах  
**Auth:** Шаардлагагүй  
**Response:**
```json
{
  "code": "OK",
  "data": {
    "valid": true,
    "user_id": 123
  }
}
```

#### POST /auth/org/change
**Тайлбар:** Байгууллага солих  
**Auth:** ✅ Required  
**Request Body:**
```json
{
  "org_id": 5
}
```

---

### 3. User Management (`/user`)

#### GET /user/me
**Тайлбар:** Одоогийн хэрэглэгчийн мэдээлэл  
**Auth:** ✅ Required  
**Response:**
```json
{
  "code": "OK",
  "data": {
    "id": 123,
    "reg_no": "УА12345678",
    "first_name": "Бат",
    "last_name": "Болд",
    "email": "bat@example.com",
    "phone_no": "99119911"
  }
}
```

#### GET /user
**Тайлбар:** Хэрэглэгчдийн жагсаалт (paginated)  
**Auth:** ✅ Required  
**Query Parameters:**
- `page` (optional, default=1): Хуудасны дугаар
- `size` (optional, default=20, max=500): Нэг хуудсанд харуулах тоо
- `q` (optional): Хайх текст
- `sort` (optional): Эрэмбэлэх талбар (жишээ: `created_at:desc`)
- `created_from` (optional): YYYY-MM-DD
- `created_to` (optional): YYYY-MM-DD

**Response:** Paginated list

#### POST /user
**Тайлбар:** Хэрэглэгч үүсгэх  
**Auth:** ✅ Required  
**Request Body:**
```json
{
  "id": 123,
  "civil_id": 456,
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

#### PUT /user/:id
**Тайлбар:** Хэрэглэгч засварлах  
**Auth:** ✅ Required  
**URL Parameters:**
- `id` (required): User ID

**Request Body:** Same as POST /user

#### DELETE /user/:id
**Тайлбар:** Хэрэглэгч устгах  
**Auth:** ✅ Required  
**URL Parameters:**
- `id` (required): User ID

#### POST /user/find-from-core
**Тайлбар:** Core системээс хэрэглэгч хайх  
**Auth:** ✅ Required  
**Request Body:**
```json
{
  "search_text": "УА12345678"
}
```

#### GET /user/profile
**Тайлбар:** Хэрэглэгчийн profile + байгууллагууд  
**Auth:** ✅ Required  
**Response:**
```json
{
  "code": "OK",
  "data": {
    "login_account_info": {
      "user_id": 123,
      "login": "bat@example.com",
      "name": "Бат Болд",
      "email": "bat@example.com",
      "image_url": "https://..."
    },
    "citizen_info": {...},
    "verifications": {
      "citizen_id": 456,
      "is_reg_no_verified": true,
      "is_dan_verified": true,
      "is_phone_verified": true,
      "is_email_verified": true
    }
  }
}
```

#### GET /user/profile/sso
**Тайлбар:** SSO-оос profile авах  
**Auth:** ✅ Required

#### GET /user/organizations
**Тайлбар:** Хэрэглэгчийн байгууллагын жагсаалт  
**Auth:** ✅ Required

---

### 4. User-Role Management (`/user-role`)

#### GET /user-role/users
**Тайлбар:** Тодорхой эрх бүхий хэрэглэгчдийн жагсаалт  
**Auth:** ✅ Required  
**Query Parameters:**
- `role_id` (required): Role ID

#### GET /user-role/roles
**Тайлбар:** Хэрэглэгчийн эрхүүдийн жагсаалт  
**Auth:** ✅ Required  
**Query Parameters:**
- `user_id` (required): User ID

#### POST /user-role
**Тайлбар:** Хэрэглэгчид эрх олгох  
**Auth:** ✅ Required  
**Request Body:**
```json
{
  "user_id": 123,
  "role_id": 5
}
```

#### DELETE /user-role
**Тайлбар:** Хэрэглэгчээс эрх хасах  
**Auth:** ✅ Required  
**Request Body:**
```json
{
  "user_id": 123,
  "role_id": 5
}
```

---

### 5. System Management (`/system`)

#### GET /system
**Тайлбар:** Системийн жагсаалт  
**Auth:** ✅ Required  
**Query Parameters:** Standard pagination

#### GET /system/:id
**Тайлбар:** Системийн дэлгэрэнгүй  
**Auth:** ✅ Required  
**URL Parameters:**
- `id` (required): System ID

#### POST /system
**Тайлбар:** Систем үүсгэх  
**Auth:** ✅ Required  
**Request Body:**
```json
{
  "code": "CORE",
  "name": "Core System",
  "description": "Main system",
  "is_active": true
}
```

#### PUT /system/:id
**Тайлбар:** Систем засварлах  
**Auth:** ✅ Required

#### DELETE /system/:id
**Тайлбар:** Систем устгах  
**Auth:** ✅ Required

#### GET /system/by-role
**Тайлбар:** Эрхээр системийн жагсаалт  
**Auth:** ✅ Required  
**Query Parameters:**
- `role_id` (required): Role ID

---

### 6. Module Management (`/module`)

#### GET /module
**Тайлбар:** Модулийн жагсаалт (menu)  
**Auth:** ✅ Required

#### POST /module
**Тайлбар:** Модуль үүсгэх  
**Auth:** ✅ Required

#### PUT /module/:id
**Тайлбар:** Модуль засварлах  
**Auth:** ✅ Required

#### DELETE /module/:id
**Тайлбар:** Модуль устгах  
**Auth:** ✅ Required

#### GET /module/by-role
**Тайлбар:** Эрхээр модулийн жагсаалт  
**Auth:** ✅ Required  
**Query Parameters:**
- `role_id` (required): Role ID

#### GET /module/by-org-admin
**Тайлбар:** Байгууллагын админы модулууд  
**Auth:** ✅ Required

---

### 7. Permission Management (`/permission`)

#### GET /permission
**Тайлбар:** Зөвшөөрлийн жагсаалт  
**Auth:** ✅ Required

#### POST /permission
**Тайлбар:** Зөвшөөрөл үүсгэх  
**Auth:** ✅ Required

#### PUT /permission/:id
**Тайлбар:** Зөвшөөрөл засварлах  
**Auth:** ✅ Required

#### DELETE /permission/:id
**Тайлбар:** Зөвшөөрөл устгах  
**Auth:** ✅ Required

---

### 8. Role Management (`/role`)

#### GET /role
**Тайлбар:** Эрхийн жагсаалт  
**Auth:** ✅ Required  
**Query Parameters:**
- `system_id` (optional): System ID filter
- `is_active` (optional): Boolean filter
- Standard pagination params

#### POST /role
**Тайлбар:** Эрх үүсгэх  
**Auth:** ✅ Required  
**Request Body:**
```json
{
  "system_id": 1,
  "code": "ADMIN",
  "name": "Администратор",
  "description": "Системийн админ",
  "is_active": true
}
```

#### PUT /role/:id
**Тайлбар:** Эрх засварлах  
**Auth:** ✅ Required

#### DELETE /role/:id
**Тайлбар:** Эрх устгах  
**Auth:** ✅ Required

#### GET /role/permissions
**Тайлбар:** Эрхийн зөвшөөрлүүд  
**Auth:** ✅ Required  
**Query Parameters:**
- `role_id` (required): Role ID

#### POST /role/permissions
**Тайлбар:** Эрхэд зөвшөөрөл олгох  
**Auth:** ✅ Required  
**Request Body:**
```json
{
  "role_id": 5,
  "permission_ids": [1, 2, 3, 5]
}
```

---

### 9. Organization Management (`/organization`)

#### GET /organization/find
**Тайлбар:** Core-оос байгууллага хайх  
**Auth:** ✅ Required  
**Query Parameters:**
- `search_text` (required): Регистрийн дугаар (7 орон)

#### GET /organization
**Тайлбар:** Байгууллагын жагсаалт  
**Auth:** ✅ Required

#### POST /organization
**Тайлбар:** Байгууллага үүсгэх  
**Auth:** ✅ Required  
**Request Body:**
```json
{
  "reg_no": "1234567",
  "name": "ХХК Компани",
  "short_name": "Компани",
  "type_id": 1,
  "phone_no": "75001122",
  "email": "info@company.mn",
  "longitude": 106.9177,
  "latitude": 47.9187,
  "is_active": true,
  "aimag_id": 1,
  "sum_id": 1,
  "bag_id": 1,
  "address_detail": "Баянзүрх дүүрэг"
}
```

#### PUT /organization/:id
**Тайлбар:** Байгууллага засварлах  
**Auth:** ✅ Required

#### DELETE /organization/:id
**Тайлбар:** Байгууллага устгах  
**Auth:** ✅ Required

#### GET /organization/tree
**Тайлбар:** Байгууллагын модон бүтэц  
**Auth:** ✅ Required  
**Query Parameters:**
- `org_id` (required): Root organization ID

---

### 10. Organization User (`/orguser`)

#### GET /orguser
**Тайлбар:** Байгууллага-хэрэглэгчийн холбоосын жагсаалт  
**Auth:** ✅ Required  
**Query Parameters:**
- `org_id` (optional): Organization ID
- `user_id` (optional): User ID
- `name` (optional): Name search
- Standard pagination

#### GET /orguser/users
**Тайлбар:** Байгууллагын хэрэглэгчид  
**Auth:** ✅ Required  
**Query Parameters:**
- `org_id` (required): Organization ID

#### GET /orguser/organizations
**Тайлбар:** Хэрэглэгчийн байгууллагууд  
**Auth:** ✅ Required  
**Query Parameters:**
- `user_id` (required): User ID

#### POST /orguser
**Тайлбар:** Байгууллагад хэрэглэгч нэмэх  
**Auth:** ✅ Required  
**Request Body:**
```json
{
  "org_id": 10,
  "user_id": 123
}
```

#### DELETE /orguser
**Тайлбар:** Байгууллагаас хэрэглэгч хасах  
**Auth:** ✅ Required  
**Request Body:**
```json
{
  "org_id": 10,
  "user_id": 123
}
```

---

### 11. Organization Type (`/orgtype`)

#### GET /orgtype
**Тайлбар:** Байгууллагын төрлийн жагсаалт  
**Auth:** ✅ Required

#### POST /orgtype
**Тайлбар:** Байгууллагын төрөл үүсгэх  
**Auth:** ✅ Required  
**Request Body:**
```json
{
  "code": "COMPANY",
  "name": "Компани",
  "description": "Аж ахуйн нэгж"
}
```

#### PUT /orgtype/:id
**Тайлбар:** Байгууллагын төрөл засварлах  
**Auth:** ✅ Required

#### DELETE /orgtype/:id
**Тайлбар:** Байгууллагын төрөл устгах  
**Auth:** ✅ Required

#### GET /orgtype/system
**Тайлбар:** Байгууллагын төрлийн системүүд  
**Auth:** ✅ Required  
**Query Parameters:**
- `type_id` (required): Organization type ID

#### POST /orgtype/system
**Тайлбар:** Төрөлд систем нэмэх  
**Auth:** ✅ Required  
**Request Body:**
```json
{
  "type_id": 1,
  "system_ids": [1, 2, 3]
}
```

---

### 12. Terminal Management (`/terminal`)

#### GET /terminal
**Тайлбар:** Терминалын жагсаалт  
**Auth:** ✅ Required

#### POST /terminal
**Тайлбар:** Терминал үүсгэх  
**Auth:** ✅ Required

#### PUT /terminal/:id
**Тайлбар:** Терминал засварлах  
**Auth:** ✅ Required

#### DELETE /terminal/:id
**Тайлбар:** Терминал устгах  
**Auth:** ✅ Required

---

### 13. OAuth Client Management (`/client`)

#### GET /client
**Тайлбар:** OAuth client-ийн жагсаалт  
**Auth:** ✅ Required

#### GET /client/scope
**Тайлбар:** OAuth scope-ийн жагсаалт  
**Auth:** ✅ Required

#### POST /client/scope
**Тайлбар:** Scope үүсгэх  
**Auth:** ✅ Required

#### DELETE /client/scope
**Тайлбар:** Scope устгах  
**Auth:** ✅ Required

---

### 14. App Service Icon (`/app-service-icon`)

#### GET /app-service-icon
**Тайлбар:** App service icon-ы жагсаалт  
**Auth:** ✅ Required

#### POST /app-service-icon
**Тайлбар:** Icon үүсгэх  
**Auth:** ✅ Required

#### PUT /app-service-icon/:id
**Тайлбар:** Icon засварлах  
**Auth:** ✅ Required

#### DELETE /app-service-icon/:id
**Тайлбар:** Icon устгах  
**Auth:** ✅ Required

---

### 15. App Service Group (`/app-service-group`)

#### GET /app-service-group
**Тайлбар:** Service group-ийн жагсаалт  
**Auth:** ✅ Required

#### GET /app-service-group/with-icons
**Тайлбар:** Icon-тай service group-ийн жагсаалт  
**Auth:** ✅ Required

#### POST /app-service-group
**Тайлбар:** Group үүсгэх  
**Auth:** ✅ Required

#### PUT /app-service-group/:id
**Тайлбар:** Group засварлах  
**Auth:** ✅ Required

#### DELETE /app-service-group/:id
**Тайлбар:** Group устгах  
**Auth:** ✅ Required

---

### 16. Platform Desktop Icon (`/app-desktop-icon`)

#### GET /app-desktop-icon
**Тайлбар:** Desktop icon модон бүтэц  
**Auth:** ✅ Required

#### POST /app-desktop-icon
**Тайлбар:** Desktop icon үүсгэх  
**Auth:** ✅ Required

#### PUT /app-desktop-icon/:id
**Тайлбар:** Desktop icon засварлах  
**Auth:** ✅ Required

#### DELETE /app-desktop-icon/:id
**Тайлбар:** Desktop icon устгах  
**Auth:** ✅ Required

---

### 17. Platform Business Icon (`/app-business-icon`)

#### GET /app-business-icon
**Тайлбар:** Business icon модон бүтэц  
**Auth:** ✅ Required

#### POST /app-business-icon
**Тайлбар:** Business icon үүсгэх  
**Auth:** ✅ Required

#### PUT /app-business-icon/:id
**Тайлбар:** Business icon засварлах  
**Auth:** ✅ Required

#### DELETE /app-business-icon/:id
**Тайлбар:** Business icon устгах  
**Auth:** ✅ Required

---

### 18. File Management (`/file`)

#### GET /file/list
**Тайлбар:** Файлын жагсаалт  
**Auth:** ✅ Required

#### POST /file/upload
**Тайлбар:** Файл upload хийх  
**Auth:** ✅ Required  
**Content-Type:** `multipart/form-data`  
**Request Body:**
- `file`: File binary

**Response:**
```json
{
  "code": "OK",
  "data": {
    "uuid": "file-uuid",
    "filename": "document.pdf",
    "size": 1024000,
    "url": "https://api.example.com/file/file-uuid"
  }
}
```

#### DELETE /file
**Тайлбар:** Файл устгах  
**Auth:** ✅ Required  
**Request Body:**
```json
{
  "uuid": "file-uuid"
}
```

#### GET /file/:uuid
**Тайлбар:** Файл татах (public)  
**Auth:** Шаардлагагүй  
**URL Parameters:**
- `uuid` (required): File UUID

**Response:** File binary stream

---

### 19. Notification Management (`/notification`)

#### GET /notification
**Тайлбар:** Мэдэгдлийн жагсаалт  
**Auth:** ✅ Required

#### GET /notification/groups
**Тайлбар:** Мэдэгдлийн бүлгүүд  
**Auth:** ✅ Required

#### POST /notification
**Тайлбар:** Мэдэгдэл илгээх  
**Auth:** ✅ Required  
**Request Body:**
```json
{
  "title": "Шинэ мэдэгдэл",
  "message": "Мэдэгдлийн агуулга",
  "user_ids": [1, 2, 3],
  "type": "info"
}
```

#### POST /notification/read
**Тайлбар:** Мэдэгдэл уншсан гэж тэмдэглэх  
**Auth:** ✅ Required  
**Request Body:**
```json
{
  "notification_id": 123
}
```

#### POST /notification/read-all
**Тайлбар:** Бүх мэдэгдэл уншсан гэж тэмдэглэх  
**Auth:** ✅ Required

---

### 20. News Management (`/news`)

#### GET /news
**Тайлбар:** Мэдээний жагсаалт (public)  
**Auth:** Шаардлагагүй

#### GET /news/get/:id
**Тайлбар:** Мэдээний дэлгэрэнгүй (public)  
**Auth:** Шаардлагагүй  
**URL Parameters:**
- `id` (required): News ID

#### POST /news
**Тайлбар:** Мэдээ үүсгэх  
**Auth:** ✅ Required

#### PUT /news/:id
**Тайлбар:** Мэдээ засварлах  
**Auth:** ✅ Required

#### DELETE /news/:id
**Тайлбар:** Мэдээ устгах  
**Auth:** ✅ Required

---

### 21. Verification (`/verify`)

#### GET /verify/dan
**Тайлбар:** DAN баталгаажуулалт  
**Auth:** ✅ Required

#### POST /verify/email
**Тайлбар:** Email баталгаажуулах код илгээх  
**Auth:** ✅ Required  
**Request Body:**
```json
{
  "email": "user@example.com"
}
```

#### POST /verify/email/confirm
**Тайлбар:** Email баталгаажуулах  
**Auth:** ✅ Required  
**Request Body:**
```json
{
  "email": "user@example.com",
  "code": "123456"
}
```

#### POST /verify/phone
**Тайлбар:** Утас баталгаажуулах код илгээх  
**Auth:** ✅ Required  
**Request Body:**
```json
{
  "phone_no": "99119911"
}
```

#### POST /verify/phone/confirm
**Тайлбар:** Утас баталгаажуулах  
**Auth:** ✅ Required  
**Request Body:**
```json
{
  "phone_no": "99119911",
  "code": "123456"
}
```

---

### 22. Video Conference Room (`/room`)

#### GET /room
**Тайлбар:** Өрөөний жагсаалт  
**Auth:** ✅ Required

#### GET /room/token
**Тайлбар:** Өрөөнд орох token үүсгэх  
**Auth:** ✅ Required  
**Query Parameters:**
- `room_id` (required): Room ID

**Response:**
```json
{
  "code": "OK",
  "data": {
    "token": "jwt-token",
    "room_id": "room-123",
    "expires_at": "2025-12-08T10:00:00Z"
  }
}
```

#### POST /room
**Тайлбар:** Өрөө үүсгэх  
**Auth:** ✅ Required  
**Request Body:**
```json
{
  "name": "Хурлын өрөө",
  "description": "Долоо хоногийн хурал",
  "max_participants": 10,
  "scheduled_at": "2025-12-10T14:00:00Z"
}
```

#### POST /room/join
**Тайлбар:** Өрөөнд орох  
**Auth:** ✅ Required  
**Request Body:**
```json
{
  "room_id": "room-123"
}
```

#### POST /room/:id/users
**Тайлбар:** Өрөөнд хэрэглэгч нэмэх  
**Auth:** ✅ Required  
**URL Parameters:**
- `id` (required): Room ID

**Request Body:**
```json
{
  "user_ids": [1, 2, 3]
}
```

#### DELETE /room/:id
**Тайлбар:** Өрөө устгах  
**Auth:** ✅ Required

#### DELETE /room/:id/users/:user_id
**Тайлбар:** Өрөөнөөс хэрэглэгч хасах  
**Auth:** ✅ Required  
**URL Parameters:**
- `id` (required): Room ID
- `user_id` (required): User ID

---

### 23. Terminal Payment (`/tpay`)

#### GET /tpay/accounts/me
**Тайлбар:** Миний дансны жагсаалт  
**Auth:** ✅ Required

#### PUT /tpay/accounts/set-default
**Тайлбар:** Үндсэн данс тохируулах  
**Auth:** ✅ Required  
**Request Body:**
```json
{
  "account_id": "acc-123"
}
```

#### GET /tpay/accounts/statement
**Тайлбар:** Дансны хуулга  
**Auth:** ✅ Required  
**Query Parameters:**
- `account_id` (required): Account ID
- `from_date` (optional): YYYY-MM-DD
- `to_date` (optional): YYYY-MM-DD

#### POST /tpay/accounts/:account_id/qr
**Тайлбар:** QR код үүсгэх  
**Auth:** ✅ Required  
**URL Parameters:**
- `account_id` (required): Account ID

**Request Body:**
```json
{
  "amount": 50000,
  "description": "Төлбөрийн тайлбар"
}
```

**Response:**
```json
{
  "code": "OK",
  "data": {
    "qr_code": "base64-qr-image",
    "qr_string": "qr-text-value",
    "expires_at": "2025-12-08T15:00:00Z"
  }
}
```

#### POST /tpay/transaction/qr_pay
**Тайлбар:** QR төлбөр хийх  
**Auth:** ✅ Required  
**Request Body:**
```json
{
  "qr_string": "qr-text-value",
  "account_id": "acc-123"
}
```

#### POST /tpay/p2p
**Тайлбар:** Данс хоорондын шилжүүлэг  
**Auth:** ✅ Required  
**Request Body:**
```json
{
  "from_account_id": "acc-123",
  "to_account_id": "acc-456",
  "amount": 10000,
  "description": "Шилжүүлгийн утга"
}
```

#### GET /tpay/card/list
**Тайлбар:** Картын жагсаалт  
**Auth:** ✅ Required

#### POST /tpay/card/create
**Тайлбар:** Карт нэмэх  
**Auth:** ✅ Required  
**Request Body:**
```json
{
  "card_number": "1234567812345678",
  "card_holder": "BAT BOLD",
  "expire_date": "12/25"
}
```

#### POST /tpay/card/confirm
**Тайлбар:** Карт баталгаажуулах  
**Auth:** ✅ Required  
**Request Body:**
```json
{
  "card_id": "card-123",
  "otp_code": "123456"
}
```

#### GET /tpay/card/send_otp
**Тайлбар:** Картын OTP илгээх  
**Auth:** ✅ Required  
**Query Parameters:**
- `card_id` (required): Card ID

#### POST /tpay/verify_card
**Тайлбар:** Карт шалгах  
**Auth:** ✅ Required  
**Request Body:**
```json
{
  "card_number": "1234567812345678"
}
```

---

### 24. Chat Management (`/chat`)

#### GET /chat
**Тайлбар:** Chat item-ийн жагсаалт  
**Auth:** ✅ Required

#### POST /chat
**Тайлбар:** Chat item үүсгэх  
**Auth:** ✅ Required

#### PUT /chat/:id
**Тайлбар:** Chat item засварлах  
**Auth:** ✅ Required

#### DELETE /chat/:id
**Тайлбар:** Chat item устгах  
**Auth:** ✅ Required

#### POST /chat/key
**Тайлбар:** Key-ээр chat item хайх  
**Auth:** ✅ Required  
**Request Body:**
```json
{
  "key": "chat-key-123"
}
```

---

## Common Data Models

### User Model
```json
{
  "id": 123,
  "civil_id": 456,
  "reg_no": "УА12345678",
  "family_name": "Болд",
  "last_name": "Бат",
  "first_name": "Дорж",
  "gender": 1,
  "birth_date": "1990-01-01",
  "phone_no": "99119911",
  "email": "bat@example.com",
  "created_at": "2025-01-01T10:00:00Z",
  "updated_at": "2025-01-01T10:00:00Z"
}
```

### Organization Model
```json
{
  "id": 10,
  "reg_no": "1234567",
  "name": "ХХК Компани",
  "short_name": "Компани",
  "type_id": 1,
  "phone_no": "75001122",
  "email": "info@company.mn",
  "longitude": 106.9177,
  "latitude": 47.9187,
  "is_active": true,
  "aimag_id": 1,
  "sum_id": 1,
  "bag_id": 1,
  "address_detail": "Баянзүрх дүүрэг",
  "parent_id": null
}
```

### Role Model
```json
{
  "id": 5,
  "system_id": 1,
  "code": "ADMIN",
  "name": "Администратор",
  "description": "Системийн админ",
  "is_active": true,
  "created_at": "2025-01-01T10:00:00Z"
}
```

### Permission Model
```json
{
  "id": 15,
  "code": "USER_CREATE",
  "name": "Хэрэглэгч үүсгэх",
  "description": "Хэрэглэгч үүсгэх эрх",
  "resource": "user",
  "action": "create"
}
```

---

## Error Codes

| Code | HTTP Status | Тайлбар |
|------|-------------|---------|
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

## Request Headers

### Required Headers (Protected endpoints)
```http
Cookie: sid=session-id-here
Content-Type: application/json
```

### Optional Headers
```http
X-Request-ID: uuid-v4          # Request tracking
Accept-Language: mn            # Language preference (mn, en)
User-Agent: MyApp/1.0          # Client info
```

---

## Response Headers

```http
Content-Type: application/json; charset=utf-8
X-Request-ID: uuid-v4
X-Response-Time: 45ms
```

---

## Rate Limiting

- **Default:** 100 requests/minute per IP
- **Authenticated:** 1000 requests/minute per user
- **File Upload:** 10 requests/minute per user

### Rate Limit Headers
```http
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 95
X-RateLimit-Reset: 1702023600
```

---

## Timeout Configuration

- **Default timeout:** 5 seconds
- **File upload:** 30 seconds
- **Long-running operations:** 60 seconds

---

## Versioning

API version нь URL-д тусгагдахгүй (v1 default).
Хэрэв шинэ version ирвэл `/v2/` prefix ашиглана.

---

## Support & Contact

- **Technical Support:** support@example.com
- **API Issues:** api@example.com
- **Documentation:** https://docs.example.com

---

## Changelog

### Version 1.0.0 (2025-12-08)
- Initial API contract
- All core endpoints documented
- Authentication & authorization flow defined
- Common data models standardized

---

**Last Updated:** 2025-12-08  
**Document Version:** 1.0.0  
**Maintained by:** Gerege Core Team

