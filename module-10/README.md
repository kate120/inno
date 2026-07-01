# Цель модуля: Перейти от Ad-Hoc команд (для быстрых проверок) к Playbooks (для реального развертывания).

## 1. Задача (Теория): Что такое Playbook?
    playbook - файл, в котором описано желаемое состояние сервер, состояит из плеев и задач. Если после -name идет hosts, то это плей. Он состоит из задач (они тоже указываются после -name)
Контекст: Это YAML-файл, который декларативно описывает желаемое состояние сервера (список "пьес" (plays) и "задач" (tasks)).

## 2. Задача (Создание Playbook): Создать nginx-playbook.yml.

- hosts: local

become: true

tasks:

```bash
katar@Tecno:/mnt/c/Users/katar/Desktop/СТАЖИРОВКА/devops_tr/module-10$ ls
README.md  nginx-playbook.yml
```

## 3. Задача (Таск 1: apt): Добавить таск "Update apt cache": name: Update apt cache \n apt: { update_cache: yes }.
```yaml
- hosts: local
  become: true # права суперпользователя
  tasks:
    - name: Update apt cache
      apt: 
        update_cache: yes 
```

## 4. Задача (Таск 2: apt): Добавить таск "Install Nginx": name: Install Nginx \n apt: { name: nginx, state: present }.
```yml
- hosts: local
  become: true
  tasks:
    - name: Update apt cache
      apt:
        update_cache: yes
    
    - name: Install nginx 
      apt:
        name: nginx 
        state: present 
```
## 5. Задача (Таск 3: service): Добавить таск "Start Nginx": name: Start Nginx \n service: { name: nginx, state: started, enabled: yes }.
```yaml
- hosts: local
  become: true
  tasks:
    - name: Update apt cache
      apt:
        update_cache: yes
    
    - name: Install nginx 
      apt:
        name: nginx 
        state: present 

    - name: Start Nginx
      service:
        name: nginx
        state: started 
        enabled: yes
```

## 6. Задача (Запуск Playbook): Запустить: ansible-playbook -i inventory.ini nginx-playbook.yml.
```yml
PLAY RECAP *************************************************************************************
localhost                  : ok=4    changed=1    unreachable=0    failed=0    skipped=0    rescued=0    ignored=0  
```
Критерии: Playbook успешно прошел.

## 7. Задача (Переменные): Добавить vars: в playbook (e.g., nginx_port: 80).
```yaml
- hosts: local
  become: true
  vars:
    nginx_port: 80
  tasks:
    - name: Update apt cache
      apt:
        update_cache: yes

    - name: Install nginx
      apt:
        name: nginx
        state: present

    - name: Start Nginx
      service:
        name: nginx
        state: started
        enabled: yes
```
Критерии: Переменная определена.

## 8. Задача (Шаблоны - template): Использовать Jinja2-шаблон. Создать nginx.conf.j2. Скопировать его: template: { src: nginx.conf.j2, dest: /etc/nginx/nginx.conf }.
```yaml
- hosts: local
  become: true
  vars:
    nginx_port: 80
  tasks:
    - name: Update apt cache
      apt:
        update_cache: yes

    - name: Install nginx
      apt:
        name: nginx
        state: present

    - name: Start Nginx
      service:
        name: nginx
        state: started
        enabled: yes

    - name: Configure Nginx
      template:
        src: templates/nginx.conf.j2
        dest: /etc/nginx/nginx.conf
```
### nginx.conf.j2
```
events {}

http {
    server {
        listen {{ nginx_port }};

        location / {
            return 200 'Hello!';
        }
    }
}
```

Контекст: В nginx.conf.j2 использовать {{ nginx_port }}.

## 9. Задача (Обработчики - handlers): Добавить notify: Restart Nginx к таску template. Определить handlers: в плейбуке.
```yml
  - hosts: local
  become: true
  vars:
    nginx_port: 80
  tasks:
    - name: Update apt cache
      apt:
        update_cache: yes

    - name: Install nginx
      apt:
        name: nginx
        state: present

    - name: Start Nginx
      service:
        name: nginx
        state: started
        enabled: yes

    - name: Configure Nginx
      template:
        src: templates/nginx.conf.j2
        dest: /etc/nginx/nginx.conf
      notify: Restart Nginx

  handlers:
    - name: Restart Nginx
      service:
        name: nginx
        state: restarted
```

Критерии: Nginx перезапускается (handler срабатывает) только если nginx.conf.j2 изменился.
```bash
katar@Tecno:/mnt/c/Users/katar/Desktop/СТАЖИРОВКА/devops_tr/module-10$ sudo cat /etc/nginx/nginx.conf
events {}

http {
    server {
        listen 80;

        location / {
            return 200 'Hello!';
        }
    }

```
## 10. Задача (Роли - roles): (Теория). Изучить структуру Ansible Roles.

Задание: Написать эссе, объясняющее, как roles/ (e.g., roles/nginx/tasks/main.yml, roles/postgres/tasks/main.yml) позволяют переиспользовать код.

  В процессе автоматизации серверов часто возникает одна и та же проблема: вы настраиваете Nginx на трёх серверах и пишете для этого три почти одинаковых плейбука. Затем добавляете PostgreSQL — и снова копируете код. Такой подход быстро превращает проект в хаос, где одно изменение конфига заставляет править файлы в пяти разных местах.

  В Ansible эта проблема решается с помощью механизма Roles (ролей). Role — это просто папка с определённой структурой, в которую складывается всё, что относится к настройке одного конкретного сервиса. Например, внутри roles/nginx/tasks/main.yml лежат задачи по установке и настройке Nginx, а в roles/postgres/tasks/main.yml — всё для PostgreSQL.

### вид Playbook
```yml
- hosts: webservers
  become: true
  roles:
    - nginx
```
### Пример роли roles/nginx/tasks/main.yml
```yml
---
- name: Install nginx
  apt:
    name: nginx
    state: present

- name: Copy nginx config
  template:
    src: nginx.conf.j2
    dest: /etc/nginx/nginx.conf
  notify: restart nginx

- name: Start and enable nginx
  service:
    name: nginx
    state: started
    enabled: yes
```

Критерии: Эссе написано.