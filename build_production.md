# Сборка и развёртывание проекта в продакшн (VPS на jino.ru, Ubuntu)

Документ описывает пошаговую сборку и установку проекта на **виртуальный сервер
Джино (классический VPS, OpenVZ) с Ubuntu**. Ключевое отличие от «голого»
сервера: **TLS терминирует веб-сервер панели Джино** — сертификат выпущен и
обновляется через админ-панель хостинг-провайдера, на сервере certbot нет и не
нужен. Наш собственный nginx работает только на внутреннем порту (по
умолчанию `8080`) и не имеет никаких `ssl_certificate`.

Регистрация новых пользователей **выключена** (feature-flag
`registration_enabled: false` на сервере и `NUXT_PUBLIC_REGISTRATION_ENABLED=false`
на клиенте).

> Документ заполнен реальными значениями:
>
> - пользователь системы: `raiyin` (sudo-доступ; вход по SSH — под `root`)
> - канонический домен: `pertsukova.ru` (без `www`, на него 301-редиректом
>   переходит `www.pertsukova.ru`)
> - каталог развёртывания: `/opt/artserver`
> - панель управления: `https://cp.jino.ru` → раздел **VPS**
> - сервер имеет **выделенный IP** (без него 80/443 пришлось бы шарить с
>   другими клиентами Джино); A-записи `pertsukova.ru` и `www.pertsukova.ru`
>   указывают на этот IP
> - внутренний порт нашего nginx: **8080** (80 и 443 снаружи заняты панелью)

---

## 0. Что уже сделано в панели Джино (сделать до установки)

Панель делает за вас часть работы, которую на «голом» сервере делает certbot.
Пройдитесь по чек-листу, иначе сайт не откроется.

| # | Что сделать | Где в панели | Зачем |
|---|-------------|--------------|-------|
| 1 | Создать/убедиться в VPS на Ubuntu 22.04/24.04, x86_64 | `cp.jino.ru` → VPS → Создать сервер | Нужна amd64-машина: Go с CGO собирается нативно, нативные бинарники `sharp` под linux-x64 |
| 2 | Задать пароль root | VPS → Управление → Основные настройки → «Задать пароль root» | Без него не войти по SSH |
| 3 | Узнать хост и порт SSH | VPS → Управление → Доступ | Порт нестандартный, в rsync/ssh его нужно указывать явно |
| 4 | Подключить **выделенный IP** | VPS → Управление → Выделенные IP | Без выделенного IP сайт живёт на общем адресе и требует «Перенаправления портов» с указанием порта в URL |
| 5 | Привязать домен к серверу | VPS → Управление → Домены (или раздел «Домены» аккаунта) | Панель должна знать, за какие домены отвечать |
| 6 | Выпустить/привязать SSL-сертификат | Раздел сертификатов в панели | **TLS обслуживается здесь**, а не certbot на сервере. Выпустите на `pertsukova.ru` + `www.pertsukova.ru` |
| 7 | Включить HTTPS→HTTP-проксирование на внутренний порт | VPS → Управление → Перенаправление портов → «Проксирование HTTP(S)» | Панель приняла запрос на 443, сняла TLS и передала **обычный HTTP** на порт `8080` нашего nginx |
| 8 | Включить переход HTTP → HTTPS (если панель предлагает) | Там же / в настройках сайта | Панель делает редирект с 80 на 443 до того, как запрос дойдёт до нас |

> Названия пунктов меню в панели Джино могут немного отличаться от приведённых —
> ориентируйтесь на смысл. Итоговая схема проверяется командами из
> раздела 9.2: если они дают ожидаемый результат, схема собрана верно.

> **Важно про «Проксирование HTTP(S)».** Для выделенного IP проксирование можно
> включать и выключать. Пока оно **выключено**, запросы на 80/443 приходят
> напрямую в вашу виртуальную машину — и тогда TLS обязан делать ваш nginx, а
> сертификат из панели не подойдёт. В нашей схеме проксирование **включено**,
> внутренний порт — `8080`.

> **Порт 25/tcp на VPS Джино заблокирован провайдером** (антиспам). Своя
> отправка писем с сервера невозможна — SMTP отправляйте через внешний сервис
> по порту 587/465, как и предполагает пустая секция `smtp` в `config.yaml`
> (раздел 6.2).

---

## 1. Обзор архитектуры и схема продакшна

| Компонент | Технология | Порт | Что делает |
|-----------|-----------|------|------------|
| **Панель Джино (внешний фронт)** | веб-сервер Джино | 80/443 | TLS, сертификат, редирект HTTP→HTTPS, проброс на внутренний порт |
| **Web (SSR)** | Nuxt 4 (Node) | 3000 | Отдаёт страницы сайта, статику `/_nuxt/` |
| **API** | Go + Gin | 8000 | Бэкенд: галерея, магазин, новости, админка, платежи, чат |
| **Reverse proxy** | nginx | 8080 | Внутренний: раздача `/content/`, прокси на 8000 и 3000. **Без TLS** |
| **БД** | SQLite (WAL) | файл | `server/db/db.sqlite` |
| **Файлы** | локальная ФС | — | изображения галереи/магазина/новостей, аватары |

Схема запросов (все — на один домен, API прячется за префиксом `/api-server/`):

```
Браузер
  │
  ├── https://pertsukova.ru/…     →  панель Джино: TLS (сертификат панели)
  │                                        │
  │        ┌─ http://<ip>:8080/… ─────────┘   ← X-Forwarded-Proto: https
  │        │
  │        └── nginx :8080  (наш, только HTTP)
  │              ├── /api-server/*  → срезает префикс → Go   :8000
  │              ├── /content/*     → статический файл (изображения работ/продаж/новостей)
  │              ├── /hero-images/* → статический файл
  │              └── /, /_nuxt/*    → Node (Nuxt SSR) :3000
  │
  └── http://pertsukova.ru/        →  панель Джино: 301 на https (редирект делает панель)
```

**Почему nginx слушает 8080, а не 80/443.** Сертификат выпущен и обновляется
панелью Джино, значит TLS держит она, а не мы: 80/443 снаружи обслуживает
веб-сервер панели. Наш nginx — внутренний слой, который принимает уже
«голый» HTTP и раскладывает запросы по приложениям. Плюс так не нужно
ни certbot, ни `ssl_certificate` в конфиге, и мы не зависим от того, куда
панель положит файлы сертификата.

**Почему префикс `/api-server/`.** У фронтенда и API совпадают пути
(`/news`, `/admin`, `/profile` и т.д.), поэтому API и сайт не могут жить на одном
и том же пути без конфликта. Префикс снимается в nginx командой `rewrite` —
Go-сервер не знает о нём и работает как обычно. Клиентская переменная
`NUXT_PUBLIC_SERVER_URL` указывает на `https://pertsukova.ru/api-server/`
(обязательно с завершающим слэшем).

**Что делать с `X-Forwarded-Proto`.** Приложениям нужно знать исходную схему
(`https`), хотя nginx видит только `http`. Панель обычно передаёт заголовок
`X-Forwarded-Proto`; если нет — считаем `http`. Из этого значения выводится
переменная `$arts_proto` (раздел 9.3), она же используется для «страхового»
редиректа на HTTPS. Клиент (Nuxt) собирает все URL из абсолютного
`NUXT_PUBLIC_SERVER_URL`, а cookie ставит без флага `Secure`, поэтому
разрыв TLS на нашей стороне на работу сайта не влияет.

---

## 2. Требования

| Зависимость | Версия | Зачем |
|-------------|--------|-------|
| Панель Джино | классический VPS (OpenVZ) | хостинг; выдаёт выделенный IP и сертификат |
| Ubuntu | 22.04 / 24.04, **x86_64 (amd64)** | ОС сервера |
| SSH-доступ | root, нестандартный порт из панели | rsync, systemd, сборка Go |
| nginx | последняя из репозитория | внутренний reverse proxy, раздача `/content/` |
| Node.js | 22 LTS | запуск Nuxt SSR (`.output`) |
| pnpm | 10.x (в проекте `pnpm@10.33.0`) | пакетный менеджер клиента — **только на машине разработки** (клиент собирается локально) |
| Go | 1.23.x (в `go.mod` — `1.23.2`) | сборка сервера |
| gcc / build-essential | — | нужен для CGO (`mattn/go-sqlite3`) |
| sqlite3 | CLI (необязательно) | резервные копии и админ-пользователь |
| ~~certbot~~ | **не нужен** | TLS обслуживает панель Джино (раздел 0, п. 6) |

> **RAM.** Минимальный тариф (2 ГБ) для Go + Node SSR хватает, но впритык:
> добавьте swap (раздел 3.4) и не собирайте клиент на сервере (раздел 7).
> На OpenVZ объём памяти — жёсткий лимит контейнера, при его превышении ядро
> убивает процессы (`OOMKilled` в логах).

---

## 3. Подготовка сервера

### 3.1. Первый вход и пользователь

На VPS Джино основная учётная запись — **root** (пароль задаётся в панели).
Хост и порт берутся из панели: *Управление → Доступ*.

```bash
ssh -p SSH_ПОРТ_ИЗ_ПАНЕЛИ root@ВЫДЕЛЕННЫЙ_IP
uname -m            # должно быть x86_64
cat /etc/os-release # Ubuntu 22.04/24.04
```

Сервисы должны работать от непривилегированного пользователя (root не нужен).
Пользователя `raiyin` на свежем VPS обычно нет — создайте его (иначе rsync из
раздела 4 не сможет писать в `/opt/artserver`):

```bash
adduser --disabled-password --gecos "" raiyin
usermod -aG sudo raiyin
install -d -m 700 -o raiyin -g raiyin /home/raiyin/.ssh
```

Положите свой SSH-ключ, чтобы rsync работал без пароля:

```bash
# локально, один раз
cat ~/.ssh/id_ed25519.pub

# на сервере
echo 'ssh-ed25519 AAAA…ваш_ключ… raiyin' >> /home/raiyin/.ssh/authorized_keys
chown raiyin:raiyin /home/raiyin/.ssh/authorized_keys
chmod 600 /home/raiyin/.ssh/authorized_keys
```

Беспарольный sudo пригодится скрипту бэкапа (раздел 13) и установке пакетов:

```bash
echo 'raiyin ALL=(ALL) NOPASSWD:ALL' > /etc/sudoers.d/raiyin
chmod 440 /etc/sudoers.d/raiyin
```

> **Защита SSH.** Порт 22 на VPS Джино обслуживает множество клиентов, поэтому
> его регулярно атакуют. Панель сама предлагает смену порта (сделайте это), а
> дополнительно можно поставить `fail2ban` — инструкция Джино:
> <https://jino.ru/spravka/articles/f2b.html>.

### 3.2. Модули ядра и файрвол

Если планируете включать `ufw` или `fail2ban`, панель должна отдать контейнеру
модуль ядра `iptables`: *Управление → Основные настройки → «подключить модули
ядра»* (fuse, tun/tap, **iptables**, ppp). Проверить:

```bash
lsmod | grep -E '^(iptable_filter|iptable_nat|nf_conntrack)' || \
  sudo iptables -L -n >/dev/null && echo 'iptables работает'
```

Сам файрвол на свежем VPS обычно выключен, и это нормально: 80/443 к вам
приходят от веб-сервера панели. **Внутренний порт 8080 должен быть доступен
извне** (его использует фронт панели). Если всё же включаете `ufw`:

```bash
sudo ufw allow SSH_ПОРТ_ИЗ_ПАНЕЛИ/tcp   # иначе потеряете доступ
sudo ufw allow 8080/tcp                  # внутренний порт для панели
sudo ufw allow 443/tcp
sudo ufw enable
sudo ufw status
```

> Не открывайте наружу 3000 и 8000 — они работают только на `127.0.0.1`.
> Порт 25/tcp заблокирован провайдером и открыть его нельзя.

### 3.3. Установка системных пакетов

> **Почему `build-essential`, а не голый `gcc`.** Go-сервер использует драйвер
> `mattn/go-sqlite3` (CGO), которому при компиляции нужны не только `gcc`, но и
> заголовки Си-стандартной библиотеки (`stdio.h` и др.) из пакета `libc6-dev`.
> `build-essential` тянет всё необходимое разом: `gcc`, `g++`, `make`, `libc6-dev`.
> На минимальных образах Ubuntu голый `gcc` не ставит `libc6-dev`, и сборка
> падает с `fatal error: stdio.h: No such file or directory`. `ca-certificates`
> нужен, чтобы `go build` мог достучаться до `proxy.golang.org`.

```bash
sudo apt update
sudo apt install -y nginx build-essential sqlite3 curl git ca-certificates

# Node.js 22 LTS (nodesource)
curl -fsSL https://deb.nodesource.com/setup_22.x | sudo -E bash -
sudo apt install -y nodejs

# pnpm — НЕ нужен на сервере (клиент собирается локально, раздел 7).
# Устанавливаем только если решите собирать клиент прямо на сервере.
# sudo npm install -g pnpm@10.33.0

# Go 1.23 для x86_64 (в системных репозиториях Ubuntu — старая версия)
cd /tmp
curl -fsSL https://go.dev/dl/go1.23.2.linux-amd64.tar.gz -o go.tgz
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go.tgz
echo 'export PATH=$PATH:/usr/local/go/bin' | sudo tee /etc/profile.d/go.sh
source /etc/profile.d/go.sh
go version   # должно показать go1.23.x
```

> **certbot не ставим.** Раньше здесь был шаг `apt install certbot` +
> `certbot certonly --standalone`. Он больше не нужен и даже вреден: панель Джино
> держит 80/443 и обновляет сертификат сама, а standalone-выпуск всё равно
> требует остановить её веб-сервер. Продление сертификата — в панели
> (раздел 9.6).

### 3.4. Swap

На тарифах с 2 ГБ RAM добавьте swap-файл, иначе OOM-killer убьёт Node или
Go-процесс под нагрузкой. Если в панели (*Статистика → объём файла SWAP*)
есть регулировка — задайте 1–2 ГБ там. Иначе вручную:

```bash
sudo fallocate -l 2G /swapfile
sudo chmod 600 /swapfile
sudo mkswap /swapfile
sudo swapon /swapfile
echo '/swapfile none swap sw 0 0' | sudo tee -a /etc/fstab
free -h
```

### 3.5. Структура каталогов

> **Почему `/opt/artserver` — лучший выбор.** По FHS `/opt` предназначен для
> самостоятельных (add-on) приложений, разворачиваемых вручную и запускаемых как
> сервисы. Приложение целиком живёт в одном месте — бинарник, конфиг, БД и
> контент — поэтому его проще бэкапить, обновлять и менять владельца.
> Альтернативы хуже: `/srv/artserver` рассчитан на данные, которые раздаёт
> система (а здесь ещё и код приложения), `/home/raiyin/artserver` смешивает
> приложение с пользовательскими данными и может попасть под чужие квоты и права.

```bash
sudo mkdir -p /opt/artserver/server/db
sudo mkdir -p /opt/artserver/server/storage/avatars
sudo mkdir -p /opt/artserver/server/bin
sudo mkdir -p /opt/artserver/client
sudo chown -R raiyin:raiyin /opt/artserver
```

---

## 4. Перенос проекта на сервер

Полный код нужен на сервере: сервер (Go) там компилируется, а `.output` клиента
переносится отдельно после локальной сборки (раздел 7.3). `node_modules`, `.nuxt`,
`.output`, `db/*.sqlite` и `storage/` переносить не нужно (они тяжёлые и/или
gitignored).

Пропишите подключение один раз в `~/.ssh/config` (на машине разработки) —
тогда не придётся добавлять `-e 'ssh -p …'` в каждую команду:

```
Host artserver
    HostName ВЫДЕЛЕННЫЙ_IP
    User root
    Port SSH_ПОРТ_ИЗ_ПАНЕЛИ
```

Проверьте, что ключ пускает root:

```bash
ssh artserver 'id root; uname -m'
```

Локально выполните:

```bash
rsync -avz --delete -e 'ssh -p SSH_ПОРТ_ИЗ_ПАНЕЛИ' \
  --exclude '.git' \
  --exclude '.DS_Store' \
  --exclude 'screenshots' \
  --exclude 'plans' \
  --exclude 'test_data' \
  --exclude 'sourcecraft' \
  --exclude 'node_modules' \
  --exclude '.nuxt' \
  --exclude '.output' \
  --exclude 'client/.env' \
  --exclude 'client/public/content' \
  --exclude 'server/config.yaml' \
  --exclude 'server/.env' \
  --exclude 'server/db' \
  --exclude 'server/storage' \
  --exclude 'server/bin' \
  --exclude 'server/server' \
  ./ \
  root@ВЫДЕЛЕННЫЙ_IP:/opt/artserver/

# и сразу выдаём права пользователю сервисов
ssh artserver 'chown -R raiyin:raiyin /opt/artserver'
```

> **Что и почему исключено (важно при повторных заливках с `--delete`).**
> `--delete` удаляет на сервере всё, чего нет локально, поэтому всё, что
> создаётся/хранится **на сервере**, из переноса исключено, иначе повторный rsync
> либо затрёт живое состояние, либо снесёт его совсем:
>
> - `server/config.yaml`, `server/.env` — продакшен-конфиг и секреты создаются
>   на сервере (разделы 6.2, 6.3); локальные дев-версии в репозитории не совпадают
>   с ними (пути `/Users/...`, `registration_enabled: true`).
> - `server/db`, `server/storage`, `client/public/content` — база данных,
>   аватары и загруженные через админку изображения. Живут только на сервере;
>   `server/db` и `client/public/content` переносятся отдельно (раздел 5.2)
>   только когда вы это явно делаете.
> - `server/bin`, `server/server` — бинарники Go (локальные артефакты сборки,
>   бессмысленные на сервере).
> - `node_modules`, `.nuxt`, `.output` — тяжёлые/пересобираемые; `.output`
>   переносится отдельно в разделе 7.3.

> **Пользователь в rsync.** В команде выше `root@…` — потому что на VPS Джино
> изначально есть только root. Если хотите заливать сразу от `raiyin`, замените
> пользователя и выполните `chown` (в этой же серии команд). Не забудьте
> `-e 'ssh -p …'`: порт у Джино нестандартный, при его потере rsync молча
> зависнет или получит `Connection refused`.

> **Сертификаты сюда не попадают.** Файлов `.pem`/`.key` в проекте нет и быть не
> должно: сертификат хранится и обновляется на стороне панели Джино.

---

## 5. База данных

### 5.1. Где лежит БД

Путь к БД захардкожен в `server/cmd/server/main.go:50` относительно рабочего
каталога процесса:

```go
db, err := sql.Open("sqlite3", "./db/db.sqlite")
```

Поэтому **рабочий каталог сервиса обязан быть `/opt/artserver/server`**, а файл
БД живёт в `/opt/artserver/server/db/db.sqlite`. При старте создаются файлы
`db.sqlite-shm` и `db.sqlite-wal` (режим WAL).

### 5.2. Перенос существующих данных (рекомендуется)

Если у вас уже есть наполнение (работы, новости, продажи, пользователи), перенесите
**БД вместе с изображениями** — пути в БД относительные (`/content/works/...`):

```bash
# локально: база
rsync -avz -e 'ssh -p SSH_ПОРТ_ИЗ_ПАНЕЛИ' server/db/db.sqlite root@ВЫДЕЛЕННЫЙ_IP:/opt/artserver/server/db/db.sqlite
# локально: изображения (галерея, магазин, новости)
rsync -avz -e 'ssh -p SSH_ПОРТ_ИЗ_ПАНЕЛИ' client/public/content/ root@ВЫДЕЛЕННЫЙ_IP:/opt/artserver/client/public/content/
```

Перед копированием локальной БД сделайте checkpoint WAL, иначе данные в `-wal`
могут потеряться:

```bash
sqlite3 server/db/db.sqlite "PRAGMA wal_checkpoint(TRUNCATE);"
```

### 5.3. Чистая установка (без переноса данных)

Достаточно создать пустую директорию `/opt/artserver/server/db/` — таблицы создаёт
сам сервер при первом запуске. Схему см. в `server/internal/integration/setup_test.go`
(эталон). Придётся вручную создать пользователя-админа — см. раздел 5.4.

### 5.4. Админ-пользователь

Регистрация выключена, поэтому аккаунт администратора должен существовать в БД
**до** первого логина. В перенесённой из разработки БД уже есть пользователь
`admin` (роль `admin`).

Если нужен новый/чистый админ, создайте хеш пароля и вставьте запись:

```bash
# 1) сгенерировать bcrypt-хеш пароля — временный помощник
mkdir -p /opt/artserver/server/cmd/genhash
cat > /opt/artserver/server/cmd/genhash/main.go <<'EOF'
package main

import (
    "fmt"
    "os"

    "golang.org/x/crypto/bcrypt"
)

func main() {
    if len(os.Args) < 2 {
        fmt.Println("usage: go run ./cmd/genhash <password>")
        os.Exit(1)
    }
    h, err := bcrypt.GenerateFromPassword([]byte(os.Args[1]), bcrypt.DefaultCost)
    if err != nil {
        panic(err)
    }
    fmt.Println(string(h))
}
EOF
cd /opt/artserver/server
go run ./cmd/genhash 'СИЛЬНЫЙ_ПАРОЛЬ'
# → строка вида $2a$10$...

# 2) вставить пользователя (сервер при этом должен быть остановлен)
sudo -u raiyin sqlite3 /opt/artserver/server/db/db.sqlite \
  "INSERT INTO users (username, password_hash, role, email, name, created_at, updated_at, email_verified)
   VALUES ('admin', 'ХЕШ_ИЗ_ШАГА_1', 'admin', 'admin@example.com', 'Admin',
           datetime('now'), datetime('now'), 1);"

# 3) удалить помощник
rm -rf /opt/artserver/server/cmd/genhash
```

> Если пароль `admin` в перенесённой БД неизвестен, сбросьте его тем же способом —
> `UPDATE users SET password_hash = 'Новый_хеш' WHERE username = 'admin';`.

---

## 6. Сервер (Go)

### 6.1. Каталог и файлы

Итоговая структура `/opt/artserver/server/`:

```
/opt/artserver/server/
├── bin/artserver          # скомпилированный бинарник
├── config.yaml            # конфигурация (обязателен, ищется в рабочем каталоге)
├── .env                   # секреты (JWT_SECRET и др.)
├── db/db.sqlite           # база данных SQLite (WAL)
└── storage/avatars/       # аватары пользователей
```

### 6.2. Конфигурация `config.yaml`

**Флаг регистрации выключен**, пути хранения — абсолютные, указывают на публичный
каталог клиента. Создайте файл **на сервере** командой ниже (раздел 4 исключает
`server/config.yaml` из rsync, поэтому файл не затирается при повторных переносах
кода — это единственное место, где он создаётся):

```bash
sudo -u raiyin tee /opt/artserver/server/config.yaml >/dev/null <<'EOF'
app:
  name: "Vera Art"
  port: 8000

cors:
  allowed_origins:
    - "https://pertsukova.ru"
    - "https://www.pertsukova.ru"
  debug: false

# Где сервер сохраняет загруженные изображения.
# Директория = корень статики клиента (client/public),
# оттуда их раздаёт nginx по URL /content/...
directories:
  abs_works_dir: "/opt/artserver/client/public"
  rel_works_dir: "/content/works/"
  abs_sales_dir: "/opt/artserver/client/public"
  rel_sales_dir: "/content/sales/"
  abs_news_dir: "/opt/artserver/client/public"
  rel_news_dir: "/content/news/"
  abs_avatars_dir: "/opt/artserver/server/storage/avatars/"

# Feature flags — регистрация ВЫКЛЮЧЕНА
features:
  registration_enabled: false

# SMTP: заполните при реальной рассылке писем.
# При host="" письма не отправляются, а логируются в консоль.
# Свой SMTP-сервер поднять нельзя: порт 25/tcp у Джино заблокирован.
smtp:
  host: ""
  port: 587
  username: ""
  password: ""
  from: "noreply@vera-art.com"
  from_name: "Vera Art"
EOF
sudo chmod 644 /opt/artserver/server/config.yaml
```

> `cors.allowed_origins` перечисляет **`https://`**-адреса независимо от того, где
> терминируется TLS: браузер присылает `Origin` именно как в адресной строке.

### 6.3. Секреты `.env`

```bash
sudo -u raiyin tee /opt/artserver/server/.env >/dev/null <<'EOF'
JWT_SECRET=сгенерируйте-длинную-случайную-строку-не-короче-32-символов
EOF
sudo chmod 600 /opt/artserver/server/.env
```

> `.env` тоже исключён из rsync (раздел 4) — при повторной заливке кода
> продакшен-секрет не затрётся локальным дев-значением из `server/.env` машины
> разработки.

Пример генерации: `openssl rand -base64 48`.

> Если `JWT_SECRET` не задан, сервер использует дефолтный `mydevsecret` и пишет
> предупреждение. В продакшне обязательно задайте свой.

### 6.4. Сборка

CGO обязателен (драйвер `mattn/go-sqlite3`), поэтому нужны `gcc` + заголовки C
`libc6-dev` — их даёт `build-essential` из раздела 3.3:

```bash
cd /opt/artserver/server
mkdir -p bin
CGO_ENABLED=1 go build -trimpath -ldflags "-s -w" -o bin/artserver ./cmd/server
```

> **Первый `go build` скачивает модули** (`go.mod`/`go.sum` уже перенесены в
> разделе 4) из `proxy.golang.org` — нужен интернет и настройка прокси не должна
> мешать. Сборка маленького проекта идёт секунды и в OOM не упирается (в отличие
> от клиента — см. 7.2).

> **Архитектура.** Сервер — Ubuntu x86_64 (amd64), поэтому бинарник собирается
> **на самом сервере** нативно и получается `linux/amd64` сам собой.
> Компилировать Go на машине разработки (Apple Silicon, arm64) не нужно:
> кросс-сборка с `CGO_ENABLED=1` и `mattn/go-sqlite3` требует отдельного
> amd64-кросс-компилятора gcc — это лишняя возня. Проверьте архитектуру, если
> образ сервера вдруг другой:

```bash
uname -m   # должно быть x86_64
```

Проверка:

```bash
cd /opt/artserver/server
timeout 3 ./bin/artserver || true   # стартует, читает config.yaml/.env; сам остановится через 3 c
# в логе должны быть "Server starting ..." и "Registration is DISABLED by feature flag"
file bin/artserver                  # должно показать x86-64
```

---

## 7. Клиент (Nuxt 4, SSR)

> Клиент **собирается на машине разработки (локально)**, на сервер переносится
> уже готовый `.output`. Серверу не нужны ни `node_modules`/pnpm, ни память на
> сборку: на VPS с 2 ГБ RAM клиентская сборка падает по OOM (см. 7.2).

### 7.1. Переменные окружения для сборки (локально)

`NUXT_PUBLIC_SERVER_URL` и `NUXT_PUBLIC_REGISTRATION_ENABLED` встраиваются в
клиентский JS-бандл на этапе сборки, поэтому задать их нужно **до** `pnpm build`
на своей машине. Проще всего передать их прямо в команду — переменные окружения
перекрывают значения из `client/.env`:

```bash
cd client
NUXT_PUBLIC_SERVER_URL=https://pertsukova.ru/api-server/ \
NUXT_PUBLIC_LIMIT=9 \
NUXT_PUBLIC_REL_WORKS_DIR=/content/works/ \
NUXT_PUBLIC_REL_SALES_DIR=/content/sales/ \
NUXT_PUBLIC_REGISTRATION_ENABLED=false \
pnpm build
```

> `NUXT_PUBLIC_SERVER_URL` обязан заканчиваться на `/` — фронтенд конкатенирует
> его с путями API (`serverUrl + 'news/'` и т.д.).
>
> Здесь указывается **внешний** `https://` адрес, хотя наш nginx слушает
> `http://127.0.0.1:8080`: клиент ходит в API снаружи, через TLS панели Джино.
>
> Альтернатива длинной команде — временно прописать те же значения в
> `client/.env`, собрать, затем вернуть прежние.

### 7.2. Установка зависимостей и сборка (локально)

```bash
cd client
pnpm install --frozen-lockfile
pnpm build    # с переменными из 7.1
```

> **Память при сборке.** Nuxt + Tailwind v4 требуют много памяти. В
> `client/.npmrc` уже прописан `node-options=--max-old-space-size=4096`, который
> pnpm подхватывает автоматически; если вылезет
> `FATAL ERROR: JavaScript heap out of memory` — увеличьте до 8192. Ошибку
> `[commonjs--resolver] The service is no longer running: write EPIPE` вызывает
> OOM-killer ядра на VPS с малым объёмом RAM (лимит контейнера OpenVZ) — именно
> поэтому собираем локально, а не на сервере. Предупреждение
> `Sourcemap is likely to be incorrect` безобидно.

> **Нативные бинарники (sharp).** SSR-часть `.output/server` содержит нативный
> модуль **sharp** (`@nuxt/image` для ресайза картинок) — бинарники привязаны к
> ОС и архитектуре. Чтобы сборка с Mac (arm64) деплоилась на Ubuntu x64, в
> `client/.npmrc` уже добавлен `supportedArchitectures` — pnpm ставит бинарники
> сразу для linux-x64 и darwin-arm64. Запись — в ини-синтаксисе `.npmrc` через
> скобки `[os][]` (НЕ в виде YAML-блока! pnpm такой блок не читает — проверьте
> командой `pnpm config get supportedArchitectures`, она должна вывести значения,
> а не `undefined`):
>
> ```
> supportedArchitectures[os][]=darwin
> supportedArchitectures[os][]=linux
> supportedArchitectures[cpu][]=arm64
> supportedArchitectures[cpu][]=x64
> supportedArchitectures[libc][]=glibc
> ```
>
> **Важно: после любого изменения `.npmrc` (например, добавления блока
> `supportedArchitectures`) обычный `pnpm install` зависимостей НЕ
> переустанавливает уже существующие пакеты** — бинарники остаются старыми
> (только darwin-arm64), а от старых установок в `node_modules/.pnpm/sharp@...`
> висят битые симлинки. Это приводит к такой ошибке при сборке:
>
> ```
> ERROR  Error: ENOENT: no such file or directory, realpath
> '.../node_modules/.pnpm/sharp@0.34.5/node_modules/@img/sharp-wasm32'
> ```
>
> Лечится принудительной переустановкой, после которой в store появятся пакеты
> для всех платформ из блока:
>
> ```bash
> cd client
> pnpm install --frozen-lockfile    # первый раз / после изменений package.json
> pnpm install --force              # если менялся .npmrc (supportedArchitectures)
> pnpm build
> ```
>
> `pnpm install --force` на чистом проекте безвреден — просто пересоздаёт
> `node_modules` из lockfile.
>
> **Про сообщение `sharp binaries have been included in your build for
> darwin-arm64`.** Оно безобидно и не значит, что в сборку попала «не та»
> платформа. Это `info`-лог самого `@nuxt/image`: в хуке `nitro:compiled` модуль
> печатает `platform-arch` **машины, на которой идёт сборка** (ваш Mac), а не то,
> что реально лежит в `.output`. Наличие каталога `@img` проверяется, но его
> содержимое — нет, поэтому это сообщение вводит в заблуждение. Ошибка на сервере
> была бы только в том случае, если бы в бандл не попали *linux-x64* бинарники.
> Убеждаться нужно не по этому логу, а по содержимому сборки (проверка ниже).
>
> **Проверка, что в бандле действительно есть x64-бинарники.** Мало увидеть
> каталог `sharp-linux-x64` — убедитесь, что внутри настоящий ELF под x86-64; если
> `supportedArchitectures` не применился (бинарники не подтянулись после первого
> добавления блока), каталог может быть пустым или вообще отсутствовать — тогда
> сделайте `pnpm install --force` (обычный `pnpm install` пропускает
> переустановку) и пересоберите:
>
> ```bash
> file client/.output/server/node_modules/@img/sharp-linux-x64/lib/*.node
> # должно показать: ELF ... x86-64
> ls client/.output/server/node_modules/@img/ | grep sharp-linux-x64
> ls client/.output/server/node_modules/@img/ | grep sharp-libvips-linux-x64
> ```

Результат: `client/.output/` (SSR-сервер в `.output/server/index.mjs` и статика
в `.output/public/`).

### 7.3. Перенос результата сборки на сервер

```bash
# локально, из корня проекта
rsync -avz --delete -e 'ssh -p SSH_ПОРТ_ИЗ_ПАНЕЛИ' \
  --exclude 'node_modules' \
  --exclude 'public/content' \
  ./client/.output/ \
  root@ВЫДЕЛЕННЫЙ_IP:/opt/artserver/client/.output/

# на сервере
sudo chown -R raiyin:raiyin /opt/artserver/client/.output
ls /opt/artserver/client/.output/server/node_modules/@img/ | grep sharp-linux-x64   # должна быть строка
file /opt/artserver/client/.output/server/node_modules/@img/sharp-linux-x64/lib/*.node   # ELF ... x86-64
```

> `public/content` исключаем: живые изображения (см. 7.4) раздаёт nginx из
> `client/public/content`, а в `.output/public` лежит только их снимок на момент
> сборки. Статику `public/` (hero-images, favicon, site.webmanifest) отдельно
> переносить не нужно — она уже уехала в разделах 4 и 5.2.

### 7.4. Где лежат изображения (галерея, магазин, новости)

Живые (загружаемые через админку) изображения пишет **Go-сервер** в
`client/public/content/`, а nginx раздаёт их по `/content/...`:

```
/opt/artserver/client/public/
├── content/
│   ├── works/   # галерея  → URL /content/works/...
│   ├── sales/   # магазин  → URL /content/sales/...
│   └── news/    # новости  → URL /content/news/YYYY/MM/DD/...
├── hero-images/ # статичные баннеры главной
├── favicon*.png
└── site.webmanifest
```

Этот каталог **не пересобирается** при `pnpm build` (в `.output/public` попадает
только то, что было на момент сборки), поэтому его раздаёт nginx напрямую.
Каталог `content/` и всё, что ниже, должно быть доступно на запись процессу
сервера `artserver`:

```bash
sudo chown -R raiyin:raiyin /opt/artserver/client/public
```

> Панель Джино кэширует и сжимает статику только на своём фронте; файлы для
> `/content/` отдаёт наш nginx напрямую с диска, поэтому загруженная через
> админку картинка видна сразу (кроме кэша панели — при необходимости
> обновите страницу с очисткой кэша).

---

## 8. systemd-сервисы

### 8.1. API (Go)

`/etc/systemd/system/artserver-server.service`:

```ini
[Unit]
Description=Vera Art API server (Go)
After=network.target

[Service]
Type=simple
User=raiyin
Group=raiyin
# Важно: рабочий каталог — там, где лежат config.yaml и db/
WorkingDirectory=/opt/artserver/server
EnvironmentFile=/opt/artserver/server/.env
ExecStart=/opt/artserver/server/bin/artserver
Restart=on-failure
RestartSec=5
NoNewPrivileges=true
PrivateTmp=true

[Install]
WantedBy=multi-user.target
```

### 8.2. Web (Nuxt SSR)

`/etc/systemd/system/artserver-web.service`:

```ini
[Unit]
Description=Vera Art web (Nuxt SSR)
After=network.target

[Service]
Type=simple
User=raiyin
Group=raiyin
WorkingDirectory=/opt/artserver/client
Environment=PORT=3000
Environment=HOST=127.0.0.1
# Дублируем флаги продакшна (перекрывают значения из бандла на время SSR)
Environment=NUXT_PUBLIC_SERVER_URL=https://pertsukova.ru/api-server/
Environment=NUXT_PUBLIC_REGISTRATION_ENABLED=false
Environment=NUXT_PUBLIC_LIMIT=9
ExecStart=/usr/bin/node /opt/artserver/client/.output/server/index.mjs
Restart=on-failure
RestartSec=5
NoNewPrivileges=true

[Install]
WantedBy=multi-user.target
```

> `HOST=127.0.0.1` — Node не должен слушать внешний интерфейс: наружу торчит
> только nginx на 8080. (В предыдущей версии было `0.0.0.0`; при включённом
> `ufw` без правила на 3000 разницы бы не было, но так безопаснее.)

### 8.3. Запуск

```bash
sudo systemctl daemon-reload
sudo systemctl enable --now artserver-server artserver-web
sudo systemctl status artserver-server artserver-web
```

Проверить, что API слушает 8000, а Node — 3000:

```bash
curl -s http://127.0.0.1:8000/news | head -c 200
curl -s -o /dev/null -w "%{http_code}\n" -H 'X-Forwarded-Proto: https' http://127.0.0.1:3000/
```

---

## 9. nginx (внутренний reverse proxy, порт 8080) + HTTPS на панели Джино

### 9.1. Кто слушает 80 и 443

Прежде чем писать конфиг, убедитесь, что эти порты действительно держит панель,
а не вы — иначе `nginx` не стартует с `bind() to 0.0.0.0:80 failed`.

```bash
ss -tlnp | grep -E ':(80|443)\b' || echo '80/443 на сервере свободны'
systemctl list-units --type=service --state=running | grep -E 'apache|nginx|httpd'
```

Ожидаемый результат: порты заняты веб-сервером панели (Apache/`httpd`), либо
не заняты вовсе (панель терминирует TLS у себя, до виртуальной машины). Наш
nginx на 80/443 **не вешаем**.

### 9.2. Схема «TLS на панели» — как она выглядит снаружи

```bash
# кто отвечает по HTTPS
curl -sI https://pertsukova.ru/ | grep -iE '^(server|x-served-by):'

# чей это сертификат и до какого он действует
echo | openssl s_client -connect pertsukova.ru:443 -servername pertsukova.ru 2>/dev/null \
  | openssl x509 -noout -subject -issuer -dates
```

> Сертификат в `Issuer`/`Subject` — от панели (или коммерческого CA, купленного
> в панели). Он обновляется там же, нашему серверу не нужен. Если в выводе
> `Can't connect`, проверьте, что панель выпустила сертификат именно на
> `pertsukova.ru` (раздел 0, п. 6), и что A-записи указывают на ваш IP:
> `dig +short pertsukova.ru` должен совпадать с IP из панели.

### 9.3. Конфиг nginx

Сначала — вспомогательная переменная в контексте `http` (директива `map`
недоступна в `sites-available`, поэтому отдельным файлом):

`/etc/nginx/conf.d/artserver-proto.conf`:

```nginx
# Исходная схема запроса. Панель Джино терминирует TLS, поэтому nginx видит
# только http — берём схему из заголовка, который прислала панель, а если его
# нет, считаем, что это http.
map $http_x_forwarded_proto $arts_proto {
    default $http_x_forwarded_proto;
    ''      $scheme;
}
```

`/etc/nginx/sites-available/artserver`:

```nginx
# Внутренний reverse proxy. Наружу (80/443) работает веб-сервер панели Джино:
# он снимает TLS (сертификат панели) и присылает нам обычный HTTP на 8080.
# Поэтому здесь НЕТ listen 80/443 и НЕТ ssl_certificate.
#
# default_server: панель не всегда подставляет исходное имя хоста,
# иначе запрос ушёл бы в другой server-блок и получил 404.

server {
    listen 8080 default_server;
    # Если в контейнере настроен IPv6 и вы предпочитаете dual-stack, добавьте
    # вторую строку: listen [::]:8080 default_server;
    server_name pertsukova.ru www.pertsukova.ru;

    # Загрузка изображений в админке
    client_max_body_size 50m;

    # ── Страховка: если запрос всё-таки пришёл по http, шлём на https ──
    # Панель обычно делает этот редирект сама. Строки можно убрать, если вас
    # смущает лишний заголовок от панели.
    if ($arts_proto != "https") {
        return 301 https://$host$request_uri;
    }

    # ── API Go-сервера: срезаем префикс /api-server/ ─────────────
    location ^~ /api-server/ {
        rewrite ^/api-server/(.*)$ /$1 break;
        proxy_pass http://127.0.0.1:8000;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $arts_proto;
    }

    # ── Изображения галереи/магазина/новостей (живая статика) ────
    location ^~ /content/ {
        alias /opt/artserver/client/public/content/;
        expires 30d;
        add_header Cache-Control "public, immutable";
        access_log off;
    }

    # ── Статичные баннеры и фавиконки ────────────────────────────
    location ^~ /hero-images/ {
        alias /opt/artserver/client/public/hero-images/;
        expires 30d;
        add_header Cache-Control "public, immutable";
        access_log off;
    }
    location = /favicon.ico         { root /opt/artserver/client/public; access_log off; }
    location = /favicon-16x16.png   { root /opt/artserver/client/public; access_log off; }
    location = /favicon-32x32.png   { root /opt/artserver/client/public; access_log off; }
    location = /favicon-128x128.png { root /opt/artserver/client/public; access_log off; }
    location = /site.webmanifest    { root /opt/artserver/client/public; access_log off; }

    # ── Собранные ассеты Nuxt ─────────────────────────────────────
    location ^~ /_nuxt/ {
        proxy_pass http://127.0.0.1:3000;
        expires 1y;
        add_header Cache-Control "public, immutable";
    }

    # ── Остальное — страницы сайта (Nuxt SSR) ─────────────────────
    location / {
        proxy_pass http://127.0.0.1:3000;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $arts_proto;
    }
}
```

Включить и применить:

```bash
sudo ln -sf /etc/nginx/sites-available/artserver /etc/nginx/sites-enabled/artserver
sudo rm -f /etc/nginx/sites-enabled/default
sudo nginx -t
sudo systemctl enable --now nginx
sudo systemctl reload nginx
ss -tlnp | grep 8080   # должен слушать nginx
```

> **Если Apache панели проксирует локально** (вариант Б ниже), ограничьте порт
> localhost-ом: замените `listen 8080 default_server;` на
> `listen 127.0.0.1:8080;` — снаружи он тогда не нужен.
>
> **Если хотите смотреть логи nginx**, раскомментируйте в
> `/etc/nginx/nginx.conf` внутри `http { … }` строки `access_log`/`error_log`
> с путями в `/var/log/nginx/`. Полезно при разборе «страница отдаётся, а API
> нет».

### 9.4. Проброс из панели Джино

**Вариант А (основной): «Проксирование HTTP(S)» в панели.**

*Управление → Перенаправление портов → Проксирование HTTP(S)*:

| Параметр | Значение |
|----------|----------|
| IP | ваш выделенный IP |
| 80/tcp | `8080` |
| 443/tcp | `8080` |
| Домен (если панель спрашивает порт отдельно для домена) | `pertsukova.ru`, `8080` |

Сохраните и **включите** проксирование для этого IP. Если панель спрашивает
внутренний порт только числом, любой свободный подойдёт — он должен совпадать с
`listen` в конфиге 9.3. Некоторые панели требуют одинаковый порт для 80 и 443;
если так — оставьте `8080` для обоих (мы не различаем эти случаи: оба приходят
одним и тем же HTTP).

**Вариант Б: веб-сервер панели стоит на самом VPS (Apache).**

Если `ss -tlnp` показал, что 80/443 держит Apache Джино, настройте его
vhost-конфиг (панель: *Управление → Настройки веб-сервера*, или
`/etc/apache2/sites-available/`), и оставьте наш nginx на 8080:

```apache
<VirtualHost *:80>
    ServerName pertsukova.ru
    ServerAlias www.pertsukova.ru
    ProxyPreserveHost On
    # TLS снят панелью — сообщаем приложениям исходную схему
    RequestHeader set X-Forwarded-Proto "https"
    ProxyPass        / http://127.0.0.1:8080/
    ProxyPassReverse / http://127.0.0.1:8080/
</VirtualHost>
```

```bash
sudo a2ensite pertsukova.conf
sudo apachectl configtest && sudo systemctl reload apache2
```

Редирект `http → https` в этом варианте делает Apache или панель; наш nginx
подстрахуется проверкой `$arts_proto` из 9.3.

### 9.5. Проверка сквозного пути

```bash
# 1) наш nginx напрямую (с эмуляцией схемы https от панели)
curl -s -o /dev/null -w "nginx:8080  %{http_code}\n" \
  -H 'X-Forwarded-Proto: https' http://127.0.0.1:8080/

# 2) без заголовка — должен сработать страховочный редирект на https
curl -s -o /dev/null -w "без proto:  %{http_code}\n" http://127.0.0.1:8080/

# 3) снаружи, через панель
curl -s -o /dev/null -w "http → https: %{http_code}\n" http://pertsukova.ru/
curl -s -o /dev/null -w "https:         %{http_code}\n" https://pertsukova.ru/
curl -s -o /dev/null -w "api:           %{http_code}\n" \
  https://pertsukova.ru/api-server/news
```

Ожидается: `200`, `301`, `301` (или `200`, если панель не редиректит и наш
nginx сделал это сам), `200`, `200`.

### 9.6. Сертификат: выпуск и продление

- **Выпуск/привязка** — в панели Джино (раздел 0, п. 6) на
  `pertsukova.ru` + `www.pertsukova.ru`. Получив сертификат, панель
  автоматически начинает отдавать сайт по HTTPS.
- **Продление** — тоже в панели; отдельного cron на сервере нет и не нужно.
  Проверить срок:

  ```bash
  echo | openssl s_client -connect pertsukova.ru:443 -servername pertsukova.ru 2>/dev/null \
    | openssl x509 -noout -dates
  ```

- **Проверить заранее**, что сертификат покрывает `www`:

  ```bash
  echo | openssl s_client -connect pertsukova.ru:443 -servername www.pertsukova.ru 2>/dev/null \
    | openssl x509 -noout -subject
  ```

  Иначе `www` отдаст предупреждение браузера, хотя по HTTP уйдёт 301 на
  канонический домен — посетитель увидит ошибку до редиректа.

- **Никаких файлов `.pem`/`.key` в проекте и на сервере не нужно.** Если вдруг
  панель всё же положит их на диск (например, при переносе на другого
  провайдера) — найти можно так:

  ```bash
  sudo find / -xdev \( -name 'fullchain.pem' -o -name '*.crt' \) 2>/dev/null | head
  ```

  В текущей схеме найденные файлы игнорируйте: TLS на нашей стороне не нужен.

---

## 10. Отключение регистрации — чек-лист

Регистрация выключена **на двух уровнях**:

| Уровень | Где | Значение |
|---------|-----|----------|
| API-эндпоинт | `server/config.yaml` → `features.registration_enabled` | `false` (роут `POST /register` не монтируется) |
| Клиент | `client/.env` → `NUXT_PUBLIC_REGISTRATION_ENABLED` | `false` (скрывается кнопка входа/«Услуги» в шапке) |
| Клиент (SSR) | systemd `artserver-web` → env | `NUXT_PUBLIC_REGISTRATION_ENABLED=false` |

После включения конфига сервера перезапустите `artserver-server`. В логе при
старте должна появиться строка `Registration is DISABLED by feature flag`.

Проверка того, что `/register` недоступен:

```bash
curl -s -o /dev/null -w "%{http_code}\n" -X POST https://pertsukova.ru/api-server/register \
  -H 'Content-Type: application/json' -d '{}'
# ожидается 404 (маршрут не зарегистрирован)
```

---

## 11. Платежи (YooKassa)

В текущей версии платёжная интеграция — заглушка (см.
`server/internal/service/payment_service.go`): возвращается тестовый URL
подтверждения, ключи YooKassa не нужны. Webhook слушает путь
`POST /payments/webhook`, снаружи он будет доступен по адресу:

```
https://pertsukova.ru/api-server/payments/webhook
```

Проверьте, что вебхуки YooKassa доходят: в панели Джино убедитесь, что на
`443 → 8080` нет ограничений по методу/размеру запроса, а в nginx не стоит
`limit_req` без нужды. Входящие POST-вебхуки идут по тому же пути
`/api-server/`, отдельных правил не требуется.

Когда появится реальная интеграция, ключи размещайте в `/opt/artserver/server/.env`
и передавайте в сервис — не в `config.yaml` и не в коде.

---

## 12. Проверка работоспособности

```bash
# Страница сайта
curl -s -o /dev/null -w "site: %{http_code}\n" https://pertsukova.ru/

# API: список новостей
curl -s https://pertsukova.ru/api-server/news | head -c 300; echo

# API: публичные работы
curl -s "https://pertsukova.ru/api-server/works" | head -c 300; echo

# Логин админа (должен вернуть access_token/refresh_token)
curl -s -X POST https://pertsukova.ru/api-server/login \
  -H 'Content-Type: application/json' \
  -d '{"username":"admin","password":"ВАШ_ПАРОЛЬ"}'

# Регистрация должна вернуть 404
curl -s -o /dev/null -w "register: %{http_code}\n" -X POST \
  https://pertsukova.ru/api-server/register -H 'Content-Type: application/json' -d '{}'

# Изображение из галереи (подставьте реальный файл)
curl -s -o /dev/null -w "content: %{http_code}\n" -I \
  https://pertsukova.ru/content/works/<каталог_работы>/<файл>

# Кто терминирует TLS и когда истекает сертификат
echo | openssl s_client -connect pertsukova.ru:443 -servername pertsukova.ru 2>/dev/null \
  | openssl x509 -noout -issuer -dates

# Слушающие порты на сервере
ss -tlnp | grep -E ':(8080|3000|8000)\b'
```

Журналы:

```bash
sudo journalctl -u artserver-server -f
sudo journalctl -u artserver-web -f
sudo tail -f /var/log/nginx/error.log   # если логирование включено (9.3)
sudo tail -f /var/log/httpd/error_log   # логи веб-сервера панели (вариант Б)
```

---

## 13. Резервное копирование

У Джино есть собственные бэкапы VPS: автоматические с периодичностью 1 раз в
2–11 дней, хранятся 30 дней; в панели можно создать **неудаляемую** копию
(до трёх на сервер) и расписание. Этого достаточно, чтобы откатиться к рабочему
состоянию контейнера, но **копии лежат у провайдера** — на случай потери
доступа к аккаунту нужен ещё и off-site бэкап.

Бэкап = SQLite-файл + каталог `content/` + конфиг и секреты (`config.yaml`,
`.env` — без них не развернуть сервер заново). Скрипт `/opt/artserver/backup.sh`
(запускать от root, например из cron — поэтому он умеет делать checkpoint от
`raiyin`):

```bash
#!/usr/bin/env bash
set -euo pipefail
TS=$(date +%F_%H%M)
DEST=/var/backups/artserver
mkdir -p "$DEST"
sudo -u raiyin sqlite3 /opt/artserver/server/db/db.sqlite \
  "PRAGMA wal_checkpoint(TRUNCATE);"
tar -czf "$DEST/db_$TS.tar.gz" -C /opt/artserver/server db config.yaml .env
tar -czf "$DEST/content_$TS.tar.gz" -C /opt/artserver/client/public content
# хранить N последних копий
ls -1t "$DEST"/db_*.tar.gz | tail -n +8 | xargs -r rm -f
ls -1t "$DEST"/content_*.tar.gz | tail -n +8 | xargs -r rm -f
```

В cron (ежедневно):

```cron
0 3 * * * /opt/artserver/backup.sh >> /var/log/artserver-backup.log 2>&1
```

> **Сертификат в бэкап не входит** и не должен: его источник — панель Джино.
> При переезде к другому провайдеру сертификат выпускается заново (раздел 9.6).

---

## 14. Частые проблемы

| Симптом | Причина / решение |
|---------|-------------------|
| `bind() to 0.0.0.0:80 failed (98: Address already in use)` | Панель уже держит 80/443 — так и должно быть. Слушайте `8080` (раздел 9.1). |
| Сайт не открывается, панель показывает 502/503 | Панель проксирует не на `8080` либо nginx не запущен. Сверьте порт в панели и `ss -tlnp \| grep 8080` (9.3, 9.4). |
| Запросы не доходят до nginx | Для выделенного IP выключено «Проксирование HTTP(S)» — включите и укажите внутренний порт `8080` (раздел 0, п. 7). |
| `502 Bad Gateway` на `/api-server/*` | Go-сервис не запущен (`systemctl status artserver-server`) либо порт 8000 не слушается. |
| `error reading config file` | `config.yaml` не найден — рабочий каталог systemd-юнита не `/opt/artserver/server`. |
| `unable to open database file` | Нет каталога `db/` или нет прав на запись у `raiyin`. |
| `CGO_ENABLED` / ошибки sqlite при `go build` | Нет `gcc` — поставьте `build-essential`. |
| Предупреждение браузера о сертификате | Сертификат в панели Джино истёк/не выпущен/не покрывает `www` — проверьте `openssl s_client` (9.6). certbot на сервере не поможет. |
| `https` не работает, `http` работает | «Проксирование HTTP(S)» настроено только для 80, либо сертификат не привязан к домену в панели. |
| Вечные 302 на `/login` после логина | Проверьте, что клиент собран с `NUXT_PUBLIC_SERVER_URL=https://pertsukova.ru/api-server/` (раздел 7.1) — иначе cookie ставится не на тот домен. |
| API отвечает 404 на `/api-server/*` | Префикс не срезан — проверьте `location ^~ /api-server/` и `rewrite`. |
| Изображения 404 после загрузки в админке | Проверьте владельца `/opt/artserver/client/public/content` (нужна запись у `raiyin`) и, что nginx отдаёт `/content/`. |
| `database is locked` | Обычно решается `busy_timeout`; при частых ошибках сделайте `PRAGMA wal_checkpoint(TRUNCATE);` и перезапустите сервер. |
| Регистрация видна на сайте | `NUXT_PUBLIC_REGISTRATION_ENABLED=false` не было в переменных на момент локального `pnpm build` (раздел 7.1) — пересоберите клиент. |
| Процессы убиваются (`OOMKilled` в `journalctl`) | Мало RAM на тарифе. Добавьте swap (3.4), поднимите тариф, не собирайте клиент на сервере (7.2). |
| `[commonjs--resolver] The service is no longer running: write EPIPE` | OOM-killer на VPS с малым объёмом RAM — собирайте клиент локально (раздел 7) и переносите `.output`. |
| `ERROR ENOENT ... realpath .../@img/sharp-wasm32` при `pnpm build` | `supportedArchitectures` добавлен в `.npmrc`, но `node_modules` не переустанавливались — остались битые симлинки `@img/*`. Сделайте `pnpm install --force` и пересоберите (раздел 7.2). |
| `rsync` зависает / `Connection refused` | Не указан нестандартный SSH-порт Джино. Добавьте `-e 'ssh -p …'` или запись `Host artserver` в `~/.ssh/config` (раздел 4). |
| Письма не уходят | Порт 25/tcp у Джино заблокирован провайдером — используйте внешний SMTP по 587/465 (6.2, 11). |

---

## 15. Альтернатива: API на поддомене

Если не хотите префикс `/api-server/`, можно вынести API на поддомен
`api.pertsukova.ru`:

- В `client/.env`: `NUXT_PUBLIC_SERVER_URL=https://api.pertsukova.ru/`.
- В nginx: отдельный `server`-блок на 8080 с `server_name api.pertsukova.ru`,
  проксирующий всё на `127.0.0.1:8000` (без `rewrite`, префикс не срезается).
- Панель Джино: привязать `api.pertsukova.ru` к тому же серверу, выпустить на
  него сертификат (или включить его в выпуск на оба имени) и задать ему тот же
  внутренний порт `8080` в «Проксировании HTTP(S)».
- DNS: A-запись `api.pertsukova.ru` на тот же IP.
- CORS: Go-сервер отражает любой `Origin`, так что разрешения в `config.yaml`
  строго не требуются, но добавьте `https://pertsukova.ru` и
  `https://www.pertsukova.ru` в `cors.allowed_origins` для порядка.

Минусы: лишний домен и лишняя сервисная настройка в панели. Основной вариант с
префиксом проще.
