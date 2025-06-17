# PanicShield Back-End API Documentation

## 공통 사항
- **Base URL**: `/api`
- 모든 요청/응답 바디는 `application/json`
- **공통 응답 포맷**  
  - 성공:
    ```json
    {
      "code": 0,
      "message": "OK",
      "data": { ... }
    }
    ```
  - 오류:
    ```json
    {
      "code": 1001,
      "message": "Error message"
    }
    ```
- **인증 헤더**:  
  ```
  Authorization: Bearer <access_token>
  ```

---

## 1. 인증·인가 (Auth)

### 1.1 회원가입
- **메서드**: `POST`
- **URL**: `/api/auth/register`
- **Request Body**:
  ```json
  {
    "username": "jihyo",
    "password": "secret123",
    "phone_number": "01012341234",
    "speaking_style": "반말",
    "tone": "유머"
  }
  ```
- **Responses**:
  - `201 Created`
    ```json
    {
      "code": 0,
      "message": "User registered",
      "data": { "user_id": 1 }
    }
    ```
  - `400 Bad Request` 입력 검증 오류
  - `409 Conflict` 중복 사용자명/전화번호

---

### 1.2 로그인
- **메서드**: `POST`
- **URL**: `/api/auth/signin`
- **Request Body**:
  ```json
  {
    "username": "jihyo",
    "password": "secret123"
  }
  ```
- **Responses**:
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
  - `400 Bad Request` 입력 검증 오류
  - `401 Unauthorized` 자격 증명 실패

---

### 1.3 전화번호 인증
- **메서드**: `POST`
- **URL**: `/api/auth/verify-phone`
- **Request Body**:
  ```json
  {
    "user_id": 1,
    "code": "123456"
  }
  ```
- **Responses**:
  - `200 OK`
    ```json
    { "code": 0, "message": "Phone number verified" }
    ```
  - `400 Bad Request` 입력 검증 오류
  - `401 Unauthorized` 코드 불일치

---

### 1.4 토큰 갱신
- **메서드**: `POST`
- **URL**: `/api/auth/refresh`
- **Request Body**:
  ```json
  { "refresh_token": "<jwt_refresh>" }
  ```
- **Responses**:
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
  - `400 Bad Request` 입력 검증 오류
  - `401 Unauthorized` 토큰 유효성 실패

---

## 2. 유저 프로필 (User)

### 2.1 내 프로필 조회
- **메서드**: `GET`
- **URL**: `/api/users/me`
- **인증 필요**: 예
- **Responses**:
  - `200 OK`
    ```json
    {
      "code": 0,
      "message": "OK",
      "data": {
        "id": 1,
        "username": "jihyo",
        "phone_number": "01012341234",
        "speaking_style": "반말",
        "tone": "유머",
        "verified": true,
        "created_at": "2025-06-18T10:00:00Z"
      }
    }
    ```

### 2.2 내 프로필 수정
- **메서드**: `PUT`
- **URL**: `/api/users/me`
- **인증 필요**: 예
- **Request Body**: (예시)
  ```json
  {
    "speaking_style": "존댓말",
    "tone": "진지함"
  }
  ```
- **Responses**:
  - `200 OK`: 업데이트된 프로필 반환
  - `400 Bad Request`: 입력 검증 오류

### 2.3 회원 탈퇴
- **메서드**: `DELETE`
- **URL**: `/api/users/me`
- **인증 필요**: 예
- **Responses**:
  - `200 OK`
    ```json
    { "code": 0, "message": "User deleted" }
    ```

---

## 3. 관심사·세부 관심사 (Interest)

### 3.1 전체 관심사 조회
- **메서드**: `GET`
- **URL**: `/api/interests`
- **Responses**:
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

### 3.2 세부 관심사 조회
- **메서드**: `GET`
- **URL**: `/api/interests/{interest_id}/subs`
- **Responses**:
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

### 3.3 내 관심사 등록/해제
- **메서드**: `POST`
- **URL**: `/api/users/me/interests`
- **Request Body**:
  ```json
  { "interest_id": 2 }
  ```
- **Responses**:
  - `200 OK`  
    ```json
    { "code": 0, "message": "Interest added" }
    ```
- **메서드**: `DELETE`
- **URL**: `/api/users/me/interests/{interest_id}`
- **Responses**:
  - `200 OK`  
    ```json
    { "code": 0, "message": "Interest removed" }
    ```

---

## 4. 챗봇 (Chat)

### 4.1 대화 요청
- **메서드**: `POST`
- **URL**: `/api/chat`
- **인증 필요**: 예
- **Request Body**:
  ```json
  { "message": "오늘 기분이 어때?" }
  ```
- **Responses**:
  ```json
  {
    "code": 0,
    "message": "OK",
    "data": { "reply": "안녕하세요! 오늘 기분이 어떤가요?" }
  }
  ```

---

## 5. 바이탈 (Vital Signs)

### 5.1 바이탈 기록 등록
- **메서드**: `POST`
- **URL**: `/api/vitals`
- **인증 필요**: 예
- **Request Body**:
  ```json
  {
    "heart_rate": 72,
    "breath_rate": 16,
    "stress_level": 30
  }
  ```
- **Responses**:
  - `201 Created`
    ```json
    { "code": 0, "message": "Vital record created" }
    ```

### 5.2 바이탈 기록 조회
- **메서드**: `GET`
- **URL**: `/api/vitals?user_id={user_id}`
- **Responses**:
  ```json
  {
    "code": 0,
    "message": "OK",
    "data": [
      {
        "id": 1,
        "user_id": 1,
        "heart_rate": 72,
        "breath_rate": 16,
        "stress_level": 30,
        "measured_at": "2025-06-18T10:05:00Z"
      }
    ]
  }
  ```

---

## 6. 공황가이드 (Panic Guides)

### 6.1 전체 가이드 조회
- **메서드**: `GET`
- **URL**: `/api/panic-guides`
- **Responses**:
  ```json
  {
    "code": 0,
    "message": "OK",
    "data": [
      { "id": 1, "title": "심호흡 연습", "description": "깊게 숨을 들이쉬고 천천히 내쉬세요." }
    ]
  }
  ```

### 6.2 가이드 등록
- **메서드**: `POST`
- **URL**: `/api/panic-guides`
- **인증 필요**: 예
- **Request Body**:
  ```json
  {
    "title": "심호흡 연습",
    "description": "깊게 숨을 들이쉬고 천천히 내쉬세요."
  }
  ```
- **Responses**:
  - `201 Created`
    ```json
    {
      "code": 0,
      "message": "Panic guide created",
      "data": {
        "id": 2,
        "title": "심호흡 연습",
        "description": "깊게 숨을 들이쉬고 천천히 내쉬세요."
      }
    }
    ```

### 6.3 즐겨찾기 추가
- **메서드**: `POST`
- **URL**: `/api/panic-guides/bookmark`
- **인증 필요**: 예
- **Request Body**:
  ```json
  { "user_id": 1, "panic_guide_id": 2 }
  ```
- **Responses**:
  ```json
  { "code": 0, "message": "Bookmark added" }
  ```

### 6.4 내 즐겨찾기 조회
- **메서드**: `GET`
- **URL**: `/api/panic-guides/bookmarks?user_id={user_id}`
- **Responses**:
  ```json
  {
    "code": 0,
    "message": "OK",
    "data": [
      { "id": 2, "title": "심호흡 연습", "description": "깊게 숨을 들이쉬고 천천히 내쉬세요." }
    ]
  }
  ```

---

```

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