# DailyTracker - Go + Docker

Aplikacja do trackowania produktywności - oceny dla pracy i życia prywatnego (0-5).

## Stack
- Backend: Go
- Database: PostgreSQL 16
- Frontend: Vanilla JS
- Deploy: Docker Compose

## Quick Start

### 1. Sklonuj repozytorium i przejdź do katalogu

```bash
cd dailytracker-go
```

### 2. Utwórz plik .env

```bash
cp .env.example .env
```

Edytuj `.env` i ustaw bezpieczne hasło:
```
POSTGRES_PASSWORD=twoje_bezpieczne_haslo
```

### 3. Uruchom z Docker Compose

```bash
docker-compose up -d
```

Aplikacja będzie dostępna na `http://localhost:8080`

### 4. Zatrzymaj aplikację

```bash
docker-compose down
```

Aby usunąć także dane z bazy:
```bash
docker-compose down -v
```

## Deployment na VPS

### 1. Zainstaluj Docker i Docker Compose na serwerze

```bash
# Ubuntu/Debian
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh
sudo usermod -aG docker $USER
```

### 2. Prześlij pliki na serwer

```bash
rsync -avz --exclude 'node_modules' --exclude '.git' \
  ./ user@your-vps-ip:/home/user/dailytracker/
```

Lub użyj git:
```bash
# Na serwerze
git clone <your-repo-url> /home/user/dailytracker
cd /home/user/dailytracker
```

### 3. Skonfiguruj środowisko

```bash
cd /home/user/dailytracker
cp .env.example .env
nano .env  # Ustaw bezpieczne hasło
```

### 4. Uruchom aplikację

```bash
docker-compose up -d
```

### 5. Skonfiguruj nginx jako reverse proxy (opcjonalne)

Jeśli chcesz wystawić aplikację na domenie z HTTPS:

```nginx
# /etc/nginx/sites-available/dailytracker
server {
    listen 80;
    server_name tracker.twoja-domena.pl;

    location / {
        proxy_pass http://localhost:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_cache_bypass $http_upgrade;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

Następnie:
```bash
sudo ln -s /etc/nginx/sites-available/dailytracker /etc/nginx/sites-enabled/
sudo certbot --nginx -d tracker.twoja-domena.pl  # Dla HTTPS
sudo systemctl reload nginx
```

## Komendy Docker

```bash
# Sprawdź logi
docker-compose logs -f

# Sprawdź logi tylko aplikacji
docker-compose logs -f app

# Sprawdź logi tylko bazy danych
docker-compose logs -f postgres

# Restart aplikacji
docker-compose restart app

# Rebuild po zmianach w kodzie
docker-compose up -d --build

# Wejdź do kontenera aplikacji
docker-compose exec app sh

# Wejdź do bazy danych
docker-compose exec postgres psql -U dailytracker -d dailytracker
```

## Backup bazy danych

```bash
# Eksport
docker-compose exec postgres pg_dump -U dailytracker dailytracker > backup.sql

# Import
cat backup.sql | docker-compose exec -T postgres psql -U dailytracker dailytracker
```

## Development lokalny (bez Dockera)

### 1. Uruchom PostgreSQL lokalnie

```bash
# macOS
brew install postgresql@16
brew services start postgresql@16

# Linux
sudo apt install postgresql-16
sudo systemctl start postgresql
```

### 2. Utwórz bazę danych

```bash
createdb dailytracker
```

### 3. Ustaw zmienne środowiskowe

```bash
export POSTGRES_URL="postgresql://localhost/dailytracker?sslmode=disable"
export PORT=8080
```

### 4. Uruchom aplikację

```bash
go run main.go
```

## API Endpoints

- `GET /api/entries` - ostatnie 30 wpisów
- `POST /api/entries` - dodaj wpis
- `GET /api/entries/:date` - wpis dla daty (YYYY-MM-DD)
- `PUT /api/entries/:date` - edytuj wpis
- `GET /api/stats/weekly` - statystyki z ostatnich 7 dni

## Frontend

- `public/index.html` - UI aplikacji
- Serwowane jako static files przez Go web server

## Zmiana portu

W `.env`:
```
PORT=3000
```

W `docker-compose.yml` zmień mapping portów:
```yaml
ports:
  - "3000:3000"
```

## Troubleshooting

### Port 8080 zajęty
Zmień port w `.env` i `docker-compose.yml`

### Baza danych nie startuje
```bash
docker-compose logs postgres
```

### Aplikacja nie łączy się z bazą
Sprawdź `POSTGRES_URL` w `.env` i upewnij się, że hasło się zgadza
