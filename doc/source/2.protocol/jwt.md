# JWT

[JSON Web Tokens](https://jwt.io/) are an open, industry standard RFC 7519 method for representing claims securely between two parties.

JWT (Json Web Token) 是一种在双方之间传输的紧凑的，URL安全的令牌表示方式
JWT 表示为编码的 JSON 对象， 可用于 JWS 的的负载， 或者 JWE 的明文
签过名的 JWT 称 JWS(Json Web Signature)
加过密的 JWT 称 JWE(Json Web Encryption)

最终的 JWS 形为 ``` ${Header}.${Payload}.${Signature} ```

它包含三部分

## 1. Header 头部
表示基本算法和格式的 JSON 对象, 例如

```
{
  "alg": "HS256", //签名算法是 HMAC SHA256
  "typ": "JWT" // token 格式是 JWT
}
```

## 2. Payload 正文
以一个 JSON 对象表示 JWT 中所包含的内容, 字段可自由定义, 以下 7 个字段是官方定义的

* iss (issuer)：签发人
* exp (expiration time)：过期时间
* sub (subject)：主题
* aud (audience)：受众
* nbf (Not Before)：生效时间
* iat (Issued At)：签发时间
* jti (JWT ID)：编号

我们也可以加上自己定义的内容, 例如

```
{
"orgId": "$orgId",
"userId": "$userId",
"roles": ["admin"],
"scopes": ["create", "write"]
}
```

## 3. Signature 签名
以头部所声明的算法将头部和正文部分的内容进行签名, 防止被人篡改.

```
HMACSHA256(base64UrlEncode(header) + base64UrlEncode(payload), secret) 
```

## 实例
写一个简单的例子来演示一下 JWT的编码和解码


```java

import io.jsonwebtoken.Claims;
import io.jsonwebtoken.JwtBuilder;
import io.jsonwebtoken.Jwts;
import io.jsonwebtoken.SignatureAlgorithm;
import io.jsonwebtoken.io.Encoders;
import io.jsonwebtoken.lang.Collections;
import lombok.extern.slf4j.Slf4j;

import javax.crypto.spec.SecretKeySpec;
import javax.xml.bind.DatatypeConverter;
import java.security.Key;
import java.util.Arrays;
import java.util.Date;
import java.util.HashMap;
import java.util.Map;
import java.util.UUID;

/**
 * @Author: Walter Fan
 * @Date: 16/6/2019, Sun
 **/
@Slf4j
public class JwtUtil {
    private static final Long ONE_HOUR_MS = 60 * 60 * 100L;

    public static String createJws(String subject, Map<String, Object> claims, long liveSeconds, String apiKey) {
        SignatureAlgorithm signatureAlgorithm = SignatureAlgorithm.HS256;
        byte[] apiKeySecretBytes = DatatypeConverter.parseBase64Binary(apiKey);
        Key signingKey = new SecretKeySpec(apiKeySecretBytes, signatureAlgorithm.getJcaName());


        long expMillis = System.currentTimeMillis() + liveSeconds * 1000;
        Date expireDate = new Date(expMillis);

        JwtBuilder jwtBuilder = Jwts.builder()
                .setId(UUID.randomUUID().toString())
                .setSubject(subject)
                .setIssuedAt(new Date())
                .setExpiration(expireDate)
                .signWith(signingKey);

        if(!Collections.isEmpty(claims)) {
            claims.entrySet().stream().forEach( x -> jwtBuilder.claim(x.getKey(), x.getValue()));
        }

        return jwtBuilder.compact();
    }


    public static Claims parseJws(String compactJws, String apiKey) {
        SignatureAlgorithm signatureAlgorithm = SignatureAlgorithm.HS256;
        byte[] apiKeySecretBytes = DatatypeConverter.parseBase64Binary(apiKey);
        Key signingKey = new SecretKeySpec(apiKeySecretBytes, signatureAlgorithm.getJcaName());

        Claims ret = Jwts.parser()
                .setSigningKey(signingKey)
                .parseClaimsJws(compactJws)
                .getBody();
        return ret;

    }


    public static void main(String[] args) {
        String proverb = "A journey of a thousand miles begins with a single step";
        String apiKey = Encoders.BASE64.encode(proverb.getBytes());
        Map<String, Object> map = new HashMap<>();
        map.put("orgId", UUID.randomUUID());
        map.put("userId", UUID.randomUUID());
        map.put("roles", Arrays.asList("admin"));
        map.put("scopes", Arrays.asList("read", "write"));


        String jws= createJws("potato", map, 300, apiKey);

        log.info("jws = {}", jws);

        Claims claims =  parseJws(jws, apiKey);
        claims.entrySet().stream().forEach(x ->  log.info("{} = {}", x.getKey(), x.getValue()));


    }
}


```

输出如下

```
# 编码过的内容
JwtUtil - jws = eyJhbGciOiJIUzM4NCJ9.eyJqdGkiOiJkOTJkNGRlZC00MjVjLTRhMGItYjVkZi05OGYzN2FiNDVmMzkiLCJzdWIiOiJwb3RhdG8iLCJpYXQiOjE1NjA2NjgwMjIsImV4cCI6MTU2MDY2ODMyMiwicm9sZXMiOlsiYWRtaW4iXSwic2NvcGVzIjpbInJlYWQiLCJ3cml0ZSJdLCJ1c2VySWQiOiIyNzYwMWE1Zi0zZDllLTQ3NWMtOGQzNS0wYTdlZGNlMzZhOGYiLCJvcmdJZCI6IjQ2ZDZhYjM1LTM5MmUtNDE4OS1iMGJjLTk3YzhkMzM2MzA2YSJ9.XydcyqdZD18Kgt0YqpOD_vD_NreHtNh_kyW8ZqxfWbhT4iNYOquQxsxyE8b2Xrul

# 解码过的内容
JwtUtil - jti = d92d4ded-425c-4a0b-b5df-98f37ab45f39
JwtUtil - sub = potato
JwtUtil - iat = 1560668022
JwtUtil - exp = 1560668322
JwtUtil - roles = [admin]
JwtUtil - scopes = [read, write]
JwtUtil - userId = 27601a5f-3d9e-475c-8d35-0a7edce36a8f
JwtUtil - orgId = 46d6ab35-392e-4189-b0bc-97c8d336306a
```

在 Web Service 实际应用中, 经常把 JWT 放到 HTTP 的 Authorization 头域中, 客户端从认证服务器那里申请一个 token , 此 token 也就是认证服务器所签发的 JWT, 用这个 token 来访问资源.
资源服务器会校验这个 token , 也会到认证服务器那里会验证。

## 利用 JWT 进行认证的过程

### 简单流程
* client 提供 username 和 password 向 server 器申请 auth token
* server 校验通过后，返回一个 auth token
* client 保存 auth token, 之后在后续的 API 请求中带着这个 token(HTTP 协议用 Authorization 头)
* server 从 Authorization 头中解码这个 token 并校验这个token


## 参考资料
* [https://jwt.io/](https://jwt.io/)
