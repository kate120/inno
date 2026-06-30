# Цель модуля: Понять, что такое "Управление конфигурациями" (IaC) и научиться использовать Ansible для декларативного описания состояния сервера.

## 1. Задача (Теория): Что такое Ansible?

Контекст: Декларативный (я хочу nginx), Агентless (работает по SSH), Идемпотентный (повторный запуск не ломает).
    Ansible - инструмент для работы с конфигурациями. Он:
* Декларативный - в конфигах описывается желамое состояние инфраструктуры, а не шаги по достижению этого сотсояния. Ansible сам решает как достичь этого состояния;
* Агентless - на серверан не надо устанавливать ничего, достатосно ssh подклюения;
* Идемпотентный - при повторном запуске ничего не ломает, если запустить конфиг, в котором описано текущее состояние, то ничего не изменится.
Задание: Написать эссе "Почему Ansible, а не Bash-скрипты?".
```
Использование ansible предотвращает повторяющиеся задачи, такие как подключение к серверу, запуск bash скрипта. Не нужно прописывать команды, достаточно написать желаемое состояние. При повторном запуске конфига, где прописано желаемое состояние, которое является текущим, никакие изменения не произойдут
```
Критерии: Выделили Идемпотентность и Декларативность.

## 2. Задача (Установка): Установить Ansible. (e.g., pip install ansible).
```bash
katar@Tecno:~$ ansible --version
ansible [core 2.16.3]
  config file = None
  configured module search path = ['/home/katar/.ansible/plugins/modules', '/usr/share/ansible/plugins/modules']
  ansible python module location = /usr/lib/python3/dist-packages/ansible
  ansible collection location = /home/katar/.ansible/collections:/usr/share/ansible/collections
  executable location = /usr/bin/ansible
  python version = 3.12.3 (main, Jan 22 2026, 20:57:42) [GCC 13.3.0] (/usr/bin/python3)
  jinja version = 3.1.2
  libyaml = True
katar@Tecno:~$
```
Критерии: ansible --version работает.

## 3. Задача (Inventory): Создать inventory.ini.

Контекст: Это список серверов, которыми мы управляем.
```ini
[local]
localhost ansible_connection=local
```

Задание: Создать [local] группу: localhost ansible_connection=local.

Критерии: Инвентарь создан.

## 4. Задача (Ad-Hoc - ping): Проверить связь: ansible local -i inventory.ini -m ping.
```bash
localhost | SUCCESS => {
    "ansible_facts": {
        "discovered_interpreter_python": "/usr/bin/python3"
    },
    "changed": false,
    "ping": "pong"
}
```

Критерии: Ответ SUCCESS (зеленый).

## 5. адача (Ad-Hoc - setup): Собрать "факты" о машине: ansible local -i inventory.ini -m setup.
```bash
ansible local -i inventory.ini -m setup
localhost | SUCCESS => {
    "ansible_facts": {
        "ansible_all_ipv4_addresses": [
            "172.29.139.69",
            "172.17.0.1",
            "10.255.255.254"
        ],
        "ansible_all_ipv6_addresses": [
```
Далее в playbook можно указать в какими именно серверами надо работать 
```yml
 when: ansible_os_family == "Debian"
```
Критерии: Вы видите гигантский JSON со всей информацией о вашей системе.

## 6. Задача (Ad-Hoc - apt): Установить Nginx: ansible local -i inventory.ini -m apt -a "name=nginx state=present" --become.

Контекст: --become — это "стать" sudo (root).
```bash
localhost | CHANGED => {
    "ansible_facts": {
        "discovered_interpreter_python": "/usr/bin/python3"
    },
    "cache_update_time": 1782826343,
    "cache_updated": false,
    "changed": true,
    "stderr": "",
    "stderr_lines": [],
```
Критерии: changed=1 ... SUCCESS. Nginx установлен.

## 7. Задача (Идемпотентность): Запустить команду из Задачи 6 еще раз.
```bash 
katar@Tecno:/mnt/c/Users/katar/Desktop/СТАЖИРОВКА/devops_tr/module-09$ ansible local -i inventory.ini -m apt -a "name=nginx state=present" --become
localhost | SUCCESS => {
    "ansible_facts": {
        "discovered_interpreter_python": "/usr/bin/python3"
    },
    "cache_update_time": 1782826343,
    "cache_updated": false,
    "changed": false
}
```

Критерии: changed=0 ... SUCCESS (зеленый). Ansible увидел, что Nginx уже установлен, и ничего не делал.

## 8. Задача (Ad-Hoc - service): Убедиться, что сервис запущен: ansible local -i inventory.ini -m service -a "name=nginx state=started enabled=yes" --become.
```bash
katar@Tecno:/mnt/c/Users/katar/Desktop/СТАЖИРОВКА/devops_tr/module-09$ ansible local -i inventory.ini -m service -a "name=nginx state=started enabled=yes" --become
localhost | SUCCESS => {
    "ansible_facts": {
        "discovered_interpreter_python": "/usr/bin/python3"
    },
    "changed": false,
    "enabled": true, # сервис автомат запускается при запуске 
    "name": "nginx",
    "state": "started", # сервис должен быть запущен прямо сейчас
    "status": {
        "ActiveEnterTimestamp": "Tue 2026-06-30 16:54:50 +03",
        "ActiveEnterTimestampMonotonic": "116624092997",
        "ActiveExitTimestampMonotonic": "0",
```

## 9. Задача (Ad-Hoc - file): Создать папку: ansible local -i inventory.ini -m file -a "path=/tmp/test-ansible state=directory".
```bash
katar@Tecno:/mnt/c/Users/katar/Desktop/СТАЖИРОВКА/devops_tr/module-09$ ansible local -i inventory.ini -m file -a "path=/tmp/test-ansible state=directory"
localhost | CHANGED => {
    "ansible_facts": {
        "discovered_interpreter_python": "/usr/bin/python3"
    },
    "changed": true,
    "gid": 1000,
    "group": "katar",
    "mode": "0755",
    "owner": "katar",
    "path": "/tmp/test-ansible",
    "size": 4096,
    "state": "directory",
    "uid": 1000
}
```

## 10. Задача (Ad-Hoc - copy): Скопировать файл: ansible local -i inventory.ini -m copy -a "src=inventory.ini dest=/tmp/test-ansible/inventory.bak".
```bash
katar@Tecno:/mnt/c/Users/katar/Desktop/СТАЖИРОВКА/devops_tr/module-09$ ansible local -i inventory.ini -m copy -a "src=inventory.ini dest=/tmp/test-ansible/inventory.bak"
localhost | CHANGED => {
    "ansible_facts": {
        "discovered_interpreter_python": "/usr/bin/python3"
    },
    "changed": true,
    "checksum": "cda59fcaddb24b6d1a522b17b9c9ec90e2183597",
    "dest": "/tmp/test-ansible/inventory.bak",
    "gid": 1000,
    "group": "katar",
    "md5sum": "b8c076350724fa8d56fe8c1707811be2",
    "mode": "0644",
    "owner": "katar",
    "size": 43,
    "src": "/home/katar/.ansible/tmp/ansible-tmp-1782828952.1266928-10682-64657544541215/source",
    "state": "file",
    "uid": 1000
}
```
```bash
katar@Tecno:/mnt/c/Users/katar/Desktop/СТАЖИРОВКА/devops_tr/module-09$ ls -la /tmp/test-ansible/
total 12
drwxr-xr-x 2 katar katar 4096 Jun 30 17:15 .
drwxrwxrwt 9 root  root  4096 Jun 30 17:16 ..
-rw-r--r-- 1 katar katar   43 Jun 30 17:15 inventory.bak
```
