### Hands on profile apps with pprof

#### O que é profile de aplicações?
  1. Usar ferramentas de criação de perfil para procurar gargalos potenciais durante o desenvolvimento pode reduzir significativamente o número de problemas que aparecem mais tarde.
  2. Saber como criar o perfil de um aplicativo e entender quais são os possíveis problemas o ajudará a escrever um código melhor. Testar rotineiramente a funcionalidade que você escreveu usando um criador de perfil e procurar os gargalos e problemas comuns permitirá que você encontre e corrija muitos problemas menores que, de outra forma, se tornariam problemas maiores mais tarde.
  3. As informações de perfil servem para auxiliar na otimização do programa.
---
#### O que é o pprof?
  1. pprof é uma ferramenta para visualização e análise de dados de criação de perfil.
  1. Existem muitos comandos disponíveis na linha de comando pprof. Os comandos comumente usados incluem "top" "traces", que imprime um resumo dos principais pontos de acesso do programa, e "web", que abre um gráfico interativo de pontos de acesso e seus gráficos de chamadas.
  2. https://github.com/google/pprof/blob/master/doc/README.md
---
#### Como fazer a instrumentação
  1. via http (https://golang.org/pkg/net/http/pprof/)
  2. via runtime (https://golang.org/pkg/runtime/pprof/)
  3. O código de instrumentação que mede o aplicativo é adicionado a ele no tempo de execução, o que pode fornecer resultados muito detalhados e precisos, mas também vem com uma grande sobrecarga.
---
#### Exemplos

- Memory leak
  1.  Programa de computador gerencia incorretamente as alocações de memória, forma que a memória que não é mais necessária não seja liberada.
  2. 
    - Criação de substrings e sublices.
    - Uso incorreto da defer.
    - Corpos de resposta HTTP não fechados (ou recursos não fechados(CLose) em geral).
    - Goroutines leak.
    - Variáveis globais.
  3. O uso de memória aumenta lentamente com o tempo 
  4. Degrada o desempenho
  5. O aplicativo irá travar(freeze/crash), exigindo uma reinicialização
  6. Após reiniciar está ok novamente, e o ciclo se repete

- Goroutine leak
  1. Cada goroutine ocupa 2 KB
  2. O espaço ocupado pela própria pilha do goroutine causa vazamentos de memória.
  3. A memória de heap ocupada por variáveis em goroutine leva a vazamentos de memória de heap, que podem ser refletidos no perfil de heap.
  4. Um tipo comum de vazamento de memória é o vazamento de Goroutines. Se você iniciar um Goroutine que espera terminar, mas nunca termina, então ele vazou.
  5. Você inicia uma goroutina, mas ela nunca termina, ocupando para sempre uma memória que ela reservou. 1.
  6. Um vazamento de goroutine acontece quando as goroutines nunca terminam a execução, o que faz com que sua pilha permaneça no heap e nunca seja coletada pelo lixo, o que eventualmente leva a uma exceção de falta de memória.
  7. https://www.ardanlabs.com/blog/2018/11/goroutine-leaks-the-forgotten-sender.html
----
#### Informação da interface web
- Goroutine
  1. o perfil Goroutine relata os rastreamentos de pilha de todos os goroutines atuais.

- Allocs(Amostra de alocações de memoria)
  1. alloc - alocações totais desde o início da execução do programa, independentemente se a memória foi liberada ou não.

- Memory/Heap (Amostra de alocações de memorias ainda persistidas)
  1.  Esta é a memória que mais tarde obtém o lixo coletado pelo Go. O heap não é o único lugar onde as alocações de memória acontecem,
  1. Pode ser código que não está mais em execução
  1. Coisas que não liberaram memoria
  2. o perfil de heap relata as alocações ativas atualmente; usado para monitorar o uso de memória atual ou verificar se há vazamentos de memória.

- Profile (CPU) 
  1. o perfil da CPU determina onde um programa gasta seu tempo enquanto consome ativamente os ciclos da CPU (ao contrário enquanto dorme ou espera por E / S).

- Trace (Execução do programa)

---
- Perfil e rastreamento de CPU
  1. http://localhost:6060/debug/pprof/profile
  1. http://localhost:6060/debug/pprof/trace?seconds=5
  1. go tool trace---
- flat: representa a memória alocada por uma função e ainda mantida por essa função.
- cum: representa a memória alocada por uma função ou qualquer outra função chamada na pilha.

- flat: quanta memória é alocada por esta função
- cum: quanta memória cumulativa é alocada por esta função ou uma função que chamou para baixo da pilha

---
- inuse_space: Significa pprof está mostrando a quantidade de memória alocada e ainda não liberada.
- inuse_objects: O meio pprof está mostrando a quantidade de objetos alocados e ainda não liberados.
- alloc_space: Significa que pprof está mostrando a quantidade de memória alocada, independentemente se foi liberada ou não.
- alloc_objects: O meio pprofestá mostrando a quantidade de objetos alocados, independentemente se foram liberados ou não.

---
1. HeapAlloc - tamanho de heap atual.
2. HeapSys - tamanho total do heap.
3. HeapObjects - número total de objetos no heap.
4. HeapReleased - quantidade de memória liberada para o sistema operacional; liberações de tempo de execução para a memória do sistema operacional sem uso por 5 minutos, você pode forçar este processo com runtime / debug.FreeOSMemory.
5. Sys - quantidade total de memória alocada do sistema operacional.
6. Sys-HeapReleased - consumo efetivo de memória do programa.
7. StackSys - memória consumida por pilhas goroutine (observe que algumas pilhas são alocadas do heap e não são contabilizadas aqui, infelizmente não há como obter o tamanho total das pilhas (https://code.google.com/p/go/ questões / detalhes? id = 7468)).
8. MSpanSys / MCacheSys / BuckHashSys / GCSys / OtherSys - quantidade de memória alocada por tempo de execução para vários fins auxiliares; geralmente não são interessantes, a menos que sejam muito altos.
9. PauseNs - durações das últimas coletas de lixo

https://blog.detectify.com/2019/09/05/how-we-tracked-down-a-memory-leak-in-one-of-our-go-microservices/
https://go101.org/article/memory-leaking.html
