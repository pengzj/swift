package cron

import (
	"testing"
	"time"
)

func TestStr2Int(t *testing.T) {
	one := "123"
	oneInt, err := Str2Int(one)
	if err != nil {
		t.Fatal(err)
	} else {
		t.Log(oneInt)
	}

	one = "12F3"
	oneInt, err = Str2Int(one)
	if err != nil {
		t.Log(err)
	} else {
		t.Fatal(oneInt)
	}

	one = "123$"
	oneInt, err = Str2Int(one)
	if err != nil {
		t.Log(err)
	} else {
		t.Error(oneInt)
	}

	one = "T123"
	oneInt, err = Str2Int(one)
	if err != nil {
		t.Log(err)
	} else {
		t.Error(oneInt)
	}

	one = "T123D"
	oneInt, err = Str2Int(one)
	if err != nil {
		t.Log(err)
	} else {
		t.Error(oneInt)
	}
}

func TestStr2Values(t *testing.T) {
	str := "1,3,5,15-18"
	arrays, err := Str2Values(str)
	if err != nil {
		t.Fatal(err)
	} else {
		t.Log(arrays)
	}
}

func TestStr2Any(t *testing.T) {
	str := "*"
	bools := Str2Any(str)
	if bools == false {
		t.Fatal("incorrect match")
	} else {
		t.Log(bools)
	}
}

func TestStr2Repeat(t *testing.T) {
	str := "*/1"
	repat, err := Str2Repeat(str)
	if err != nil {
		t.Fatal(err)
	} else {
		t.Log(repat)
	}
}

func TestParse(t *testing.T) {
	oneSchedule, err := Parse("* * * * 0,3 *")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(oneSchedule.Week)

	AddJob("*/2 * * * * *", func() {
		t.Log(time.Now(), "hello cron")
	})
}

func TestAddJob(t *testing.T) {
	AddJob("55 * * * * *", func() {
		t.Log(time.Now(), "hello, 55 job")
	})
}
