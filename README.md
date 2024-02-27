# O que é isso?

Este é um serviço recebe um CEP brasileiro e retorna temperaturas (celsius, fahrenheit e kelvin) 

# Como executá-lo?

Execute com docker compose:


```bash
docker compose up
```

Execute a solicitação atraves do endpoint http://localhost:8080/?cep=69400970 utilizando o Verbo Get

Ex: 

```bash
curl --request GET --url 'http://localhost:8080/?cep=69400970'
```

# Na produção usar o enviroment 

No arquivo `.env` e adicione sua chave de API Weather (https://www.weatherapi.com/).

# Google Cloud Run

Este serviço está hospedado no serveless do Google cloud run e estará disponível online por tempo limitado. 