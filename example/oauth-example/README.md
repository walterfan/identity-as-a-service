
##


## prepare

1. 在 GitHub 创建 OAuth App
登录 GitHub 后，访问 Settings -> Developer settings -> OAuth Apps

点击 "New OAuth App"

填写信息：

Application name: 你的应用名称

Homepage URL: http://localhost:8081 (Vue.js 开发服务器)

Authorization callback URL: http://localhost:8081/callback (或你的后端回调端点)

注册后，记下 Client ID 和生成一个 Client Secret


2. write backend by fastapi


3. write frontend by vue

   


## backend flow

```plantuml

@startuml

autonumber

frontend -> backend: GET /auth/github
backend -> frontend:  307 location=https://github.com/login/oauth/authorize?client_id=xxx&redirect_uri=http://localhost:8000/auth/github/callback
frontend -> github: https://github.com/login/oauth/authorize?client_id=xxx&redirect_uri=http://localhost:8000/auth/github/callback
note right frontend: user login into github and give permission
github -> frontend: 302 location=http://localhost:8000/auth/github/callback?code=yyy
frontend -> backend: /auth/github/callback?code=yyy
note right frontend: backend ask github's access token by the authorization code
backend -> github: POST https://github.com/login/oauth/access_token
github --> backend: access token
backend --> frontend: 308, location=http://localhost:5173?access_token=zzz
frontend -> backend: /?access_token=zzz
note right frontend: backend get  github's user info by client_id, client_secret and access token
backend --> github:  get https://api.github.com/user with header Authorization Bearer: {token}
github --> backend: user info
backend --> frontend: user info
frontend -> frontend: welcome walter, etc.

@enduml
```