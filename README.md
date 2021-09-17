# Automatic configuration NGINX

Автоматические создание файлов конфигураций из шаблонов на основании параметров полученных из yaml файла.  
Применяется для конфигурирования reverse proxy на базе Nginx для Docker контейнеров.

## Использование

Параметры запуска:  
`-input-dir` - папка с конфигурационными файлами (осуществялет поск по подпапкам);  
`-output-dir` - папка куда генерируются конфигурационные файлы исходя из полученной конфигурации;  
`-config-name` - имя конфигруационного файла для сравнения при поиске (по умолчанию `nginx.yaml`);  

На входе yaml файл  
````yaml
conf:
  host: HOST_NAME
  container: CONTAINER_NAME
  port: PORT_NUMBER
  ssl: 0/1
  sslNameCert: SSL_NAME_CERT
  sslNameKey: SSL_NAME_KEY
````
`HOST_NAME` - имя хоста  
`CONTAINER_NAME` - имя контейнера для доступа по порту
`PORT_NUMBER` - номер внешнего порта контейнера
`ssl` - генерировать шаблон с учетом ssl (1) или без него (0)
`sslNameCert` - имя сертификата
`sslNameKey` - привватный ключ  

## Примеры
> Без SSL
````yaml
conf:
  host: myhost.local
  container: myhost
  port: 5050
  ssl: 0
````
> С учетом SSL
````yaml
conf:
  host: myhost.local
  container: myhost
  port: 5050
  ssl: 1
  sslNameCert: myhost.crt
  sslNameKey: myhost.key
````

## Стрктура проекта

 ````
 - src (расположение проекта)
 -- template (папка с шаблонами для конфигуратора)
 --- nginx.ssl.template (шаблон с ssl)
 --- nginx.nonssl.template (шаблон без ssl)
 -- main.go (файл с проектом)
````