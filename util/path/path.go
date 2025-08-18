package path

//多层级用户关系
import (
	"errors"
	"strconv"
	"strings"
)

// 格式  `/[x]/yyy/zzz/` 最后有斜线

type Path struct {
	opts *Options
}

func NewPath(options ...Option) *Path {
	opts := DefaultOptions()
	return &Path{opts: opts}
}

// 把uint32的uid转换成长度短一些的字符串
func (e *Path) Uid2Path(etype string, uid ...uint32) (string, error) {
	if len(uid) > e.opts.maxLevel {
		return "", errors.New("max level limit")
	}

	path := "/" + etype + "/"
	for _, v := range uid {
		if v != 0 {
			path += strconv.FormatUint(uint64(v), e.opts.base) + "/"
		}
	}
	return path, nil
}

func (e *Path) Path2Uid(path string) ([]uint32, error) {
	ids := []uint32{}
	//默认是'/'
	if len(path) < 3 {
		return ids, nil
	}
	arr := strings.SplitSeq(path[3:len(path)-1], "/")
	for v := range arr {
		uid, err := strconv.ParseUint(v, e.opts.base, 64)
		if err != nil {
			return nil, err
		}
		ids = append(ids, uint32(uid))
	}
	return ids, nil
}

// 再原来的path上面再追加uid
func (e *Path) UidPathAppend(path string, uid ...uint32) (string, error) {
	arr := strings.Split(path[3:len(path)-1], "/")
	if len(arr)+len(uid) > e.opts.maxLevel {
		return "", errors.New("max level limit")
	}
	for _, v := range uid {
		if v > 0 {
			path += strconv.FormatUint(uint64(v), e.opts.base) + "/"
		}
	}
	return path, nil
}
