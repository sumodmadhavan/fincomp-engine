{-# LANGUAGE DataKinds #-}
{-# LANGUAGE TypeOperators #-}
{-# LANGUAGE DeriveGeneric #-}
{-# LANGUAGE BangPatterns #-}
{-# LANGUAGE UnboxedTuples #-}

module Main where

import Text.Printf (printf)
import Data.Time.Clock.POSIX (getPOSIXTime)
import Servant
import Network.Wai.Handler.Warp (run)
import GHC.Generics
import Data.Aeson

data FinancialParams = FinancialParams
  { numYears :: !Int
  , auHours :: !Double
  , initialTSN :: !Double
  , rateEscalation :: !Double
  , aic :: !Double
  , hsitsn :: !Double
  , overhaulTSN :: !Double
  , hsiCost :: !Double
  , overhaulCost :: !Double
  }

calculateFinancials :: Double -> FinancialParams -> Double
calculateFinancials rate params = go 1 0
  where
    go !year !acc
      | year > numYears params = acc
      | otherwise =
          let !tsn = initialTSN params + auHours params * fromIntegral year
              !escalatedRate = rate * (1 + rateEscalation params / 100) ^ (year - 1)
              !engineRevenue = auHours params * escalatedRate
              !aicRevenue = engineRevenue * aic params / 100
              !totalRevenue = engineRevenue + aicRevenue
              !hsi = tsn >= hsitsn params && (year == 1 || tsn - auHours params < hsitsn params)
              !overhaul = tsn >= overhaulTSN params && (year == 1 || tsn - auHours params < overhaulTSN params)
              !hsiCost' = if hsi then hsiCost params else 0
              !overhaulCost' = if overhaul then overhaulCost params else 0
              !totalCost = hsiCost' + overhaulCost'
              !totalProfit = totalRevenue - totalCost
              !newAcc = acc + totalProfit
          in go (year + 1) newAcc

goalSeek :: Double -> FinancialParams -> Double -> (# Double, Int #)
goalSeek targetProfit params initialGuess =
  let objective rate = calculateFinancials rate params - targetProfit
      derivative rate =
        let epsilon = 1e-6
            f1 = objective (rate + epsilon)
            f2 = objective rate
        in (f1 - f2) / epsilon
      go !x !iter
        | iter >= 100 = (# x, iter #)
        | abs (objective x) < 1e-8 = (# x, iter + 1 #)
        | otherwise =
            let !fx = objective x
                !f'x = derivative x
                !newX = if f'x /= 0 then x - fx / f'x else x
            in go newX (iter + 1)
  in go initialGuess 0

data ApiInput = ApiInput
  { initialRate :: Double
  , targetProfit :: Double
  } deriving (Generic, Show)

data ApiOutput = ApiOutput
  { optimalRate :: Double
  , iterations :: Int
  , finalCumulativeProfit :: Double
  , computationTimeMicros :: Int
  } deriving (Generic, Show)

instance ToJSON ApiInput
instance FromJSON ApiInput
instance ToJSON ApiOutput
instance FromJSON ApiOutput

type API = "calculate" :> ReqBody '[JSON] ApiInput :> Post '[JSON] ApiOutput

server :: Server API
server = calculateHandler

calculateHandler :: ApiInput -> Handler ApiOutput
calculateHandler input = do
  start <- liftIO getPOSIXTime
  let params = FinancialParams
        { numYears = 10
        , auHours = 450
        , initialTSN = 100
        , rateEscalation = 5
        , aic = 10
        , hsitsn = 1000
        , overhaulTSN = 3000
        , hsiCost = 50000
        , overhaulCost = 100000
        }
      (# optRate, iters #) = goalSeek (targetProfit input) params (initialRate input)
      finalProfit = calculateFinancials optRate params
  end <- liftIO getPOSIXTime
  let elapsedMicros :: Int
      !elapsedMicros = round $ (end - start) * 1000000
  return $ ApiOutput optRate iters finalProfit elapsedMicros

app :: Application
app = serve (Proxy :: Proxy API) server

main :: IO ()
main = do
  putStrLn "Starting server on port 8080..."
  run 8080 app
