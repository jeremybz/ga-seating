# ga-seating
Find good-enough solutions to simple-enough problems of seating tables to maximize individual meetups.
This is my first project in golang, so a lot of it is wrong.

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes.

### Prerequisites

```
[golang](https://golang.org/doc/install)
```

### Installing and Running

Install golang

```
$ sudo apt install golang
```

Set up golang's directory structure
```
$ mkdir -p ~/go/{src/ga-seating,bin,pkg}
$ export GOPATH=~/go
```

Copy the source
```
$ cp ga-seating.go ~/go/src/ga-seating/
$ cd ~/go/src/ga-seating
$ go build
```

Run
```
$ ./ga-seating
```

## Problem Statement

Given a group of eighteen individuals, and five tables to sit them at, find a schedule of seatings such that each individual is seated with every other individual at least once.  The tables can seat four, four, four, three, and three individuals, respectively.


## Algorithm

1. Start with a randomly generated set of schedules
1. Evaluate each schedule according to how many individuals each individual sits with during the course of the scheduled seatings
1. Order the schedules according to their evaluated fitness
1. Replace the weakest schedules with a new set of schedules generated based on the existing population
1. Evaluate the new schedules
1. Repeat from step 3


## Data

* **seating**: a seating plan: matches individuals to seats.  To avoid permutation complexity, individuals are sorted by table: `EFGH ABCD IJKL MNO PQR` is a valid seating, but `BACD EFGH IJKL MNO PQR` is not.
* **schedule**: a list of `rotations` seatings: ```ABCD  EFGH  IJKL  MNO  PQR
CGIM  FJKR  ABNQ  EHL  DOP
DFLM  BCHQ  GKNP  IJO  AER
GLQR  DEIN  CHKO  AFP  BJM
AHIP  EJNQ  BFLO  DGK  CMR
CEOR  AFLN  BGIP  DHJ  KMQ
DHNR  AGMO  CJLP  BEK  FIQ
AGJK  EHMP  CFLN  BIR  DOQ
```


## Messing Around

Fun variables to play with include:
* **tables**: list of tables with number of seats.  Order is arbitrary.
* **rotations**: number of times to change seating.  My guess is that 8 is the minimum with a full solution.
* **pop_size**: total number of schedules to keep alive
* **gen_size**: number of new schedules to create each round
* **mutation_permille**: chance in a thousand of a new schedule being mutated


## Built With

* [vim](https://www.vim.org/)


## Contributing

PRs welcome.


## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details


## Acknowledgments

* Inspiration from Rey
* README taken from [PurpleBooth](https://gist.github.com/PurpleBooth/109311bb0361f32d87a2)
* I learned GAs a long time ago from [Peter G. Anderson](https://www.cs.rit.edu/~pga/)
