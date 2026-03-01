# Back-end

**เวลา:** งานนี้ควรทำเสร็จภายใน 2–3 ชั่วโมง

**ส택ที่ใช้ได้:** Go (แนะนำ: Go Fiber) หรือ Node (แนะนำ: Express และ Elysia บนรันไทม์ใดก็ได้)

**เครื่องมือ:** อนุญาตและสนับสนุนให้ใช้ AI Assistants (ทั้งแท็บ autocomplete และ coding assistant) แต่คุณต้องเป็นเจ้าของและเข้าใจโค้ดที่ส่ง

**ไฟล์เสริม:** https://github.com/CLOUDFOREST-CO-TH/interview-supplemental-files

---

# ภาพรวมงาน

คุณจะสร้าง RESTful API Service สำหรับจัดการ Cloud Servers ที่ประสานงานระหว่าง data store ภายในกับ infrastructure micro-service ภายนอก

![image.png](./image.png)

## งานที่ 1: **การออกแบบฐานข้อมูล**

สร้างไฟล์ Markdown อธิบาย data schema ของแอปนี้ คุณกำหนดฟิลด์และชนิดข้อมูลได้ตามต้องการ แต่ต้องเป็นไปตามข้อกำหนดดังนี้:

- **Users:** ต้องเก็บข้อมูลสำหรับการยืนยันตัวตน (อีเมลและรหัสผ่าน)
- **Servers:** ต้องเก็บสถานะการเปิดเครื่อง (เช่น On/Off) และ SKU (เช่น "C1-R1GB-D40GB" สำหรับ 1 core CPU, 1 GB RAM, 40 GB Disk)
    - **ข้อกำหนดสำคัญ:** เมื่อ provision server แล้ว infrastructure micro-service จะออก **Infrastructure Resource ID** ที่ไม่ซ้ำ คุณต้องออกแบบ schema ให้เก็บ ID นี้เพื่อจัดการ server ในภายหลัง Infrastructure Resource ID นี้ไม่ควรใช้เป็น Server ID ฝั่ง API Service
- **ActivityLogs:** ต้องทำหน้าที่เป็น audit trail สำหรับการกระทำผ่าน API

<aside>
ℹ️

**หมายเหตุเรื่องการ implement:** ในส่วนเขียนโค้ด (งานที่ 2 และ 3) คุณไม่ต้องใช้ฐานข้อมูลจริง จะ implement design นี้ด้วย in-memory storage แบบง่าย (เช่น Go Structs หรือ JS Objects) ออกแบบ schema ให้เหมาะกับความง่ายนี้

</aside>

## งานที่ 2: การ implement Endpoints

สร้าง RESTful API ด้วย Go (แนะนำ: Fiber) หรือ Node (แนะนำ: Express หรือ Elysia)

- **การเก็บข้อมูล:** ใช้ in-memory storage แบบง่าย (เช่น Global Structs/Objects/Arrays) **ไม่ต้องใช้ฐานข้อมูลจริง**
- **Middleware:** implement JWT Authentication และ CORS (อนุญาตเฉพาะ `localhost:3000`)

**ข้อกำหนดก่อนเริ่ม:**

ก่อนเริ่ม ต้องดาวน์โหลด **Mocked Infrastructure Microservice** จาก https://github.com/CLOUDFOREST-CO-TH/interview-supplemental-files (โฟลเดอร์ /backend) แล้วรันบนเครื่องที่พอร์ต `:8081` (`go run backend/go/server.go` หรือ `node backend/node/server.mjs`) API ของคุณจะต้องส่ง HTTP request ไปยัง service นี้เพื่อ provision และจัดการ servers ดู `swagger.yaml` สำหรับ API docs

**Endpoints:**

1. **Authentication**

    **POST /auth/login**

    - **พฤติกรรม:** ยืนยันตัวตนผู้ใช้และออก JWT token ใน cookie แบบ `HttpOnly`
    - **Seed Data:** แอปต้องโหลดผู้ใช้นี้ไว้ล่วงหน้า:
        - **User ID:** `123123123`
        - **Email:** `john.smith@gmail.com`
        - **Password:** `not-so-secure-password`

2. **Server Management (Protected)**

    **Endpoints เหล่านี้ต้องตรวจสอบ JWT ที่ valid และบังคับใช้ CORS**

    **GET /servers**

    - **คำอธิบาย:** แสดงรายการ servers ทั้งหมดของผู้ใช้ที่ล็อกอินแล้ว

    **POST /servers**

    - **คำอธิบาย:** Provision server ใหม่ ต้องเรียก Mock Infrastructure Service เพื่อรับ Infrastructure Resource ID แล้วบันทึกรายละเอียด server ลง in-memory store
    - **Request Body:** `{ "sku": "C1-R1GB-D40GB" }`
    - **Response:** `{ "success": true, "id": "<Infrastructure Resource ID>" }`
    - **การจัดการข้อผิดพลาด:**
        - ตรวจสอบ sku กับรายการจาก microservice API `/v1/skus` คืน HTTP status ที่เหมาะสมเมื่อข้อมูลผิดรูปแบบหรือไม่ถูกต้อง

    **POST /servers/:server-id/power**

    - **คำอธิบาย:** เปลี่ยนสถานะการเปิดเครื่องของ server ต้องค้นหา Infrastructure Resource ID ของ server แล้วเรียก Mock Service เพื่อ apply การเปลี่ยนแปลง
    - **Request Body:** `{ "action": "on" }` หรือ `{ "action": "off" }`
    - **Response:** `{ "success": true, "state": "on" }`
    - **การจัดการข้อผิดพลาด:**
        - จัดการกรณีผิดพลาด เช่น ไม่พบ server ID โดยใช้ HTTP status ที่เหมาะสม

## งานที่ 3: กรณีใช้งานจริง

คุณอาจสังเกตว่า **Mocked Infrastructure Microservice** ถูกออกแบบให้ไม่เสถียรโดยตั้งใจ บางครั้งตอบช้าหรือคืน error ในโลกจริงเราต้องสร้าง service ที่ทนทานและจัดการความล้มเหลวเหล่านี้ได้อย่างเหมาะสม

**ข้อกำหนด:**

ปรับ implementation ให้จัดการความล้มเหลวจากภายนอก **หรือ** ถ้าเวลาจำกัด ให้เขียน design doc เป็น Markdown แทน

- คุณอาจเพิ่ม/แก้ schema ฐานข้อมูล logic ของ endpoint หรือรูปแบบ response ตามความเหมาะสม
- วิธีแก้ควรแข็งแรงแต่ไม่ over-engineered (เช่น ไม่จำเป็นต้องใช้ message queue แบบเต็มอย่าง Kafka แต่ควรอธิบายเหตุผลได้)

**คำถามที่ควรคิดขณะ implement:**

- ถ้า external service ล่มระหว่าง provisioning จะเกิดอะไรขึ้น?
- จะทำให้ผู้ใช้ไม่ต้องรอไม่รู้จบได้อย่างไร?

# วิธีส่งงาน

มีสองวิธีส่งงาน

1. Push โค้ดไปยัง repository สาธารณะบน GitHub หรือ GitLab แล้วส่งลิงก์ทางอีเมล
    1. ต้องมี `README.md` ที่ชัดเจน มี:
        - วิธีติดตั้ง dependencies และรันโปรเจกต์บนเครื่องตัวเอง รวมถึง Markdown ของงานที่ 1
        - (ถ้าทำ) ถ้าเลือกทาง "Design Doc" สำหรับงานที่ 3 ให้ใส่ไว้ที่นี่
2. หรือส่งทางอีเมล: ไฟล์ .zip แล้วส่งกลับมาที่อีเมล
