package ports

// An interface which defines from how many ways the user will have acces to our core application
// for ex: through init command

type Commands interface{
	// right now we don't know what this function will take and return
	Init()
}