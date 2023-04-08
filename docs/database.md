# Сервис базы данных

Сервис отвечает за структуру, а также за добавление, обновление и получение записей из базы данных.

Общая схема базы данных
```mermaid
erDiagram
  GAMES <|-- GAMES2USERS
  GAMES2USERS <|-- USERS
  MOVES <|-- GAMES
  MOVES {
    int id PK
    int game_id FK "many-to-many"
    int move_number
    string move "e7e8Q is allowed"
    datetime time
    int branch
  }
  GAMES {
    int id PK
    string board_fen
    datetime created_at
    datetime updated_at
    bool is_active
    int winner
    bool white_turn
  }
  GAMES2USERS {
     int id PK
     int user_id FK "many-to-many"
     int game_id FK "one-to-many"
     string color "white, black allowed"
   }
  USERS {
    string username UK
    string email "example@example.example is allowed"
    string password_hash
    datetime created_at
    string profile_info
  }
```

## Требования к API

### 1. Пользователи

#### 1.1. Создание пользователя (PUT):

Метод создает запись о новой игровой сессии в таблице users. Автоматически выставляется время created_at.

##### Входящие параметры:
* username
* password_hash
* email
* profile_info

##### Возвращает:
* Header:
  * `201 Created` - запись создана в БД
  * `403 Bad Request` - данные некорректны, поэтому не записаны в БД
  * `405 Method Not Allowed` - неправильный метод
  * `500 Internal Server Error` - данные корректны, ошибка записи в БД

* JSON Body:
  * id пользователя

#### 1.2. Обновление информации о пользователе (PATCH):

Метод изменяет запись пользователя в таблице users, обновляя соответствующие поля.

##### Входящие параметры:
* username
* profile_info

##### Возвращает:
* Header:
  * `200 OK` - данные записаны в БД
  * `403 Bad Request` - данные некорректны, поэтому не записаны в БД
  * `405 Method Not Allowed` - неправильный метод
  * `500 Internal Server Error` - данные корректны, ошибка записи в БД

#### 1.3. Получить информацию о пользователе (GET):

Метод извлекает информацию о пользователе из таблицы users.

##### Входящие параметры:
* username

##### Возвращает:
* Header:
  * `200 OK` - данные извлечены из БД
  * `403 Forbidden` - данные некорректны, поэтому не извлечены из БД
  * `400 Bad Request` - данные корректны, ошибка доступа к БД
  * `405 Method Not Allowed` - неправильный метод
* JSON Body:
  * created_at, profile_info, win_count, draw_count, loss_count


### 2. Игровые сессии

#### 2.1. Создание игровой сессии (PUT):

Метод создает запись о новой игровой сессии в таблице games. Автоматически выставляется время created_at, updated_at, winner=None, white_turn=true, is_active=true.

##### Входящие параметры:
* user_id_white
* user_id_black
* board_fen

##### Возвращает:
* Header:
  * `201 Created` - запись создана в БД
  * `403 Bad Request` - данные некорректны, поэтому не записаны в БД
  * `405 Method Not Allowed` - неправильный метод
  * `500 Internal Server Error` - данные корректны, ошибка записи в БД
* JSON Body:
  * id игровой сессии

#### 2.2. Обновление (изменение очередности хода) игровой сессии (PATCH):

Метод создает изменяет запись игровой сессии по id в таблице games. Обновляются поля board_fen, updated_at, white_turn. При значении outcome != 0, обновляется поле winner и выставляется is_active = false.

##### Входящие параметры:
* id
* board_fen
* outcome (0 игра в процессе, 1 выигрыш белых, 2 выигрыш черных, 3 ничья)

##### Возвращает:
* Header:
  * `200 OK` - данные записаны в БД
  * `403 Bad Request` - данные некорректны, поэтому не записаны в БД
  * `405 Method Not Allowed` - неправильный метод
  * `500 Internal Server Error` - данные корректны, ошибка записи в БД

#### 2.3. Обновление статусов игровых сессий (PATCH):

Метод изменяет записи игровых сессии в таблице games, присваивая значение is_active = false тем, для которых поле updated_at не удовлетворяет нужному временному интервалу time_limit. Возвращается массив значений id сессий, которые были переведены в статус Неактивно.

##### Входящие параметры:
* time_limit

##### Возвращает:
* Header:
  * `200 OK` - данные записаны в БД
  * `403 Bad Request` - данные некорректны, поэтому не записаны в БД
  * `405 Method Not Allowed` - неправильный метод
  * `500 Internal Server Error` - данные корректны, ошибка записи в БД
* JSON Body:
  * []id

#### 2.4. Получить позицию на доске для игровой сессии (GET):

Метод извлекает текущую позицию на доске board_fen для записи игровой сессии с указанным id из таблице games.

##### Входящие параметры:
* id

##### Возвращает:
* Header:
  * `200 OK` - данные извлечены из БД
  * `403 Bad Request` - данные некорректны, поэтому не извлечены в БД
  * `405 Method Not Allowed` - неправильный метод
  * `500 Internal Server Error` - данные корректны, ошибка доступа в БД
* JSON Body:
  * board_fen


### 3. История ходов

#### 3.1. Запись хода в таблицу (PUT):

Метод создает запись о новом ходе игровой сессии game_id в таблицу moves.

##### Входящие параметры:
* game_id
* move_number
* move
* time
* branch

##### Возвращает:
* Header:
  * `201 Created` - запись создана БД
  * `403 Bad Request` - данные некорректны, поэтому не записаны в БД
  * `405 Method Not Allowed` - неправильный метод
  * `500 Internal Server Error` - данные корректны, ошибка записи в БД

#### 3.2. Получение полной истории ходов для сессии (GET):

Метод извлекает полную историю ходов для игровой сессии game_id из таблицы moves.

##### Входящие параметры:
* game_id

##### Возвращает:
* Header:
  * `200 OK` - данные извлечены из БД
  * `403 Bad Request` - данные некорректны, поэтому не извлечены в БД
  * `405 Method Not Allowed` - неправильный метод
  * `500 Internal Server Error` - данные корректны, ошибка доступа в БД
* JSON Body:
  * []move
