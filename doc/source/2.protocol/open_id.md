# OIDC

OpenID Connect（OIDC）是一个基于OAuth 2.0协议的身份认证协议，它允许客户端应用程序通过第三方身份提供者（Identity Provider, IdP）验证用户身份并获取用户的基本信息。OpenID Connect建立在OAuth 2.0的授权框架之上，但它专注于身份认证，而OAuth 2.0主要用于授权。

## OpenID 协议介绍

OpenID 是一个基于 OAuth2 的身份验证协议，允许用户使用单一身份（如 Google、Facebook 等）登录多个网站或应用，而无需为每个网站创建独立的账号。OpenID 的核心目标是实现**单点登录（SSO, Single Sign-On）**。

OpenID Connect（OIDC）是 OpenID 的扩展，基于 OAuth2 协议，增加了身份验证功能。它通过 ID Token（JWT 格式）传递用户身份信息。

---

### OpenID Connect 的核心角色

1. **用户（End-User）**：拥有身份信息的用户。
2. **客户端（Client）**：请求用户身份信息的第三方应用。
3. **身份提供者（Identity Provider, IdP）**：负责验证用户身份并颁发 ID Token。
4. **授权服务器（Authorization Server）**：与身份提供者集成，负责颁发访问令牌和 ID Token。

---

### OpenID Connect 的核心概念

1. **ID Token**：一个 JWT（JSON Web Token），包含用户的身份信息（如用户 ID、姓名、邮箱等）。
2. **UserInfo Endpoint**：一个 API 端点，客户端可以使用访问令牌获取用户的详细信息。
3. **Discovery Document**：一个 JSON 文件，包含 OpenID 服务的配置信息（如端点 URL、支持的算法等）。

---

### OpenID Connect 的授权流程

OpenID Connect 基于 OAuth2 的授权码模式（Authorization Code Flow），并增加了身份验证功能。

#### 授权码模式（Authorization Code Flow）

**步骤：**

1. **用户访问客户端**：用户通过浏览器访问客户端应用。
2. **重定向到授权服务器**：客户端将用户重定向到授权服务器，请求授权和身份验证。
3. **用户登录并授权**：用户在授权服务器上登录并授权客户端访问其身份信息。
4. **授权码返回**：授权服务器将授权码通过重定向返回给客户端。
5. **客户端请求令牌**：客户端使用授权码向授权服务器请求访问令牌和 ID Token。
6. **颁发令牌**：授权服务器验证授权码并颁发访问令牌和 ID Token。
7. **验证 ID Token**：客户端验证 ID Token 的签名和内容。
8. **访问用户信息**：客户端使用访问令牌向 UserInfo Endpoint 请求用户的详细信息。

---

### 示例

#### 示例场景

假设用户使用 Google 登录一个第三方应用（客户端）。

1. **用户访问客户端**：用户打开第三方应用，点击“使用 Google 登录”。
2. **重定向到 Google 授权服务器**：客户端将用户重定向到 Google 的授权页面。
3. **用户登录并授权**：用户在 Google 上登录，并同意授权客户端访问其基本信息。
4. **授权码返回**：Google 将授权码通过重定向返回给客户端。
5. **客户端请求令牌**：客户端使用授权码向 Google 请求访问令牌和 ID Token。
6. **颁发令牌**：Google 验证授权码并颁发访问令牌和 ID Token。
7. **验证 ID Token**：客户端验证 ID Token 的签名和内容，确保其来自 Google。
8. **访问用户信息**：客户端使用访问令牌向 Google 的 UserInfo Endpoint 请求用户的详细信息（如姓名、邮箱等）。

---

#### 流程图

```plaintext
+--------+                                   +---------------+
|        |                                   |               |
|  User  |                                   |  Authorization|
|        |                                   |     Server    |
|        |                                   |  (Google IdP) |
+---|----+                                   +-------|-------+
    |                                                 |
    | (1) User Accesses Client                        |
    |------------------------------------------------>|
    |                                                 |
    | (2) Redirect to Authorization Server            |
    |<------------------------------------------------|
    |                                                 |
    | (3) User Logs in and Authorizes Client          |
    |------------------------------------------------>|
    |                                                 |
    | (4) Authorization Code Returned                 |
    |<------------------------------------------------|
    |                                                 |
    | (5) Request Access Token and ID Token           |
    |------------------------------------------------>|
    |                                                 |
    | (6) Access Token and ID Token Issued            |
    |<------------------------------------------------|
    |                                                 |
    | (7) Validate ID Token                           |
    |                                                 |
    | (8) Request UserInfo                            |
    |------------------------------------------------>|
    |                                                 |
    | (9) UserInfo Returned                           |
    |<------------------------------------------------|
    |                                                 |
+---|----+                                   +-------|-------+
|        |                                   |               |
| Client |                                   |  Resource     |
|        |                                   |    Server     |
+--------+                                   +---------------+
```

---

### ID Token 示例

ID Token 是一个 JWT，包含以下信息：

```json
{
  "iss": "https://accounts.google.com", // 签发者
  "sub": "1234567890", // 用户唯一标识
  "aud": "client-id", // 客户端 ID
  "exp": 1672444800, // 过期时间
  "iat": 1672441200, // 签发时间
  "name": "John Doe", // 用户姓名
  "email": "john.doe@example.com" // 用户邮箱
}
```

---

### UserInfo Endpoint 示例

客户端使用访问令牌向 UserInfo Endpoint 请求用户信息：

**请求：**
```http
GET /userinfo HTTP/1.1
Host: https://accounts.google.com
Authorization: Bearer <access_token>
```

**响应：**
```json
{
  "sub": "1234567890",
  "name": "John Doe",
  "email": "john.doe@example.com",
  "picture": "https://example.com/john.jpg"
}
```

---

### 总结

OpenID Connect 是基于 OAuth2 的身份验证协议，通过 ID Token 和 UserInfo Endpoint 提供用户身份信息。它简化了单点登录的实现，广泛应用于现代互联网服务。通过 OpenID Connect，用户可以使用单一身份登录多个应用，而无需重复注册和登录。