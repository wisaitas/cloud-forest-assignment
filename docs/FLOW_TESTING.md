# คู่มือการเทสตามโจทย์ (Task 2)

API ใช้ **base path** `/api/v1` และรันที่ port **8080**  
Infrastructure Mock Service ต้องรันที่ port **8081**

---

## 1. เตรียมสภาพแวดล้อม

### 1.1 รัน Infrastructure Mock Service (port 8081)

```bash
make run-infra-service
# หรือ
cd backend/go/cmd/infraservice && go run main.go
```

รอจนเห็นว่า service พร้อม (เช่น listen :8081)

### 1.2 รัน Interview API Service (port 8080)

เปิด terminal อีก tab/window:

```bash
make gen
make run
# หรือ
cd backend/go && go run ./cmd/interviewservice
```

---

## 2. เทสตามโจทย์

### 2.0
เปิด swagger ที่ http://localhost:8080/swagger/index.html

### 2.1 Authentication — POST /auth/login

**โจทย์:** Authenticate ด้วย seed user แล้วได้ JWT ใน HttpOnly cookie

**Seed user (ต้องมีในระบบ):**
- Email: `john.smith@gmail.com`
- Password: `not-so-secure-password`

| รายการ | ค่า |
|--------|-----|
| URL | `POST http://localhost:8080/api/v1/auth/login` |
| Body (JSON) | `{"email":"john.smith@gmail.com","password":"not-so-secure-password"}` |

**ผลที่คาดหวัง:**
- Status **200**
- Response มี `access_token` (และอาจมี `refresh_token`)
- Response header มี **Set-Cookie** ที่มี `access_token` และ **HttpOnly**

---

### 2.2 Server Management (Protected)

ทุก endpoint ด้านล่าง **ต้องส่ง JWT** ผ่าน:
Header ใน swagger
ถ้าไม่ส่งหรือ token ไม่ถูกต้อง → **401 Unauthorized**

---

#### GET /servers — รายการ servers ของ user

| รายการ | ค่า |
|--------|-----|
| URL | `GET http://localhost:8080/api/v1/servers` |
| Header | `Authorization: <access_token>` |

**ผลที่คาดหวัง:**
- Status **200**
- Body มี `servers` (array) — ครั้งแรกอาจเป็น `[]`

---

#### POST /servers — Provision server ใหม่

**โจทย์:** ต้อง validate SKU กับ infra แล้วได้ Infrastructure Resource ID กลับมา

| รายการ | ค่า |
|--------|-----|
| URL | `POST http://localhost:8080/api/v1/servers` |
| Header | `Authorization: Bearer <access_token>` |
| Body (JSON) | `{"sku":"C1-R1GB-D20GB"}` | (ต้องเป็น SKU ที่มีใน infraservice เช่น C1-R1GB-D20GB)

**ผลที่คาดหวัง:**
- Status **200**
- Body: `{"success": true, "id": "<Infrastructure Resource ID>"}`  
  (เช่น `"id": "i-xxxx"` จาก mock infra)

**เทส error ตามโจทย์:**
- SKU ไม่มีใน `/v1/skus` หรือ body ไม่ถูกต้อง → **400**
- Infra service ล้มหรือ error → **502** (หรือ 5xx ตามที่ implement)

**หลัง provision สำเร็จ:** ใช้ **Server ID** (ที่ API สร้างเอง ไม่ใช่ Infrastructure Resource ID) สำหรับ power endpoint — ดูได้จาก `GET /servers` ว่าแต่ละตัวมี `id` อะไร

---

#### POST /servers/:server-id/power — เปิด/ปิด server

**โจทย์:** ใช้ **Server ID** (จาก GET /servers) ไม่ใช่ Infrastructure Resource ID

| รายการ | ค่า |
|--------|-----|
| URL | `POST http://localhost:8080/api/v1/servers/<server-id>/power` |
| Header | `Authorization: Bearer <access_token>` |
| Body (JSON) | `{"action":"on"}` หรือ `{"action":"off"}` |

**ผลที่คาดหวัง:**
- Status **200**
- Body: `{"success": true, "state": "on"}` หรือ `"state": "off"`

**เทส error ตามโจทย์:**
- Server ID ไม่มีในระบบ / ไม่ใช่ของ user นี้ → **404**
- Body ไม่มี `action` หรือไม่ใช่ `on`/`off` → **400**

---

## 3. ลำดับเทสแนะนำ (Manual flow)

1. **รัน infra + API** (ตามข้อ 1)
2. **Login**  
   `POST /api/v1/auth/login` ด้วย seed user → เก็บ `access_token` จาก response หรือใช้ cookie
3. **List servers**  
   `GET /api/v1/servers` + Bearer token → ควรได้ `[]` หรือรายการเดิม
4. **Provision**  
   `POST /api/v1/servers` body `{"sku":"C1-R1GB-D20GB"}` → เก็บ `id` จาก response (Infra Resource ID) และจาก `GET /servers` เก็บ **Server ID** ของ server ที่สร้าง
5. **Power**  
   `POST /api/v1/servers/<server-id>/power` body `{"action":"off"}` แล้ว `{"action":"on"}` → ตรวจ response `success` และ `state`
6. **(Optional) เทส CORS**  
   เปิดแอปที่รันบน `http://localhost:3000` แล้วเรียก API — ควรผ่าน; เรียกจาก origin อื่น ควรถูก block

---
