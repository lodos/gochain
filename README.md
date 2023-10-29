# GoChain

## Цели проекта
1. Создание собственного блокчейн-проекта для сохранения данных и их производных из различных приложений LineClub R&D.
2. Проект в эксплуатации должен иметь закрытый код и работать максимально быстро и обладать максимальной
   масштабируемостью
3. Все инструменты среды разработки бесплатны, основаны на открытом коде и носит академический характер.
4. Безопасность проекта - важный фактор!

## Описание проекта

Проект реализован на Go и SQLLite3. База данных хранится в файле blocks.db. Логин и пароль могут быть установлены
опционально

## Как использовать
0. Установить Go (https://go.dev/doc/install)
1. Клонировать репозиторий
2. Запустить проект

Для тестирования, компилляции кода и запуска
```
go run .
```

Для запуска скомпиллированного приложения
```shell
./build/go_build_blockchain
```

Для сборки приложения с помощью docker compose

```shell
docker-compose up
```

GoChain будет запущен по адресу: http://127.0.0.1:3000 с помощью Fiber

Далее, необходимо передать данные в БД, чтобы сохранить их в блокчейн.

Обязательными параметрами являются:
1. server_key = EDLJDhHD9mzbloENzA9pjdteAhgE4VdClUWR2SPP7tfkkQFHb9APBW4STOTTlM7S
2. data = JSON-объект со следующей структурой:
```json
   {
   "method": "add",
   "filename": "filename.mp4",
   "encrypted": 1,
   "key": "123123123123",
   "username": "anyuser",
   "url": "https://any.domain.ru/tmp/filename.mp4"
}
```
В этом ключе может быть любой объем данных. Именно из него высчитывается hash.

Пример запроса с помощью POST
```
POST http://127.0.0.1:3000/blocks/?data=%7B%22method%22%3A%22add%22,%22filename%22%3A%22123.mp4%22,%22encrypted%22%3A1,%22key%22%3A%22123123123123%22,%22username%22%3A%22anyuser%22,%22url%22%3A%22https%3A%2F%2Fstats.lineclub.ru%2Ftmp%2F123.mp4%22%7D&server_key=EDLJDhHD9mzbloENzA9pjdteAhgE4VdClUWR2SPP7tfkkQFHb9APBW4STOTTlM7S
```

Получить данные обо всех записях в блокчейн можно с помощью команды GET
```
http://127.0.0.1:3000/blocks
```

Получить данные о блоке по его хешу:
```
http://127.0.0.1:3000/blocks/023918e3bb6467425b15c876c9aac4d987312d06718adb63f63c61b01c05e668
```


#### Данный код разработан в рамках проекта LineClub R&D и является открытым, свободнораспространяемым ПО. При использовании кода без изменений, копирайт в коде с указанием проекта LineClub.RU обязателен