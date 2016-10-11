package predicate

type Predicate func() bool

func All(predicates ...Predicate) Predicate {
	return func() bool {
		for _, p := range predicates {
			if !p() {
				return false
			}
		}
		return true
	}
}

func Any(predicates ...Predicate) Predicate {
	return func() bool {
		for _, p := range predicates {
			if p() {
				return true
			}
		}
		return false
	}
}

func Not(predicate Predicate) Predicate {
	return func() bool {
		return !predicate()
	}
}
