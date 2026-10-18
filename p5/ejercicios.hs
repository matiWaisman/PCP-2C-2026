import Control.Concurrent (forkIO)
import Control.Concurrent.STM
import System.IO (hSetBuffering, stdout, BufferMode(LineBuffering))
import Control.Concurrent
import Control.Concurrent.STM
import Control.Concurrent.Async
import System.IO
import Data.Maybe (isJust)


-- Ejercicio 1 
-- Versiones no bloqueantes
get :: TVar(Maybe t) -> STM (Maybe t) 
get buf = do
    contenido <- readTVar buf
    case contenido of
        Nothing -> return Nothing 
        Just elemento -> do
            writeTVar buf Nothing
            return (Just elemento)

put :: TVar(Maybe t) -> t -> STM Bool
put buf e = do 
    contenido <- readTVar buf 
    case contenido of 
        Nothing -> do 
            writeTVar buf (Just e)
            return True
        Just elemento -> return False 


ejercicio1b :: IO ()
ejercicio1b = do 
    buffer <- newTVarIO (Just 0)
    escritor <- async (atomically (put buffer 1))
    lector   <- async ( 
        do
            numero <- atomically (get buffer)    
            print numero 
        )
    wait escritor 
    wait lector 
    return ()

-- Versiones bloqueantes 
getBloqueante :: TVar (Maybe b) -> STM b
getBloqueante buf = do
    valorBuf <- readTVar buf
    case valorBuf of
        Nothing -> retry
        Just valor -> do
            writeTVar buf Nothing
            return valor

putBloqueante :: TVar(Maybe t) -> t -> STM()
putBloqueante buf e = do 
    valorBuf <- readTVar buf 
    case valorBuf of 
        Just _ -> retry 
        Nothing -> do 
            writeTVar buf (Just e)


demoBufferBloqueante :: IO ()
demoBufferBloqueante = do
    hSetBuffering stdout LineBuffering
    -- Se crea una sola vez y se comparte entre ambos threads.
    buffer <- newTVarIO Nothing
    
    productor1 <- async (atomically (putBloqueante buffer 10))
    productor2 <- async (atomically (putBloqueante buffer 20))

    consumidor1 <- async (do
        primero <- atomically (getBloqueante buffer)
        putStrLn ("Soy el consumidor 1 y lei " ++ show primero)
        )

    consumidor2 <- async (do
        segundo <- atomically (getBloqueante buffer)
        putStrLn ("Soy el consumidor 2 y lei " ++ show segundo)
        )
    -- Esperar evita que main termine antes que los threads.
    wait productor1
    wait productor2
    wait consumidor1
    wait consumidor2

read :: TVar (Maybe t) -> STM t
read buf = do
    valorBuf <- readTVar buf
    case valorBuf of
        Nothing -> retry
        Just valor -> return valor

-- Ejercicio 2 
acquireSem :: TVar(Int) -> STM()
acquireSem permisos = do 
    valor <- readTVar permisos 
    check (valor > 0)
    writeTVar permisos (valor - 1)

releaseSem :: TVar(Int) -> STM()
releaseSem permisos = do 
    valor <- readTVar permisos 
    writeTVar permisos (valor + 1)

demoSemaforo :: IO ()
demoSemaforo = do 
    hSetBuffering stdout LineBuffering
    semaforo <- newTVarIO 0
    consumidor <- async $ do
        putStrLn ("Soy el consumidor 1 y todavia no lei nada ")
        atomically $ do
            acquireSem semaforo
            acquireSem semaforo
        putStrLn ("Soy el consumidor 1 y ya tome los dos permisos ")

    consumidor2 <- async $ do
        putStrLn ("Soy el consumidor 2 y todavia no tome ningun permiso")
        atomically (acquireSem semaforo)
        putStrLn ("Soy el consumidor 2 y ya tome un permiso")

    productor1 <- async (do
        atomically(releaseSem semaforo)
        putStrLn("Soy el productor 1 y acabo de liberar un permiso")
        )
    
    productor2 <- async (do
        atomically(releaseSem semaforo)
        putStrLn("Soy el productor 2 y acabo de liberar un permiso")
        )

    productor3 <- async (do
        atomically(releaseSem semaforo)
        putStrLn("Soy el productor 3 y acabo de liberar un permiso")
        )

    wait productor1
    wait productor2
    wait productor3
    wait consumidor
    wait consumidor2
    return()
     
-- Ejercicio 3 
empezarAJugar :: TVar(Int) -> STM()
empezarAJugar capacidad = do 
    valor <- readTVar capacidad 
    writeTVar capacidad (valor + 1)

terminarDeJugar :: TVar Int -> TVar Int -> STM ()
terminarDeJugar bolsa capacidad = do
    bolitas <- readTVar bolsa
    generadores <- readTVar capacidad
    check (generadores > 0 && bolitas <= generadores - 1)
    writeTVar capacidad (generadores - 1)

consumir :: TVar Int -> STM Int
consumir bolsa = do
    bolitas <- readTVar bolsa
    check (bolitas >= 2)
    writeTVar bolsa (bolitas - 2)
    return 1

producir :: TVar Int -> TVar Int -> STM Int
producir bolsa capacidad = do
    bolitas <- readTVar bolsa
    generadores <- readTVar capacidad
    check (bolitas < generadores)
    writeTVar bolsa (bolitas + 1)
    return 1

demoBolitas :: IO ()
demoBolitas = do
    hSetBuffering stdout LineBuffering
    capacidad <- newTVarIO 0
    bolsa <- newTVarIO 0

    generador1 <- async $ do
        atomically (empezarAJugar capacidad)
        putStrLn "A entro al juego"
        threadDelay 100000
        atomically (producir bolsa capacidad)
        putStrLn "A produjo su primera bolita"
        putStrLn "A intenta producir otra bolita"
        atomically (producir bolsa capacidad)
        putStrLn "A produjo su segunda bolita: tiene 2 puntos"

    generador2 <- async $ do
        atomically (empezarAJugar capacidad)
        putStrLn "B entro al juego"
        threadDelay 100000
        atomically (producir bolsa capacidad)
        putStrLn "B produjo su primera bolita"
        putStrLn "B intenta producir otra bolita"
        atomically (producir bolsa capacidad)
        putStrLn "B produjo su segunda bolita"
        threadDelay 100000
        atomically (producir bolsa capacidad)
        putStrLn "B produjo su tercera bolita: tiene 3 puntos"

    consumidor1 <- async $ do
        threadDelay 500000
        putStrLn "X intenta consumir un par"
        atomically (consumir bolsa)
        putStrLn "X consumio un par: tiene 1 punto"

    consumidor2 <- async $ do
        threadDelay 500000
        putStrLn "Y intenta consumir un par"
        atomically (consumir bolsa)
        putStrLn "Y consumio su primer par"
        threadDelay 200000
        putStrLn "Y intenta consumir otro par"
        atomically (consumir bolsa)
        putStrLn "Y consumio su segundo par: tiene 2 puntos"

    generador3 <- async $ do
        threadDelay 200000
        atomically (empezarAJugar capacidad)
        putStrLn "C entro: ahora hay 3 generadores"
        atomically (producir bolsa capacidad)
        putStrLn "C produjo una bolita: tiene 1 punto. Intenta retirarse"
        atomically (terminarDeJugar bolsa capacidad)
        putStrLn "C pudo retirarse"

    wait consumidor1
    wait consumidor2
    wait generador1
    wait generador2
    wait generador3

    -- A y B permanecen registrados hasta que termina el consumo.
    atomically (terminarDeJugar bolsa capacidad)
    atomically (terminarDeJugar bolsa capacidad)
    putStrLn "A y B se retiraron. Fin: se produjeron 6 bolitas y se consumieron 3 pares"

-- Ejercicio 4
data Mesa = Mesa {
    hayTabaco :: Bool, 
    hayPapel :: Bool, 
    hayFosforos :: Bool
}

tabaco :: TVar Mesa -> STM ()
tabaco disp = do 
    valores <- readTVar disp 
    check (hayPapel valores && hayFosforos valores)
    writeTVar disp (Mesa False False False)

papel :: TVar Mesa -> STM ()
papel disp = do 
    valores <- readTVar disp 
    check (hayTabaco valores && hayFosforos valores)
    writeTVar disp (Mesa False False False)

fosforos :: TVar Mesa -> STM ()
fosforos disp = do 
    valores <- readTVar disp 
    check (hayTabaco valores && hayPapel valores)
    writeTVar disp (Mesa False False False)

nuevaMesa :: STM Mesa 
nuevaMesa = do 
    return (Mesa False False False)



-- Ejercicio 5
data Barra = Barra {
    cantidadChicas :: Int, 
    cantidadGrandes :: Int
}

depositarChica :: TVar Barra -> STM()
depositarChica barra = do 
    cantidadPizzas <- readTVar barra 
    writeTVar barra (Barra (cantidadChicas cantidadPizzas + 1) (cantidadGrandes cantidadPizzas))

depositarGrande :: TVar Barra -> STM()
depositarGrande barra = do 
    cantidadPizzas <- readTVar barra 
    writeTVar barra (Barra (cantidadChicas cantidadPizzas) (cantidadGrandes cantidadPizzas + 1))

tomarComida :: TVar Barra -> STM ()
tomarComida barra = do
    cantidadPizzas <- readTVar barra
    let grandes = cantidadGrandes cantidadPizzas
    let chicas = cantidadChicas cantidadPizzas
    check (grandes > 0 || chicas > 1)
    case () of
        _ | grandes > 0 ->
                writeTVar barra (Barra chicas (grandes - 1))
          | chicas > 1 ->
                writeTVar barra (Barra (chicas - 2) grandes)
          | otherwise -> retry

-- Ejercicio 6
-- Version con lista
type PilaLista = TVar [Int]

pushLista :: PilaLista -> Int -> STM ()
pushLista pilaLista e = do
    elementos <- readTVar pilaLista
    writeTVar pilaLista (e : elementos)

popLista :: PilaLista -> STM Int
popLista pilaLista = do
    elementos <- readTVar pilaLista
    case elementos of 
        [] -> retry 
        (x:xs) -> do 
            writeTVar pilaLista xs
            return x

-- Version con nodos
data Nodo = Nodo Int (TVar Nodo) | NodoVacio

type PilaNodo = TVar Nodo

pushNodo :: PilaNodo -> Int -> STM ()
pushNodo pilaNodo valor = do
    tope <- readTVar pilaNodo
    siguiente <- newTVar tope
    let nodoTope = Nodo valor siguiente
    writeTVar pilaNodo nodoTope

popNodo :: PilaNodo -> STM Int
popNodo pilaNodo = do
    tope <- readTVar pilaNodo
    case tope of
        NodoVacio -> retry
        Nodo valor siguiente -> do
            nuevoTope <- readTVar siguiente
            writeTVar pilaNodo nuevoTope
            return valor

-- Ejercicio 7
data Cola a = Cola (TVar[a])

encolar :: Cola a -> a -> STM ()
encolar (Cola referencia) e = do
    elementos <- readTVar referencia
    writeTVar referencia (elementos ++ [e])

desencolar :: Cola a -> STM a 
desencolar (Cola referencia) = do 
    elementos <- readTVar referencia 
    check (length elementos > 0)
    let (x:xs) = elementos 
    writeTVar referencia xs 
    return x

desencolarSi :: Cola a -> STM (Maybe a)
desencolarSi (Cola referencia) = do 
    elementos <- readTVar referencia 
    case elementos of 
        [] -> return Nothing 
        (x:xs) -> do
            writeTVar referencia xs 
            return (Just x)

encolarDesencolar :: Cola a -> Cola a -> IO()
encolarDesencolar aEncolar aDesencolar = do 
    atomically $ do 
        elementoAEncolar <- desencolar aDesencolar
        encolar aEncolar elementoAEncolar
    
data Prioridad = Alta | Media | Baja
data ColaPrioridad a = ColaPrioridad (Cola a) (Cola a) (Cola a)

encolarP :: ColaPrioridad a -> a -> Prioridad -> STM ()
encolarP (ColaPrioridad colaAlta colaMedia colaBaja) e prio = do 
    case prio of 
        Alta -> encolar colaAlta e 
        Media -> encolar colaMedia e 
        Baja -> encolar colaBaja e

tieneElementos :: Cola a -> STM(Bool)
tieneElementos (Cola referencia) = do 
    elementos <- readTVar referencia 
    case elementos of 
        [] -> return False 
        (_:_) -> return True 

desencolarP :: ColaPrioridad a -> STM a
desencolarP (ColaPrioridad colaAlta colaMedia colaBaja) = do 
    tieneElementosColaAlta <- tieneElementos colaAlta
    tieneElementosColaMedia <- tieneElementos colaMedia 
    tieneElementosColaBaja <- tieneElementos colaBaja
    case (tieneElementosColaAlta, tieneElementosColaMedia, tieneElementosColaBaja) of 
        (True, _, _) -> desencolar colaAlta
        (_, True, _) -> desencolar colaMedia 
        (_, _, True) -> desencolar colaBaja 
        _ -> retry 

-- Ejercicio 8
data RWLock = RWLock(TVar Int) (TVar Bool)

newRWLock :: IO RWLock
newRWLock = do
    varLectores <- newTVarIO 0
    flagEscritor <- newTVarIO False
    return $ RWLock varLectores flagEscritor

acquireLectura :: RWLock -> STM ()
acquireLectura (RWLock referenciaCantidadLectores referenciaHayEscritor) = do 
    cantidadLectores <- readTVar referenciaCantidadLectores
    hayEscritor <- readTVar referenciaHayEscritor
    check(hayEscritor == False)
    writeTVar referenciaCantidadLectores (cantidadLectores + 1)

releaseLectura :: RWLock -> STM ()
releaseLectura (RWLock referenciaCantidadLectores _) = do 
    cantidadLectores <- readTVar referenciaCantidadLectores 
    writeTVar referenciaCantidadLectores (cantidadLectores - 1)

acquireEscritura :: RWLock -> STM ()
acquireEscritura (RWLock referenciaCantidadLectores referenciaHayEscritor) = do 
    cantidadLectores <- readTVar referenciaCantidadLectores
    hayEscritor <- readTVar referenciaHayEscritor
    check(hayEscritor == False && cantidadLectores == 0)
    writeTVar referenciaHayEscritor True

releaseEscritura :: RWLock -> STM ()
releaseEscritura (RWLock _ referenciaHayEscritor) = do 
    writeTVar referenciaHayEscritor False

threadLectura :: RWLock -> IO()
threadLectura lock = do 
    putStrLn "Lector: intento adquirir el lock"
    atomically(acquireLectura lock)
    putStrLn "Lector: leyendo"
    threadDelay 500000 -- Simula una lectura que tarda medio segundo.
    atomically(releaseLectura lock)
    putStrLn "Lector: libere el lock"

threadEscritura :: RWLock -> IO ()
threadEscritura lock = do
    putStrLn "Escritor: intento adquirir el lock"
    atomically (acquireEscritura lock)
    putStrLn "Escritor: escribiendo"
    threadDelay 500000
    atomically (releaseEscritura lock)
    putStrLn "Escritor: libere el lock"

demoRWLock :: IO ()
demoRWLock = do
    hSetBuffering stdout LineBuffering
    lock <- newRWLock

    lector1 <- async (threadLectura lock)
    lector2 <- async (threadLectura lock)
    escritor <- async (threadEscritura lock)
    lector3 <- async (threadLectura lock)

    wait lector1
    wait lector2
    wait lector3
    wait escritor

-- Ejercicio 9
