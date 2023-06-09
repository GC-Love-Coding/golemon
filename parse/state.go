package parse

// Each state of the genrated parser's finite state machine
// is encoded as an instance of the following structure
type State struct {
	bp         *Config  // The basis configuration for this state
	cfp        *Config  // ALl configurations in this set
	index      int      // Sequencial number for this satte
	ap         []Action // Array of actions for this state
	nTknAct    int      // Number of actions on terminals
	nNtAct     int      // Number of actions on nonterminals
	iTknOffset int      // yy_action[] offset for terminals
	iNtOfst    int      // yy_action[] offset for nonterminals
	iDefAction int      // Default action
}
