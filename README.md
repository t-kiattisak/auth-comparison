# 🔐 Auth Comparison in Go (JWT, PASETO, Sessions, Cookies)

## ✅ Implemented Methods

| Method            | Description                                      | Endpoint Example                 |
| ----------------- | ------------------------------------------------ | -------------------------------- |
| JWT (Bearer)      | JSON Web Token ที่ client ถือ token เอง          | `/login` → `/me`                 |
| PASETO (Bearer)   | Token แบบปลอดภัยกว่า JWT                         | `/paseto-login` → `/me-paseto`   |
| Sessions (Cookie) | ใช้ session store ฝั่ง server + sessionId cookie | `/session-login` → `/me-session` |
| JWT (Cookie)      | เก็บ JWT ใน HTTPOnly Cookie                      | `/cookie-login` → `/me-cookie`   |

---

## ⚙️ Tech Stack

- Go + Fiber
- Clean Architecture
- JWT (`github.com/golang-jwt/jwt/v5`)
- PASETO (`github.com/o1egl/paseto`)
- Fiber Session Middleware
