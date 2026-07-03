# Цель модуля: Научиться развертывать stateful-приложения (БД) и передавать в них конфигурацию (ConfigMaps) и пароли (Secrets).

## 1. Задача (Теория): ConfigMap vs Secret.

Контекст: ConfigMap — для не-секретных конфигов (в plain text). Secret — для секретных (пароли, ключи API), хранятся в base64.

    Огромному количеству приложений необходимо зачитывать файлы для работы, это могут быть файлы конфигурации, публичные или приватные TLS ключи, текстовые шаблоны и многое другое. Самый популярный файл, который используется разработчиками - это config.yaml. 
```yml
## пример config.yaml

app:
  name: my-application
  version: 1.0.0

server:
  port: 8080
  host: 0.0.0.0

database:
  host: localhost
  port: 5432
  user: db_user
  password: db_password
  name: my_database

logging:
  level: info
```
## 2. Задача (Secret): Создать postgres-secret.yml (kind: Secret).

stringData: { POSTGRES_PASSWORD: "mysecret" } (K8s сам закодирует в base64).

Критерии: kubectl apply ...
```yml
apiVersion: v1
kind: Secret
metadata:
  name: postgres-secret
stringData:
  POSTGRES_PASSWORD: "mysecret"
```
secret/postgres-secret created

## 3. Задача (ConfigMap): Создать postgres-configmap.yml (kind: ConfigMap).
```yml
apiVersion: v1
kind: ConfigMap
metadata:
  name: postgres-config
data:
  POSTGRES_DB: "devopsdb"
  POSTGRES_USER: "devops"
```
data: { POSTGRES_DB: "devopsdb", POSTGRES_USER: "devops" }

Критерии: kubectl apply ...

    configmap/postgres-config created

## 4. Задача (Теория): PersistentVolume (PV) и PersistentVolumeClaim (PVC).
```yml
# Файл: nfs-storage.yml
apiVersion: v1
kind: PersistentVolume
metadata:
  name: nfs-pv
spec:
  capacity:
    storage: 20Gi
  accessModes:
    - ReadWriteMany
  nfs:
    server: 192.168.1.100   # замените на ваш NFS-сервер
    path: "/exports/data"
  persistentVolumeReclaimPolicy: Retain
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: nfs-pvc
spec:
  accessModes:
    - ReadWriteMany
  resources:
    requests:
      storage: 5Gi
```
  PV - в случае с миникубом, это часть памяти на диске (зарезервированная миникубом), которая может использоваться потом, если тому надо место для записи данных, например, под БД. PVC это посредник между ними, т е под запрашивает сначала PVC, а тот уже дает PV. При удалении пода, его данные (данные БД) сохраняются в PV  

Контекст: PV — это "диск" (e.g., SSD). PVC — это "запрос" (e.g., "Хочу 1 ГБ").

## 5. Задача (PVC): Создать postgres-pvc.yml (kind: PersistentVolumeClaim).

spec: { accessModes: [ReadWriteOnce], resources: { requests: { storage: 1Gi } } }
```yml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: postgres-pvc
spec:
  accessModes: 
    - ReadWriteOnce
  resourсes:
    requests:
      storage: 1Gi

```
Критерии: kubectl apply ... -> kubectl get pvc (статус Pending, т.к. Minikube создаст PV автоматически).
### PVC
```bash
NAME           STATUS   VOLUME                                     CAPACITY   ACCESS MODES   STORAGECLASS   VOLUMEATTRIBUTESCLASSAGE
postgres-pvc   Bound    pvc-3f5b4f64-5ca7-4d09-8632-3e2ad81c42e1   1Gi        RWO            standard       <unset>54s
```
### PV
```bash
PS C:\Users\katar\Desktop\СТАЖИРОВКА\devops_tr\module-13> kubectl get pv 
NAME                                       CAPACITY   ACCESS MODES   RECLAIM POLICY   STATUS   CLAIM                  STORAGECLASS VOLUMEATTRIBUTESCLASS   REASON   AGE
pvc-3f5b4f64-5ca7-4d09-8632-3e2ad81c42e1   1Gi        RWO            Delete           Bound    default/postgres-pvc   standard <unset>  
```
## 6. Задача (YAML - Deployment Postgres): Создать postgres-deployment.yml. (Использовать image: postgres:14-alpine).
```yml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: postgres-deployment
spec:
  replicas: 1
  selector:
    matchLabels:
      postgres: my-bd
  template:
    metadata:
      labels:
        postgres: my-bd
    spec:
      containers:
        - name: my-bd
          image: postgres:14-alpine
```
## 7. Задача (Инъекция Config): В postgres-deployment.yml в spec.template.spec.containers добавить envFrom:

envFrom: [ { configMapRef: { name: postgres-configmap } }, { secretRef: { name: postgres-secret } } ]
```yml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: postgres-deployment
spec:
  replicas: 1
  selector:
    matchLabels:
      postgres: my-bd
  template:
    metadata:
      labels:
        postgres: my-bd
    spec:
      containers:
        - name: my-bd
          image: postgres:14-alpine
          envFrom:
            configMapRef:
              name: postgres-configmap
            secretRef:
              name: postgres-secret
```

Критерии: Postgres "увидит" POSTGRES_USER, _PASSWORD, _DB как переменные окружения.

## 8. Задача (Монтирование Volume): В postgres-deployment.yml "примонтировать" PVC.

volumes: [ { name: pg-data, persistentVolumeClaim: { claimName: postgres-pvc } } ]

volumeMounts: [ { name: pg-data, mountPath: /var/lib/postgresql/data } ]
```yml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: postgres-deployment
spec:
  replicas: 1
  selector:
    matchLabels:
      postgres: my-bd
  template:
    metadata:
      labels:
        postgres: my-bd
    spec:
      containers:
        - name: my-bd
          image: postgres:14-alpine
          envFrom:
            - configMapRef:
                name: postgres-configmap
            - secretRef:
                name: postgres-secret
          volumeMounts:
            - name: pg-data
              mountPath: /var/lib/postgresql/data
      volumes:
        - name: pg-data
          persistentVolumeClaim:
            claimName: postgres-pvc
```
Критерии: Данные Postgres будут храниться на PVC.

## 9. Задача (Service Postgres): Создать postgres-service.yml (type: ClusterIP, имя postgres-db).
```yml
apiVersion: v1
kind: Service
metadata:
  name: postgres-db
spec:
  type: ClusterIP
  selector:
    postgres: my-bd
  ports:
    - port: 5432
      targetPort: 5432
```
Критерии: K8s создаст DNS-имя postgres-db, доступное внутри кластера.

## 10. Задача (Обновление my-app): В ConfigMap для my-app прописать DATABASE_URL (e.g., postgresql://devops:mysecret@postgres-db:5432/devopsdb). Перезапустить my-app.
```yml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: my-app-deployment
spec:
  replicas: 3
  selector:
    matchLabels:
      app: my-app
  template:
    metadata:
      labels:
        app: my-app
    spec:
      containers:
        - name: my-app
          image: my-new-app:k8s
          imagePullPolicy: IfNotPresent
          env:
            - name: DATABASE_URL
              value: "postgresql://devops:mysecret@postgres-db:5432/devopsdb"
```
Критерии: my-app (frontend) успешно подключается к postgres-db (backend) через K8s DNS.
```bash
PS C:\Users\katar\Desktop\СТАЖИРОВКА\devops_tr\module-13> kubectl get pods
NAME                                   READY   STATUS    RESTARTS   AGE
my-app-deployment-78fd9bb8b8-4ldtz     1/1     Running   0          31m
my-app-deployment-78fd9bb8b8-7r4gl     1/1     Running   0          31m
my-app-deployment-78fd9bb8b8-lmzsj     1/1     Running   0          31m
postgres-deployment-6d4f46bccd-rxhfg   1/1     Running   0          31m
```


