# Tài Liệu API - Module Identity (Quản Lý Người Dùng)

## 1. Tổng Quan

Module **Identity** chịu trách nhiệm quản lý thông tin định danh của người dùng trong hệ thống. Tài liệu này mô tả chi tiết các điểm cuối (endpoints) liên quan đến việc tạo và quản lý người dùng.

Tài liệu được thiết kế để phù hợp cho cả:
*   **Người dùng kỹ thuật (Developers/QA):** Cần biết cấu trúc dữ liệu, phương thức HTTP, và luồng xử lý.
*   **Người dùng phi kỹ thuật (Product Managers/Business):** Cần hiểu chức năng nghiệp vụ mà API cung cấp.

---

## 2. API Chi Tiết: Tạo Người Dùng Mới (Create User)

### 2.1. Mô Tả Nghiệp Vụ
API này cho phép đăng ký một tài khoản người dùng mới vào hệ thống. Khi gọi API này, hệ thống sẽ kiểm tra tính hợp lệ của thông tin đầu vào (như email, mật khẩu) và lưu trữ thông tin người dùng vào cơ sở dữ liệu nếu mọi thứ chính xác.

*   **Mục đích:** Thêm thành viên mới vào hệ thống.
*   **Đối tượng sử dụng:** Ứng dụng di động, Web Frontend, hoặc các hệ thống bên ngoài cần tích hợp thêm user.

### 2.2. Thông Tin Kỹ Thuật

| Thuộc tính | Giá trị |
| :--- | :--- |
| **Tên API** | `Add User` |
| **Phương thức** | `POST` |
| **Đường dẫn (Path)** | `/identity/v1/user` |
| **Accept** | `application/json` |
| **Produce** | `application/json` |
| **Tags** | `identity` |

### 2.3. Tham Số Đầu Vào (Request)

#### Query Parameters (Tham số truy vấn)
Các tham số được gửi kèm theo đường dẫn URL.

| Tên | Kiểu | Bắt buộc | Mô tả | Giá trị mặc định |
| :--- | :--- | :--- | :--- | :--- |
| `lang` | string | Không | Ngôn ngữ hiển thị thông báo lỗi/thành công | `vi` |
| | | | *Giá trị chấp nhận:* `en`, `vi` | |

#### Request Body (Nội dung yêu cầu)
Dữ liệu được gửi dưới dạng JSON trong thân yêu cầu.

| Trường | Kiểu | Bắt buộc | Mô tả | Ví dụ |
| :--- | :--- | :--- | :--- | :--- |
| `email` | string | Có | Địa chỉ email của người dùng | `user@example.com` |
| `password` | string | Có | Mật khẩu đăng nhập | `Password123!` |
| `full_name` | string | Có | Họ và tên đầy đủ | `Nguyễn Văn A` |
| `phone` | string | Không | Số điện thoại liên hệ | `0901234567` |

> **Lưu ý:** Các trường bắt buộc phải tuân thủ quy tắc xác thực (validation) như định dạng email, độ dài mật khẩu tối thiểu.

### 2.4. Phản Hồi Đầu Ra (Response)

#### Trường hợp Thành Công (Status Code: 200 OK)
Trả về thông tin chi tiết của người dùng vừa được tạo.

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "id": "uuid-hoặc-id-so",
    "email": "user@example.com",
    "full_name": "Nguyễn Văn A",
    "created_at": "2023-10-27T10:00:00Z"
  }
}
```

#### Trường hợp Thất Bại (Status Code: 400 Bad Request)
Xảy ra khi dữ liệu không hợp lệ (ví dụ: email đã tồn tại, mật khẩu yếu).

```json
{
  "code": 400,
  "message": "not_allow",
  "error": "Email đã được sử dụng"
}
```

### 2.5. Ví Dụ CURL

Dưới đây là lệnh mẫu để gọi API từ dòng lệnh (Terminal):

```bash
curl --location 'https://api.example.com/identity/v1/user?lang=vi' \
--header 'Content-Type: application/json' \
--data '{
    "email": "lich.tv@example.com",
    "password": "SecurePass123!",
    "full_name": "Lịch TV",
    "phone": "0912345678"
}'
```

### 2.6. Sơ Đồ Luồng Xử Lý (Sequence Diagram)

Sơ đồ dưới đây minh họa cách hệ thống xử lý yêu cầu từ khi nhận được request cho đến khi trả về kết quả.

```mermaid
sequenceDiagram
    participant Client as Client (App/Web)
    participant Handler as UserHandler (HTTP)
    participant Validator as UserValidation
    participant Service as UserService
    participant DB as Database

    Client->>Handler: POST /identity/v1/user (JSON)
    activate Handler
    
    Handler->>Handler: Bind JSON Data
    alt Binding Failed
        Handler-->>Client: 400 Bad Request
        deactivate Handler
    end

    Handler->>Validator: Validate Input (Add)
    activate Validator
    alt Validation Failed
        Validator-->>Handler: Error Message
        Handler-->>Client: 400 Validation Error
        deactivate Validator
        deactivate Handler
    else Validation Passed
        Validator-->>Handler: Valid
        deactivate Validator
        
        Handler->>Service: Create User
        activate Service
        Service->>DB: Insert User Record
        activate DB
        DB-->>Service: Success
        deactivate DB
        Service-->>Handler: User Object
        deactivate Service
        
        Handler-->>Client: 200 OK + User Detail
        deactivate Handler
    end
```

### 2.7. Lưu Ý Quan Trọng

1.  **Bảo mật:** Mật khẩu (`password`) nên được mã hóa trước khi truyền tải qua mạng (HTTPS) và luôn được băm (hash) trước khi lưu vào cơ sở dữ liệu.
2.  **Ngôn ngữ:** Nếu không truyền tham số `lang`, hệ thống mặc định trả về thông báo tiếng Việt (`vi`). Để có trải nghiệm quốc tế, hãy truyền `lang=en`.
3.  **Xác thực:** API này thường dành cho việc đăng ký mới. Nếu người dùng đã tồn tại, hệ thống sẽ trả về lỗi cụ thể thay vì ghi đè thông tin cũ.
4.  **Versioning:** Đường dẫn API bao gồm phiên bản `/v1/`. Trong tương lai, nếu có thay đổi lớn về cấu trúc dữ liệu, phiên bản sẽ được nâng lên `/v2/` để đảm bảo tính tương thích ngược.