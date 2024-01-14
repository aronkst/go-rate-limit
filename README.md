# Go Rate Limiter

## Visão Geral

Este projeto envolve o desenvolvimento de um limitador de taxa em Go, que pode ser configurado para limitar o número máximo de requisições por segundo com base em um endereço IP específico ou um token de acesso. O objetivo do limitador é controlar o tráfego de serviços web de forma eficiente.

## Características

- **Limitação por Endereço IP**: Restringe requisições de um único endereço IP dentro de um intervalo de tempo definido.
- **Limitação por Token de Acesso**: Limita requisições com base em tokens de acesso únicos, permitindo diferentes limites de tempo de expiração para diferentes tokens. Os tokens devem ser fornecidos no cabeçalho como API_KEY.
- **Configurações Sobrepostas**: As configurações do token de acesso têm prioridade sobre as configurações de endereço IP.
- **Integração com Middleware**: Funciona como um middleware injetado no servidor web.
- **Configuração de Limite de Requisições**: Permite definir o número máximo de requisições por tempo.
- **Duração do Bloqueio**: Opção para definir a duração do tempo de bloqueio para IP ou Token após exceder os limites de requisição.
- **Configurações de Variáveis de Ambiente**: As configurações de limite podem ser feitas por meio de variáveis de ambiente ou um arquivo .env na pasta raiz.
- **Resposta HTTP ao Exceder o Limite**: Responde com o código HTTP 429 e uma mensagem indicando a superação do número máximo de requisições.
- **Banco de Dados Redis**: Utiliza Redis para armazenar e consultar informações do limitador.
- **Estratégia de Persistência Flexível**: Padrão de estratégia para alternar facilmente entre o Redis e outros mecanismos de persistência.
- **Lógica do Limiter Separada**: A lógica do limitador é independente do middleware.
- **Configurações de Quantidade de Requisições para Lista de IPs ou Tokens**: Opção para definir um IP ou Token individual o seu número máximo de requisições por tempo.
- **Configurações de Duração do Bloqueio para Lista de IPs ou Tokens**: Opção para definir um IP ou Token individual o seu tempo de bloqueio.

## Exemplos de Uso

- **Exemplo de Limitação por IP**: Se configurado para um máximo de 5 requisições por segundo por IP, a 6ª requisição do IP 192.168.1.1 dentro de um segundo deve ser bloqueada.
- **Exemplo de Limitação por Token**: Se um token `abc123` estiver configurado com um limite de 10 requisições por segundo, a 11ª requisição dentro desse segundo deve ser bloqueada.
- **Tempo de Expiração**: Após atingir o limite, novas requisições do mesmo IP ou token só são possíveis após o tempo de expiração.
- **Configuração Personalizada por IP**: Suponha que o `CUSTOM_MAX_REQ_PER_SEC` esteja configurado com "192.168.1.2=2", isso significa que o IP 192.168.1.2 terá um limite personalizado de 2 requisições por tempo, independente do limite padrão para outros IPs.
- **Configuração Personalizada por Token**: Se `CUSTOM_MAX_REQ_PER_SEC` incluir "token123=8", isso indica que o token `token123` tem um limite personalizado de 8 requisições por tempo, que pode ser diferente do limite padrão para outros tokens.
- **Duração de Bloqueio Personalizada**: Usando `CUSTOM_BLOCK_DURATION` com "192.168.1.2=30s", o IP 192.168.1.2 será bloqueado por 30 segundos após exceder seu limite de requisições, que é uma configuração específica diferente do bloqueio padrão.
- **Bloqueio Personalizado para Token**: Por exemplo, com "token123=1m" em `CUSTOM_BLOCK_DURATION`, o token `token123` enfrentará um bloqueio de 1 minuto após atingir seu limite de requisições.

## Variáveis de Ambiente

A configuração do rate limiter é gerenciada através das seguintes variáveis de ambiente. Cada uma desempenha um papel crucial no controle e personalização do comportamento do limitador de taxa:

- **REDIS_ADDRESS**: Define o endereço do servidor Redis utilizado pelo limitador de taxa.
- **REDIS_PASSWORD**: Senha para autenticação no servidor Redis.
- **REDIS_DB**: Número do banco de dados Redis a ser utilizado pelo aplicativo.
- **DEFAULT_IP_MAX_REQ_PER_SEC**: Define o limite padrão de requisições por segundo por endereço IP. Este valor é aplicado a todos os IPs, a menos que uma configuração personalizada seja especificada.
- **DEFAULT_TOKEN_MAX_REQ_PER_SEC**: Estabelece o limite padrão de requisições por segundo por token de acesso. Esse limite é aplicado a todos os tokens, exceto aqueles com configurações personalizadas.
- **DEFAULT_IP_BLOCK_DURATION**: Duração do bloqueio padrão para um endereço IP que excede seu limite de requisições. Especificado em um formato de duração, como 10s para dez segundos.
- **DEFAULT_TOKEN_BLOCK_DURATION**: Duração do bloqueio padrão para um token que excede seu limite de requisições, especificado no mesmo formato de duração que o bloqueio de IP.
- **CUSTOM_MAX_REQ_PER_SEC**: Permite a definição de limites de requisição personalizados para IPs ou tokens específicos. Formato esperado: `ip_ou_token=valor;outro_ip_ou_token=valor`, por exemplo, `127.0.0.1=2;abc123=10`.
- **CUSTOM_BLOCK_DURATION**: Configura durações de bloqueio personalizadas para IPs ou tokens específicos. O formato é similar ao CUSTOM_MAX_REQ_PER_SEC, por exemplo, `127.0.0.1=30s;abc123=1m`.

Estas variáveis são fundamentais para a flexibilidade e eficácia do limitador, permitindo uma adaptação às necessidades específicas. É importante definir essas variáveis de forma apropriada para garantir que o sistema funcione como esperado.

## Funcionamento do Rate Limit

O Rate Limiter é implementado como um middleware no servidor HTTP, permitindo que ele intercepte e controle as requisições.

### Middleware de Rate Limiting

- O middleware `RateLimiterMiddleware` é aplicado a cada requisição recebida pelo servidor.
- Ele identifica cada requisição por um identificador único, que pode ser um token de acesso (se presente no cabeçalho API_KEY) ou o endereço IP do solicitante.
- Após identificar a requisição, o middleware consulta o `RateLimit`, uma estrutura que contém as configurações de limitação e o Redis para verificar se o limite foi excedido.

### Verificação e Controle de Limite

- A função `IsLimitExceeded` da estrutura `RateLimit` é responsável por determinar se uma requisição excede o limite configurado.
- Ela verifica o número atual de requisições feitas pelo identificador (IP ou token) no Redis.
- Com base no identificador, a função determina o limite máximo de requisições permitidas por segundo `maxReqPerSec` e a duração do bloqueio `blockDuration`. Esses valores podem ser configurados de forma personalizada para cada IP ou token ou utilizar os valores default.
- Se o número de requisições já feitas for maior ou igual ao limite permitido, a função indica que o limite foi excedido.

### Resposta a Requisições Excessivas

- Se o limite for excedido, o middleware responde imediatamente com um erro HTTP 429 "you have reached the maximum number of requests or actions allowed within a certain time frame", indicando ao cliente que o limite de requisições foi atingido.
- Em caso de erros internos (por exemplo, falha ao acessar o armazenamento), o middleware responde com um erro HTTP 500.

### Incremento de Contagem de Requisição

- Para cada requisição válida que não excede o limite, o contador de requisições no armazenamento é incrementado. Este contador é usado para rastrear o número de requisições feitas pelo identificador dentro do intervalo de tempo especificado.

## Pré-requisitos

Certifique-se de ter o Docker instalado no seu sistema. Você pode baixar e instalar o Docker a partir do [site oficial do Docker](https://www.docker.com/).

## Configuração

- Clone o repositório.
- Configure o arquivo `.env` ou defina as variáveis de ambiente conforme necessário.

## Executando o Projeto

```bash
docker compose up
```
