# PanicShield API Documentation

## Base URL

```
https://panicshield.ngrok.dev/api
```

All endpoints accept and return `application/json`. Protected endpoints require `Authorization: Bearer <ACCESS_TOKEN>` header.

---

## 1. Health Check

```bash
curl -i "$BASE/api/health"
```

**Response**  
```json
{"status":"ok"}
```

---

## 2. Register

```bash
REG_RESP=$(curl -s -X POST "$BASE/api/auth/register" \
  -H "Content-Type: application/json" \
  -d '{
    "username":"testuser",
    "password":"password123",
    "phone_number":"01012345678",
    "speaking_style":"반말",
    "tone":"유머"
  }')
echo "Register response: $REG_RESP"
ACCESS_TOKEN=$(echo "$REG_RESP" | jq -r '.access_token')
REFRESH_TOKEN=$(echo "$REG_RESP" | jq -r '.refresh_token')
echo "Access Token: $ACCESS_TOKEN"
echo "Refresh Token: $REFRESH_TOKEN"
```

**Response**  
```json
{
  "message": "회원가입 성공",
  "access_token": "<ACCESS_TOKEN>",
  "refresh_token": "<REFRESH_TOKEN>",
  "user_id": 1,
  "speaking_style": "반말",
  "tone": "유머"
}
```

---

## 3. Refresh Token

```bash
NEW_TOKEN=$(curl -s -X POST "$BASE/api/auth/refresh" \
  -H "Content-Type: application/json" \
  -d "{\"refresh_token\":\"$REFRESH_TOKEN\"}" \
  | jq -r '.access_token')
echo "New Access Token: $NEW_TOKEN"
```

**Response**  
```json
{"access_token":"<NEW_ACCESS_TOKEN>"}
```

---

## 4. Get Profile

```bash
curl -i "$BASE/api/users/me?user_id=1" \
  -H "Authorization: Bearer $NEW_TOKEN"
```

**Response**  
```json
{
  "id": 1,
  "username": "testuser",
  "phone_number": "01012345678",
  "verified": false,
  "speaking_style": "반말",
  "tone": "유머",
  "created_at": "2025-07-24T02:17:37.007865+09:00"
}
```

---

## 5. Interests

### 5.1 Add Interest

```bash
curl -i -X POST "$BASE/api/interests" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $NEW_TOKEN" \
  -d '{"user_id":1,"interest":"독서"}'
```

### 5.2 List Interests

```bash
curl -i "$BASE/api/interests" \
  -H "Authorization: Bearer $NEW_TOKEN"
```

**Response**  
```json
{"interests":[{"id":1,"name":"독서"}]}
```

---

## 6. Vitals

### 6.1 Create Vital

```bash
curl -i -X POST "$BASE/api/vitals" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $NEW_TOKEN" \
  -d '{"user_id":1,"heart_rate":75,"breath_rate":18,"stress_level":30}'
```

### 6.2 List Vitals

```bash
curl -i "$BASE/api/vitals?user_id=1" \
  -H "Authorization: Bearer $NEW_TOKEN"
```

**Response**  
```json
{
  "data":[
    {
      "id":1,
      "user_id":1,
      "heart_rate":75,
      "breath_rate":18,
      "stress_level":30,
      "measured_at":"2025-07-24T02:17:38.544263+09:00",
      "created_at":"2025-07-24T02:17:38.549801+09:00"
    }
  ]
}
```

---

## 7. Chatbot

### 7.1 Chat

```bash
curl -i -X POST "$BASE/api/chat" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $NEW_TOKEN" \
  -d '{"user_id":1,"message":"안녕하세요"}'
```

**Response**  
```json
{"reply":"안녕하세요! 무엇을 도와드릴까요?"}
```

### 7.2 History

```bash
curl -i "$BASE/api/chat?user_id=1" \
  -H "Authorization: Bearer $NEW_TOKEN"
```

**Response**  
```json
{"data":[{"id":1,"user_id":1,"message":"안녕하세요","reply":"안녕하세요!","created_at":"..."}]}
```

---

## 8. Panic Guides

### 8.1 List Guides

```bash
curl -i "$BASE/api/panic-guides" \
  -H "Authorization: Bearer $NEW_TOKEN"
```

### 8.2 Create Guide

```bash
curl -i -X POST "$BASE/api/panic-guides" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $NEW_TOKEN" \
  -d '{"user_id":1,"title":"심호흡","description":"깊게 숨을 들이마시고 천천히 내쉬세요."}'
```

### 8.3 Bookmark Guide

```bash
curl -i -X POST "$BASE/api/panic-guides/bookmark" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $NEW_TOKEN" \
  -d '{"user_id":1,"panic_guide_id":1}'
```

### 8.4 List Bookmarks

```bash
curl -i "$BASE/api/panic-guides/bookmarks?user_id=1" \
  -H "Authorization: Bearer $NEW_TOKEN"
```
