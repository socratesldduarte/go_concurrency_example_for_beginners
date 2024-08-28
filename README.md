## A SIMPLE EXAMPLE OF GO ROUTINES (CONCURRENCY) FOR BEGINNERS

<p>O contexto desse execmplo é o seguinte: fazer um código que tenha um loop externo (de 10 execuções por exemplo) e um loop interno (outras 10 execuções) que insira dados em um SQLite.</p>

### PRIMEIRA FORMA: SEM USO DE GO ROUTINES (SEM CONCORRÊNCIA)

<p>Na primeia implementação, eu simplesmente realizei os dois loops, sem nenhum controle de concorrência - dessa forma, os inserts são executados sequencialmente, i.e., (1, 1), (1, 2), (1, 3), ..., (1, 10), (2, 1), (2, 2), ..., ..., (9, 8), (9, 9), (9, 10)</p>

### SEGUNDA FORMA: USANDO GO ROUTINES, MAS SEM CHECK DO FINAL DA EXECUÇÃO DAS ROUTINAS PELA MAIN

<p>Nessa segunda implementação, faremos as chamadas do loop interno utilizando go routines - dessa forma, diversas execuções da função create serão chamadas concorrentemente (NÃO É NECESSARIAMENTE EM PARALELO) - e o go controla essa concorrência e distribuição entre os processadores (ou seja, provavelmente haverá paralelismo, mas controlado pelo go e não pela programação).<br>
Se não colocarmos nenhum controle que permita "esperar" a execução das go routines, provavelmente elas não serão executadas a tempo (todas ou parte delas) - fazendo com que as go routines ainda não finalizadas quando a main finalizar sejam descartadas (não executadas ou finalizadas sem completar a execução).<br>
Além disso, não podemos prever a ordem de execução dos go routines, i.e., os dados inseridos não serão exatamente na sequência "padrão" de chamada.<br>
O sleep no final da segunda forma foi adicionado para permitir aumentar / diminuir o tempo que a main vai aguardar para encerrar (o que permitirá que mais ou menos go routines sejam executadas - brinque com esses valores para ver o resultado).</p>

### TERCEIRA FORMA: USANDO GO ROUTINES, E COM VERIFICAÇÃO DO FINAL DA EXECUÇÃO DAS ROUTINAS PELA MAIN

<p>Nessa terceira forma, adicionamos um conceito importante - de wait groups. Se previamente eu sei quantas chamdas de go routines eu farei, posso utilizar essa estratégia para fazer com que a main somente se encerre APÓS o término de execução de todas as go routines.<br>
Para tanto, eu declarei um wg no escopo do pacote (estará visível para todas as funções desse pacote) e inicializei com o número de execuções que terei.<br>
Dentro da função que será chamada com go routine, inseri um wg.Done() ao final da execução da routine chamada (com um defer - que faz com que seja executado somente ao final).<br>
De maneira análoga à segunda forma, não podemos prever a ordem de execução dos go routines, i.e., os dados inseridos não serão exatamente na sequência "padrão" de chamada.</p>
