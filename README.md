
# PanicShield Back-End API Documentation

## Error Codes
| Code  | HTTP Status | Description                       |
|-------|-------------|-----------------------------------|
| 0     | 200/201     | Success                           |
| 1001  | 400         | Validation error                  |
| 1002  | 401         | Unauthorized (invalid credentials / token) |
| 1003  | 403         | Forbidden                         |
| 1004  | 404         | Not Found                         |
| 1005  | 409         | Conflict (duplicate resource)     |
| 1006  | 500         | Internal server error             |

---

## Common

- **Base URL**: `https://api.example.com/api`
- **Headers**:
  - `Content-Type: application/json`
  - `Authorization: Bearer <access_token>` (required for protected endpoints)

- **Response Envelope**  
  ```json
  {
    "code": <int>,
    "message": "<string>",
    "data": <object|null>
  }
  ```

---

## 1. Authentication

### 1.1 Register
- **POST** `/api/auth/register`
- **Body Parameters**:
  | Name            | Type   | Required | Description                   |
  |-----------------|--------|----------|-------------------------------|
  | `username`      | string | yes      | 2-32 characters, unique       |
  | `password`      | string | yes      | 6+ characters                 |
  | `phone_number`  | string | yes      | e.g. "01012341234", unique    |
  | `speaking_style`| string | yes      | User’s default speaking style |
  | `tone`          | string | yes      | User’s default tone           |

- **cURL Example**:
  ```bash
  curl -X POST https://api.example.com/api/auth/register \
    -H "Content-Type: application/json" \
    -d '{
      "username":"jane",
      "password":"pass1234",
      "phone_number":"01011112222",
      "speaking_style":"formal",
      "tone":"friendly"
    }'
  ```

- **Success (201)**:
  ```json
  {
    "code": 0,
    "message": "User registered",
    "data": { "user_id": 42 }
  }
  ```
- **Error (409 Conflict)**:
  ```json
  {
    "code": 1005,
    "message": "Username already taken"
  }
  ```

---

### 1.2 Sign In
- **POST** `/api/auth/signin`
- **Body Parameters**:
  | Name       | Type   | Required | Description           |
  |------------|--------|----------|-----------------------|
  | `username` | string | yes      |                       |
  | `password` | string | yes      |                       |

- **cURL Example**:
  ```bash
  curl -X POST https://api.example.com/api/auth/signin \
    -H "Content-Type: application/json" \
    -d '{"username":"jane","password":"pass1234"}'
  ```

- **Success (200)**:
  ```json
  {
    "code": 0,
    "message": "Login successful",
    "data": {
      "access_token":"<jwt_access>",
      "refresh_token":"<jwt_refresh>"
    }
  }
  ```
- **Error (401 Unauthorized)**:
  ```json
  {
    "code":1002,
    "message":"Invalid credentials"
  }
  ```

---

### 1.3 Verify Phone
- **POST** `/api/auth/verify-phone`
- **Body Parameters**:
  | Name    | Type   | Required | Description         |
  |---------|--------|----------|---------------------|
  | `user_id` | uint | yes      | Returned from register |
  | `code`  | string | yes (6) | OTP code from SMS    |

- **cURL Example**:
  ```bash
  curl -X POST https://api.example.com/api/auth/verify-phone \
    -H "Authorization: Bearer <token>" \
    -d '{"user_id":42,"code":"123456"}'
  ```

- **Success (200)**:
  ```json
  { "code":0,"message":"Phone number verified","data":null }
  ```

---

### 1.4 Refresh Token
- **POST** `/api/auth/refresh`
- **Body Parameters**:
  | Name           | Type   | Required | Description     |
  |----------------|--------|----------|-----------------|
  | `refresh_token`| string | yes      | Refresh JWT     |

- **cURL Example**:
  ```bash
  curl -X POST https://api.example.com/api/auth/refresh \
    -H "Content-Type: application/json" \
    -d '{"refresh_token":"<jwt_refresh>"}'
  ```

- **Success (200)**:
  ```json
  {
    "code":0,
    "message":"Token refreshed",
    "data":{"access_token":"<new>","refresh_token":"<new>"}
  }
  ```

---

## 2. User Profile

### 2.1 Get My Profile
- **GET** `/api/users/me`
- **Query Parameters**: none
- **Headers**: Authorization required
- **cURL Example**:
  ```bash
  curl https://api.example.com/api/users/me \
    -H "Authorization: Bearer <token>"
  ```

- **Success (200)**:
  ```json
  {
    "code":0,
    "message":"OK",
    "data": {
      "id":42,
      "username":"jane",
      "phone_number":"01011112222",
      "speaking_style":"formal",
      "tone":"friendly",
      "verified":true,
      "created_at":"2025-06-18T10:00:00Z"
    }
  }
  ```

### 2.2 Update My Profile
- **PUT** `/api/users/me`
- **Body Parameters**:
  | Name            | Type   | Required | Description             |
  |-----------------|--------|----------|-------------------------|
  | `speaking_style`| string | no       |                         |
  | `tone`          | string | no       |                         |

- **Success (200)**:
  ```json
  {
    "code":0,
    "message":"Profile updated",
    "data":{ ... updated user ... }
  }
  ```

### 2.3 Delete My Profile
- **DELETE** `/api/users/me`
- **Success (200)**:
  ```json
  { "code":0, "message":"User deleted","data":null }
  ```

---

## 3. Interests & Sub-Interests

### 3.1 List All Interests
- **GET** `/api/interests`
- **Query Parameters**:
  | Name | Type | Description |
  |------|------|-------------|
  | `page` | int | optional, default=1 |
  | `size` | int | optional, default=20 |

- **Success (200)**:
  ```json
  {
    "code":0,
    "message":"OK",
    "data":{
      "items":[ { "id":1,"name":"운동" }, ... ],
      "page":1,"size":20,"total":5
    }
  }
  ```

### 3.2 List Sub-Interests
- **GET** `/api/interests/{id}/subs`
- **Path Parameters**: `id` (interest ID)
- **Success (200)**:
  ```json
  {
    "code":0,"message":"OK",
    "data":[ { "id":1,"interest_id":1,"name":"축구" }, ... ]
  }
  ```

### 3.3 Add/Remove My Interest
- **POST** `/api/users/me/interests`
  ```json
  { "interest_id":2 }
  ```
- **DELETE** `/api/users/me/interests/{id}`

---

*(Continue similarly for Chat, Vitals, Panic Guides with pagination, detailed params, and cURL examples.)*
# PanicShield API 문서

## 1. 공통 사항

- **Base URL**:
  ```
  https://panicshield.ngrok.dev/api
  ```
- 모든 요청/응답의 `Content-Type`은 `application/json`입니다.
- 보호된 엔드포인트는 `Authorization: Bearer <ACCESS_TOKEN>` 헤더가 필요합니다.

---

## 2. Health Check

```bash
curl -i "https://panicshield.ngrok.dev/api/health"
```
- 서버 상태 확인용 엔드포인트입니다.
- **Response**
  ```json
  { "status": "ok" }
  ```

---

## 3. 인증 (Auth)

### 3.1 회원가입 (Register)
```bash
curl -i -X POST "https://panicshield.ngrok.dev/api/auth/register" \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"password123","phone_number":"01012345678","speaking_style":"반말","tone":"유머"}'
```
#### Request Body
```json
{
  "username": "testuser",
  "password": "password123",
  "phone_number": "01012345678",
  "speaking_style": "반말",
  "tone": "유머"
}
```
#### Response
- `201 Created`
  ```json
  {
    "code": 0,
    "message": "User registered",
    "data": { "user_id": 1 }
  }
  ```
- 오류 코드: 1001(입력 오류), 1005(중복)

---

### 3.2 로그인 (Signin)
```bash
curl -i -X POST "https://panicshield.ngrok.dev/api/auth/signin" \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"password123"}'
```
#### Request Body
```json
{
  "username": "testuser",
  "password": "password123"
}
```
#### Response
- `200 OK`
  ```json
  {
    "code": 0,
    "message": "Login successful",
    "data": {
      "access_token": "<jwt_access>",
      "refresh_token": "<jwt_refresh>"
    }
  }
  ```
- 오류 코드: 1002(인증 실패)

---

### 3.3 토큰 갱신 (Refresh Token)
```bash
curl -i -X POST "https://panicshield.ngrok.dev/api/auth/refresh" \
  -H "Content-Type: application/json" \
  -d '{"refresh_token":"<jwt_refresh>"}'
```
#### Request Body
```json
{ "refresh_token": "<jwt_refresh>" }
```
#### Response
- `200 OK`
  ```json
  {
    "code": 0,
    "message": "Token refreshed",
    "data": {
      "access_token": "<new_access>",
      "refresh_token": "<new_refresh>"
    }
  }
  ```
- 오류 코드: 1002

---

## 4. 유저 (User)

### 4.1 내 프로필 조회 (Get Profile)
```bash
curl -i "https://panicshield.ngrok.dev/api/users/me?user_id=1" \
  -H "Authorization: Bearer <ACCESS_TOKEN>"
```
- 인증 필요
- **Response**
  ```json
  {
    "code": 0,
    "message": "OK",
    "data": {
      "id": 1,
      "username": "testuser",
      "phone_number": "01012345678",
      "speaking_style": "반말",
      "tone": "유머",
      "verified": true,
      "created_at": "2025-06-18T10:00:00Z"
    }
  }
  ```

### 4.2 내 프로필 수정 (Update Profile)
```bash
curl -i -X PUT "https://panicshield.ngrok.dev/api/users/me" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <ACCESS_TOKEN>" \
  -d '{"speaking_style":"존댓말","tone":"진지함"}'
```
- 인증 필요
- **Request Body**
  ```json
  {
    "speaking_style": "존댓말",
    "tone": "진지함"
  }
  ```
- **Response**
  - `200 OK`: 업데이트된 프로필 반환
  - 오류 코드: 1001

### 4.3 회원 탈퇴 (Delete Profile)
```bash
curl -i -X DELETE "https://panicshield.ngrok.dev/api/users/me" \
  -H "Authorization: Bearer <ACCESS_TOKEN>"
```
- 인증 필요
- **Response**
  ```json
  { "code": 0, "message": "User deleted" }
  ```

---

## 5. 관심사 (Interests)

### 5.1 전체 관심사 조회 (List All Interests)
```bash
curl -i "https://panicshield.ngrok.dev/api/interests" \
  -H "Authorization: Bearer <ACCESS_TOKEN>"
```
- **Response**
  ```json
  {
    "code": 0,
    "message": "OK",
    "data": [
      { "id": 1, "name": "운동" },
      { "id": 2, "name": "음악" }
    ]
  }
  ```

### 5.2 세부 관심사 조회 (List Sub-Interests)
```bash
curl -i "https://panicshield.ngrok.dev/api/interests/1/subs" \
  -H "Authorization: Bearer <ACCESS_TOKEN>"
```
- **Response**
  ```json
  {
    "code": 0,
    "message": "OK",
    "data": [
      { "id": 1, "interest_id": 1, "name": "축구" },
      { "id": 2, "interest_id": 1, "name": "달리기" }
    ]
  }
  ```

### 5.3 내 관심사 등록 (Add My Interest)
```bash
curl -i -X POST "https://panicshield.ngrok.dev/api/interests" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <ACCESS_TOKEN>" \
  -d '{"user_id":1,"interest":"독서"}'
```
- **Response**
  ```json
  { "code": 0, "message": "Interest added" }
  ```

### 5.4 내 관심사 전체 조회 (List My Interests)
```bash
curl -i "https://panicshield.ngrok.dev/api/interests" \
  -H "Authorization: Bearer <ACCESS_TOKEN>"
```
- **Response**
  ```json
  {
    "code": 0,
    "message": "OK",
    "data": [
      { "id": 1, "user_id": 1, "interest": "독서" }
    ]
  }
  ```

---

## 6. 바이탈 (Vitals)

### 6.1 바이탈 기록 등록 (Create Vital Record)
```bash
curl -i -X POST "https://panicshield.ngrok.dev/api/vitals" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <ACCESS_TOKEN>" \
  -d '{"user_id":1,"heart_rate":75,"breath_rate":18,"stress_level":30}'
```
- **Response**
  - `201 Created`
    ```json
    { "code": 0, "message": "Vital record created" }
    ```

### 6.2 바이탈 기록 조회 (List Vital Records)
```bash
curl -i "https://panicshield.ngrok.dev/api/vitals?user_id=1" \
  -H "Authorization: Bearer <ACCESS_TOKEN>"
```
- **Response**
  ```json
  {
    "code": 0,
    "message": "OK",
    "data": [
      {
        "id": 1,
        "user_id": 1,
        "heart_rate": 75,
        "breath_rate": 18,
        "stress_level": 30,
        "measured_at": "2025-06-18T10:05:00Z"
      }
    ]
  }
  ```

---

## 7. 챗봇 (Chatbot)

### 7.1 대화 요청 (Chat)
```bash
curl -i -X POST "https://panicshield.ngrok.dev/api/chat" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <ACCESS_TOKEN>" \
  -d '{"user_id":1,"message":"안녕하세요"}'
```
- **Response**
  ```json
  {
    "code": 0,
    "message": "OK",
    "data": { "reply": "안녕하세요! 무엇을 도와드릴까요?" }
  }
  ```

### 7.2 대화 히스토리 조회 (Chat History)
```bash
curl -i "https://panicshield.ngrok.dev/api/chat?user_id=1" \
  -H "Authorization: Bearer <ACCESS_TOKEN>"
```
- **Response**
  ```json
  {
    "code": 0,
    "message": "OK",
    "data": [
      { "role": "user", "message": "안녕하세요" },
      { "role": "bot", "message": "안녕하세요! 무엇을 도와드릴까요?" }
    ]
  }
  ```

---

## 8. 공황가이드 (Panic Guides)

### 8.1 전체 가이드 조회 (List Panic Guides)
```bash
curl -i "https://panicshield.ngrok.dev/api/panic-guides" \
  -H "Authorization: Bearer <ACCESS_TOKEN>"
```
- **Response**
  ```json
  {
    "code": 0,
    "message": "OK",
    "data": [
      { "id": 1, "title": "심호흡", "description": "깊게 숨을 들이마시고 천천히 내쉬세요." }
    ]
  }
  ```

### 8.2 가이드 등록 (Create Panic Guide)
```bash
curl -i -X POST "https://panicshield.ngrok.dev/api/panic-guides" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <ACCESS_TOKEN>" \
  -d '{"user_id":1,"title":"심호흡","description":"깊게 숨을 들이마시고 천천히 내쉬세요."}'
```
- **Response**
  - `201 Created`
    ```json
    {
      "code": 0,
      "message": "Panic guide created",
      "data": {
        "id": 2,
        "title": "심호흡",
        "description": "깊게 숨을 들이마시고 천천히 내쉬세요."
      }
    }
    ```

### 8.3 즐겨찾기 추가 (Bookmark Panic Guide)
```bash
curl -i -X POST "https://panicshield.ngrok.dev/api/panic-guides/bookmark" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <ACCESS_TOKEN>" \
  -d '{"user_id":1,"panic_guide_id":1}'
```
- **Response**
  ```json
  { "code": 0, "message": "Bookmark added" }
  ```

### 8.4 내 즐겨찾기 조회 (List My Bookmarks)
```bash
curl -i "https://panicshield.ngrok.dev/api/panic-guides/bookmarks?user_id=1" \
  -H "Authorization: Bearer <ACCESS_TOKEN>"
```
- **Response**
  ```json
  {
    "code": 0,
    "message": "OK",
    "data": [
      { "id": 1, "title": "심호흡", "description": "깊게 숨을 들이마시고 천천히 내쉬세요." }
    ]
  }
  ```

---

## 9. 에러 코드

| 코드   | HTTP 상태    | 설명                                |
|--------|--------------|-------------------------------------|
| 0      | 200/201      | 성공                                |
| 1001   | 400          | 입력값 오류                         |
| 1002   | 401          | 인증 실패 (토큰/자격증명 불일치)    |
| 1003   | 403          | 권한 없음                           |
| 1004   | 404          | 리소스 없음                         |
| 1005   | 409          | 중복 리소스                         |
| 1006   | 500          | 서버 내부 오류                      |