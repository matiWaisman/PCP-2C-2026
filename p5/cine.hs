-- INCOMPLETO, Falta completar R3 

import Control.Concurrent.STM
    ( STM
    , TVar
    , atomically
    , check
    , modifyTVar'
    , newTVarIO
    , readTVar
    , retry
    , writeTVar, newTVar
    )
import Control.Monad (forever)
import Data.List (isInfixOf)

type Butaca = TVar Bool -- True = ocupada
type Fila = [Butaca]
type Sala = [Fila] -- de la fila más buscada a la menos buscada

inicializarSala :: Int -> Int -> IO Sala 
inicializarSala 1 columnas = do 
    fila <- inicializarFila columnas
    return [fila]
inicializarSala filas columnas = do 
    nuevaFila <- inicializarFila columnas 
    resto <- inicializarSala (filas - 1) columnas 
    return (nuevaFila : resto)


inicializarFila :: Int -> IO Fila 
inicializarFila 1 = do 
    butaca <- newTVarIO False 
    return [butaca]
inicializarFila columnas = do 
    nuevaButaca <- newTVarIO False 
    resto <- inicializarFila (columnas - 1)
    return (nuevaButaca : resto)

hayKContiguasEnFila :: Int -> Fila -> STM Bool
hayKContiguasEnFila k fila = do
    estados <- mapM readTVar fila
    return (replicate k False `isInfixOf` estados)

hayKContiguas :: Int -> Sala -> STM Bool 
hayKContiguas _ [] = return False
hayKContiguas k (fila:filas) = do 
    hayKContiguasEnEstaFila <- hayKContiguasEnFila k fila 
    hayKContiguasEnElResto <- hayKContiguas k filas 
    return (hayKContiguasEnEstaFila || hayKContiguasEnElResto)

cantidadButacasLibres :: Fila -> STM Int 
cantidadButacasLibres [] = return 0 
cantidadButacasLibres (butaca:butacas) = do 
    valor <- readTVar butaca 
    cantidadButacasRestantes <- cantidadButacasLibres butacas 
    if valor then 
        return cantidadButacasRestantes 
    else 
        return (1 + cantidadButacasRestantes) 
        
hayK :: Int -> Sala -> STM Bool
hayK k _ | k <= 0 = return True
hayK _ [] = return False
hayK k (fila:filas) = do
    libres <- cantidadButacasLibres fila
    hayK (k - libres) filas

-- Tiene la precondicion de que sabemos que en esta fila hay k espacios libres contiguos
reservarKContiguasEnFila :: Int -> Fila -> STM [Butaca]
reservarKContiguasEnFila k fila = do
    let bloque = take k fila
    estados <- mapM readTVar bloque
    if or estados then -- Si hay algun true tomamos mal los asientos, asi que probamos desde el siguiente
        reservarKContiguasEnFila k (tail fila)
    else do
        mapM_ (\butaca -> writeTVar butaca True) bloque
        return bloque

reservarKContiguas :: Int -> Sala -> STM [Butaca]
reservarKContiguas k (fila:filas) = do
    hayKContiguasEnEstaFila <- hayKContiguasEnFila k fila
    case (hayKContiguasEnEstaFila) of 
        False -> do 
            butacas <- reservarKContiguas k filas 
            return butacas 
        True -> reservarKContiguasEnFila k fila

-- Reserva hasta k posiciones de una fila, si es que las hay. Si no reserva las que pueda
reservarKEnFila :: Int -> Fila -> STM [Butaca]
reservarKEnFila k _ | k <= 0 = return []
reservarKEnFila _ [] = return []
reservarKEnFila k (butaca:butacas) = do 
    valorActual <- readTVar butaca 
    case valorActual of 
        True -> reservarKEnFila k butacas 
        False -> do 
            writeTVar butaca True 
            resto <- reservarKEnFila (k - 1) butacas
            return (butaca : resto) 

-- Reserva las primeras k butacas que encuentra, sabiendo que hay k butacas libres en la sala
reservarK :: Int -> Sala -> STM [Butaca]
reservarK k _ | k <= 0 = return []
reservarK k (fila:filas) = do 
    butacasEnEstaFila <- reservarKEnFila k fila 
    resto <- reservarK (k - length butacasEnEstaFila) filas 
    return (butacasEnEstaFila ++ resto)

-- Devuelve las butacas asignadas para que despues cuando se levanten las liberen
reservar :: Int -> Sala -> STM [Butaca]
reservar k sala = do 
    hayKContiguasEnAlgunLado <- hayKContiguas k sala 
    hayKEnAlgunLado <- hayK k sala
    case (hayKContiguasEnAlgunLado, hayKEnAlgunLado) of 
        (False, False) -> retry 
        (True, _) -> reservarKContiguas k sala 
        (_, True) -> reservarK k sala 

liberar :: [Butaca] -> STM ()
liberar butacas =
    mapM_ (\butaca -> writeTVar butaca False) butacas