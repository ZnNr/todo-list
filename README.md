# todo-list
a simple api for the todo list application for creating tasks for the day


Разработать простое api для приложение todo list для создание задач на день 
1. Функции api
   
a. Создается задача с полями 
    i. Заголовок 
   ii. Описание 
  iii. Дата на которую заводится задача 
   iv. Выполнено/не выполнено 

b. Для задачи реализовать CRUDL
c. Для списков предусмотреть возможность пагинации со статусом "Не выполнено" либо "Выполнено"
d. Предусмотреть возможность выдавать задачи по дате со статусом "Не выполнено" либо "Выполнено"
e. swagger

2. Приложение запускается в docker 
3. База для приложения postgres 


Руководство по запуску приложения


1. Запуск кода локально:
- Скопировать файлы приложения локально
  к примеру в терминале выполнить команду git clone https://github.com/ZnNr/todo-list.git

- После сохранения всех изменений, вы можете запустить  Go-приложение go-todo локально с помощью команды go run.

Например, для запуска приложения main.go:
*из директории cmd
- go run main.go

2. Выполнение тестов (используйте ветку для SQLite version):
  git checkout SQLiteVersion

- Для запуска тестов в Go, используйте команду go test.

- Go будет автоматически находить и выполнять все файлы с тестами.

- Пример выполнения всех тестов в проекте: go test ./...
* перед запуском тестов приложение main.go должно быть запущено локально и доступно по адресу http://localhost:7540/


* Для пользователей из России Docker Hub заблокирован для пользователей с IP адресами из РФ
необходимо добавить иформациб о зеркалах в фаил docker/daemon.

Чтобы собрать и запустить приложение в Docker, используйте следующие команды:

1. Сборка Docker-образа: 

- Для сборки Docker-образа на основе  docker-compose проекта выполните команду docker build.

Пример: docker-compose up .
 
2. Запуск Docker-контейнера:

- После успешной сборки образа вы можете запустить Docker-контейнер с  Go-приложением todo-list.

- Пример: docker run -p 7540:7540 go-todo

После выполнения этих шагов, Go-приложение go-todo будет запущено в Docker-контейнере и будет доступно по указанному порту.

после успешного запуска контейнера приложение становится доступно по адресу http://localhost:7540/

Протестировать приложение запущенное в контейнере можно вручную.
