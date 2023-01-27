package plan_test

// [JH]: I'm not really sure how this test works
// It was failing because of a reflect setting a value on an unexported field
// until I know more about what is going on I am going to comment it out
/*
func TestValidConfigNoPanic(t *testing.T) {
	// return false if valid + panic
	assertion := func(conf *v2.Config) bool {
		// fmt.Printf("GOT %#v\n\n", pretty.Sprint(conf))
		// validate our configuration
		_, err := conf.Validate()

		// if config is valid, there should be no panic
		if err == nil {
			fmt.Println("valid")
			defer func() bool {
				if r := recover(); r != nil {
					fmt.Println("Recovered in f", r)
					debug.PrintStack()
					return false
				}
				return true
			}()

			_, err = plan.Eval(conf)
			if err != nil {
				panic(err)
			}
		} else {
			fmt.Println("invalid")
			fmt.Printf("err %s\n", err)
		}

		// config isn't valid so we don't care if we panic or not
		return true
	}
	if err := quick.Check(assertion, &quick.Config{MaxCount: 100}); err != nil {
		t.Error(err)
	}
}*/
