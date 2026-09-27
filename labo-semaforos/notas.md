Para hacer un thread una forma es con herencia 
Hay que overridear el run 

Otra forma es con un runnable, vamos a usar esa. 

Un join de un thread padre es como un wait a un hijo. 

Yield hace que un thread se suspenda solo/ delegue su quantum de ejecución. 

No podemos garantizar que la escritura y la lectura se hagan en el orden correcto (en especial la escritura)

Variables volátiles garantiza visibilidad, cuando escribo en esa variable se va a ver en todos lados, das una garantía de orden. 

En Java los mutex son ReentrantLock