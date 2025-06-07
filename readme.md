# Калькулятор

### Описание
Проект представляет собой сервис-калькулятор, способный вычислять математические выражения. 
Реализация включает два микросервиса: GRPC и HTTP API.

### Технологический стек
Язык программирования: Go 1.24.3

Использованные библиотеки:
- `google.golang.org/grpc` для реализации gRPC API
- `github.com/gin-gonic/gin` для реализации HTTP API
- `github.com/swaggo/swag` как инструмент документации API
- `github.com/stretchr/testify` для модульного тестирования

Проект реализован без ориентации на какую-то конкретную архитектуру, делал, как казалось удобнее:
- Отделил лигику получения *"сырого"* запроса(controllers) от доменной логики (services) для удобства 
переиспользования домена с разными видами API,
- Выделил часть сервисов под интерфейсы для упрощения unit-тестирования,
- Сделал единую точку конструкции объектов в main.go, получилось правильное внедрение зависимостей.

Параллельный расчет сделан с помощью сборки инструкций в дерево выражений, затем вычисления выражений через горутины. 
Кол-во горутин ограничено 10-ю через WaitGroup. 

### Как запустить

1. Переходим в корень проекта
2. `docker-compose up`

### Обращение к HTTP API

Адрес метода вычисления выражения: http://localhost:8080/http://0.0.0.0:8080/api/v1/solve

Адрес документации: http://localhost:8080/swagger/index.html

Пример запроса:
```json 
[
{ "type": "calc", "op": "+", "var": "x",        "left": 10,   "right": 2    },
{ "type": "calc", "op": "*", "var": "y",        "left": "x",  "right": 5    },
{ "type": "calc", "op": "-", "var": "q",        "left": "y",  "right": 20   },
{ "type": "calc", "op": "+", "var": "unusedA",  "left": "y",  "right": 100  },
{ "type": "calc", "op": "*", "var": "unusedB",  "left": "unusedA", "right": 2 },
{ "type": "print",             "var": "q"                        },
{ "type": "calc", "op": "-", "var": "z",        "left": "x",  "right": 15   },
{ "type": "print",             "var": "z"                        },
{ "type": "calc", "op": "+", "var": "ignoreC",  "left": "z",  "right": "y"  },
{ "type": "print",             "var": "x"                        }
]
```

Пример ответа:
```json
{
    "items": [
        {
            "var": "x",
            "value": 12
        },
        {
            "var": "z",
            "value": -3
        },
        {
            "var": "q",
            "value": 40
        }
    ]
}
```

Проверить можно через curl командой:
`curl -X POST http://localhost:8080/api/v1/solve -H "Content-Type: application/json" -d "[{\"type\":\"calc\",\"op\":\"+\",\"var\":\"x\",\"left\":10,\"right\":2},{\"type\":\"calc\",\"op\":\"*\",\"var\":\"y\",\"left\":\"x\",\"right\":5},{\"type\":\"calc\",\"op\":\"-\",\"var\":\"q\",\"left\":\"y\",\"right\":20},{\"type\":\"calc\",\"op\":\"+\",\"var\":\"unusedA\",\"left\":\"y\",\"right\":100},{\"type\":\"calc\",\"op\":\"*\",\"var\":\"unusedB\",\"left\":\"unusedA\",\"right\":2},{\"type\":\"print\",\"var\":\"q\"},{\"type\":\"calc\",\"op\":\"-\",\"var\":\"z\",\"left\":\"x\",\"right\":15},{\"type\":\"print\",\"var\":\"z\"},{\"type\":\"calc\",\"op\":\"+\",\"var\":\"ignoreC\",\"left\":\"z\",\"right\":\"y\"},{\"type\":\"print\",\"var\":\"x\"}]"`

### Обращение к gRPC API

Адрес вычисления выражения: http://localhost:8081

Пример запроса:
```json
{
  "instructions": [
    {
      "type": "calc",
      "var": "x",
      "op": "+",
      "left_num": 10,
      "right_num": 2
    },
    {
      "type": "calc",
      "var": "y",
      "op": "*",
      "left_var": "x",
      "right_num": 5
    },
    {
      "type": "calc",
      "var": "q",
      "op": "-",
      "left_var": "y",
      "right_num": 20
    },
    {
      "type": "calc",
      "var": "unusedA",
      "op": "+",
      "left_var": "y",
      "right_num": 100
    },
    {
      "type": "calc",
      "var": "unusedB",
      "op": "*",
      "left_var": "unusedA",
      "right_num": 2
    },
    {
      "type": "print",
      "var": "q"
    },
    {
      "type": "calc",
      "var": "z",
      "op": "-",
      "left_var": "x",
      "right_num": 15
    },
    {
      "type": "print",
      "var": "z"
    },
    {
      "type": "calc",
      "var": "ignoreC",
      "op": "+",
      "left_var": "z",
      "right_var": "y"
    },
    {
      "type": "print",
      "var": "x"
    }
  ]
}
```

Пример ответа:
```json
{
    "results": [
        {
            "variable": "x",
            "value": "12"
        },
        {
            "variable": "z",
            "value": "-3"
        },
        {
            "variable": "q",
            "value": "40"
        }
    ],
    "error": ""
}
```

Проверить можно через Postman или BloomRPC используя .proto файл по пути `api/proto/calculator/v1/calculator.proto`
