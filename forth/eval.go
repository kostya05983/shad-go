//go:build !solution

package main

import (
	"errors"
	"strconv"
	"strings"
)

type Evaluator struct {
	stack     []int
	operators map[string]func() error
}

// NewEvaluator creates evaluator.
func NewEvaluator() *Evaluator {
	e := &Evaluator{
		stack:     make([]int, 0),
		operators: map[string]func() error{},
	}
	e.operators["+"] = e.plus
	e.operators["-"] = e.minus
	e.operators["*"] = e.multiply
	e.operators["/"] = e.divide
	e.operators["dup"] = e.dup
	e.operators["over"] = e.over
	e.operators["drop"] = e.drop
	e.operators["swap"] = e.swap

	return e
}

// Process evaluates sequence of words or definition.
//
// Returns resulting stack state and an error.
func (e *Evaluator) Process(row string) ([]int, error) {
	row = strings.ToLower(row)
	parts := strings.Split(row, " ")

	funcs := make([]func() error, 0)
	newOperator := ""

	for index, part := range parts {
		if part == ";" || part == ":" {
			continue
		}

		num, err := strconv.Atoi(part)
		isNumber := err == nil
		if index == 1 && parts[0] == ":" && isNumber {
			return nil, errors.New("can't redefine numbers")
		}

		// если число
		if err == nil {
			operator := func() error {
				e.stack = append(e.stack, num)

				return nil
			}

			funcs = append(funcs, operator)
			continue
		}

		operator, ok := e.operators[part]

		//override или новый оператор без разницы, всегда высчитываем занова
		if index == 1 && parts[0] == ":" {
			newOperator = part
			continue
		}
		if !ok {
			return nil, errors.New("can't find operator" + part)
		}

		funcs = append(funcs, operator)
	}

	if newOperator != "" && len(funcs) == 0 {
		return nil, errors.New("operator definition is wrong")
	}

	if newOperator != "" {
		e.operators[newOperator] = func() error { return composition(funcs) }

		return e.stack, nil
	}

	err := composition(funcs)
	if err != nil {
		return nil, err
	}

	return e.stack, nil
}

func composition(funcs []func() error) error {
	for _, f := range funcs {
		err := f()
		if err != nil {
			return err
		}
	}

	return nil
}

func (e *Evaluator) poppedTwo() ([]int, error) {
	if len(e.stack) < 2 {
		return nil, errors.New(wrongStateForOperation)
	}

	numbers := e.stack[len(e.stack)-2:]
	e.stack = e.stack[:len(e.stack)-2]

	return numbers, nil
}

func (e *Evaluator) plus() error {
	numbers, err := e.poppedTwo()
	if err != nil {
		return err
	}

	e.stack = append(e.stack, numbers[0]+numbers[1])

	return nil
}

func (e *Evaluator) minus() error {
	numbers, err := e.poppedTwo()
	if err != nil {
		return err
	}

	e.stack = append(e.stack, numbers[0]-numbers[1])

	return nil
}

func (e *Evaluator) multiply() error {
	numbers, err := e.poppedTwo()
	if err != nil {
		return err
	}

	e.stack = append(e.stack, numbers[0]*numbers[1])

	return nil
}

func (e *Evaluator) divide() error {
	numbers, err := e.poppedTwo()
	if err != nil {
		return err
	}
	if numbers[1] == 0 {
		return errors.New("integer division by zero")
	}

	e.stack = append(e.stack, numbers[0]/numbers[1])

	return nil
}

func (e *Evaluator) dup() error {
	if len(e.stack) == 0 {
		return errors.New(wrongStateForOperation)
	}
	e.stack = append(e.stack, e.stack[len(e.stack)-1])

	return nil
}

func (e *Evaluator) over() error {
	if len(e.stack) < 2 {
		return errors.New(wrongStateForOperation)
	}
	e.stack = append(e.stack, e.stack[len(e.stack)-2])

	return nil
}

func (e *Evaluator) drop() error {
	if len(e.stack) == 0 {
		return errors.New(wrongStateForOperation)
	}
	e.stack = e.stack[:len(e.stack)-1]

	return nil
}

func (e *Evaluator) swap() error {
	if len(e.stack) < 2 {
		return errors.New(wrongStateForOperation)
	}
	n := len(e.stack)

	tmp := e.stack[n-1]
	e.stack[n-1] = e.stack[n-2]
	e.stack[n-2] = tmp

	return nil
}
