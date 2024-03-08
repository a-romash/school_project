# Школьный проект

### Система тестирования учеников

Отправная точка: http://localhost:8080/ (если вы запускете на локалке)

## Запуск проекта: 
0) Установите [docker engine](https://docs.docker.com/engine/install/) и [docker compose](https://docs.docker.com/compose/install/) (обязательно версию 2.24.5 и новее, иначе на Windows не будет билдить)
1) `cd <"path/to/project/root/directory">`
2) `docker compose up -d` (без флага `-d` если хотите, чтобы выводились логи докера в консоль)

## Вы можете получить данную ошибку: 

`Error response from daemon: Ports are not available: exposing port TCP 0.0.0.0:5673 -> 0.0.0.0:0: listen tcp 0.0.0.0:5673: bind: address already in use`

В этом случае вы должны поставить порт на любой другой свободный [docker-compose](docker-compose.yml)

Меняйте только проброшенные порты

Например:
- Было:
```yaml
  postgresql:
    container_name: postgresql
    image: postgres
    ports:
      - 5432:5432
```
- Стало:
```yaml
  postgresql:
    container_name: postgresql
    image: postgres
    ports:
      - 8081:5432
```

## Костыли

Костыли есть, признаём. Куда же без них (и да, мы знаем что проект кривой)

# Скриншоты

in the progress, скоро будут

# Схемки

![Схема всего проекта](/docs/schema.png)