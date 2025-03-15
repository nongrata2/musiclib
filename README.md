# Music Library API

Этот проект представляет собой RESTful API для управления библиотекой песен. API позволяет добавлять, удалять, обновлять и получать информацию о песнях.

## Запуск проекта

Для запуска проекта выполните следующие шаги:

1. **Клонируйте репозиторий:**

   ```bash
   git clone https://github.com/nongrata2/musiclib.git
   cd musiclib
   ```
2. **Запустите проект с помощью Docker Compose**:

```bash
docker compose up --build
```
Эта команда соберёт Docker-образы для API и базы данных и запустит контейнеры с API и PostgreSQL.

После запуска API будет доступен по адресу http://localhost:8081. Для тестирования можно использовать curl.

## Доступные эндпоинты
### 1. Получить список песен
#### Метод: GET

#### URL: /songs

#### Параметры:
- **group_name** (опционально): Фильтр по названию группы.
- **song_name** (опционально): Фильтр по названию песни.
- **release_date** (опционально): Фильтр по дате выпуска.
- **text** (опционально): Фильтр по тексту песни.
- **link** (опционально): Фильтр по ссылке.
- **page** (опционально): Номер страницы.
- **limit** (опционально): Количество песен на странице.

для применения пагинации должны быть указаны и page, и limit

#### Пример:
```bash
curl -X GET "http://localhost:8081/songs"
```

#### Пример запроса с фильтрацией:

```bash
curl -X GET "http://localhost:8081/songs?group_name=Muse"
```

#### Пример запроса с пагинацией: 
```bash
curl -X GET "http://localhost:8081/songs?page=1&limit=3"
```

### 2. Добавить новую песню
#### Метод: PUT

#### URL: /songs

#### Тело запроса (JSON):

```json
{
  "group_name": "Muse",
  "song_name": "Supermassive Black Hole",
}
```
#### Пример:

```bash
curl -X PUT "http://localhost:8081/songs" \
     -H "Content-Type: application/json" \
     -d '{
           "group_name": "Muse",
           "song_name": "Supermassive Black Hole",
         }'
```

### 3. Обновить информацию о песне
#### Метод: PUT

#### URL: /songs/{songID}

#### Тело запроса (JSON):

```json
{
  "group_name": "New Group Name",
  "song_name": "New Song Name",
  "release_date": "2023-10-01",
  "text": "New lyrics",
  "link": "https://new-link.com"
}
```

####  Пример:
```bash
curl -X PUT "http://localhost:8081/songs/{songID}" \
     -H "Content-Type: application/json" \
     -d '{
           "group_name": "New Group Name",
           "song_name": "New Song Name",
           "release_date": "2023-10-01",
           "text": "New lyrics",
           "link": "https://new-link.com"
         }'
```
где songID - id песни, которую нужно отредактировать
### 4. Удалить песню
#### Метод: DELETE

#### URL: /songs/{songID}

#### Пример:

```bash
curl -X DELETE "http://localhost:8081/songs/{songID}"
```
где songID - id песни, которую нужно удалить

### 5. Получить текст песни

#### Метод: GET

#### URL: /songs/{id}

#### Пример:

```bash
curl -X GET "http://localhost:8081/songs/{songID}"
```
где songID - id песни, текст которой нужно получить
## Версии
- Go 1.23.6 
- PostgreSQL 13.20
- Docker 26.1.3
- Docker Compose 2.33.0
