## Building & running

```shell
docker build --tag instagram:v2.0.6 --force-rm  .
docker run --name instagram --rm --log-driver json-file --log-opt max-size=50m --log-opt max-file=5 --env-file service.env -p 80:80 -d instagram:v2.0.6
```

## Slots

По умолчанию используется файл со слотами `instagram.slots`, который указывается в качестве источника в `service.env`.

Файл содержит список хостов, разделенный `\n`.

Так же в `service.env` в качестве источника можно указать URL-запрос на удаленный http/https ресурс, который вернет список слотов, в формате, описанным выше.

Примеры:

`instagram.slots`
```text
instagram_api_8125:8125
instagram_api_8126:8126
```

`service.env`
```text
SLOTS_URI=https://instagram_api.ontec.ru/slots
SLOTS_URI=instagram.slots
```