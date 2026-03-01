# Interview Supplemental Files (Go Implementation)

โปรเจกต์นี้เป็น RESTful API สำหรับจัดการ Cloud Servers โดยมีการเชื่อมต่อไปยัง Mocked Infrastructure Microservice ภายนอก

## สิ่งที่ส่งตามโจทย์

- **Task 1 (Database Schema):** `docs/TASK_ONE.md`
- **Task 2 (Flow testing):** `docs/FLOW_TESTING.md`
- **Task 3 (Design Doc):** อยู่ในหัวข้อ `Task 3: Real-world resiliency design` ด้านล่าง

## Tech Stack

- Go (Fiber)
- In-memory repository (ไม่มีฐานข้อมูลจริง)
- JWT auth + CORS

## วิธีรันโปรเจกต์

### 1) รัน Mocked Infrastructure Service (port 8081)

```bash
make run-infra-service
```

### 2) รัน API Service (port 8080)

เปิดอีก terminal แล้วรัน:

```bash
make gen
make run
```

### 3) Swagger

เปิดที่:

`http://localhost:8080/swagger/index.html`

## Configuration

ระบบมีค่า default ในโค้ดอยู่แล้ว แต่สามารถ override ด้วย env ได้:

- `SERVICE_PORT` (default: `8080`)
- `SERVICE_ALLOWED_ORIGINS` (default: `http://localhost:3000`)
- `INFRA_SERVICE_URL` (default: `http://localhost:8081`)
- `JWT_SECRET` (default: `dev-secret-change-in-production`)

> ดูตัวอย่างไฟล์ env ได้ที่ `backend/go/.env.template`

---

## Task 3: Real-world resiliency design

### ปัญหาที่เจอจริง

Mocked Infrastructure Service ถูกออกแบบให้ flaky (ช้า/ล่ม/สุ่ม error) ทำให้ endpoint ฝั่งเรา เช่น `POST /api/v1/servers` สามารถได้ `502 Bad Gateway` ได้จริงเมื่อเรียก infra ไม่สำเร็จ

ตัวอย่างอาการที่พบ:
- ยิง `POST /api/v1/servers` แล้วได้ `502`
- สาเหตุหลักมาจากขั้นตอน validate SKU (`/v1/skus`) หรือ provision (`/v1/resources`) ที่ infra ตอบ error/timeout

### เป้าหมายการออกแบบ

1. ผู้ใช้ไม่ต้องรอจน timeout แบบไม่มีขอบเขต
2. สถานะข้อมูลฝั่งเราไม่เพี้ยน เมื่อ external call ล้มเหลวกลางทาง
3. บอกผู้ใช้ได้ชัดว่าควร retry เองหรือรอระบบ retry
4. วิธีแก้ต้องแข็งแรง แต่ไม่ over-engineered

### แนวทางที่เลือก (ไม่ over-engineered)

#### 1) กำหนด timeout ที่ชัดเจนต่อ request

- ตั้ง timeout สำหรับ call ไป infra ทุกครั้ง (โค้ดปัจจุบันมี timeout client แล้ว)
- เมื่อ timeout ให้คืน `502` พร้อมข้อความที่สื่อว่า dependency ล่ม/ช้า

เหตุผล:
- กัน request ค้างนาน
- ควบคุม latency สูงสุดฝั่ง API ได้

#### 2) Retry แบบจำกัดครั้ง เฉพาะ error ที่เหมาะสม

- retry เฉพาะ transient failures เช่น network error, timeout, `5xx`
- จำนวนครั้งแนะนำ: 2-3 ครั้ง พร้อม backoff สั้น ๆ (เช่น 100ms, 300ms)
- ไม่ retry กรณี `4xx` จาก infra เพราะเป็น logical/client error

เหตุผล:
- ลด false failure จากอาการสะดุดชั่วคราว
- ไม่เพิ่มโหลดเกินจำเป็น

#### 3) ทำ provisioning ให้เป็น state machine

เพิ่มสถานะ server ฝั่งเราให้รองรับงาน async/failure:
- `provisioning`
- `active`
- `failed`

flow ที่แนะนำ:
1. รับคำขอ `POST /servers` แล้วสร้าง record สถานะ `provisioning` ก่อน
2. เรียก infra เพื่อ provision
3. ถ้าสำเร็จ อัปเดตเป็น `active` พร้อมเก็บ `infrastructure_resource_id`
4. ถ้าล้มเหลว อัปเดตเป็น `failed` พร้อม `error_reason`

เหตุผล:
- ป้องกันกรณีล่มกลางทางแล้วข้อมูลหาย/ไม่สอดคล้อง
- ผู้ใช้เห็นสถานะปัจจุบันชัดเจนจาก `GET /servers`

#### 4) ปรับ response ให้สื่อสารได้ดีขึ้น

สำหรับ `POST /servers`:
- ถ้าจบเร็วและสำเร็จ: `200 OK` + `{ "success": true, "id": "..." }`
- ถ้าภายนอกล้มเหลวชั่วคราว: `502 Bad Gateway` + error code เช่น `INFRA_UNAVAILABLE`
- ถ้าจะรองรับ async เต็มรูปแบบ: `202 Accepted` + server id ฝั่งเราเพื่อไป poll สถานะ

เหตุผล:
- ให้ client ตัดสินใจได้ว่าจะ retry แบบไหน
- รองรับ UX ที่ดีขึ้นเมื่อระบบภายนอกไม่เสถียร

#### 5) Idempotency สำหรับการกดซ้ำ

- รองรับ `Idempotency-Key` ที่ `POST /servers`
- ถ้าผู้ใช้ส่ง key เดิมภายในช่วงเวลาเดียวกัน ให้คืนผลเดิมแทนสร้างซ้ำ

เหตุผล:
- กัน duplicate provisioning เมื่อผู้ใช้ retry เองจาก timeout/502

#### 6) Logging และ observability ขั้นพื้นฐาน

- log ต่อ request: latency, endpoint, status จาก infra, retry count, error type
- ผูกกับ `activity_logs` เพื่อ trace ว่าใครกด provision และผลลัพธ์คืออะไร

เหตุผล:
- ช่วย debug production incident ได้เร็ว
- ใช้งานจริงได้โดยไม่ต้องเพิ่มระบบใหญ่

### การเปลี่ยนแปลง schema ที่แนะนำ (ถ้าจะ implement ต่อ)

ใน `servers`:
- `status` (`provisioning|active|failed`)
- `error_reason` (nullable string)
- `retry_count` (int, default 0)
- `last_attempt_at` (timestamp, nullable)

ใน `activity_logs`:
- `error_code` (nullable string)
- `metadata` (json/string)

### การ map HTTP status ที่แนะนำ

- `400` ข้อมูล request ไม่ถูกต้อง (เช่น sku format ผิด)
- `401` ไม่มี/ใช้ token ไม่ถูกต้อง
- `404` ไม่พบ server
- `409` idempotency conflict (ถ้ามี policy เฉพาะ)
- `502` infra ล่ม/timeout/ตอบ 5xx
- `500` internal error ฝั่ง API

### เหตุผลที่ยังไม่ใช้ระบบใหญ่ (เช่น Kafka)

โจทย์กำหนดเวลา 2-3 ชั่วโมง และระบบยังขนาดเล็ก จึงเลือก:
- timeout + retry + state machine + idempotency + logging

ซึ่งได้ reliability เพิ่มขึ้นมากโดย complexity ยังพอเหมาะ และสามารถขยายไป job queue/circuit breaker ภายหลังได้

### Test plan สำหรับ Task 3

1. **Happy path**
   - `POST /servers` ด้วย sku ถูกต้อง ควรได้ `200`
2. **Infra flaky**
   - ยิงซ้ำหลายครั้งจนเจอ `502`, ตรวจว่าระบบไม่ค้างและ response time อยู่ในกรอบ timeout
3. **Retry effectiveness**
   - จำลอง network ช้า/สะดุดชั่วคราว แล้วตรวจว่า request บางส่วน recover ได้
4. **Data consistency**
   - ถ้าล้มเหลวกลางทาง ต้องไม่เกิด server ที่สถานะกำกวม
5. **Idempotency**
   - ส่ง `Idempotency-Key` เดิมซ้ำ ต้องไม่สร้าง resource ซ้ำ

