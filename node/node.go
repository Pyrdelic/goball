package node

import "github.com/hajimehoshi/ebiten/v2"

type Message struct {
	TypeStr  string
	Msg      int
	IntExtra int
}

type Node interface {
	Update() Message
	Draw(screen *ebiten.Image)
}

func Update(n Node) Message {
	return n.Update()
}

func Draw(n Node, screen *ebiten.Image) {
	n.Draw(screen)
}
