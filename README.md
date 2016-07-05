## set-token Cloud Foundry plaugin

This plugin allow to authenticate with UAA by directly sertting access and refresh tokens.

###Installation

~~~
git clone https://github.com/s-matyukevich/set-token/
cd set-token
go build -o set-token
cf install-plugin set-token
~~~

### Ussage

```
cf set-token [-a ACCESS_TOKEN] [-r REFRESH_TOKEN]
```
