package main

func (s *SpaceResponse) GetPercentUsed() int {
	if s.TotalSpace == 0 {
		return 1
	}
	return int((s.UsedSpace * 100) / s.TotalSpace)
}
