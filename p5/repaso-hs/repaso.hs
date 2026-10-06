import Prelude hiding (mapM)

-- Ejercicio 1 
longitud :: [a] -> Int 
longitud s = foldr (\_ acumulador -> acumulador + 1) 0 s

sumatoria :: [Int] -> Int 
sumatoria = foldr (\x acc -> x + acc) 0 

pertenece :: Eq a => a -> [a] -> Bool 
pertenece e = foldr (\x acc -> e == x || acc) False 

quitar :: Eq a => a -> [a] -> [a] 
quitar _ [] = []
quitar e (x:xs) 
    | e == x = xs 
    | otherwise = x : quitar e xs 

-- Ejercicio 2 
cuadrados :: [Int] -> [Int] 
cuadrados = map (\x -> x ^ 2)

positivos :: [Int] -> [Int] 
positivos = filter (\a -> a > 0)

cantidad :: (a -> Bool) -> [a] -> Int 
cantidad p = foldr (\x acc -> if p x then 1 + acc else acc) 0

-- Ejercicio 3 
data Arbol a = Hoja | Nodo ( Arbol a ) a ( Arbol a )

foldArbol :: b -> (b -> a -> b -> b) -> Arbol a -> b 
foldArbol cHoja _ Hoja = cHoja 
foldArbol cHoja cArbol (Nodo izq valor der) = cArbol (foldArbol cHoja cArbol izq) valor (foldArbol cHoja cArbol der)


insertar :: Ord a => a -> Arbol a -> Arbol a 
insertar e Hoja = Nodo Hoja e Hoja
insertar e (Nodo izq valor der)
    | e == valor = (Nodo izq valor der)
    | e < valor = (Nodo (insertar e izq) valor der)
    | otherwise = (Nodo izq valor (insertar e der))

estaEn :: Ord a => a -> Arbol a -> Bool 
estaEn _ Hoja = False 
estaEn e (Nodo izq valor der)
    | valor == e = True 
    | valor < e = estaEn e izq 
    | otherwise = estaEn e der 

enOrden :: Arbol a -> [a]
enOrden Hoja = []
enOrden (Nodo izq valor der) = enOrden izq ++ [valor] ++ enOrden der

desdeLista :: Ord a => [a] -> Arbol a
desdeLista = foldr insertar Hoja

-- Ejercicio 4 
safeHead :: [a] -> Maybe a 
safeHead [] = Nothing 
safeHead (x:_) = Just x 

buscar :: Eq k => k -> [(k, v)] -> Maybe v 
buscar _ [] = Nothing 
buscar clave ((c, v):xs) 
    | clave == c = Just v 
    | otherwise = buscar clave xs 

dividir :: Int -> Int -> Maybe Int 
dividir _ 0 = Nothing 
dividir n d = Just (n `div` d)

data Cuenta = Cuenta { titular :: String , saldo :: Int } 

depositar :: Int -> Cuenta -> Cuenta
depositar monto (Cuenta nombre saldoActual) = Cuenta nombre (saldoActual + monto)

extraer :: Int -> Cuenta -> Maybe Cuenta
extraer monto (Cuenta nombre saldoActual)
    | monto > saldoActual = Nothing
    | otherwise = Just (Cuenta nombre (saldoActual - monto))

-- Ejercicio 5 
calcularCase :: Int -> Int -> Int -> Maybe Int
calcularCase a b c =
    case dividir a b of
        Nothing -> Nothing
        Just x ->
            case dividir x c of
                Nothing -> Nothing
                Just y -> Just (y + 1)

calcularBind :: Int -> Int -> Int -> Maybe Int 
calcularBind a b c = dividir a b >>= (\x -> dividir x c) >>= (\y -> Just (y + 1))

calcularDo :: Int -> Int -> Int -> Maybe Int 
calcularDo a b c = do 
        x <- dividir a b 
        y <- dividir x c 
        return (y + 1)

transferir :: Int -> Cuenta -> Cuenta -> Maybe (Cuenta, Cuenta)
transferir monto cuentaOrigen cuentaDestino = do
    cuentaOrigenPostExtraer <- extraer monto cuentaOrigen
    let cuentaDestinoPostDeposito = depositar monto cuentaDestino
    return (cuentaOrigenPostExtraer, cuentaDestinoPostDeposito)

dividirE :: Int -> Int -> Either String Int
dividirE a b = 
    case dividir a b of 
        Nothing -> Left "No se puede dividir por cero"
        Just x -> Right x 

extraerE :: Int -> Cuenta -> Either String Cuenta
extraerE monto cuenta =
    case extraer monto cuenta of
        Nothing -> Left "No hay plata"
        Just cuentaPostExtraer -> Right cuentaPostExtraer

transferirE :: Int -> Cuenta -> Cuenta -> Either String (Cuenta, Cuenta)
transferirE monto cuentaOrigen cuentaDestino = 
    case transferir monto cuentaOrigen cuentaDestino of 
        Nothing -> Left "No hay plata"
        Just res -> Right res

-- Ejercicio 6
pares :: [a] -> [b] -> [(a, b)]
pares l1 l2 = do 
    a <- l1
    b <- l2 
    let actual = (a,b)
    return actual  

pitagoricas :: Int -> [(Int, Int, Int)]
pitagoricas n = do
    a <- [1..n]
    b <- [a..n]
    c <- [b..n]
    if a^2 + b^2 == c^2
        then return (a, b, c)
        else []

-- Ejercicio 7
when :: Monad m => Bool -> m () -> m () 
when condition m 
    | condition = m 
    | otherwise = pure ()

replicateM_ :: Monad m => Int -> m a -> m () 
replicateM_ 0 _ = pure ()
replicateM_ n m = do 
    m 
    replicateM_ (n - 1) m

mapM :: Monad m => (a -> m b) -> [a] -> m [b] 
mapM _ [] = pure []
mapM f (x:xs) = do
    y <- f x
    ys <- mapM f xs
    pure (y : ys)

forM_ :: Monad m => [a] -> (a -> m b) -> m ()
forM_ [] _ = pure ()
forM_ (x:xs) f = do
    f x
    forM_ xs f

forever :: Monad m => m a -> m b 
forever m = do 
    m 
    forever m

-- Ejercicio 8 
saludar :: IO () 
saludar = do 
    putStrLn "Hola, dame tu nombre"
    nombre <- getLine 
    putStrLn ("Hola, " ++ nombre) 

sumarHastaCero :: IO Int
sumarHastaCero = do 
    putStrLn "Ingresa un numero, si es cero la corto"
    numero <- readLn :: IO Int 
    if numero == 0 then 
        return 0 
        else do
            sumaResto <- sumarHastaCero
            return (numero + sumaResto)

-- Ejercicio 9
hola :: IO ()
hola = putStrLn " hola "

-- Imprime solo chau porque nunca se ejecuta el hola, solo se lo asigna a x
programa1 :: IO ()
programa1 = do
    let x = hola
    putStrLn " chau "

-- Se imprimen dos veces hola porque se llama dos veces a x
programa2 :: IO ()
programa2 = do
    let x = hola
    x
    x

-- Imprime 5 por el print n del final:
-- n <- return 5 obtiene el valor 5 de una accion IO y lo vincula a n; no imprime.
-- return 10 produce 10 dentro de IO, pero su resultado se descarta: no cambia n,
-- no imprime ni termina la funcion. return envuelve un valor en la monada.
-- Finalmente, print n imprime 5, que es el valor que sigue teniendo n.
programa3 :: IO ()
programa3 = do
    n <- return 5
    return 10
    print n

-- Primero imprime el resultado de quedarse con el elemento en la segunda posición
-- Dp imprime la cabeza. 
programa4 :: IO ()
programa4 = do
    let as = [putStrLn " a ", putStrLn " b ", putStrLn " c " ]
    as !! 2
    head as

-- Ejercicio 10 
adivinar :: Int -> IO ()
adivinar secreto = adivinarAux secreto 0

adivinarAux :: Int -> Int -> IO ()
adivinarAux secreto intentos = do
    putStrLn "Ingresa un numero para adivinar el secreto:"
    numero <- readLn :: IO Int
    let intentosActuales = intentos + 1
    if numero == secreto
        then putStrLn ("Acertaste en " ++ show intentosActuales ++ " intentos.")
        else do
            if numero < secreto
                then putStrLn "El numero secreto es mas alto."
                else putStrLn "El numero secreto es mas bajo."
            adivinarAux secreto intentosActuales

main :: IO ()
main = adivinar 42
