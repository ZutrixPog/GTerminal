package terminal

type TerminalRepo interface {
	GetTerminal(token string) *Terminal
	SetTerminal(token string, terminal *Terminal) error
}

type MemTerminalRepo struct {
	terminals map[string]*Terminal
}

func NewMemTerminalRepo() *MemTerminalRepo {
	return &MemTerminalRepo{
		terminals: make(map[string]*Terminal),
	}
}

func (r *MemTerminalRepo) GetTerminal(token string) *Terminal {
	return r.terminals[token]
}

func (r *MemTerminalRepo) SetTerminal(token string, terminal *Terminal) error {
	r.terminals[token] = terminal
	return nil
}
