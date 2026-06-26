# Цель модуля: Настроить production-ready docker-compose файл со связями, вольюмами, переменными и healthcheck.

## 1. Задача (Переменные окружения): Добавить environment в сервис db.

db:

...

environment:

- POSTGRES_USER=devops

- POSTGRES_PASSWORD=mysecret
```yml

version: '3.8'

services:
  app:
    build: .
    container_name: app
    ports:
      - "5000:5000"

  db:
    image: postgres:14-alpine
    container_name: bd
    ports:
      - "5432:5432"
    environment:
      - POSTGRES_USER=devops
      - POSTGRES_PASSWORD=mysecret

```

## 2. Задача (Файл .env): Создать .env файл в корне.

POSTGRES_USER=devops

POSTGRES_PASSWORD=mysecret

В docker-compose.yml использовать: - POSTGRES_USER=${POSTGRES_USER}.

Критерии: Compose "подтягивает" переменные из .env (Секреты не должны быть в docker-compose.yml).

```yaml

 db:
    image: postgres:14-alpine
    container_name: bd
    ports:
      - "5432:5432"
    environment:
      - POSTGRES_USER=${POSTGRES_USER}.
      - POSTGRES_PASSWORD=${POSTGRES_PASSWORD}

```

## 3. Задача (Связь Сервисов): Передать DATABASE_URL в app.

app:

...

environment:

- DATABASE_URL=postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@db:5432/db_name

Контекст: db — это имя сервиса PostgreSQL, которое работает как DNS-имя! Compose автоматически создает сеть.
```yml 

services:
  app:
    build: .
    container_name: app
    ports:
      - "5000:5000"
    environment:
      - DATABASE_URL=postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@db:5432/db_name

```

Критерии: Приложение app теперь может подключиться к db.

```bash

/app # POSTGRES_PGPASSWORD=mysecret  psql -h db -U devops
Password for user devops: 
psql (18.4, server 14.23)
Type "help" for help.

devops=# ls
devops-# 

```

## 4. Задача (Volumes): Добавить volumes для db.

db:

...

volumes:

- postgres-data:/var/lib/postgresql/data

volumes: (в корне)

postgres-data: (объявить named volume)

```
в корне значит на уровне service
```
```yaml

version: '3.8'

services:
  app:
    build: .
    container_name: app
    ports:
      - "5000:5000"
    environment:
      - DATABASE_URL=postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@db:5432/db_name

  db:
    image: postgres:14-alpine
    container_name: bd
    ports:
      - "5432:5432"
    environment:
      - POSTGRES_USER=${POSTGRES_USER}
      - POSTGRES_PASSWORD=${POSTGRES_PASSWORD}
    volumes:
      - postgres-data:/var/lib/postgresql/data
  
  
volumes: 
  postgres-data: 

```
Критерии: Данные Postgres теперь персистентны (не удаляются при docker-compose down).
```bash

PS C:\Users\katar\Desktop\СТАЖИРОВКА\devops_tr\module-06> docker compose down       
time="2026-06-26T13:25:30+03:00" level=warning msg="C:\\Users\\katar\\Desktop\\СТАЖИРОВКА\\devops_tr\\module-06\\docker-compose.yml: the attribute `version` is obsolete, it will be ignored, please remove it to avoid potential confusion"
[+] down 3/3
 ✔ Container bd              Removed                                                                                    0.6s
 ✔ Container app             Removed                                                                                    0.6s
 ✔ Network module-06_default Removed                                                                                    0.3s
PS C:\Users\katar\Desktop\СТАЖИРОВКА\devops_tr\module-06> docker ps
CONTAINER ID   IMAGE     COMMAND   CREATED   STATUS    PORTS     NAMES
PS C:\Users\katar\Desktop\СТАЖИРОВКА\devops_tr\module-06> docker volume ls          
DRIVER    VOLUME NAME
local     5b39787823ca7219aca20aa20c259f2906d31475d1ef2ac486501fecf2e99fdd
local     7f1a52b13829f1677708b5721611f7880b9f8ea074ffd89aecd1ff3442193567
local     26d78d9f0340cce5829d3bc6b75bb04f6aa597e7709ff6a2e8bcf804ec3f57d9
local     573f718eb1f173801782fbafdd8f872af8ba4b74097e640bfc1e29e58c705d44
local     32457d0b822f310bb7e695f1654a70dfdd44711c6d2c95c56e08650d3f1b851c
local     c27fdf8b18125db34bb9fc3b7866e438ecca306b89dce2ce98469f9813924d70
local     e66e9fc7a7351d726e84a63f87a0a948c80415c6fdf4859cf7bac7fdf432b084
local     module-06_postgres-data
local     postgres-data
local     postgres-data-new

```

## 5. Задача (Bind Mount): Добавить Bind Mount для app (для разработки).

app:

...

volumes:

- .:/app

Критерии: "Горячая" перезагрузка кода app работает.
![alt text](image.png)
![alt text](image-1.png)
```docker-compose

services:
  app:
    build: .
    ports:
      - "5000:5000"
    environment:
      - DATABASE_URL=postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@db:5432/db_name
    volumes:
      - ./for_go/:/myapp
      - go-modules:/root/go/pkg/mod
    depends_on:
      db:
        condition: service_healthy

  db:
    image: postgres:14-alpine
    container_name: bd
    ports:
      - "5432:5432"
    environment:
      - POSTGRES_USER=${POSTGRES_USER}
      - POSTGRES_PASSWORD=${POSTGRES_PASSWORD}
    volumes:
      - postgres-data:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U ${POSTGRES_USER}"]
      interval: 5s
      timeout: 5s
      retries: 5

  redis:
    image: redis:alpine
    container_name: redis
    ports:
      - "6379:6379"

volumes:
  postgres-data:


```
```
Для реализации bind mound я использовада дополнительные библиотеки с github, т к html файлы при пробросе обновлялись без проблем, но app.go (а именно внутри переменной) не отображался в браузере (т е обновление не отображалось). Поэтому отдельно был создан файл .air.toml и были внесены изменения  в docker-compose.eml 

```
.air.toml

```
root = "."
tmp_dir = "tmp"

[build]
  cmd = "go build -o ./tmp/main ."
  bin = "./tmp/main"
  include_ext = ["go", "html"]
  exclude_dir = ["tmp"]
  delay = 1000
  poll = true
  poll_interval = 1000


[screen]
  clear_on_rebuild = false

```
dockerfile
```dockerfile

FROM golang:1.26-alpine
WORKDIR /myapp
RUN go install github.com/air-verse/air@latest
COPY for_go/go.mod ./
RUN go mod download
CMD ["air"]

```
## 6. Задача (depends_on): Гарантировать порядок запуска.

app:

...

depends_on:

- db
```bash
Было:
✔ Network module-06_default Created                                                                                                                                0.1s
 ✔ Container app             Started                                                                                                                                0.9s
 ✔ Container bd              Started 

Стало:
✔ Network module-06_default Created                                                                                                                                0.1s
 ✔ Container bd              Started                                                                                                                                0.8s
 ✔ Container app             Started 
```
Критерии: db всегда стартует раньше, чем app.

## 7. Задача (healthcheck): Добавить healthcheck в db.

db:

...

healthcheck:

test: ["CMD-SHELL", "pg_isready -U devops"]

interval: 5s

timeout: 5s

Критерии: docker-compose ps показывает (healthy).
```bash

PS C:\Users\katar\Desktop\СТАЖИРОВКА\devops_tr\module-06> docker-compose ps
time="2026-06-26T14:54:55+03:00" level=warning msg="C:\\Users\\katar\\Desktop\\СТАЖИРОВКА\\devops_tr\\module-06\\docker-compose.yml: the attribute `version` is obsolete, it will be ignored, please remove it to avoid potential confusion"
NAME      IMAGE                COMMAND                  SERVICE   CREATED              STATUS                        PORTS
app       module-06-app        "go run app.go --rel…"   app       About a minute ago   Up About a minute             0.0.0.0:5000->5000/tcp, [::]:5000->5000/tcp
bd        postgres:14-alpine   "docker-entrypoint.s…"   db        About a minute ago   Up About a minute (healthy)   0.0.0.0:5432->5432/tcp, [::]:5432->5432/tcp
PS C:\Users\katar\Desktop\СТАЖИРОВКА\devops_tr\module-06> 

```

## 8. Задача (depends_on + healthcheck): Сделать app ждал здорового db.

app:

...

depends_on:

db:

condition: service_healthy
```bash

app:
    build: .
    container_name: app
    ports:
      - "5000:5000"
    environment:
      - DATABASE_URL=postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@db:5432/db_name
    volumes:
      - ./home_page.html:/app/home_page.html
    command: go run app.go --reload
    depends_on:
      - db
    condition: service_healthy

```

Критерии: app не запустится, пока pg_isready не вернет true.

## 9. Задача (Третий сервис): Добавить redis (из image: redis:alpine) в docker-compose.yml.

Критерии: docker-compose up запускает app, db и redis.
```yaml

redis:
    image: redis:alpine
    container_name: redis
    ports:
      - "6379:6379"

```

## 10. Задача (Масштабирование): Запустить 3 инстанса app.

Задание: docker-compose up -d --scale app=3.
```bash

PS C:\Users\katar\Desktop\СТАЖИРОВКА\devops_tr\module-06> docker ps
CONTAINER ID   IMAGE                COMMAND                  CREATED          STATUS                    PORTS                                           NAMES
8cc7877440bb   module-06-app        "go run app.go --hos…"   11 seconds ago   Up 4 seconds              0.0.0.0:50360->5000/tcp, [::]:50360->5000/tcp   module-06-app-3
c177f95047f4   module-06-app        "go run app.go --hos…"   11 seconds ago   Up 4 seconds              0.0.0.0:50361->5000/tcp, [::]:50361->5000/tcp   module-06-app-1
2c1520e1fe7e   module-06-app        "go run app.go --hos…"   11 seconds ago   Up 3 seconds              0.0.0.0:51259->5000/tcp, [::]:51259->5000/tcp   module-06-app-2
7708c392de6d   postgres:14-alpine   "docker-entrypoint.s…"   12 seconds ago   Up 10 seconds (healthy)   0.0.0.0:5432->5432/tcp, [::]:5432->5432/tcp     bd
31c17e9aa3b4   redis:alpine         "docker-entrypoint.s…"   12 seconds ago   Up 10 seconds             0.0.0.0:6379->6379/tcp, [::]:6379->6379/tcp     redis
PS C:\Users\katar\Desktop\СТАЖИРОВКА\devops_tr\module-06> 

```
Критерии: docker-compose ps показывает 3 инстанса app и по одному db/redis. (Примечание: ports: 5000:5000 вызовет ошибку, нужен reverse-proxy типа Nginx... это подводит к K8s).