# avito-bootcamp-assignment

[![Audit](https://github.com/Ropho/avito-bootcamp-assignment/actions/workflows/audit.yml/badge.svg?branch=master)](https://github.com/Ropho/avito-bootcamp-assignment/actions/workflows/audit.yml)

## Задание 
Задание состояло в создании сервиса домов, с помощью которого пользователь сможет продать квартиру, загрузив объявление на Авито. 

## Реализация

- Была дана openAPI спецификация, по которой генерируем код с помощью 
Так мы отвязываемся от рутинных задач по написанию ручек и посвящаем время написанию бизнес-логики, а также автоматизируем взаимодействие с заказчиком, ведь SOT (Source of truth) является спецификация сервиса, от которой сервис и должен отталкиваться.
	    
        wget -d --header="Accept: application/yaml" -O oas/openapi.yaml https://app.swaggerhub.com/apiproxy/registry/IVANKOVALENKOWORK/backend-bootcamp/1.0.0?resolved=true&flatten=true&pretty=true
