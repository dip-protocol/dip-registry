package merkle

type ProofNode struct {
	Hash     string `json:"hash"`
	Position string `json:"position"`
}

func GenerateProof(leaves []string, index int) ([]ProofNode, string) {

	var proof []ProofNode
	level := leaves

	for len(level) > 1 {

		var next []string

		for i := 0; i < len(level); i += 2 {

			if i+1 < len(level) {

				left := level[i]
				right := level[i+1]

				next = append(next, Combine(left, right))

				if i == index {
					proof = append(proof, ProofNode{Hash: right, Position: "right"})
					index = len(next) - 1
				} else if i+1 == index {
					proof = append(proof, ProofNode{Hash: left, Position: "left"})
					index = len(next) - 1
				}

			} else {

				next = append(next, level[i])

				if i == index {
					index = len(next) - 1
				}
			}
		}

		level = next
	}

	root := level[0]

	return proof, root
}