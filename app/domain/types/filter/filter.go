package filter

type Id byte

const (
	All  Id = 0xFF // 全検索条件に引っ掛けるためのフィルタ
	Turf Id = 0x01
	Dirt Id = 0x02
	Jump Id = 0x04
)

var filterIdMap = map[Id]string{
	All:  "条件なし",
	Turf: "芝",
	Dirt: "ダート",
	Jump: "障害",
}

func (i Id) Value() int {
	return int(i)
}

func (i Id) String() string {
	id, _ := filterIdMap[i]
	return id
}
