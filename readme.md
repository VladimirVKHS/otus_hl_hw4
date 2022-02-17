# Otus Highload Architect homework 4

Масштабируемая подсистема диалогов

## Run

    docker-compose up // запуск БД

    cp .env.example .env

    ./main

## .env configuration

     DB_HOST=127.0.0.1
     DB_PORT=16033
     DB_USER=otus
     DB_PASSWORD=otus
     DB_NAME=otus
     HTTP_PORT=8086


## Описание

### Функции системы

- Создание сообщений от одного пользователя к другому

  Пример запроса:

      POST http://localhost:8086/api/messages
      
      {
          "AuthorId": 10,
          "AddresseeId": 20,
          "Message": "New message"
      }

  Пример ответа:
      
       {
           "AuthorId": 10,
           "CreatedAt": "2022-02-17 22:13:36",
           "Id": "7f4a21fb-93ca-40da-9e76-dc0d9457a332",
           "Message": "New message"
       }      

- Просмотр диалога между двумя пользователями:
   
   Пример запроса:

       GET http://localhost:8086/api/messages?user1_id=10&user2_id=20    

   Пример ответа:

       {
           "items": [
               {
                   "AuthorId": 10,
                   "CreatedAt": "2022-02-17 22:13:36",
                   "Id": "7f4a21fb-93ca-40da-9e76-dc0d9457a332",
                   "Message": "New message"
               }
           ]
       }

### Структура базы данных, ключ шардирования

   Таблица messages:
   
        CREATE TABLE messages (
            id BINARY(16) PRIMARY KEY DEFAULT (UUID_TO_BIN(UUID())),
            author_id int NOT NULL,
            chat_id VARCHAR(255) NOT NULL,
            message VARCHAR(4096),
            shard_factor VARCHAR(2),
            created_at TIMESTAMP DEFAULT (NOW())
        );
        CREATE INDEX messages_chat_id_idx ON messages (chat_id);

   В качестве первичного ключа примеяется uuid, что обеспечивает возможность переноса данных между шардами.

   Ключ шардирования рассчитывается по id пользователей между которыми происходит диалог (user1Id, user2Id):
   
       shard_factor = user1Id % 10 + user2Id % 10

   Ключ шардирования всегда находится в диапаозоне от 0 до 18.
   “Эффект Леди Гаги” (один пользователь пишет сильно больше среднего) компенсируется ввиду зависимости ключа шардирования не только от id автора, но и id адресата.
   
   Параметр chat_id служит для доступа к конкретному диалогу с применением индекса и также расчитвается по user1Id и user2Id:
   
       if user1Id < user2Id {
 		   return strconv.Itoa(user1Id) + "_" + strconv.Itoa(user2Id)
       }
       return strconv.Itoa(user2Id) + "_" + strconv.Itoa(user1Id)

   Для проксирования запросов между шардами применяется proxysql.

### Решардинг

   Процедура решардинга без даунтайма:
   - определяем целевой shard_factor (размещен на шарде А)
   - определяем шард на который хотим перенсти данный shard_factor (шард Б)
   - копируем с шарда А все сообщения с данным shard_factor на шард Б
   - изменяем правила proxysql, теперь запросы с данным shard_factor идут на шард Б
   - определяем были ли в данный промежуток времени добавлены новые сообщения на шард А с данным shard_factor, копируем их на шард Б
   - удаляем с шарда А сообщения с данным shard_factor
   