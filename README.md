# holtwinters
Holt-Winters algorithm

reference by https://github.com/DmitrySerg/DataStart
and rewrite in Golang

Example
model := holtwinters.New(data, 24, 50, 0.11652680227350454, 0.002677697431105852, 0.05820973606789237, 3)
model.TripleExponentialSmoothing()
