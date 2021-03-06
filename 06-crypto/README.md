# LetoCTF Taskbot 2022 | 06-crypto

Автор: [darkside](https://github.com/darkside0000001)

## Информация

> Флаг опять украли. Но на этот раз похитители дают нам больше возможностей.

## Описание

Симметричное шифрование AES-CBC

## Статика
- [serv.py](static/serv.py) - статика

## Решение

Рассмотрим код сервера, который нам дан. 

Чтобы получить флаг, нам нужно отправить пароль, который мы можем получить в зашифрованном виде.

Видим, что используется шифрование AES в режиме CBC

```
cipher = AES.new(key, AES.MODE_CBC, iv)
```

Также замечаем, что на сервере есть возможности проверить, корретный ли padding у сообщения, которое мы отправим. 

```
Also, you can try to input your cipher and check it padding.
```

Итак, у нас есть шифрование AES-CBC и возможность проверить padding. 

Это сразу наводит мысль о Padding Oracle Attack. 

Подробнее про атаку [тут](https://habr.com/ru/post/247527/)

## Флаг

`LetoCTF{A3C_M0d3_CBC_Y0u_kn0w}`