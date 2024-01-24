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