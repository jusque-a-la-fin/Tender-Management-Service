# Tender-Management-Service
## Как запустить:
```bash
git clone git@github.com:jusque-a-la-fin/Tender-Management-Service.git && cd Tender-Management-Service && docker compose up --build
```
## Замечания
  Выполнены все дополнительные требования:
  ### Дополнительные требования

1. Расширенный процесс согласования:

   - Если есть хотя бы одно решение reject, предложение отклоняется.
   
   - Для согласования предложения нужно получить решения больше или равно кворуму.
   
   - Кворум = min(3, количество ответственных за организацию).

3. Просмотр отзывов на прошлые предложения:

   - Ответственный за организацию может просмотреть отзывы на предложения автора, который создал предложение для его тендера.

5. Оставление отзывов на предложение:

   - Ответственный за организацию может оставить отзыв на предложение.

7. Добавить возможность отката по версии (Тендер и Предложение):

   - После отката, считается новой правкой с увеличением версии.

9. Описание конфигурации линтера.
