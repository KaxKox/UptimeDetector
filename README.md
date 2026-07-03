# GoCheck – Uptime Monitor API

## Przeznaczenie

GoCheck to system do monitorowania dostępności stron internetowych. Zalogowany
użytkownik dodaje listę adresów URL, które chce obserwować, a aplikacja w tle
cyklicznie sprawdza ich dostępność (kod odpowiedzi HTTP oraz czas odpowiedzi),
zapisuje wyniki w bazie danych i udostępnia je zarówno przez REST API, jak i na
żywo przez WebSocket (dashboard).

## Zaimplementowane endpointy

| Metoda | Endpoint            | Opis                                             | Autoryzacja |
|--------|----------------------|---------------------------------------------------|-------------|
| POST   | `/api/register`      | Rejestracja nowego użytkownika                    | -           |
| POST   | `/api/login`          | Logowanie, zwraca token JWT                       | -           |
| GET    | `/api/sites`          | Lista stron monitorowanych przez usera             | JWT         |
| POST   | `/api/sites`          | Dodanie nowej strony do monitoringu                | JWT         |
| PUT    | `/api/sites/{id}`     | Aktualizacja danych monitorowanej strony           | JWT         |
| DELETE | `/api/sites/{id}`     | Usunięcie strony z monitoringu                     | JWT         |
| GET    | `/api/export`         | Eksport historii sprawdzeń do pliku CSV            | JWT         |
| POST   | `/api/import`         | Import listy stron z pliku CSV                     | JWT         |
| GET    | `/ws`                  | Kanał WebSocket z danymi live (dashboard)          

## Zakres funkcjonalny

- Rejestracja i logowanie użytkowników (hasła hashowane bcryptem, sesja oparta o JWT).
- CRUD na monitorowanych stronach, ograniczony do stron należących do zalogowanego użytkownika.
- Cykliczny worker w tle (`monitor.Start`), który co 15 sekund sprawdza wszystkie strony w bazie i zapisuje wynik (kod HTTP, czas odpowiedzi) do tabeli `checks`.
- Podgląd wyników na żywo przez WebSocket (`/ws`) – ostatnie 5 sprawdzeń, odświeżane co 5 sekund.
- Eksport historii sprawdzeń do CSV oraz import listy stron z pliku CSV.
- Logowanie żądań API (metoda, status, czas trwania, IP klienta).


## Baza danych – projekt / relacje

Baza składa się z trzech tabel powiązanych relacjami 1:N:

```
users (1) ───< (N) sites (1) ───< (N) checks
```

- **users** – konta użytkowników (`id`, `username` unikalny, `password_hash`).
- **sites** – strony monitorowane przez danego użytkownika (`id`, `user_id` → FK do `users`, `name`, `url`, `interval`). Usunięcie użytkownika kasuje kaskadowo jego strony (`ON DELETE CASCADE`).
- **checks** – historia pojedynczych sprawdzeń danej strony (`id`, `site_id` → FK do `sites`, `status_code`, `response_time_ms`, `checked_at`). Usunięcie strony kasuje kaskadowo jej historię sprawdzeń (`ON DELETE CASCADE`).

Jeden użytkownik może mieć wiele monitorowanych stron, a każda strona ma wiele
zapisanych w czasie sprawdzeń – stąd relacje 1:N w obie strony łańcucha.

## Uruchomienie

1. `docker-compose up -d` – uruchamia bazę PostgreSQL i inicjalizuje schemat ze skryptu `init.sql`.
2. `go run cmd/api/main.go` – uruchamia serwer API na `localhost:8080`.
3. Frontend `npm run dev` (`App.jsx`) łączy się z API pod `http://localhost:5173`.
