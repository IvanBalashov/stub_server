# Stub Server

---

Данный сервер предназначет для моков/стабов/эмуляции другого сервера. Основная цель его использования, повторить поведение удаленного сервера.
Удобно в случае отсусвтия доступа к сети.

## Конфигурация
Ниже предоставлен пример файла `config.json`. В нем описывается поведение сервера, на какие запросы и как отвечать.
Так же есть возможность конфигурации при помощи `yaml` файла.
```json
{
  "version": "example",
  "host": "localhost",
  "port": "8080",
  "urls": {
    "api/v1/example1": [
      {
        "method": "GET",
        "answers": {
          "answer1": {
            "http_status": 200,
            "queries": {
              "example_query_1": "example1",
              "query_data": "query_data1"
            },
            "request_headers": {
              "example_request_header": "example1"
            },
            "response_headers": {
              "example_response_header": "example1"
            },
            "data": "testing_data1\n"
          },
          "answer2": {
            "http_status": 200,
            "queries": {
              "example_query_2": "example2",
              "query_data": "query_data2"
            },
            "request_headers": {
              "example_request_header": "example2"
            },
            "response_headers": {
              "example_response_header": "example2"
            },
            "data": "testing_data2\n"
          }
        }
      },
      {
        "method": "POST",
        "answers": {
          "answer1": {
            "http_status": 200,
            "queries": {
              "example_query_1": "example",
              "example_query_2": "example",
              "example_query_3": "example"
            },
            "expected_headers": {
              "header_1": "val1",
              "header_2": "val2",
              "header_3": "val3",
              "headers_data": "data"
            },
            "post_forms": {
              "arg1": "val1",
              "arg2": "val2",
              "arg3": "val3",
              "post_data": "post_data"
            },
            "data": ""
          }
        }
      }
    ],
    "api/v1/example2": [
      {
        "method": "GET",
        "answers": {
          "answer1": {
            "http_status": 200,
            "queries": {
              "example_query_1": "example",
              "example_query_2": "example",
              "example_query_3": "example"
            },
            "expected_headers": {
              "header_1": "val1",
              "header_2": "val2",
              "header_3": "val3",
              "headers_data": "data"
            },
            "data": ""
          }
        }
      }
    ]
  }
}
```

* "version" - содержит версию конфига, позже появится возможность использовать несколько конфигов и быстро переключаться между ними , тип - `string`
* "host" - адрес сервера, тип - `string`
* "port" - порт сервера, тип - `string`
* "urls" - список уролв Rest API которые необходимо эмулировать, тип - `string`

### Urls

* "method" - Rest API метод для запроса, `GET, POST, PUT, HEAD, DELETE, PATCH`, тип - `string`
* "answers" - лист с ответами для одного типа запросов
* "[query_name]" - название запроса, тип - `string`
* "http_status" - статус ответа, тип -`int`
* "queries" - лист с `url_query` для запроса, именно по ним сервер и различает что и куда отдавать
* "example_query_1": "example1" - описание аргументов запроса. 1 - название аргумента, 2 - значение.
* "query_data" - зарезервированное поле, именно это значение будет возвращаться при аргументах описанных выше
* "request_headers" - список `Headers` которые ожидает сервер, описываются подобно аргументам
* "response_headers" - список `Headers` которые отдаст сервер, описываются подобно аргументам
* "data" - данные которые будт отданы в случе если не указанны в `query_data`
* "mime_type" - mime тип ответа, дефолтное значение `text/plain`
* "wait_time" - время ожидания сервера перед ответом, указывается в виде `1s`, тип - `string` (подбронее тут - )
* "post_arguments" - аргументы ожидаемые запросом типа `POST`, описываются подобно аргументам выше
* "data_from_file" - путь до файла который необходимо вернуть 
* "cookies" - установит cookies в ответ на запрос

#### Cookies

* "name" - имя куки
* "value" - значение куки
* "max_age" - длительность жизни куки 