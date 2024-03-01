# Телеграм-бот для введения канала

Основная функция - это подтягивания статей с указанных источников и их постинг в телеграм-канал. Также бот способен создавать краткую сводку со статьи с помощью OpenAI API.
## Бот состоит из трех компонентов:
- Самого телеграм бота и middleware, чтобы команды были доступны только администраторам тг канала.
- Fetcher, который подтягивает статьи через RSS и добавляет их в базу данных.
- Notifier, который выкладывает еще невыложенные статьи через телеграм бота. 

## Команды, доступные боту
- /start - список команд.
- /addsource arg - добавление источника, где аргументом передаются данные об источнике в формате json.
- /delete id - удаляление источника по его id.
- /listsource - вывод списка всех источников.
- /source id - вывод информации об одном из источников по его id.
- /update arg - обновление информации об источнике, где аргументом передаются данные об источнике в формате json, включая id.

## Пример работы бота можно увидеть здесь:
- https://t.me/pet_project_news_golang_bot - бот
- https://t.me/kimcodec_tg_news - канал
