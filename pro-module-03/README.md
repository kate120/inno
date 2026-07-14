# Модуль 3: GitLab CI (Альтернативный CI/CD Локально)
Цель модуля: Развернуть локально GitLab и GitLab Runner. Сравнить Jenkinsfile (Groovy) и .gitlab-ci.yml (YAML).
## 1. Задача (Развертывание GitLab): Запустить GitLab в Docker (используя docker-compose из документации GitLab).
```bash
mkdir gitlab-local
cd gitlab-local
mkdir -p data\docker\gitlab\etc\gitlab
mkdir -p data\docker\gitlab\var\opt\gitlab
mkdir -p data\docker\gitlab\var\log\gitlab
mkdir -p data\dind\docker
mkdir -p configNew-Item -Path "config\config.toml" -ItemType File
nano docker-compose.yml
docker-compose up
```
```yml
version: '3.8'

services:
  gitlab:
    image: gitlab/gitlab-ce:latest
    hostname: "localhost"
    restart: unless-stopped
    environment:
      GITLAB_OMNIBUS_CONFIG: |
        external_url 'http://localhost:9090'
        gitlab_rails['gitlab_shell_ssh_port'] = 8822
        gitlab_rails['initial_root_password'] = 'CHANGEME123'
    ports:
      - "9090:80"
      - "8822:22"
    volumes:
      - ./data/docker/gitlab/etc/gitlab:/etc/gitlab
      - ./data/docker/gitlab/var/opt/gitlab:/var/opt/gitlab
      - ./data/docker/gitlab/var/log/gitlab:/var/log/gitlab
    networks:
      - gitlab_net

  dind:
    image: docker:20-dind
    restart: always
    privileged: true
    environment:
      DOCKER_TLS_CERTDIR: ""
    command: [
      "--host", "tcp://0.0.0.0:2375",
      "--storage-driver=overlay2",
      "--tls=false"
    ]
    volumes:
      - ./data/dind/docker:/var/lib/docker
    networks:
      - gitlab_net

  gitlab-runner:
    image: gitlab/gitlab-runner:alpine
    restart: unless-stopped
    environment:
      DOCKER_HOST: "tcp://dind:2375"
    volumes:
      - ./config:/etc/gitlab-runner:z
      - /var/run/docker.sock:/var/run/docker.sock
    networks:
      - gitlab_net
    depends_on:
      - dind

networks:
  gitlab_net:
    name: gitlab_net
```
Контекст: Это "тяжелый" сервис, требует 4-8 ГБ RAM.
Критерии: UI GitLab доступен (e.g., http://localhost:10080).
## 2. Задача (Настройка): Войти (root, пароль из docker logs), создать проект, "запушить" my-app в этот локальный GitLab.

    Логин: root
    Пароль: CHANGEME123

## 3. Задача (Развертывание Runner): Запустить GitLab Runner в Docker.
Задание: docker run ... gitlab/gitlab-runner:latest.
```bash
# регитрация runner
PS C:\Users\katar\Desktop\СТАЖИРОВКА\devops_tr\pro-module-03\gitlab-local> docker ps
CONTAINER ID   IMAGE                           COMMAND                  CREATED          STATUS                    PORTS                                                                                  NAMES
7ae9f1e148f6   gitlab/gitlab-runner:alpine     "/usr/bin/dumb-init …"   24 minutes ago   Up 24 minutes                                                                                                    gitlab-local-gitlab-runner-1
cf27c83bb17a   gitlab/gitlab-ce:latest         "/assets/init-contai…"   24 minutes ago   Up 24 minutes (healthy)   0.0.0.0:9090->9090/tcp, [::]:9090->9090/tcp, 0.0.0.0:8822->22/tcp, [::]:8822->22/tcp   gitlab-local-gitlab-1
f893d50752d7   moby/buildkit:buildx-stable-1   "/usr/bin/buildkitd-…"   6 days ago       Up 4 days                                                                                                        buildx_buildkit_my-builder0
PS C:\Users\katar\Desktop\СТАЖИРОВКА\devops_tr\pro-module-03\gitlab-local> docker exec -it 7ae9f1e148f6 gitlab-runner register                  
Runtime platform                                    arch=amd64 os=linux pid=19 revision=24b9b726 version=19.1.1
Running in system-mode.                            
                                                   
Enter the GitLab instance URL (for example, https://gitlab.com/):
http://gitlab:9090
Enter the registration token:
glrt-6HXDj.01.171613vz3
Verifying runner... is valid                        correlation_id=01KXDWBW41P3Z4TV7ZJXYFSJ2C runner=6HXDjpDzV runner_name=7ae9f1e148f6
Enter a name for the runner. This is stored only in the local config.toml file:
[7ae9f1e148f6]: runner
Enter an executor: custom, instance, docker, kubernetes, parallels, docker-windows, docker-autoscaler, docker+machine, ssh, virtualbox, shell:
docker
Enter the default Docker image (for example, ruby:3.3):
alpine:latest
Runner registered successfully. Feel free to start it, but if it's running already the config should be automatically reloaded!
 
Configuration (with the authentication token) was saved in "/etc/gitlab-runner/config.toml" 

What's next:
    Try Docker Debug for seamless, persistent debugging tools in any container or image → docker debug 7ae9f1e148f6
    Learn more at https://docs.docker.com/go/debug-cli/
PS C:\Users\katar\Desktop\СТАЖИРОВКА\devops_tr\pro-module-03\gitlab-local> 
```
  Чтобы пайплайн заработал адо было изменить файл  C:\Users\katar\Desktop\СТАЖИРОВКА\devops_tr\pro-module-03\gitlab-local\config\config.toml
  
```toml
  [runners.docker]
    tls_verify = false
    image = "alpine:latest"
    privileged = false
    disable_entrypoint_overwrite = false
    oom_kill_disable = false
    disable_cache = true
    network_mode = "gitlab_net"
    volume_keep = false
    shm_size = 0
    network_mtu = 0
```
## 4. Задача (Регистрация Runner): Зарегистрировать Runner в GitLab.
Задание: docker exec [runner] gitlab-runner register (указать URL, токен).
Критерии: Runner виден в UI GitLab (в Settings -> CI/CD -> Runners).
## 5. Задача (.gitlab-ci.yml): Создать .gitlab-ci.yml в корне my-app.
```yml
stages:
  - test
  - build
```
stages: [test, build]
## 6. Задача (Job test): Описать test job.
test:
image: python:3.10-slim
script:
- pip install -r requirements.txt
- pytest
Критерии: git push -> Runner "подхватывает" job -> Тесты проходят.
## 7. Задача (Docker-in-Docker): Настройка GitLab CI для сборки Docker.
Контекст: GitLab CI использует image: docker:latest и services: [docker:dind].
Критерии: Вы понимаете, что Runner запускает dind (Docker-in-Docker) как "братский" контейнер.

    Контейнер dind нужен для того, чтобы внутри сети, созданной с помощью docker compose, контенер gitlab при сборке образов и из пуше мог работаь с docker движком, реализовать это можно через проброс сокетов или через создание нового докер контейнера dind. Runner при сборке образа обращается к dind 

## 8. Задача (Job build): Описать build job.
build:
stage: build
image: docker:latest
services: [docker:dind]
script:
- docker login ... -u $CI_REGISTRY_USER -p $CI_REGISTRY_PASSWORD $CI_REGISTRY
- docker build -t $CI_REGISTRY_IMAGE:$CI_COMMIT_SHA .
- docker push ...
Контекст: $CI_... — это встроенные переменные GitLab CI. G. 
```yml
stages:
  - test
  - build

variables:
  DOCKER_HOST: tcp://dind:2375
  DOCKER_TLS_CERTDIR: ""

run-tests:
  stage: test
  image: python:3.10-slim
  script:
    - cd module-08
    - pip install -r requirements.txt
    - pytest tests/

build-image:
  stage: build
  image: docker:latest
  services:
    - docker:dind
  variables:
    DOCKER_HOST: tcp://docker:2375
    DOCKER_TLS_CERTDIR: ""
  script:
    - cd module-08
    - docker login -u $CI_REGISTRY_USER -p $CI_REGISTRY_PASSWORD $CI_REGISTRY
    - docker build -t $CI_REGISTRY_IMAGE:$CI_COMMIT_SHA .
    - docker push $CI_REGISTRY_IMAGE:$CI_COMMIT_SHA
```
## 9. Задача (GitLab Registry): Включить встроенный Docker Registry в GitLab.
Критерии: Пайплайн (Задача 8) пушит образ во внутренний Registry GitLab'а.
## 10. Задача (Сравнение): Написать эссе "Jenkins vs. GitLab CI".
Критерии: Сравнили: Jenkinsfile (Groovy, гибкий, много плагинов) vs. .gitlab-ci.yml (YAML, простой, интегрирован с SCM/Registry).

  Jenkins и GitLab CI — это два принципиально разных подхода к CI/CD. Jenkins — это гибкий, но сложный конструктор на Groovy (Jenkinsfile) с огромной экосистемой из 1800+ плагинов, дающий полный контроль и свободу, но требующий ручного обслуживания инфраструктуры, совместимости плагинов и выделенной DevOps-команды. GitLab CI — это унифицированная платформа "всё в одном", использующая простой декларативный YAML (.gitlab-ci.yml) с глубокой интеграцией кода, Container Registry и встроенными инструментами безопасности, что обеспечивает быстрый старт и низкие операционные расходы ценой меньшей гибкости и привязки к экосистеме GitLab. Итог: Jenkins выбирают за максимальную кастомизацию и контроль, а GitLab CI — за простоту, единую среду и экономию ресурсов команды.
