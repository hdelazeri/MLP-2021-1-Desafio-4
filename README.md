# Desafio 4 - MLP

Escolha duas linguagens, dentre as seguintes:

- Rust
- Julia
- Go
- Java 8
- Scala

Implemente, nas duas linguagens escolhidas, um programa que contém produtores e consumidores de dados. Depois compare o quão fácil ou difícil é em cada uma delas. A quantidade de produtores e consumidores deve ser parametrizável, e, uma vez lançado o programa, ela não muda. Os dados a serem processados são cadeias de caracteres (uma por linha) em um único arquivo de entrada que contém cerca de 14 milhões de cadeias de caracteres (veja abaixo como obter o arquivo de entrada). O arquivo de entrada tem o nome de “rockyou.txt” e possui 14442062 palavras em 134 megabyte de dados. A função dos produtores é realizar a leitura do arquivo de entrada, colocando os dados em uma fila. Os consumidores, por outro lado, removem os dados que estão na fila e fazem o seu processamento. O processamento consiste em encontrar uma determinada cadeia de caracteres (passada como parâmetro para os consumidores) e contar quantas ocorrências dessa entrada aparecem na fila.

O arquivo rockyou.txt está disponível no seguinte link: [rockyou.txt](https://github.com/brannondorsey/naive-hashcat/releases/download/data/rockyou.txt)

# Rust

Para executar o programa é necessário executar os seguintes comandos:

```bash
cargo run -- <ARQUIVO> <TEXTO>
```

Para alterar a quantidade de produtores e consumidores os respectivos parâmetros podem ser utilizados:

```bash
cargo run -- -p <PRODUCERS> -c <CONSUMERS> <ARQUIVO> <TEXTO>
```

# Go

Para executar o programa é necessário executar os seguintes comandos:

```bash
go run main.go <ARQUIVO> <TEXTO>
```

Para alterar a quantidade de produtores e consumidores os respectivos parâmetros podem ser utilizados:

```bash
go run main.go -P <PRODUCERS> -C <CONSUMERS> <ARQUIVO> <TEXTO>
```