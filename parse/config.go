package parse

// A followset propagation link indicates that the contents of one
// configuration followset should be propagated to another whenever
// the first changes.
type PLink struct {
	cfp  *Config // The configuration to which linked
	next *PLink  // The next progagate link
}

// A configuration is a production rule of the grammar together with
// a mark (dot) showing how much of that rule has been processed so far.
// Configurations also contain a follow-set which is a list of terminal
// symbols which are allowed to immediately follow the end of the rule.
// Every configuration is recorded as an instance of the following:
type Config struct {
	rp   *Rule  // The rule upon which the configuration is based
	dot  int    // The parse point
	fws  string // Follow-set for this configuration only
	fplp *PLink // Follow-set forward propagation links
	bplp *PLink // Follow-set backward propagation links

}
