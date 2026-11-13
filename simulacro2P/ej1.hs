import Control.Concurrent.STM
    ( STM
    , TVar
    , atomically
    , check
    , modifyTVar'
    , newTVarIO
    , readTVar
    , retry
    , writeTVar
    )
import Control.Monad (forever)


data Tipo = Urgente | Comun
type Llamada = Int

data Centro = Centro {
    capacidad :: Int, 
    colaUrgente :: TVar ([Llamada]),
    colaComun :: TVar ([Llamada])
}

atender :: Llamada -> IO ()
atender llamada =
    putStrLn ("Un agente atiende la llamada " ++ show llamada)
    
nuevoCentro :: Int -> IO Centro 
nuevoCentro k = do 
    urgentes <- newTVarIO ([])
    comunes <- newTVarIO ([])
    return (Centro {capacidad = k, colaUrgente = urgentes, colaComun = comunes})

encolarse :: TVar ([Llamada]) -> Int -> Llamada -> STM ()
encolarse referenciaCola capacidad llamada = do 
    cola <- readTVar referenciaCola
    check (length cola < capacidad)
    writeTVar referenciaCola (cola ++ [llamada])
    return()
    

llamar :: Centro -> Tipo -> Llamada -> IO ()
llamar centro tipo llamada = do
    case tipo of
        Urgente -> atomically $ do 
            encolarse (colaUrgente centro) (capacidad centro) llamada
        Comun -> atomically $
            encolarse (colaComun centro) (capacidad centro) llamada
    -- Esperamos hasta que nos saquen de la lista que significa que ya atendieron nuestro pedido y termino la llamada 
    case tipo of 
        Urgente -> atomically $ do 
            colaAObservar <- readTVar (colaUrgente centro)
            check (not (elem llamada colaAObservar))

        Comun -> atomically $ do 
            colaAObservar <- readTVar (colaComun centro)
            check (not (elem llamada colaAObservar))
    return ()

agente :: Centro -> IO ()
agente centro = do
    urgentesAtendidos <- newTVarIO 0
    forever $ do 
        llamadaAAtender <- atomically $ do 
            urgentes <- readTVar (colaUrgente centro)
            comunes <- readTVar (colaComun centro)
            atendidosSeguidos <- readTVar urgentesAtendidos
            llamada <- case
                ( length urgentes > 0
                    && (atendidosSeguidos < 3 || length comunes == 0)
                , length comunes > 0
                ) of
                    (False, False) -> retry
                    (True, _) -> do
                        writeTVar (colaUrgente centro) (tail urgentes)
                        modifyTVar' urgentesAtendidos (+ 1)
                        return (head urgentes)
                    (False, True) -> do
                        writeTVar (colaComun centro) (tail comunes)
                        writeTVar urgentesAtendidos 0
                        return (head comunes)
            return(llamada)
        atender llamadaAAtender
            
    
    
