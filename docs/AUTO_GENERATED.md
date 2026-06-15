# API Tài Khoản Người Dùng: Tạo Mới (Create User)

## 📖 Tổng Quan (Overview)

API này cho phép hệ thống tạo một tài khoản người dùng mới trong cơ sở dữ liệu Identity. Đây là bước đầu tiên để một khách hàng hoặc nhân viên đăng ký vào hệ thống Clawathon 2026.

*   **Đối tượng sử dụng:** Ứng dụng Frontend (Web/Mobile), Hệ thống đối tác tích hợp.
*   **Mục đích:** Lưu trữ thông tin đăng nhập và hồ sơ cơ bản của người dùng.
*   **Quyền hạn:** Công khai (Public) - Không cần token xác thực để gọi API đăng ký.

---

## 🔌 Thông Tin Kỹ Thuật (Technical Specifications)

| Thuộc tính | Giá trị |
| :--- | :--- |
| **Endpoint** | `/identity/v1/user` |
| **HTTP Method** | `POST` |
| **Content-Type** | `application/json` |
| **Version** | v1 |
| **Module** | Identity |

### Tham Số Đầu Vào (Request Parameters)

#### 1. Query Parameters (Tham số truy vấn)

| Tên | Kiểu | Bắt Buộc | Mô Tả | Ví Dụ |
| :--- | :--- | :--- | :--- | :--- |
| `lang` | String | Không | Ngôn ngữ phản hồi lỗi/thông báo | `vi`, `en` |

#### 2. Body Parameters (Nội dung yêu cầu)

Dữ liệu gửi đi ở định dạng JSON theo cấu trúc `UserAddRequest`. Dưới đây là các trường phổ biến (tùy thuộc vào định nghĩa schema cụ thể):

```json
{
  "email": "user@example.com",
  "password": "SecurePassword123!",
  "full_name": "Nguyễn Văn A",
  "phone_number": "0901234567"
}
```

> **Lưu ý:** Các trường bắt buộc sẽ được kiểm tra bởi lớp `validation.UserValidation`. Nếu thiếu trường bắt buộc, hệ thống sẽ trả về lỗi 400.

---

## 📤 Phản Hồi (Response)

### Trường Hợp Thành Công (Success - 200 OK)

Hệ thống trả về thông tin chi tiết của người dùng vừa được tạo.

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "id": "uuid-v123-xyz",
    "email": "user@example.com",
    "full_name": "Nguyễn Văn A",
    "created_at": "2026-01-01T10:00:00Z",
    "status": "active"
  }
}
```

### Trường Hợp Thất Bại (Error)

Nếu xảy ra lỗi (ví dụ: Email đã tồn tại, mật khẩu yếu, định dạng sai), hệ thống sẽ trả về mã lỗi tương ứng.

```json
{
  "code": 400,
  "message": "not_allow",
  "data": null
}
```

---

## 🛠️ Ví Dụ Thực Hành (CURL Example)

Dưới đây là lệnh mẫu để gọi API từ dòng lệnh (Terminal/CMD):

```bash
curl --location 'https://api.clawathon.com/identity/v1/user?lang=vi' \
--header 'Content-Type: application/json' \
--data '{
    "email": "lich.tv@example.com",
    "password": "StrongPass#2026",
    "full_name": "Lich TV",
    "phone_number": "0912345678"
}'
```

---

## 🔄 Luồng Xử Lý (Sequence Diagram)

Biểu đồ dưới đây mô tả luồng xử lý nội bộ khi API nhận được yêu cầu tạo người dùng mới.

```mermaid
sequenceDiagram
    participant Client as Client (App/Web)
    participant Handler as UserHandler (HTTP)
    participant Validator as UserValidation
    participant Service as UserService
    participant DB as Database

    Client->>Handler: POST /identity/v1/user
    activate Handler
    
    Handler->>Handler: Bind JSON Request
    alt JSON Invalid
        Handler-->>Client: Return 400 Bad Request
        deactivate Handler
    else JSON Valid
        Handler->>Validator: Validate Data
        activate Validator
        alt Validation Failed
            Validator-->>Handler: Return Error
            Handler-->>Client: Return 400 Validation Error
            deactivate Validator
        else Validation Passed
            Validator-->>Handler: OK
            deactivate Validator
            
            Handler->>Service: Create User
            activate Service
            Service->>DB: Insert Record
            activate DB
            DB-->>Service: Commit Success
            deactivate DB
            Service-->>Handler: User Object
            deactivate Service
            
            Handler-->>Client: Return 200 + User Detail
        end
    end
    deactivate Handler
```

### Giải thích luồng hoạt động:

1.  **Tiếp nhận:** `UserHandler` nhận yêu cầu HTTP từ Client.
2.  **Gắn kết (Bind):** Chuyển đổi dữ liệu JSON thành object Go. Nếu thất bại, trả lỗi ngay lập tức.
3.  **Kiểm tra (Validate):** Lớp `UserValidation` đảm bảo dữ liệu hợp lệ (email đúng định dạng, mật khẩu đủ mạnh...).
4.  **Xử lý nghiệp vụ (Service):** Nếu qua bước validate, `UserService` sẽ thực hiện logic lưu trữ vào Database.
5.  **Trả kết quả:** Trả về thông tin người dùng nếu thành công hoặc mã lỗi nếu thất bại.

---

## ⚠️ Lưu Ý Quan Trọng (Notes & Best Practices)

1.  **Bảo Mật:** Mật khẩu (`password`) nên được mã hóa (hash) trước khi lưu vào database. API chỉ nhận plaintext password từ client để hash bên server.
2.  **Ngôn ngữ:** Tham số `lang` giúp hệ thống trả về thông báo lỗi phù hợp với người dùng (Tiếng Việt hoặc Tiếng Anh).
3.  **Idempotency:** API tạo người dùng không đảm bảo idempotent. Gọi nhiều lần với cùng dữ liệu có thể tạo nhiều tài khoản trùng lặp nếu không có ràng buộc unique trên Database (thường là Unique Index trên `email`).
4.  **Xác thực:** Sau khi tạo tài khoản thành công, hệ thống thường sẽ yêu cầu người dùng xác minh email hoặc đăng nhập lại để lấy Token.