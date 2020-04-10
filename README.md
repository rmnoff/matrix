### Requests

Так как настроено версионирование API, все запросы действуют из под пути /api/v1/
При обновлении API сменить v1 на v2 и так далее.

###### GET /healthcheck
- Запрос для проверки работоспособности сервера
Тело запроса:
```
-
```
Ответ:
```JSON
{"ok": true, "error": null}
```

###### POST /register
- Запрос для регистрации пользователя
Тело запроса:
```JSON
{
  "email": String,
  "password": String,
  "firstname": String,
  "lastname": String,
  "birthdate": UNIXTimestamp,
  "gender": Bool
}
```
Ответ:
```JSON
{"ok": true, "error": null, "data": "User registered"}
```

###### POST /auth
- Запрос для авторизации пользователя
Тело запроса:
```JSON
{
  "email": String,
  "password": String,
}
```
Ответ:
```JSON
{
  "ok": true,
  "error": null,
  "data": {
    "email": String,
	"firstname": String,
	"lastname": String,
	"birthdate": String
  }
}
```

###### GET /check/new/:input
- Запрос для получения списков предсказаний для пользователя по дате

Формат input: UNIXTimestamp gender personal language

Пример input: 637286400fcru

Разбор input: 637286400 - дата, f - женский пол (female), c - безличное (common), ru - русский язык

Ответ:
```JSON
{
  "ok": true,
  "error": null,
  "data": [{...}, {...}, {...}] // массив объектов слишком велик
}
```


### Utility Requests

Имеется несколько утилитарных реквестов для работы админ. панели и вывода информации

###### GET /show/types
- Запрос для получения списка существующих типов предсказаний
Ответ:
```JSON
{
  "ok": true,
  "error": null,
  "data": [{...}, {...}, {...}] // массив объектов слишком велик
}
```

###### GET /show/predictions
- Запрос для получения списка существующих предсказаний
Ответ:
```JSON
{
  "ok": true,
  "error": null,
  "data": [{...}, {...}, {...}] // массив объектов слишком велик
}
```

###### GET /show/prediction/:id
- Запрос для получения предсказания по id
Ответ:
```JSON
{
  "ok": true,
  "error": null,
  "data": {...} // объект слишком велик
}
```

###### POST /add
- Запрос для добавления нового предсказания
Тело запроса:
```JSON
{
  "ptypeid": Number, // ID типа предсказания
  "combo": String, // комбинация для предсказания (так как есть вариативность от 1 до 1-1-1 используется строковый тип)
  "prediction": String, // текст предсказания
  "personal": Bool, // вариативность предсказания между common и personal (безлич и лич)
  "language": Number // передается ID языка
}
```
Ответ:
```JSON
{
  "ok": true,
  "error": null,
  "data": Number // ID нового предсказания
}
```

###### POST /add
- Запрос для редактирования существующего предсказания
Тело запроса:
```JSON
{
  "ptypeid": Number, // ID типа предсказания
  "combo": String, // комбинация для предсказания (так как есть вариативность от 1 до 1-1-1 используется строковый тип)
  "prediction": String, // текст предсказания
  "personal": Bool, // вариативность предсказания между common и personal (безлич и лич)
  "language": Number // передается ID языка
}
```
Ответ:
```JSON
{
  "ok": true,
  "error": null,
  "data": Number // ID предсказания
}
```

### API Error Handling

Для отработки ошибок используется стандартный протокол ответа:
```JSON
{
  "ok": false,
  "data": null,
  "error": Object/String // Тут идет полный вывод содержимого ошибки
}
```

### UI Requests

Запросы для вызова страниц административной панели. Эти запросы игнорируют /api/v1 и должны вызываться напрямую

###### GET /
- Главная страница

###### GET /add
- Страница добавления новых предсказаний

###### GET /edit
- Страница списка предсказаний для редактирования

###### GET /edit/:id
- Страница редактирования предсказания по ID
