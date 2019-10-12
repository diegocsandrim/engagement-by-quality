package git

type Contributor struct {
	Id string
}

func NewContributor(id string) *Contributor {
	c := Contributor{
		Id: id,
	}
	return &c
}
