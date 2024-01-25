# CEP-Cloud: Aplicativo Go para Informações de CEP e Clima

## Visão Geral

CEP-Cloud é um aplicativo web baseado em Go que permite aos usuários recuperar informações de Código de Endereçamento
Postal (CEP) para endereços brasileiros e dados climáticos correspondentes. Este serviço está hospedado no Google Cloud
Run e pode ser acessado via uma simples API REST.

## Acessando o Serviço

### Recuperando Informações de CEP

Para obter informações sobre um CEP específico, você pode usar o seguinte formato de URL:

```https://cep-cloud-fp32umazia-uc.a.run.app/cep?cep=<Valor-do-CEP>```

Substitua `<Valor-do-CEP>` pelo CEP desejado. Por exemplo, para obter informações do CEP `01001000`, use:

Ex:

```https://cep-cloud-fp32umazia-uc.a.run.app/cep?cep=01001000```

O projeto inclui um `Dockerfile.dev` configurado para fins de desenvolvimento, incluindo depuração.


## Configuração Local

Para executar o aplicativo localmente, é necessário configurar algumas variáveis de ambiente e usar o Docker Compose.

### Pré-requisitos

- Docker e Docker Compose instalados.
- Uma chave de API válida para o Weather API.

### Configuração do Ambiente

Crie um arquivo `.env` na raiz do projeto com o seguinte conteúdo:

```
API_KEY=sua_chave_api_aqui
PORT=8080
```

Substitua `sua_chave_api_aqui` pela sua chave de API do Weather API.

### Usando Docker Compose

Para iniciar o aplicativo localmente usando Docker Compose, utilize o seguinte comando:

```
docker-compose up -d
```

Para testar o serviço localmente, você pode usar o comando `curl`. Por exemplo:

```
curl "http://localhost:8080/cep?cep=12345-678"
```

Pode ser que tenha que tentar mais de uma vez, pois a api da weather api é bem fraquinha 
retornando da seguinte forma 

```
curl: (52) Empty reply from server
```