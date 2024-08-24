<p align="center">
  <img src="/logo-rate-limiter.png" alt="Rate Limiter Logo" width="400" height="400">
</p>

<h1 style="text-align: center">Rate Limiter</h1>

### Visão Geral

O projeto "Rate Limiter" é uma aplicação desenvolvida em Golang que implementa um limitador de taxa (rate limiter) para controlar o número de requisições HTTP que podem ser feitas por clientes em um determinado período de tempo. A aplicação utiliza tokens JWT para identificar e autenticar as requisições, permitindo a personalização do limite de requisições por cliente.

### Requisitos

Antes de iniciar o projeto "Rate Limiter", certifique-se de que o seu ambiente de desenvolvimento atende aos seguintes requisitos:

- **Go**: Versão 1.22.1 ou superior instalada.
- **Docker**: Docker instalado para gerenciar contêineres.
- **Docker Compose**: Para orquestrar os contêineres necessários.

### Subindo o Projeto com Docker

Para subir o projeto, siga os passos abaixo:

1. **Clone o Repositório**:
   Se ainda não fez isso, clone o repositório do projeto para o seu ambiente local:

```shell
git clone https://github.com/seu-usuario/rate-limiter.git
cd rate-limiter
```
2. **Configuração do Arquivo .env**:
   Antes de iniciar o Docker, certifique-se de ter um arquivo `.env` na raiz `REDIS_PORT`, `REDIS_CACHE_KEY`, entre outras.
3. **Iniciar os Contêineres com Docker Compose:**:
   Utilize o Docker Compose para iniciar o ambiente. Isso irá criar e iniciar os contêineres necessários, incluindo o Redis:
```shell
docker-compose up -d
```
4. **Verificar os Contêineres:**:
Para verificar se os contêineres estão rodando corretamente, use:
```shell
docker ps
```
5. **Acessar a Aplicação:**:
   Com os contêineres em execução, a aplicação estará disponível no endereço configurado no arquivo `.env`.
### Parando os Contêineres
Para parar os contêineres, utilize o comando `docker-compose down`:
```shell
docker-compose down
```
Isso irá parar e remover todos os contêineres, redes e volumes criados pelo Docker Compose.

### Versão do Go
O projeto utiliza a versão 1.22.1 do Go. Para garantir a compatibilidade e funcionamento correto, certifique-se de que esta versão está instalada no seu ambiente de desenvolvimento:
```shell
go version
# go version go1.22.1

```
### Arquitetura

A aplicação está estruturada em vários pacotes, cada um com uma responsabilidade clara:

```go
// configpkg: Lida com o carregamento e gerenciamento das configurações da aplicação, utilizando variáveis de ambiente.
package confpkg

// cache: Fornece a interface para interação com o sistema de cache (Redis).
package cache

// repository: Implementa o repositório de requisições, responsável por verificar e registrar o número de requisições feitas por um cliente.
package repository

// webserver: Inicia e gerencia o servidor HTTP.
package webserver

// handlers: Define os handlers HTTP, onde são aplicadas as regras de rate limit e gerados os tokens JWT.
package handlers
```
### Configuração
As configurações da aplicação são carregadas pelo pacote `configpkg`, que utiliza o pacote `envsnatch` para extrair variáveis de ambiente de um arquivo`.env. Essas configurações incluem parâmetros como ambiente da aplicação, host e porta do Redis, chave para cache Redis, limite padrão de requisições por segundo, entre outros.

```go
// Estrutura de configuração usada na aplicação.
type Conf struct {
	AppEnv              string `env:"APP_ENV"`
	WSHost              string `env:"WS_HOST"`
	JWTKey              string `env:"JWT_KEY"`
	RedisHost           string `env:"REDIS_HOST"`
	RedisPort           string `env:"REDIS_PORT"`
	RedisCacheKey       string `env:"REDIS_CACHE_KEY"`
	DefaultMaxReqPerSec int    `env:"DEFAULT_MAX_REQ_PER_SEC"`
	TokenExpiresInSec   int    `env:"TOKEN_EXPIRES_IN_SEC"`
	TimeoutDuration     int    `env:"TIMEOUT_DURATION"`
}
```
### Inicialização do Servidor
No ponto de entrada da aplicação (`main.go), o servidor é configurado e iniciado após a inicialização do cache e do repositório de requisições. O servidor é responsável por escutar as requisições HTTP e aplicar as regras de rate limit.
```go
// Função main que inicializa as configurações, cache e repositório, e inicia o servidor web.
func main() {
	conf, _, err := confpkg.LoadConfig()
	if err != nil {
		log.Fatalln(err)
	}

	cacheProvider, err := cache.NewClientProvider(&cache.ClientSettings{
		Host:     conf.RedisHost,
		Port:     conf.RedisPort,
		Password: conf.RedisCacheKey,
		AppEnv:   conf.AppEnv,
	})

	if err != nil {
		log.Fatalln(err)
	}

	cacheClient := cache.NewClient(cacheProvider)
	requestRepository := repository.NewRequestRepository(cacheClient)

	webserver.Start(requestRepository)
}
```
### Repositório de Requisições
O repositório de requisições (`RequestRepository`) interage com o cache (`Redis`) para verificar e registrar o número de requisições feitas por um cliente. O método CheckRateLimit é o responsável por verificar se o cliente ultrapassou o limite de requisições permitidas por segundo.
```go
// Verifica se a requisição é permitida de acordo com o limite de taxa.
func (r *RequestRepository) CheckRateLimit(key string, limit int) (bool, error) {
	ctx := context.Background()
	current, err := r.CacheClient.Get(ctx, key).Int()
	if err != nil && !errors.Is(err, redis.Nil) {
		return false, err
	}

	if current < limit {
		delay := time.Duration(confpkg.Config.TimeoutDuration) * time.Second
		_, err := r.CacheClient.Set(ctx, key, current+1, delay).Result()
		if err != nil {
			return false, err
		}
		return true, nil
	}

	return false, nil
}
```
### Geração de Tokens JWT
Os tokens JWT são gerados e validados para autenticar as requisições e aplicar as regras de rate limit. O token inclui o IP do cliente e o limite máximo de requisições por segundo.
```go
// Gera um novo token JWT contendo o IP e o limite máximo de requisições por segundo.
func NewJWT(ip string, expirationDuration time.Duration, maxReqPerSec int) (string, error) {
	expirationTime := time.Now().Add(expirationDuration)
	claims := &Claims{
		IP:           ip,
		MaxReqPerSec: maxReqPerSec,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(JwtKey())

	return tokenString, nil
}
```
### Middlewares
Os middlewares aplicam as regras de segurança e limite de taxa. Um middleware extrai o token JWT do cabeçalho da requisição e define as claims no contexto da requisição, enquanto outro middleware impõe o limite de requisições por IP.
```go
// Middleware que extrai o token JWT e define as claims no contexto da requisição.
func (m *MiddlewarePkg) SetJWTClaimsMiddleware(next http.Handler) http.Handler {}

// Middleware que impõe o limite de requisições por IP.
func (m *MiddlewarePkg) RateLimitMiddleware(next http.Handler) http.Handler {}
```
### Execução do Servidor Web
O servidor web é iniciado com as configurações carregadas, e fica escutando requisições HTTP, aplicando as regras de rate limit definidas.
```go
// Inicia o servidor HTTP e lida com sinais do sistema para desligamento gracioso.
func Start(reqRepository entity.RequestRepositoryInterface) {
	server := &http.Server{Addr: confpkg.Config.WSHost, Handler: handlers.Handler(reqRepository)}
}
```

### Testes Unitários

O projeto "Rate Limiter" inclui uma série de testes unitários para garantir o correto funcionamento das funcionalidades de middleware e geração de tokens JWT. Abaixo está uma descrição dos principais testes realizados:

#### Testes do Middleware

- **TestNewRateLimiterMiddleware**: Verifica se o middleware de rate limiter é inicializado corretamente com o repositório de requisições.

- **TestSetJWTClaimsMiddleware**:
    - **Valid Token**: Testa se um token JWT válido é corretamente processado, extraindo as claims e definindo-as no contexto da requisição.
    - **No Token**: Verifica o comportamento do middleware quando nenhum token é fornecido, utilizando os valores padrão de configuração.
    - **Invalid Token**: Avalia a resposta do middleware ao receber um token inválido, esperando que retorne `HTTP 401 Unauthorized`.
    - **Expired Token**: Testa o comportamento do middleware com um token expirado, esperando também um retorno `HTTP 401 Unauthorized`.

- **RateLimitMiddleware**:
    - **Rate limit middleware error**: Testa o comportamento do middleware quando ocorre um erro na verificação do limite de requisições, esperando um retorno `HTTP 500 Internal Server Error`.
    - **Rate limit middleware allows request**: Verifica se o middleware permite a requisição quando o limite de taxa não foi atingido, esperando um retorno `HTTP 200 OK`.
    - **Rate limit middleware blocks request**: Avalia se o middleware bloqueia a requisição quando o limite de taxa foi atingido, esperando um retorno `HTTP 429 Too Many Requests`.

#### Testes de Geração de Tokens JWT

- **TestNewJWTGenerates**:
    - **Geração de token com valores padrão**: Testa a criação de um token JWT, verificando se as claims, como IP e `MaxReqPerSec`, são definidas corretamente.
    - **Geração de token com `MaxReqPerSec` igual a zero**: Avalia a criação de um token JWT quando o limite de requisições por segundo é definido como zero, garantindo que o comportamento esperado seja mantido.

Esses testes são fundamentais para garantir que as regras de segurança e controle de requisições implementadas no projeto "Rate Limiter" funcionem conforme o esperado em diferentes cenários.

### Conclusão
Este projeto demonstra como implementar um limitador de taxa em Golang utilizando Redis para armazenar os dados temporários de requisições e JWT para autenticação dos clientes. As principais funcionalidades incluem a criação de tokens personalizados para cada cliente e a aplicação de limites de requisição por segundo, garantindo a estabilidade e segurança da API.