package main

import (
	"fmt"
	"strings"

	"github.com/markkurossi/scheme"
)

type SDK struct {
}

// myAdd is a Go function that will be callable from Scheme
func myAdd(scm *scheme.Scheme, args []scheme.Value) (scheme.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("myAdd expects exactly 2 arguments")
	}

	// Convert the scheme values to int64
	n1, err := scheme.Int64(args[0])
	if err != nil {
		return nil, err
	}
	n2, err := scheme.Int64(args[1])
	if err != nil {
		return nil, err
	}

	// Return the result as a scheme.Value
	return scheme.Int(n1 + n2), nil
}

func main() {
	v := true
	noRun := false
	scm, err := scheme.NewWithParams(scheme.Params{
		Verbose:   v,
		NoRuntime: noRun,
	})
	if err != nil {
		panic(err)
	}

	// Define our function as a builtin
	// addBuiltin := scheme.Builtin{
	// 	Name:   "my-add",
	// 	Args:   []string{"x", "y"},
	// 	Native: myAdd,
	// 	Return: &types.Type{Kind: types.Fixed},
	// }

	// Register the builtin function
	// scm.DefineBuiltin(addBuiltin)

	// Define the complex data structure
	// _ = eval(`
	// 	(define flag-data
	// 		'((data
	// 			(key "flag2")
	// 			(on #f)
	// 			(prerequisites ())
	// 			(targets ())
	// 			(contextTargets ())
	// 			(rules ())
	// 			(fallthrough ())
	// 			(offVariation ())
	// 			(variations ())
	// 			(clientSide #f)
	// 			(salt "")
	// 			(trackEvents #f)
	// 			(trackEventsFallthrough #f)
	// 			(debugEventsUntilDate ())
	// 			(version 0)
	// 			(deleted #f))))
	// `, scm)

	// Define the assoc function first
	_ = eval(`
		(define (assoc key lst)
			(cond
				((null? lst) #f)
				((equal? key (car (car lst))) (car lst))
				(else (assoc key (cdr lst)))))
	`, scm)
	_ = eval(` (define (assoc-cdr key lst) (cdr (assoc key lst))) `, scm)

	_ = eval(`
	(define payload '(
    (flags (
        (flag2 (
					(on #t)
				))
    ))
))`,
		scm)

	// fmt.Println("~~~~~~~~~1")
	// _ = eval(`(assoc-cdr 'flags payload)`, scm)
	// fmt.Println("~~~~~~~~~2")
	// _ = eval(`(car (assoc-cdr 'flags payload))`, scm)
	fmt.Println("~~~~~~~~~3")
	_ = eval(`(assoc 'flag2 (car (assoc-cdr 'flags payload)))`, scm)
	_ = eval(`(assoc-cdr 'on (car (assoc-cdr 'flag2 (car (assoc-cdr 'flags payload)))))`, scm)
	_ = eval(`
		(define (evaluate ctx flag-key default)
			(assoc-cdr 'on (car (assoc-cdr flag-key (car (assoc-cdr 'flags payload)))))
		)
	`, scm)
	_ = eval(`(evaluate "" 'flag2 #f)`, scm)

}

func eval(s string, scm *scheme.Scheme) scheme.Value {
	payload := strings.NewReader(s)
	res, err := scm.Eval("", payload)

	if err != nil {
		panic(err)
	}
	fmt.Printf("Result: %+v\n", res)
	return res
}
