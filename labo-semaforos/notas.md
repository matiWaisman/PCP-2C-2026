Para hacer un thread una forma es con herencia 
Hay que overridear el run 

Otra forma es con un runnable, vamos a usar esa. 

Un join de un thread padre es como un wait a un hijo. 

Yield hace que un thread se suspenda solo/ delege su quantum de ejecucion. 

No podemos garantizar que la escritura y la lectura se hagan en el orden correcto (en especial la escritura)

Variables volatiles garantiza visibilidad, cuando escribo en esa variable se va a ver en todos lados, das una garantia de orden. 

En Java los mutex son ReentrantLock