package internal

type Parser struct {
	engine *Engine
}

func NewParser(engine *Engine) *Parser {
	return &Parser{
		engine: engine,
	}
}
