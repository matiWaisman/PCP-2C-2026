import Control.Concurrent (forkIO)
import Control.Concurrent.STM
import System.IO (hSetBuffering, stdout, BufferMode(LineBuffering))
import Control.Concurrent
import Control.Concurrent.STM
import Control.Concurrent.Async
import System.IO
import Control.Monad

type Butaca = TVar Bool -- True = ocupada
type Fila = [Butaca]
type Sala = [Fila] -- de la fila más buscada a la menos buscada

tomarTodas :: [Butaca] -> STM [Butaca]
tomarTodas butacas = forM butacas (\x -> do 
        ocupada <- readTVar x 
        if ocupada then retry 
        else 
            writeTVar x true 
            return x)

bloques :: Int -> Fila -> [[Butaca]]
bloques k fila = 
    if length fila < k  then [] else take k fila : bloques k (tail fila)

reservarContiguas :: Int -> Sala -> STM [Butaca]
reservarContiguas k sala = foldr orElse retry (map (contiguasEnFila k) sala)

contiguasEnFila  :: Int -> Fila -> STM [Butaca]
contiguasEnFila k fila = foldr orElse retry (map tomarTodas (bloques k fila))


{-
reservarContiguas _ [] = retry
reservarContiguas k (f : fs) = 
    contiguasEnFila k f `orElse` reservarContiguas k fs 
-}

reservar :: Int -> Sala -> STM [Butaca]
reservar k sala = reservarContiguas k sala `orElse` tomarAlgunas k (concat sala) 

tomarAlgunas :: Int -> [Butaca] -> STM [Butaca]
tomarAlgunas 0 _ = return []
tomarAlgunas _ [] = retry
tomarAlgunas k (b:bs) = do 
    ocupada <- readTVar b 
    if ocupada then tomarAlgunas k bs 
    else do 
        writeTVar b True 
        resto <- tomarAlgunas (k-1) bs 
        return (b : resto)