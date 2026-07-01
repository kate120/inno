# Цель модуля: Понять зачем нужен K8s (проблемы Docker Compose) и развернуть первый Pod (контейнер).

## 1. Задача (Теория): Проблемы docker-compose.

Контекст: Что если docker-compose up --scale app=3 (М6 З10)? Как балансировать трафик? Что если один контейнер упадет? (Compose его не перезапустит).

    Docker-compose up --scale app=3 запускает 3 app контенера, но при падении одного из них не востаналивает ничего. Docker compose не балансирует нагрузку, отвечает первый из 3х контенеров, надо использовать nginx (например).

Задание: Написать эссе "Почему Docker Compose — это не Production".
 * Self-healing - если контейнер упал, докер его не востановит
 * Scaling - не балансирует трафик между контейнерами
 * Zero-downtime deployment - при остановке контенера и повторном запуске есть downtime.

Критерии: Выделили 3 проблемы: Self-healing, Scaling, Zero-downtime deployment.

## 2. Задача (Теория): Что такое Оркестрация (Kubernetes)?

Контекст: K8s — это "дирижер" (оркестратор) для контейнеров. Он решает проблемы Compose.
    Оркестрация - управление контенерами с помощию таких оркестраторов как k8s. k8s проводит запуск контенеров, их учет, балансировку нагрузки между ними и управлению сетью.

Критерии: Вы можете объяснить, что K8s — это "ОС для датацентра".

## 3. Задача (Установка): Установить kubectl (CLI для K8s).
```bash
$ kubectl version
Client Version: v1.35.0
Kustomize Version: v5.7.1
Unable to connect to the server: dial tcp 127.0.0.1:51016: connectex: No connection could be made because the target machine actively refused it.
```
## 4. Задача (Установка): Установить minikube (Локальный K8s-кластер для разработки).
```bash
$ minikube status
minikube
type: Control Plane
host: Running
kubelet: Running
apiserver: Running
kubeconfig: Configured
```

## 5. Задача (Запуск): Запустить Minikube: minikube start.
```bash
 C:\Users\katar\Desktop\СТАЖИРОВКА\devops_tr> minikube start 
😄  minikube v1.38.1 on Microsoft Windows 11 Home Single Language 25H2
✨  Using the docker driver based on existing profile
```
## 6. Задача (Проверка): Проверить кластер: kubectl get nodes.
```bash
NAME       STATUS   ROLES           AGE     VERSION
minikube   Ready    control-plane   3m39s   v1.35.1
```
Критерии: Вы видите 1 ноду (minikube) в статусе Ready.

## 7. Задача (Теория): K8s Объекты (YAML). Изучить 4 ключевых:
* Pod -  наименьшая еденица, которой упраляет кубер, состоит из 1 или нескольких контенеров, изолированных с помощью linux ns, контенеры видят друг друга и могут делить общий volume. Чаще в одном поде один контейнер. 
* Deployment - модуль, отвечающий за колличество подов на нодах (если быть точнее, то он отвечает за replicaset, а тот уже на deployment)
* Service - модуль, отвечающий за наличии у подов стабильного ip (точнее при падении и востановлении меняется под и его ip соотвественно, чтобы обращаться к поду еслипользует sevice)
* ConfigMap / Secret - модули, которые отвечают за натсройку окружения/перемнных внутри контенера (конфигурации)

Pod: (1+ контейнер). Минимальная единица. Эфемерна.

Deployment: (Контроллер). Говорит K8s: "Я хочу, чтобы всегда было 3 Pod'а". Обеспечивает Self-healing.

Service: (Сеть). Дает стабильный IP-адрес и DNS-имя для группы Pod'ов.

ConfigMap / Secret: (Конфигурация).

## 8. Задача (Императивный запуск): Запустить Nginx без YAML: kubectl run my-nginx --image=nginx.
```bash
pod/my-nginx created
```
## 9. Задача (Управление): kubectl get pods. kubectl describe pod my-nginx.
```bash
NAME       READY   STATUS    RESTARTS   AGE
my-nginx   1/1     Running   0          36s
```
```bash
Name:             my-nginx
Namespace:        default
Priority:         0
Service Account:  default
Node:             minikube/192.168.49.2
```
## 10. Задача (Очистка): kubectl delete pod my-nginx.
```bash
pod "my-nginx" deleted from default namespace
```