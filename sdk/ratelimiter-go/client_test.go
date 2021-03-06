package ratelimiter

import (
	"testing"
	"time"

	"github.com/golang/protobuf/jsonpb"
)

var (
	rcTypeId        = []byte("zwf-test")
	quota    uint32 = uint32(100000)
	expire   int64  = int64(1000)
	cluster         = []string{
		"127.0.0.1:20000",
		"127.0.0.1:20001",
		"127.0.0.1:20002",
	}
)

func TestRegistQuota(t *testing.T) {
	l, err := New(&ClientConfig{
		Cluster: cluster,
	})
	if err != nil {
		t.Fatal(err.Error())
	}
	defer l.Close()

	if err := l.RegistQuota(rcTypeId, quota, 60); err != nil {
		t.Fatal(err.Error())
	} else {
		t.Logf("Regist quota OK")
	}
}

func TestDeleteQuota(t *testing.T) {
	l, err := New(&ClientConfig{
		Cluster: cluster,
	})
	if err != nil {
		t.Fatal(err.Error())
	}
	defer l.Close()

	if err := l.DeleteQuota(rcTypeId); err != nil {
		t.Fatal(err.Error())
	} else {
		t.Logf("Delete quota OK")
	}
}

func TestBorrow_(t *testing.T) {
	l, err := New(&ClientConfig{
		Cluster: cluster,
	})
	if err != nil {
		t.Fatal(err.Error())
	}
	defer l.Close()

	// 借用资源
	start := time.Now()
	rcId, err := l.Borrow(rcTypeId, expire)
	if err != nil {
		t.Fatal(err.Error())
	} else {
		t.Logf("Borrow OK, rcId=%s, elapse:%v\n", rcId, time.Since(start))
	}
	_ = rcId

	// 模拟做业务处理
	//time.Sleep(5 * time.Second)

	// 归还资源
	start = time.Now()
	if err = l.Return(rcId); err != nil {
		t.Fatal(err.Error())
	} else {
		t.Logf("Return OK, rcId=%s, elapse:%v\n", rcId, time.Since(start))
	}
}

func TestBorrowWithTimeout(t *testing.T) {
	l, err := New(&ClientConfig{
		Cluster: cluster,
	})
	if err != nil {
		t.Fatal(err.Error())
	}
	//defer l.Close()

	// 借用资源
	start := time.Now()
	rcId, err := l.BorrowWithTimeout(rcTypeId, expire, 10*time.Second)
	if err != nil {
		t.Fatal(err.Error())
	} else {
		t.Logf("Borrow OK, rcId=%s, elapse:%v\n", rcId, time.Since(start))
	}
}

func TestResourceList(t *testing.T) {
	l, err := New(&ClientConfig{
		Cluster: cluster,
	})
	if err != nil {
		t.Fatal(err.Error())
	}

	start := time.Now()
	rcList, err := l.ResourceList(nil)
	if err != nil {
		t.Fatal(err.Error())
	} else {
		t.Logf("List Resource OK. elapse:%v\n", time.Since(start))
		var encoder = jsonpb.Marshaler{EmitDefaults: true}
		for i, dt := range rcList {
			str, _ := encoder.MarshalToString(dt)
			t.Logf("(%d) %s\n", i, str)
		}
	}
}

func BenchmarkBorrow(b *testing.B) {
	l, err := New(&ClientConfig{
		Cluster: cluster,
	})
	if err != nil {
		b.Fatal(err.Error())
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rcId, err := l.BorrowWithTimeout(rcTypeId, expire, 10*time.Second)
		if err != nil {
			b.Fatal(err.Error())
		}
		if err = l.Return(rcId); err != nil {
			b.Fatal(err.Error())
		}
	}
}
