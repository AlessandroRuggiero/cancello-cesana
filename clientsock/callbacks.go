package clientsock

type Callbacks struct {
	Apricancello    func(*CallbackData)
	Apricancelletto func(*CallbackData)
	SendAperture    func(*CallbackData, string)
}
