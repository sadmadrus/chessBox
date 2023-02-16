
# API микросервиса сессий

Данный МС обеспечивает работу с игровыми сессиями, то есть позволяет создавать партии между двумя игроками, получать данные по партиям для конкретного игрока, а также искать партии в общей базе по некоторым ключевым полям.

## Детализация

### Модели и структуры

#### Storage - структура для обработки активностей с таблицами в БД

#### Session - игровая сессия

```go
type Session struct {
 Uid         uint64
 StartDate   time.Time
 EndDate     time.Time
 BlackPlayer uint64
 WhitePlayer uint64
 Description string
 State       uint8
 Notation    string
}
```

#### Repository - 

## Реализованный API
