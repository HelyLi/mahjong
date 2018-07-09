package mahjong

import (
	"testing"
)

func BenchmarkCheckHu(b *testing.B) {
	aiCards := []aiCard{
		{COLOR_WAN*MAHJONG_MASK + MAHJONG_1, 3},
		{COLOR_WAN*MAHJONG_MASK + MAHJONG_2, 1},
		{COLOR_WAN*MAHJONG_MASK + MAHJONG_3, 1},
		{COLOR_WAN*MAHJONG_MASK + MAHJONG_4, 1},
		{COLOR_WAN*MAHJONG_MASK + MAHJONG_5, 1},
		{COLOR_WAN*MAHJONG_MASK + MAHJONG_6, 1},
		{COLOR_WAN*MAHJONG_MASK + MAHJONG_7, 1},
		{COLOR_WAN*MAHJONG_MASK + MAHJONG_8, 1},
		{COLOR_WAN*MAHJONG_MASK + MAHJONG_9, 3},
	}
	for i := 0; i < b.N; i++ {
		for j := 0; j < len(aiCards); j++ {
			aiCards[j].num++
			CheckHu(aiCards)
			aiCards[j].num--
		}
	}
}

func BenchmarkCheckHuForLZ(b *testing.B) {
	aiCards := []aiCard{
		{COLOR_WAN*MAHJONG_MASK + MAHJONG_1, 3},
		{COLOR_WAN*MAHJONG_MASK + MAHJONG_2, 1},
		{COLOR_WAN*MAHJONG_MASK + MAHJONG_3, 1},
		{COLOR_WAN*MAHJONG_MASK + MAHJONG_4, 1},
		{COLOR_WAN*MAHJONG_MASK + MAHJONG_5, 2},
		{COLOR_WAN*MAHJONG_MASK + MAHJONG_6, 1},
		{COLOR_WAN*MAHJONG_MASK + MAHJONG_7, 1},
		{COLOR_WAN*MAHJONG_MASK + MAHJONG_8, 1},
		{COLOR_WAN*MAHJONG_MASK + MAHJONG_9, 3},
		{COLOR_TIAO*MAHJONG_MASK + MAHJONG_1, 0},
	}
	for i := 0; i < b.N; i++ {
		aiCards[len(aiCards)-1].num = 1
		for j1 := 0; j1 < len(aiCards)-1; j1++ {
			aiCards[j1].num--
			CheckHuForLZ(aiCards, COLOR_TIAO*MAHJONG_MASK+MAHJONG_1)
			aiCards[j1].num++
		}

		aiCards[len(aiCards)-1].num = 2
		for j1 := 0; j1 < len(aiCards)-1; j1++ {
			if aiCards[j1].num == 0 {
				continue
			}
			aiCards[j1].num--
			for j2 := j1; j2 < len(aiCards)-1; j2++ {
				if aiCards[j2].num == 0 {
					continue
				}
				aiCards[j2].num--
				CheckHuForLZ(aiCards, COLOR_TIAO*MAHJONG_MASK+MAHJONG_1)
				aiCards[j2].num++
			}
			aiCards[j1].num++
		}

		aiCards[len(aiCards)-1].num = 3
		for j1 := 0; j1 < len(aiCards)-1; j1++ {
			if aiCards[j1].num == 0 {
				continue
			}
			aiCards[j1].num--
			for j2 := j1; j2 < len(aiCards)-1; j2++ {
				if aiCards[j2].num == 0 {
					continue
				}
				aiCards[j2].num--
				for j3 := j2; j3 < len(aiCards)-1; j3++ {
					if aiCards[j3].num == 0 {
						continue
					}
					aiCards[j3].num--
					CheckHuForLZ(aiCards, COLOR_TIAO*MAHJONG_MASK+MAHJONG_1)
					aiCards[j3].num++
				}
				aiCards[j2].num++
			}
			aiCards[j1].num++
		}

		aiCards[len(aiCards)-1].num = 4
		for j1 := 0; j1 < len(aiCards)-1; j1++ {
			if aiCards[j1].num == 0 {
				continue
			}
			aiCards[j1].num--
			for j2 := j1; j2 < len(aiCards)-1; j2++ {
				if aiCards[j2].num == 0 {
					continue
				}
				aiCards[j2].num--
				for j3 := j2; j3 < len(aiCards)-1; j3++ {
					if aiCards[j3].num == 0 {
						continue
					}
					aiCards[j3].num--
					for j4 := j3; j4 < len(aiCards)-1; j4++ {
						if aiCards[j4].num == 0 {
							continue
						}
						aiCards[j4].num--
						CheckHuForLZ(aiCards, COLOR_TIAO*MAHJONG_MASK+MAHJONG_1)
						aiCards[j4].num++
					}
					aiCards[j3].num++
				}
				aiCards[j2].num++
			}
			aiCards[j1].num++
		}
	}
}

func BenchmarkCheckTing(b *testing.B) {
	aiCards := []aiCard{
		{COLOR_WAN*MAHJONG_MASK + MAHJONG_1, 3},
		{COLOR_WAN*MAHJONG_MASK + MAHJONG_2, 1},
		{COLOR_WAN*MAHJONG_MASK + MAHJONG_3, 1},
		{COLOR_WAN*MAHJONG_MASK + MAHJONG_4, 1},
		{COLOR_WAN*MAHJONG_MASK + MAHJONG_5, 1},
		{COLOR_WAN*MAHJONG_MASK + MAHJONG_6, 1},
		{COLOR_WAN*MAHJONG_MASK + MAHJONG_7, 1},
		{COLOR_WAN*MAHJONG_MASK + MAHJONG_8, 1},
		{COLOR_WAN*MAHJONG_MASK + MAHJONG_9, 3},
	}
	for i := 0; i < b.N; i++ {
		CheckTing(aiCards)
	}	
}

func BenchmarkCheckTingForLZ(b *testing.B) {
	aiCards := []aiCard{
		 {COLOR_WAN*MAHJONG_MASK + MAHJONG_1, 2},
		 {COLOR_WAN*MAHJONG_MASK + MAHJONG_2, 1},
		 {COLOR_WAN*MAHJONG_MASK + MAHJONG_3, 3},
		 {COLOR_WAN*MAHJONG_MASK + MAHJONG_4, 1},
		 {COLOR_WAN*MAHJONG_MASK + MAHJONG_5, 1},
		 {COLOR_WAN*MAHJONG_MASK + MAHJONG_6, 1},		
		 {COLOR_WAN*MAHJONG_MASK + MAHJONG_9, 1},
		 {COLOR_TIAO*MAHJONG_MASK + MAHJONG_1, 3},
	}
	for i := 0; i < b.N; i++ {
		CheckTingForLZ(aiCards, COLOR_TIAO*MAHJONG_MASK + MAHJONG_1)
	}	
}