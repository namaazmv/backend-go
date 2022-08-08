package lib

import "time"

type Prayer struct {
	Atolls      []*Atoll
	PrayerTimes []*PrayerTime
	Islands     []*Island
	Timings     []string
}

func (p Prayer) GetAtoll(Id int) *Atoll {
	var atoll *Atoll = nil

	for i := range p.Atolls {
		if p.Atolls[i].CategoryId == Id {
			atoll = p.Atolls[i]
		}
	}

	return atoll
}

func (p Prayer) GetIsland(Id int) *Island {
	var island *Island = nil

	for i := range p.Islands {
		if p.Islands[i].IslandId == Id {
			island = p.Islands[i]
		}
	}

	return island
}

func (p Prayer) GetEntryFromDay(day int, island *Island) *PrayerTime {
	var entry *PrayerTime = nil

	for i := range p.PrayerTimes {
		prayer := p.PrayerTimes[i]

		if prayer.Date == day && prayer.CategoryId == island.CategoryId {
			entry = prayer
		}
	}

	return entry
}

func (p Prayer) GetToday(island *Island) *PrayerTime {
	return p.GetEntryFromDay(DaysIntoYear(time.Now().Local()), island)
}
