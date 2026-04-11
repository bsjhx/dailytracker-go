# DailyTracker - Go + Vercel

Aplikacja do trackowania produktywności - oceny dla pracy i życia prywatnego (0-5).

## Stack
- Backend: Go (serverless functions)
- Database: Vercel Postgres
- Frontend: Vanilla JS
- Deploy: Vercel

## Setup Vercel

### 1. Zainstaluj Vercel CLI

```bash
npm i -g vercel
```

### 2. Login

```bash
vercel login
```

### 3. Stwórz Postgres database

W dashboardzie Vercel:
1. Idź do https://vercel.com/dashboard
2. Storage → Create Database → Postgres
3. Wybierz region (najlepiej blisko Ciebie)
4. Skopiuj `POSTGRES_URL` (będzie automatycznie dodany do projektu)

### 4. Deploy

```bash
vercel
```

Przy pierwszym deploymencie:
- Link to existing project? → No
- Project name? → dailytracker (lub własna nazwa)
- Directory? → `./` (enter)
- Vercel wykryje Go i skonfiguruje automatycznie

### 5. Podłącz Postgres do projektu

W Vercel Dashboard:
- Twój projekt → Settings → Environment Variables
- Połącz database z projektem (powinno być auto)

### 6. Deploy ponownie

```bash
vercel --prod
```

Gotowe! Aplikacja działa na `https://twoj-projekt.vercel.app`

## Lokalne testowanie

Vercel nie ma prostego lokalnego dev dla Go + Postgres, więc najlepiej:
1. Deploy do Vercel
2. Testuj na `https://twoj-projekt.vercel.app`

Albo ustaw lokalnego Postgresa i zmienne środowiskowe.

## API Endpoints

- `GET /api/entries` - ostatnie 30 wpisów
- `POST /api/entries` - dodaj wpis
- `GET /api/entries/:date` - wpis dla daty (YYYY-MM-DD)
- `PUT /api/entries/:date` - edytuj wpis
- `GET /api/stats/weekly` - statystyki z ostatnich 7 dni

## Frontend

- `public/index.html` - UI aplikacji
- Automatycznie serwowane przez Vercel
