# Privae Key

## PKCS#1 和 PKCS#8 私钥格式

**PKCS#1** 和 **PKCS#8** 是两种常见的私钥格式，它们各自具有特定的结构和字段。以下是它们的详细结构和字段解释。

---

### **1. PKCS#1 私钥格式结构**

**PKCS#1** 专门用于存储 RSA 密钥，包括 RSA 公钥和私钥。RSA 私钥是一个包含多个整数的结构，它包含了 RSA 算法的核心参数。以下是 **PKCS#1** 格式私钥的结构和主要字段：

#### **PKCS#1 私钥结构**

```text
RSAPrivateKey ::= SEQUENCE {
    version       Version,
    modulus       INTEGER,  -- n
    publicExponent INTEGER, -- e
    privateExponent INTEGER, -- d
    prime1        INTEGER,  -- p
    prime2        INTEGER,  -- q
    exponent1     INTEGER,  -- dP = d % (p - 1)
    exponent2     INTEGER,  -- dQ = d % (q - 1)
    coefficient   INTEGER,  -- qInv = (1 / q) % p
    otherPrimeInfo  SEQUENCE OF OtherPrimeInfo OPTIONAL
}
```

**字段说明**：
1. **version**：版本号，通常为 `0`（表示是 PKCS#1 版本 1 密钥）。在 PKCS#1 中，版本号是可选的，但一般使用版本号 `0`。
2. **modulus (n)**：RSA 模数，表示两个大质数 `p` 和 `q` 的乘积。这是公钥和私钥的一部分。
3. **publicExponent (e)**：公钥指数，通常设置为 `65537`（这是一个常见的值）。
4. **privateExponent (d)**：私钥指数，是解密运算的核心，满足 `e * d ≡ 1 (mod (p-1)(q-1))`。
5. **prime1 (p)** 和 **prime2 (q)**：RSA 私钥的两个质数，这两个数是 RSA 模数 `n` 的因子。
6. **exponent1 (dP)**：`d % (p - 1)`，用于 CRT (Chinese Remainder Theorem) 加速 RSA 解密过程。
7. **exponent2 (dQ)**：`d % (q - 1)`，同样是 CRT 加速的一部分。
8. **coefficient (qInv)**：`(1 / q) % p`，用于加速 RSA 解密的 CRT 参数。
9. **otherPrimeInfo**：这个字段包含额外的质因数信息（如果存在），一般用于多质数的 RSA。

#### **PKCS#1 格式私钥示例**

PEM 编码的 PKCS#1 私钥格式如下：

```
-----BEGIN RSA PRIVATE KEY-----
MIIBOwIBAAJBAKQ2JtnVNNWQxeVAO6O9l0sXxUJbHnvpuJxhHwtPq7zyA5V0JkZB
yHNBqVfEKwz3rHVVqYX+vD4izNnsnn0tfKL0Eu5Ot+yCEzQZInJWeWshQyy8t1ns
DoFgfz7n9Bwy+gj6qZlXWeDgCKlUBn8KNpZKZrGx3g7g0WOpAhO5DqfKDiPHQdfu
L7guO0dsqbD5kwFg5A5zW6CGIcG0EGkqdywVsznXBIt7ct+PqaqympGEu9tIwwTg
Qmf9WwBlaMrb4VKuyt2W7P8l6YrUmRbGVLO5g8n8ryLw+gZycqD/YQg8yveLgKjS
IWlItPeW6/1cGjzJtRsfJmc1mFqghPsbDbyGvtyq9GeXsA9jEX8kj9AeHTeblHkD
5gWxkj0B+9vfnCH1EKOsYQIH5x6rSi6Vrk5ul5yzojKzJlG4YidWfTyD7tQzMD1N
J9gCmzTmuB5ii+mPfEddH+Xq91cIduh5AKTI6ZokfwTwDLOk9kn1u9Fss5Jcswgs
cL+gF5vV5t5aZj70lrDiB1mpk1V69fI=
-----END RSA PRIVATE KEY-----
```

---

### **2. PKCS#8 私钥格式结构**

**PKCS#8** 是一个更加通用的标准，支持多种算法的私钥。与 **PKCS#1** 格式不同，**PKCS#8** 可以用于表示除了 RSA 以外的其他密钥类型（例如 DSA、ECDSA 等）。PKCS#8 格式的私钥包含 **私钥算法标识符** 和 **私钥数据**。

#### **PKCS#8 私钥结构**

```text
PrivateKeyInfo ::= SEQUENCE {
    version         Version,
    privateKeyAlgorithm PrivateKeyAlgorithmIdentifier,
    privateKey      OCTET STRING
}

PrivateKeyAlgorithmIdentifier ::= SEQUENCE {
    algorithm       OBJECT IDENTIFIER,
    parameters      ANY DEFINED BY algorithm OPTIONAL
}

Version ::= INTEGER { v1(0) }

PrivateKey ::= OCTET STRING
```

**字段说明**：
1. **version**：版本号。`v1(0)` 表示 PKCS#8 版本 1，当前标准中大多使用版本号 `0`。
2. **privateKeyAlgorithm**：一个结构，包含密钥算法标识符和可选的参数：
   - **algorithm**：算法标识符，例如 RSA 使用的 `1.2.840.113549.1.1.1`，表示 RSA。
   - **parameters**：可选的参数，通常为空或与算法相关的参数（例如 EC 密钥会有椭圆曲线参数）。
3. **privateKey**：包含实际的私钥数据，经过编码并存储为 OCTET STRING。这些数据是由特定算法（例如 RSA）私钥的实际数字组成。

#### **PKCS#8 格式私钥示例**

PEM 编码的 PKCS#8 格式私钥示例如下：

```
-----BEGIN PRIVATE KEY-----
MIIEvgIBADANBgkqhkiG9w0BAQEFAASCATwggE4AgEAAkEAmZo0BqNzyP7kEo9tA
zyvF2Y+Xh+P0JHcQpBHzmLgf8yz0uNjmf9QOeug9oFgMZm0HhOYx1x/x6/ZTQrkM
2obf4HJ9gKcFnFvK9o0mPSYPgN3gQhFqZEnK0U9TxTj9p8pLTKzYr34ml66y5R5f
yXhANMFc1FfcY7Ev3gpgjrZT1hGww3v/8tO+P2WFTKmlfISIn+d5cVA5S+Mb06hrk
FHfR1VqxJlqOW3wMN41w9d8M5leudVgpfxI5/Fl3o5Gq3lX3GoIWZyFO4I4wftdM
3Rqa8Xz9mQDbVrr73d5JwYQK+1AGKbFsdvlydhplBbgI1mRPGxL8ocT7z9T8Ekwr
j4sBghld+ZZ5vVfqMBRY9d5pq7tFEJbxmwtr9kTuaPzO3IEm4Bfd0VrgTGHsgALH
Zc00a6iHtVK84BswS+a82N3x9Q9t4Az9l9RMRdMv60pM5PlBhhIhAOwP/kl5m6Wz
2eWyLS3G5WauK8Sjd+Z5ob0ppZlIojN+gtK8CHc3RPsmXybgzVzSk6OQdx5A29LO
5pfxwR4LfgS2wCKOZtPV0h9g2JkDHHlnZ4fgP3kbysb3ihrE8yNsf50ZmwX7MQO6
ax9hVg==
-----END PRIVATE KEY-----
```

---

### **PKCS#1 和 PKCS#8 的关键区别**

| **字段/特性**         | **PKCS#1**                                       | **PKCS#8**                                      |
|--------------------|-------------------------------------------------|------------------------------------------------|
| **支持的算法**        | 仅支持 RSA 算法                                   | 支持多种算法（如 RSA、DSA、ECDSA 等）            |
| **结构**             | 直接存储 RSA 密钥的各个部分，如 `n`、`e`、`d` 等      | 包含一个 **算法标识符** 和 **私钥数据**，格式更通用     |
| **用途**             | 专门为 RSA 密钥设计，通常用于存储 RSA 私钥和公钥      | 适用于存储多种算法的私钥，可用于 RSA、DSA、ECDSA 等算法 |
| **可扩展性**          | 不支持其他算法，只适用于 RSA                      | 支持多种算法，具有更强的通用性                     |

