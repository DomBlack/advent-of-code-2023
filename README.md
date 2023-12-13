# Advent of Code λ2023

This repository contains my solutions for the [Advent of Code 2023](https://adventofcode.com/2023), written in [Go](https://goland.org).

## Project Structure

The project is structured in a way that each day of the Advent of Code challenge has its own package under the `internal`
directory. For example, the solution for Day 10 can be found in [`internal/day10/day10.go`](internal/day10/day10.go).

The [`cmd/aoc2023`](cmd/aoc2023) package is the main entry point of the application.
It combines all the days into a single application.

All my personalised inputs are stored in the [`inputs`](inputs) directory, and each day has a corresponding text file
in there. I've also included expected outputs for each day in the day's package using a call to `WithExpectedAnswers`,
this is to allow for regression testing after refactoring.

The `pkg` directory contains various packages that are used across the solutions.

- [`pkg/runner`](pkg/runner) contains a generic runner for registering, running and testing the solutions for each day.
  Every day uses this without fail. It will be at the top of each day's go file, and pretty much the only thing in the days
  test file.
- [`pkg/stream`](pkg/stream) contains various stream functors and sinks used across the solutions.
- [`pkg/alogrithms`](pkg/algorithms) contains various common algorithms used across the solutions.
- [`pkg/datastructures`](pkg/datastructures) contains various common data structures used across the solutions.

## Running the Code

To run the code, you need to have Go installed on your system. You can download it from the [official Go website](https://golang.org/dl/).

Once you have Go installed, you can run the code for all days by navigating to the root directory of the project and running:

```bash
go run ./cmd/aoc2023
```

If you want to run the code for a specific day only, you can do so by providing the `--day` flag followed by the day number. For example, to run the code for Day 10, you would do:

```bash
go run ./cmd/aoc2023 --day 10
```

## Contributing

While this is primarily a personal project, contributions are welcome. If you see an issue or have a suggestion for improvement, feel free to open an issue or submit a pull request.

## License

This project is licensed under the MIT License. See the [`LICENSE`](LICENSE) file for more details.
