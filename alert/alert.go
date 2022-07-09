package alert

type Alert interface {
	Push(title, body string)
}
