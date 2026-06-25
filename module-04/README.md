# Цель модуля: Научиться сохранять данные (Volumes) и связывать контейнеры (Networking), чтобы App мог "увидеть" DB.

## Задача (Теория): Почему данные пропадают?
Контекст: Файловая система контейнера эфемерна. При docker rm все данные внутри него удаляются.

Задание: docker run -it ubuntu bash, создать touch /data/test.txt, exit. docker rm [container]. docker run ... (тот же образ) — файла /data/test.txt нет.

Критерии: Вы практически убедились в эфемерности.
``` bash  
PS C:\Users\katar\Desktop\СТАЖИРОВКА\devops_tr> docker run -it ubuntu bash 
root@08f8cb861839:/# mkdir data
root@08f8cb861839:/# cd data 
root@08f8cb861839:/data# touch test.txt
root@08f8cb861839:/data# ls
test.txt
root@08f8cb861839:/data# exit
exit
PS C:\Users\katar\Desktop\СТАЖИРОВКА\devops_tr> docker ps 
CONTAINER ID   IMAGE           COMMAND                  CREATED       STATUS       PORTS     NAMES
d0a4463e14fa   java-main:1.0   "/__cacert_entrypoin…"   3 hours ago   Up 3 hours             objective_wilson
PS C:\Users\katar\Desktop\СТАЖИРОВКА\devops_tr> docker stop d0a4463e14fa 
d0a4463e14fa
PS C:\Users\katar\Desktop\СТАЖИРОВКА\devops_tr> docker rm d0a4463e14fa  
d0a4463e14fa
PS C:\Users\katar\Desktop\СТАЖИРОВКА\devops_tr> docker run -it ubuntu bash     
root@64e135cb4b67:/# ls
bin  boot  dev  etc  home  lib  lib64  media  mnt  opt  proc  root  run  sbin  srv  sys  tmp  usr  var
root@64e135cb4b67:/# 

```
## 2. Задача (Теория): Bind Mount vs. Volume.
Контекст: Bind Mount (флаг -v /host/path:/container/path) — "горячая" проброска папки с хоста, идеальна для разработки (hot-reload). Volume (флаг -v volume-name:/container/path) — управляется Docker, для production-данных (БД).
```
Bind Mount идеале для разработки, т к измение файла на хосте приведет к измению файла внутри компьютера. При удалении соответвенно тоже самое. За работу с файлом отвеачает челокек. А при работе с volume данные из контнера храняться в файле, который управляется Docker, человек не может случайно поменять данные, хорошо для БД.
```
Задание: Нарисовать схему, объясняющую разницу.
## 3. Задача (Bind Mount): "Горячая" перезагрузка (Hot Reload).
Задание: Запустить Python/Node приложение, используя Bind Mount: docker run -p 5000:5000 -v $(pwd):/app my-app:1.0 (может требоваться uvicorn --reload или nodemon).
``` Python
from fastapi import FastAPI
import datetime

app = FastAPI()

@app.get("/")
def root():
    return {
        "message": "Hello from Python and Kate with hot-reload!",
        "time": str(datetime.datetime.now())
    }

@app.get("/test")
def test():
    return {"status": "ok", "data": "you can change me!"}

if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=5000)
```
``` Dockerfile
FROM python:3.9-slim
WORKDIR /app
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt
CMD ["uvicorn", "app:app", "--host", "0.0.0.0", "--port", "5000", "--reload"]
```
``` bash
PS C:\Users\katar\Desktop\СТАЖИРОВКА\devops_tr\module-04> docker run -p 5000:5000 -v C:\Users\katar\Desktop\СТАЖИРОВКА\devops_tr\module-04\my-app.py:/app/app.py my-app:1.0
INFO:     Will watch for changes in these directories: ['/app']
INFO:     Uvicorn running on http://0.0.0.0:5000 (Press CTRL+C to quit)
INFO:     Started reloader process [1] using StatReload
INFO:     Started server process [8]
INFO:     Waiting for application startup.
INFO:     Application startup complete.
INFO:     172.17.0.1:55870 - "GET / HTTP/1.1" 200 OK
INFO:     172.17.0.1:55870 - "GET /favicon.ico HTTP/1.1" 404 Not Found
WARNING:  StatReload detected changes in 'app.py'. Reloading...
INFO:     Shutting down
INFO:     Waiting for application shutdown.
INFO:     Application shutdown complete.
INFO:     Finished server process [8]
INFO:     Started server process [10]
INFO:     Waiting for application startup.
INFO:     Application startup complete.
INFO:     172.17.0.1:44586 - "GET / HTTP/1.1" 200 OK

```
Критерии: Изменения кода на хост-машине немедленно отражаются в контейнере без пересборки.

## 4. Задача (Volume): Создать "именованный" Volume: docker volume create postgres-data.
``` bash
PS C:\Users\katar\Desktop\СТАЖИРОВКА\devops_tr\module-04> docker volume create postgres-data.
postgres-data.
PS C:\Users\katar\Desktop\СТАЖИРОВКА\devops_tr\module-04> docker volume ls
DRIVER    VOLUME NAME
local     postgres-data.
PS C:\Users\katar\Desktop\СТАЖИРОВКА\devops_tr\module-04> 

```
Критерии: docker volume ls показывает postgres-data.

## 5. Задача (Volume): Запустить PostgreSQL с Volume.

Задание: docker run -d -p 5432:5432 -v postgres-data:/var/lib/postgresql/data -e POSTGRES_PASSWORD=mysecret postgres.
``` bash

CONTAINER ID   IMAGE         COMMAND                  CREATED         STATUS         PORTSNAMES
f466b9009373   postgres:15   "docker-entrypoint.s…"   4 seconds ago   Up 3 seconds   0.0.0.0:5433->5432/tcp, [::]:5433->5432/tcpmy-postgres

```
Критерии: Postgres запущен.

## 6. Задача (Проверка Volume): docker stop/rm контейнер Postgres. Запустить команду из Задачи 5 еще раз.
``` bash
PS C:\Users\katar\Desktop\СТАЖИРОВКА\devops_tr\module-04> docker run -d --name my-postgres-new1 -p 5433:5432 -v postgres-data-new:/var/lib/postgresql/data -e POSTGRES_PASSWORD=mysecret postgres:15
4ad2b2c8aea7994e6b192a256242f7264e08c39c38ca38e5a152f9524c20489d
PS C:\Users\katar\Desktop\СТАЖИРОВКА\devops_tr\module-04> docker ps
CONTAINER ID   IMAGE         COMMAND                  CREATED         STATUS         PORTSNAMES
4ad2b2c8aea7   postgres:15   "docker-entrypoint.s…"   3 seconds ago   Up 2 seconds   0.0.0.0:5433->5432/tcp, [::]:5433->5432/tcpmy-postgres-new

```
Критерии: Postgres стартует мгновенно и не жалуется на отсутствие БД. Данные остались на месте, т.к. они жили в Volume.

## 7. Задача (Теория Сети): Сеть bridge (по умолчанию).

Контекст: Контейнеры могут общаться по IP, но IP меняются.

Задание: Запустить 2 ubuntu контейнера (sleep 1000). docker inspect и найти их IP. Зайти в первый (docker exec -it [id] bash) и ping второй по IP.
``` bash
$ docker ps
CONTAINER ID   IMAGE     COMMAND        CREATED              STATUS              PORTS     NAMES
3afa9609a97c   ubuntu    "sleep 1000"   57 seconds ago       Up 56 seconds                 ubuntu2
389790dba973   ubuntu    "sleep 1000"   About a minute ago   Up About a minute             ubuntu1

katar@Tecno MINGW64 ~/Desktop/СТАЖИРОВКА/devops_tr/module-04 (feature/module-04-docker-networking-volumes)
$ docker inspect ubuntu1 | grep IPAddress
                    "IPAddress": "172.17.0.2",

katar@Tecno MINGW64 ~/Desktop/СТАЖИРОВКА/devops_tr/module-04 (feature/module-04-docker-networking-volumes)
$ docker inspect ubuntu2 | grep IPAddress
                    "IPAddress": "172.17.0.3",

katar@Tecno MINGW64 ~/Desktop/СТАЖИРОВКА/devops_tr/module-04 (feature/module-04-docker-networking-volumes)
$

```

Критерии: Пинг проходит.
``` bash 

ping 172.17.0.3
PING 172.17.0.3 (172.17.0.3) 56(84) bytes of data.
64 bytes from 172.17.0.3: icmp_seq=1 ttl=64 time=0.497 ms
64 bytes from 172.17.0.3: icmp_seq=2 ttl=64 time=0.045 ms
64 bytes from 172.17.0.3: icmp_seq=3 ttl=64 time=0.041 ms
64 bytes from 172.17.0.3: icmp_seq=4 ttl=64 time=0.059 ms
64 bytes from 172.17.0.3: icmp_seq=5 ttl=64 time=0.038 ms

```
## 8. Задача (Сеть bridge): Создать свою bridge сеть: docker network create my-app-net.
``` bash

PS C:\Users\katar\Desktop\СТАЖИРОВКА\devops_tr\module-04> docker network create my-app-net
fc86587128655e3aaec74f3a05ba9788dc14e14b401857c12b7875d178d3211c
PS C:\Users\katar\Desktop\СТАЖИРОВКА\devops_tr\module-04> docker network ls
NETWORK ID     NAME         DRIVER    SCOPE
4c68ebf3c6f3   bridge       bridge    local
cf97a18bd6d7   host         host      local
fc8658712865   my-app-net   bridge    local
8a063acecb95   none         null      local
PS C:\Users\katar\Desktop\СТАЖИРОВКА\devops_tr\module-04> 

```
Критерии: Сеть создана.

## 9. Задача (Подключение к сети): Запустить Postgres в этой сети (--network my-app-net --name pg-db).
``` bash
PS C:\Users\katar\Desktop\СТАЖИРОВКА\devops_tr\module-04> docker run -d --network my-app-net --name pg-db -p 5432:5432  -v postgres-data:/var/lib/postgresql/data -e POSTGRES_PASSWORD=mysecret postgres:15 
30c6b9f81f5b14b5009d7971df0c5615ab6fe47c28ac24f8bbae8ac4f1720bc1
PS C:\Users\katar\Desktop\СТАЖИРОВКА\devops_tr\module-04> docker ps                                                                
CONTAINER ID   IMAGE         COMMAND                  CREATED          STATUS          PORTS  NAMES
30c6b9f81f5b   postgres:15   "docker-entrypoint.s…"   4 seconds ago    Up 3 seconds    0.0.0.0:5432->5432/tcp, [::]:5432->5432/tcp  pg-db

```
Критерии: Postgres в my-app-net.

## 10. Задача (DNS-Discovery): Запустить ubuntu в этой же сети (--network my-app-net). apt update && apt install dnsutils -y. Выполнить ping pg-db.
```bash

 C:\Users\katar\Desktop\СТАЖИРОВКА\devops_tr\module-04> docker ps
CONTAINER ID   IMAGE         COMMAND                  CREATED          STATUS          PORTS  NAMES
de6ff3af7326   ubuntu        "sleep 1000"             2 minutes ago    Up 2 minutes  ubuntu4

```
```bash

ping pg-db
PING pg-db (172.18.0.2) 56(84) bytes of data.
64 bytes from pg-db.my-app-net (172.18.0.2): icmp_seq=1 ttl=64 time=0.324 ms
64 bytes from pg-db.my-app-net (172.18.0.2): icmp_seq=2 ttl=64 time=0.112 ms
64 bytes from pg-db.my-app-net (172.18.0.2): icmp_seq=3 ttl=64 time=0.120 ms
64 bytes from pg-db.my-app-net (172.18.0.2): icmp_seq=4 ttl=64 time=0.129 ms
64 bytes from pg-db.my-app-net (172.18.0.2): icmp_seq=5 ttl=64 time=0.081 ms

```
Критерии: Пинг по имени контейнера (pg-db) проходит! Docker обеспечивает DNS-обнаружение внутри кастомной сети.