[![Golang](https://img.shields.io/badge/Go-v1.22-EEEEEE?logo=go&logoColor=white&labelColor=00ADD8)](https://go.dev/)


<div align="center">
    <h1>Сервис баннеров</h1>
    <h5>
        Cервис, который позволяет показывать пользователям баннеры, в зависимости от требуемой фичи и тега пользователя, а также управлять баннерами и связанными с ними тегами и фичами.
    </h5>
</div>

---

## Используемые технологии:
- [Golang](https://go.dev), [PostgreSQL](https://www.postgresql.org/)
- [Docker](https://www.docker.com/),
- [REST](https://ru.wikipedia.org/wiki/REST), [Swagger UI](https://swagger.io/tools/swagger-ui/)


---

## Описание системы
**Сервис позволяет создавать различные баннеры с привязкой к конкретному идентификатору фичи и тэга.
Есть возможность хранить до 3 версий баннера, пользователь может ознакомиться со всеми версиями или указать конкретную от 1 до 3.
По необходимости системой предусмотрено получение баннера актуальной версии, либо сроком хранения до 5 минут.
Данная возможность реализована с помощью `cron`. 
Сервис так же имеет возможность получить конкретный баннер, просмотреть список всех баннер, внести каки-либо изменения,
а так же позволит вам удалить баннер по его идентификатору, тэгу или фиче.
В том случае, если вы укажете тэг или фичу будут удалены все баннеры, имеющие к ним непосредственное отношение**

## Установка
```shell
git clone git@github.com:Klansander/avito-test.git
```

---

## Начало работы


1. **Запуск сервиса:**
```shell
make
```

2. **Для работоспособности сервиса на устройстве должны быть свободны порты: 8000, 5432,6379(5433,6378 для запуска тестов)**

3. **Чтобы протестировать работу сервиса, можно перейти по адресу
   http://localhost:8000/docs/index.html для получения Swagger документации.**

4. **Для токена авторизации необходимо добавить заголовок "token"**


## Дополнительные возможности

1. **Отображение списка доступных команд**
```shell
make help
```

2. **Сборка приложения с помощью Docker Compos**
```shell
make build
```

3. **Запуск приложения с помощью Docker Compose**
```shell
make up
```

4. **Обновление swagger**
```shell
make swag
```

5. **Остановка всех запущенных контейнеров**
```shell
make down
```

6. **Остановка и удаление всех запущенных контейнеров**
```shell
make clean
```

7. **Запуск тестов**
```shell
make test.integration
```

8. **Запуск линтера**
```shell
make lint
```


## Примеры запросов

### Получение баннера

#### Успешный запрос баннера 5 минутной давности
![get_request](materials/get200.png)
#### Успешный запрос актуального баннера
![get_request](materials/get200(1).png)
#### Нет токена 
В Заголовке "to
![get_request](materials/get401.png)
#### Нет прав доступа 
![get_request](materials/get403.png)
#### Некорректные данные
![get_request](materials/get400.png)
#### Баннер не найден
![get_request](materials/get404.png)

### Получение списка баннеров

#### Успешный запрос 
![list_request](materials/list200.png)
#### Нет токена 
![list_request](materials/list401.png)
#### Нет прав доступа 
![list_request](materials/list403.png)



### Создание баннера

#### Успешный запрос 
![post_request](materials/post201.png)
#### Нет токена 
![post_request](materials/post401.png)
#### Нет прав доступа 
![post_request](materials/post403.png)
#### Некорректные данные
![post_request](materials/post400.png)


### Обновление баннера

#### Успешный запрос
![patch_request](materials/patch200.png) 
#### Нет токена 
![patch_request](materials/patch401.png)
#### Нет прав доступа 
![patch_request](materials/patch403.png)
#### Некорректные данные
![patch_request](materials/patch400.png)
#### Баннер не найден
![patch_request](materials/patch404.png)

### Удаление баннера

#### Успешный запрос 
![del_request](materials/del204.png)
#### Нет токена 
![del_request](materials/del401.png)
#### Нет прав доступа 
![del_request](materials/del403.png)
#### Некорректные данные
![del_request](materials/del400.png)
#### Баннер не найден
![del_request](materials/del404.png)

### Удаление баннера по тэгу или фиче

#### Успешный запрос 
![delby_request](materials/delby204.png)
#### Нет токена 
![delby_request](materials/delby401.png)
#### Нет прав доступа 
![delby_request](materials/delby403.png)
#### Некорректные данные
![delby_request](materials/delby400.png)
#### Баннер не найден
![delby_request](materials/delby404.png)

### Получение версий баннера

#### Получение всех версий
![ver_request](materials/ver_full.png)
#### Получение конкретной версии
![ver_request](materials/ver_2.png)
#### Нет токена 
![ver_request](materials/ver401.png)
#### Нет прав доступа 
![ver_request](materials/ver401.png)
#### Некорректные данные
![ver_request](materials/ver400.png)
#### Баннер не найден
![ver_request](materials/ver404.png)


