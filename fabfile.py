import os
import sys
import http.server
import socketserver
import codecs

from fabric import task
from fabric import Connection

BASE_PATH = os.path.dirname(__file__)
default_hosts = ["localhost"]
conn = Connection('walter@www.fanyamin.com')

def run_cmd(c, cmd):
    print(cmd)
    c.local(cmd)

@task(hosts=['localhost'])
def md2rst(c, src, dest=None):
    if not dest:
        dest = src[:-3] + ".rst";
    cmd = "pandoc --to RST --reference-links {} > {}".format(src, dest)
    c.local(cmd)


@task(hosts=default_hosts)
def rst2md(c, src, dest=None):
    if not dest:
        dest = src[:-4] + ".md";
    cmd = "pandoc {} -f rst -t markdown -o {}".format(src, dest)
    c.local(cmd)

@task(hosts=default_hosts)
def make_note(c):
    build_cmd = 'cd doc && make clean html'
    c.local(build_cmd)

@task(hosts=default_hosts)
def init_backend(c):
    cmds = """
    mkdir -p ./backend
    mkdir -p ./backend/internal/api/handlers
    mkdir -p ./backend/internal/api/middleware
    mkdir -p ./backend/internal/config
    mkdir -p ./backend/internal/models
    mkdir -p ./backend/cmd/server
    mkdir -p ./backend/pkg/auth
    mkdir -p ./backend/pkg/logger
    mkdir -p ./backend/pkg/validator
    touch mkdir -p ./backend/.env
    """
    for cmd in cmds.strip().split('\n'):
        run_cmd(c, cmd)
    libs = """
    go mod init github.com/walterfan/idaas/backend
    go get -u github.com/gin-gonic/gin
    go get -u gorm.io/gorm
    go get -u gorm.io/driver/sqlite
    go get -u github.com/golang-jwt/jwt/v5
    go get -u github.com/gin-contrib/cors
    go get -u github.com/go-playground/validator/v10
    """
    if not os.path.exists('./backend/go.mod'):
        with c.cd('./backend'):
            for lib_cmd in libs.strip().split('\n'):
                c.local(lib_cmd)

@task(hosts=default_hosts)
def init_frontend(c):
    cmds = """
    npm create vite@latest frontend -- --template vue
    npm install
    npm install axios tailwindcss postcss autoprefixer @headlessui/vue @heroicons/vue
    npm install -D @tailwindcss/forms @tailwindcss/typography
    npx tailwindcss init -p
    """
    for cmd in cmds.strip().split('\n'):
        run_cmd(c, cmd)

@task(hosts=default_hosts)
def run_backend(c):
    cmd = "cd ./backend && go run cmd/server/main.go"
    c.local(cmd)

@task(hosts=default_hosts)
def run_frontend(c):
    cmd = "cd ./frontend && npm run dev"
    c.local(cmd)

@task(hosts=default_hosts)
def build_frontend(c):
    cmd = "cd ./frontend && npm run build"
    c.local(cmd)

@task(hosts=default_hosts)
def build_backend(c):
    cmd = "cd ./backend && go build -o bin/server cmd/server/main.go"
    c.local(cmd)

@task(hosts=default_hosts)
def test_backend(c):
    cmd = "cd ./backend && go test ./..."
    c.local(cmd)

@task(hosts=default_hosts)
def deploy(c):
    # Build frontend
    build_frontend(c)
    
    # Build backend
    build_backend(c)
    
    # Deploy to server
    c.local("rsync -avz --delete ./frontend/dist/ walter@www.fanyamin.com:/var/www/html/identity-service/")
    c.local("rsync -avz ./backend/bin/server walter@www.fanyamin.com:/opt/identity-service/")
    
    # Restart service on remote server
    conn.run("sudo systemctl restart identity-service")

@task(hosts=default_hosts)
def setup_db(c):
    # Create database schema
    cmd = "cd ./backend && go run cmd/migrations/main.go"
    c.local(cmd)

@task(hosts=default_hosts)
def generate_api_docs(c):
    cmd = "cd ./backend && swag init -g cmd/server/main.go -o api/docs"
    c.local(cmd)

@task(hosts=default_hosts)
def clean(c):
    # Clean build artifacts
    c.local("rm -rf ./frontend/dist")
    c.local("rm -rf ./backend/bin")
