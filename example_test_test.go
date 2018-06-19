package rego

// Indicates how to use the testing utility provided in the package
func ExampleRun() {

	hippo := struct {
		Name string `json:"name"`
		Age int `json:"age"`
		Friends []string `json:"friends"`
	}{
		Name: "jim",
		Age: 123,
		Friends: []string{"tommy", "ben"},
	}

	test := TestCase{
		Target: "hippo",
		Rules: []string{
			`hippo = {"name":"jim", "age": 123, "friends":["tom", "ben"]} {true}`,
		},
		Expected: hippo,
	}

	// would put a testing.t here
	test.Run(nil, nil, nil)
}
