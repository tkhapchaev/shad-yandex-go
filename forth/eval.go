//go:build !solution

package main

import (
	"errors"
	"strconv"
	"strings"
)

type (
	Stack struct {
		top    *node
		length int
	}
	node struct {
		value int
		prev  *node
	}
)

func New() *Stack {
	return &Stack{nil, 0}
}

func (s *Stack) Pop() interface{} {
	if s.length == 0 {
		return nil
	}

	n := s.top
	s.top = n.prev
	s.length--

	return n.value
}

func (s *Stack) Push(value int) {
	n := &node{value, s.top}
	s.top = n
	s.length++
}

func (s *Stack) GetValues() []int {
	if s.length == 0 {
		return []int{}
	}

	var values []int
	var result []int

	for s.top != nil {
		values = append(values, s.top.value)
		s.Pop()
	}

	for i := len(values) - 1; i >= 0; i-- {
		result = append(result, values[i])
	}

	return result
}

type Evaluator struct {
	Memory           Stack
	Commands         map[string][]string
	Operators        map[string]func(int, int) int
	BuiltInOperators map[string]func(int, int) int
}

func NewEvaluator() *Evaluator {
	operators := map[string]func(int, int) int{
		"+": func(a, b int) int { return a + b },
		"-": func(a, b int) int { return a - b },
		"*": func(a, b int) int { return a * b },
		"/": func(a, b int) int { return a / b },
	}

	return &Evaluator{Commands: make(map[string][]string), Operators: operators, BuiltInOperators: operators}
}

func (e *Evaluator) Process(row string) ([]int, error) {
	input := strings.Split(row, " ")

	for i := 0; i < len(input); i++ {
		command := input[i]
		command = strings.ToLower(command)

		switch command {
		case "+":
			if e.Memory.length < 2 {
				return e.Memory.GetValues(), errors.New("not enough elements in memory stack")
			}

			node1 := e.Memory.top
			e.Memory.Pop()
			node2 := e.Memory.top
			e.Memory.Pop()

			operation := e.Operators["+"]
			e.Memory.Push(operation(node1.value, node2.value))
		case "-":
			if e.Memory.length < 2 {
				return e.Memory.GetValues(), errors.New("not enough elements in memory stack")
			}

			node1 := e.Memory.top
			e.Memory.Pop()
			node2 := e.Memory.top
			e.Memory.Pop()

			operation := e.Operators["-"]
			e.Memory.Push(operation(node2.value, node1.value))
		case "*":
			if e.Memory.length < 2 {
				return e.Memory.GetValues(), errors.New("not enough elements in memory stack")
			}

			node1 := e.Memory.top
			e.Memory.Pop()
			node2 := e.Memory.top
			e.Memory.Pop()

			operation := e.Operators["*"]
			e.Memory.Push(operation(node1.value, node2.value))
		case "/":
			if e.Memory.length < 2 {
				return e.Memory.GetValues(), errors.New("not enough elements in memory stack")
			}

			node1 := e.Memory.top
			e.Memory.Pop()
			node2 := e.Memory.top
			e.Memory.Pop()

			if node1.value == 0 {
				return e.Memory.GetValues(), errors.New("division by zero")
			}

			operation := e.Operators["/"]
			e.Memory.Push(operation(node2.value, node1.value))
		case "dup":
			if cmd, exists := e.Commands[command]; exists {
				part1 := input[:i]
				part1 = append(part1, cmd...)
				part2 := input[i+1:]
				input = append(part1, part2...)
				i -= 1
			} else {
				if e.Memory.length == 0 {
					return e.Memory.GetValues(), errors.New("memory stack is empty")
				}

				e.Memory.Push(e.Memory.top.value)
			}
		case "over":
			if cmd, exists := e.Commands[command]; exists {
				part1 := input[:i]
				part1 = append(part1, cmd...)
				part2 := input[i+1:]
				input = append(part1, part2...)
				i -= 1
			} else {
				if e.Memory.length < 2 {
					return e.Memory.GetValues(), errors.New("not enough elements in memory stack")
				}

				node1 := e.Memory.top
				e.Memory.Pop()
				node2 := e.Memory.top
				e.Memory.Pop()
				e.Memory.Push(node2.value)
				e.Memory.Push(node1.value)
				e.Memory.Push(node2.value)
			}
		case "drop":
			if cmd, exists := e.Commands[command]; exists {
				part1 := input[:i]
				part1 = append(part1, cmd...)
				part2 := input[i+1:]
				input = append(part1, part2...)
				i -= 1
			} else {
				if e.Memory.length == 0 {
					return e.Memory.GetValues(), errors.New("memory stack is empty")
				}

				e.Memory.Pop()
			}
		case "swap":
			if cmd, exists := e.Commands[command]; exists {
				part1 := input[:i]
				part1 = append(part1, cmd...)
				part2 := input[i+1:]
				input = append(part1, part2...)
				i -= 1
			} else {
				if e.Memory.length < 2 {
					return e.Memory.GetValues(), errors.New("not enough elements in memory stack")
				}

				node1 := e.Memory.top
				e.Memory.Pop()
				node2 := e.Memory.top
				e.Memory.Pop()
				e.Memory.Push(node1.value)
				e.Memory.Push(node2.value)
			}
		case ":":
			definition := input[2 : len(input)-1]
			subcommand := input[1]

			_, err := strconv.Atoi(subcommand)

			if err == nil {
				return e.Memory.GetValues(), errors.New("cannot redefine a number")
			}

			if subcommand == "dup" || subcommand == "drop" || subcommand == "over" || subcommand == "swap" {
				if len(definition) != 1 {
					return e.Memory.GetValues(), nil
				}

				if definition[0] != "dup" && definition[0] != "drop" && definition[0] != "over" && definition[0] != "swap" {
					return e.Memory.GetValues(), nil
				}
			}

			if subcommand == "+" || subcommand == "-" || subcommand == "*" || subcommand == "/" {
				e.Operators[subcommand] = e.BuiltInOperators[definition[0]]

				return e.Memory.GetValues(), nil
			}

			for j := 0; j < len(definition); j++ {
				if definition[j] == subcommand {
					part1 := definition[:j]
					part1 = append(part1, e.Commands[subcommand]...)
					part2 := definition[j+1:]
					definition = append(part1, part2...)
				}
			}

			// cringe
			for j := 0; j < len(definition); j++ {
				_, exists := e.Commands[definition[j]]

				if exists {
					part1 := definition[:j]
					part1 = append(part1, e.Commands[definition[j]]...)
					part2 := definition[j+1:]
					definition = append(part1, part2...)
				}
			}

			e.Commands[strings.ToLower(subcommand)] = definition

			return e.Memory.GetValues(), nil
		default:
			if cmd, exists := e.Commands[command]; exists {
				part1 := input[:i]
				part1 = append(part1, cmd...)
				part2 := input[i+1:]
				input = append(part1, part2...)
				i -= 1
			} else {
				integer, err := strconv.Atoi(command)

				if err == nil {
					e.Memory.Push(integer)
				} else {
					return e.Memory.GetValues(), err
				}
			}
		}
	}

	return e.Memory.GetValues(), nil
}
