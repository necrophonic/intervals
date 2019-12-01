# Intervals

Intervals is a beeper for exercise intervals. Start it up with a set of intervals you wish to use and it'll beep as each times out

## Usage

```shell
intervals [--sets|-s n] [--rest|-r n] --intervals|-i <durations> 

  -d, --delay      amount of seconds delay before the first set starts
  -i, --intervals  (optional) number of sets (default [])
  -r, --rest       amount of seconds rest between sets
  -s, --sets       the number of times to repeat the interval set (default 1)
```

Example: 3 intervals of 10, 20 then 15 seconds

```shell
> intervals -i 10s,20s,15s
```

Example: Repeat set of intervals ( 3x 15 15 25 )

```shell
> intervals --sets 3 -i 15s,15s,25
```

Example: Repeat sets with rest

```shell
> intervals -s 3 -r 30s -i 30s,30s,30s
```

## License

Copyright (c) 2019 J Gregory

Released under MIT license, see [LICENSE](LICENSE) for details.
