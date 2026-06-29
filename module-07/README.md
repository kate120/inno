# Цель модуля: Понять концепцию CI/CD (Автоматизация) и настроить первый пайплайн для тестирования кода (Linting, Unit Tests).

## 1. Задача (Теория): Что такое CI (Continuous Integration)?

Контекст: Автоматическая сборка и тестирование кода при каждом git push (или PR).

    Continuous Integration (CI) — это практика разработки, при которой разработчики регулярно (обычно несколько раз в день) интегрируют свой код в общий репозиторий, и каждая интеграция автоматически проверяется сборкой и тестами.

Критерии: Вы понимаете, что CI — это "защита" ветки main.

## 2. Задача (Теория): Что такое CD (Continuous Delivery/Deployment)?

Контекст: Автоматическая выкатка (деплой) приложения, если CI прошел успешно.

* Continuous Delivery (CD) — это практика, при которой код автоматически подготавливается к релизу в продакшн после прохождения всех проверок. Ключевое слово здесь — подготавливается. Код готов к деплою в любой момент, но само решение о релизе принимает человек.
* Continuous Deployment — это практика, при которой каждое изменение кода, прошедшее все автопроверки, автоматически попадает в продакшн без какого-либо ручного вмешательства. Написал код, запушил, прошли тесты — и через 15-30 минут изменения уже доступны потребителям.

  ### Окружения в CD
CD подразумевает наличие нескольких окружений, через которые проходит код на пути к продакшну. Типичная схема выглядит так:

1. Development (dev) — самое нестабильное окружение. Сюда деплоится каждый коммит в основную ветку. Здесь могут быть баги, здесь тестируются новые идеи. Разработчики используют его для проверки интеграции.

2. Testing/QA — окружение для тестировщиков. Более стабильное, чем dev. Сюда попадают изменения, которые прошли базовые проверки. Тестировщики проводят здесь ручное и автоматизированное тестирование.

3. Staging — копия продакшна. Такое же железо, такие же данные (обезличенные), такие же настройки. Здесь проводится финальная проверка перед релизом. Если работает на стейджинге, должно работать и на продакшне.

4. Production — боевое окружение, где работают реальные пользователи. Сюда попадает только проверенный код после всех стадий тестирования.

Критерии: Вы понимаете разницу (Delivery — готово к деплою, Deployment — задеплоено).

## 3. Задача (Теория): Анатомия GitHub Actions (или GitLab CI).

Задание: Изучить: Workflow (файл), Event (e.g., on: push), Runner (виртуалка), Job (e.g., build), Step (команда).
* Workflow - файл с инструкциями по пути .github/workflows/*.yaml. В одном репозитории может быть много файлов workflow.
```yaml
# .github/workflows/ci.yml
name: My CI Pipeline
```
* Event - событие, после наступления которого запускается workflow (пайплайн)
```yaml
on:
  push:
    branches: [ main ]
  pull_request:
    branches: [ main ]
  workflow_dispatch:   # запуск вручную через UI
  schedule:
    - cron: '0 9 * * 1'   # каждый понедельник в 9:00
```
* Runner - изолированнная среда для запуска пайплайна (виртуалка)
```yaml
jobs:
  build:
    runs-on: ubuntu-latest   # ← выбор runner'а
```
* Job - действия, из которйх состоит workflow, по умолчанию выполняются одновременно, набор шагов
```yaml
jobs:
  build:
    runs-on: ubuntu-latest

  deploy:
    runs-on: ubuntu-latest
    needs: build          # ← сначала дождётся build
    if: github.ref == 'refs/heads/main'
```
* Step - конкретные дейтсивия (команды) из котторых состоят job. Выполрняются последовательно
```yaml
steps:
  - name: Checkout code
    uses: actions/checkout@v4    # готовое Action из маркетплейса

  - name: Install deps
    run: npm install             # обычная shell-команда

  - name: Run tests
    run: npm test
```
### workflow:
```yaml

name: Проверка кода
on: [push]  # <-- Это Event
jobs:       # <-- Список задач
  build:    # <-- Это Job
    runs-on: ubuntu-latest # <-- Это Runner
    steps:  # <-- Это Steps
      - name: Скачать код
        run: git clone ...
      - name: Установить зависимости
        run: npm install
      - name: Запустить тесты
        run: npm test
    
```

Критерии: Вы понимаете иерархию.

## 4. Задача (Создание Workflow): Создать папку .github/workflows/. В ней — ci.yml. (Или .gitlab-ci.yml в корне).
```bash

$ find -name ci.yml
./.github/workflows/ci.yml

```

## 5. Задача (Триггер): Настроить Event (триггер): on: pull_request: { branches: [main] }. (GitLab: only: [merge_requests]).
``` yaml

name: code check
on: 
  pull_request: 
    branches: [main] 

```

## 6. Задача (Jobs): Настроить Job (задачу):

jobs:

test:

runs-on: ubuntu-latest
``` yaml

name: code check
on: 
  pull_request: 
    branches: [main] 
jobs:
  test:
    runs-on: ubuntu-latest

```

## 7. Задача (Steps - checkout): Добавить первый шаг — получение кода: uses: actions/checkout@v4. (GitLab: script:).
```yaml

  name: code check
  on: 
    pull_request: 
      branches: [main] 
  jobs:
    test:
      runs-on: ubuntu-latest
      steps:
        - name: checkout 
          uses: actions/checkout@v4

```

## 8. Задача (Steps - setup): Настроить среду (e.g., Python): uses: actions/setup-python@v5 (или image: python:3.10 в GitLab).
```yaml
  name: code check
  on: 
    pull_request: 
      branches: [main] 
  jobs:
    test:
      runs-on: ubuntu-latest
      steps:
        - name: checkout 
          uses: actions/checkout@v4
        - name: setups
          uses: actions/setup-python@v5
```

## 9. Задача (Steps - run): Запустить тесты:

run: pip install -r requirements.txt

run: pytest (Предполагается, что у вас есть pytest и tests/).
```yaml
name: code check
on:
  pull_request: 
    branches: [main] 
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - name: checkout 
        uses: actions/checkout@v4
      - name: setups
        uses: actions/setup-python@v5
      - name: install dependencies
        run: pip install -r requirements.txt
      - name: run tests
        run: pytest
```

## 10. Задача (Практика): Создать PR и увидеть, как GitHub Actions запускает пайплайн. Увидеть "зеленую галочку" или "красный крестик".

Критерии: CI-пайплайн автоматически проверяет ваш код.