package mecachis

type Cacher interface {
	Insert(k interface{}, value interface{}) bool
	Upsert(k interface{}, value interface{})
	Access(k interface{}, value interface{}) (interface{}, error)
	Size() int
	Dump() string
}
