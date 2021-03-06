package mahjong

import (
	"sort"
	"util"
)

type ChiCard struct {
	CardId    int32
	CardType  int32
	FromIndex int32
	ToIndex   int32
	ChiPosBit int32 //吃位置的二进制标示，位运算
}

func NewChiCard() *ChiCard {
	c := &ChiCard{
		CardId:    0,
		CardType:  0,
		FromIndex: 0,
		ToIndex:   0,
		ChiPosBit: 0,
	}
	return c
}

func (c *ChiCard) SelectFirstChi() {
	if c.CardType == MJ_CHI_CHI {
		for i := uint32(1); i <= 3; i++ {
			if c.ChiPosBit&(0x01<<i) != 0 {
				c.ChiPosBit = 0x01 << i
				break
			}
		}
	}
}

type IMahjong interface {
	//是否癞子胡
	IsHasLaiziHu() bool
	//获取癞子牌
	GetLaiziCard() []int32
	//设置癞子牌
	SetLaiziCard(lzCards []int32)
	//是否为癞子牌
	IsLaiziCard(card int32) bool
	//是否胡七对
	IsHasQidui() bool
	//获取吃牌数组
	GetChiCard(chiCard *ChiCard) []int32
	//获取胡牌类型及分组信息
	CheckHuType(chiCards []*ChiCard, holdCards []int32, card int32, flag uint32) (uint64, [][][3]int32)
	//获取听牌信息
	CheckTing(chiCards []*ChiCard, holdCards []int32, all bool) []int32
	//获取向听数及对应的听牌信息
	CheckNTing(chiCards []*ChiCard, holdCards []int32, N int32, all bool) (int32, []int32)
}

type RuleMahjong struct {
	rType   int32
	lzCards []int32
}

func NewRuleMahjong(rType int32) *RuleMahjong {
	return &RuleMahjong{
		rType: rType,
	}
}

func (rule *RuleMahjong) IsHasLaiziHu() bool {
	return rule.rType == RULE_HN_MAHJONG_HONGZHONG
}

func (rule *RuleMahjong) GetLaiziCard() []int32 {
	return rule.lzCards
}

func (rule *RuleMahjong) SetLaiziCard(lzCards []int32) {
	rule.lzCards = lzCards
}

func (rule *RuleMahjong) IsLaiziCard(card int32) bool {
	for _, v := range rule.lzCards {
		if v == card {
			return true
		}
	}
	return false
}

func (rule *RuleMahjong) IsHasQidui() bool {
	return rule.rType == RULE_MAHJONG_GUOBIAO || rule.rType == RULE_SC_MAHJONG_XUELIU || rule.rType == RULE_SC_MAHJONG_XUEZHAN
}

func (rule *RuleMahjong) GetChiCard(chiCard *ChiCard) []int32 {
	var cards []int32
	if chiCard.CardType == MJ_CHI_GANG || chiCard.CardType == MJ_CHI_GANG_WAN || chiCard.CardType == MJ_CHI_GANG_AN {
		cards = []int32{chiCard.CardId, chiCard.CardId, chiCard.CardId, chiCard.CardId}
	} else if chiCard.CardType == MJ_CHI_PENG {
		cards = []int32{chiCard.CardId, chiCard.CardId, chiCard.CardId}
	} else if chiCard.CardType == MJ_CHI_CHI {
		c := chiCard.CardId
		c1, c2, c3 := int32(0), int32(0), int32(0)
		if chiCard.ChiPosBit&(0x01<<1) != 0 {
			c1, c2, c3 = c-2, c-1, chiCard.CardId
		} else if chiCard.ChiPosBit&(0x01<<2) != 0 {
			c1, c2, c3 = c-1, chiCard.CardId, c+1
		} else if chiCard.ChiPosBit&(0x01<<3) != 0 {
			c1, c2, c3 = chiCard.CardId, c+1, c+2
		}
		cards = []int32{c1, c2, c3}
	}
	return cards
}

func (rule *RuleMahjong) CheckHuType(chiCards []*ChiCard, holdCards []int32, card int32, flag uint32) (uint64, [][][3]int32) {
	huType := MJ_HU_TYPE_NONE
	var groups [][][3]int32
	if card != 0 && len(holdCards)%3 != 1 || card == 0 && len(holdCards)%3 != 2 {
		return huType, groups
	}

	aiCards, lzNum := ConvertSliceToAICard(holdCards, card, rule.lzCards)

	if rule.IsHasQidui() && len(chiCards) == 0 {
		if ok, group := CheckQiDuiHu(aiCards, lzNum, flag&MJ_CHECK_HU_FLAG_GROUP != 0); ok {
			huType |= MJ_HU_TYPE_QIDUI
			groups = append(groups, group)
		}
	}
	if ok, group := CheckCommonHu(aiCards, lzNum, flag&MJ_CHECK_HU_FLAG_GROUP != 0); ok {
		huType |= MJ_HU_TYPE_COMMON
		groups = append(groups, group...)
	}
	return huType, groups
}

func (rule *RuleMahjong) CheckTing(chiCards []*ChiCard, holdCards []int32, all bool) []int32 {
	tingInfo := make([]int32, 0, 34)

	aiCards, lzNum := ConvertSliceToAICard(holdCards, 0, rule.lzCards)

	if rule.IsHasQidui() && len(chiCards) == 0 {
		tingInfo = append(tingInfo, CheckQiDuiTing(aiCards, lzNum)...)
	}
	if len(tingInfo) > 0 && (!all || tingInfo[0] == MAHJONG_ANY) {
		return tingInfo
	}

	tingInfo = append(tingInfo, CheckCommonTing(aiCards, lzNum, all)...)
	return tingInfo
}

/* 检查N向听牌，N: 0,1,2,3...，返回实际向听数及相应的向听牌
** N为0，等价于CheckTing函数，N为1，一上一听，依次类推，时间复杂度为组合数C{M,N}*Ot，M手牌数，Ot为检查听牌的时间复杂度
** 手牌数为14，最坏情况为四上一听，手牌数为17，最坏情况为五上一听，依次类推
 */
func (rule *RuleMahjong) CheckNTing(chiCards []*ChiCard, holdCards []int32, N int32, all bool) (int32, []int32) {
	if len(holdCards)%3 != 1 || N < 0 {
		return -1, []int32{}
	}

	cards := make([]int32, len(holdCards))
	copy(cards, holdCards)
	if len(cards) > 3 {
		cards = rule.removeSolKeSeq(cards)
	}
	sort.Slice(cards, func(i, j int) bool { return cards[i] < cards[j] })

	indexs := make([]int, 0, len(cards))
	for i := 0; i < len(cards); i++ {
		if rule.IsLaiziCard(cards[i]) || i > 0 && cards[i] == cards[i-1] {
			continue
		}
		indexs = append(indexs, i)
	}

	//对indexs按权重预排序会提高剪枝效率，能加快搜索树的返回，对于实际听数小于检测听数的剪枝效率很高，其它情况效率一样
	cardVal := GetCardValue(rule, cards, 0)
	mapVal := make(map[int32]float32)
	for _, v := range cardVal {
		mapVal[v.Card] = v.Val
	}
	sort.Slice(indexs, func(i, j int) bool { return mapVal[cards[indexs[i]]] < mapVal[cards[indexs[j]]] })

	baseN, tings := rule.doCheckNTing(chiCards, cards, N, indexs, 0, all)
	if len(tings) == 0 {
		return -1, tings
	}
	return N - baseN, tings
}

func (rule *RuleMahjong) doCheckNTing(chiCards []*ChiCard, cards []int32, N int32, indexs []int, minIndex int, all bool) (int32, []int32) {
	tings := rule.CheckTing(chiCards, cards, all)
	if len(tings) > 0 || N == 0 {
		return N, tings
	}

	//由于麻将牌型的复杂性，目前只能抽取牌，然后加上癞子牌，取听牌的并集
	maxBaseN, lzCard := int32(0), int32(0)
	if len(rule.lzCards) > 0 {
		lzCard = rule.lzCards[0]
	}
	for k := minIndex; k < len(indexs); k++ {
		bkc := cards[indexs[k]]
		cards[indexs[k]] = lzCard
		baseN, subTings := rule.doCheckNTing(chiCards, cards, N-maxBaseN-1, indexs, k+1, all)
		if baseN > 0 { //加快向听搜索
			maxBaseN += baseN
			tings = tings[0:0]
		}
		tings = append(tings, subTings...)
		cards[indexs[k]] = bkc
	}

	if maxBaseN == 0 && N == 2 {
		duiIndex := make([]int, 0, len(indexs))
		for k := 0; k < len(indexs); k++ {
			if cards[indexs[k]] != 0 && rule.IsLaiziCard(cards[indexs[k]]) == false &&
				indexs[k] < len(cards)-1 && cards[indexs[k]+1] == cards[indexs[k]] &&
				(indexs[k] > len(cards)-3 || cards[indexs[k]+1] != cards[indexs[k]+2]) {
				duiIndex = append(duiIndex, k)
			}
		}
		if len(duiIndex) >= 5 { //当对子数超过5个时，需要尝试拆对子来检查听
			for _, v := range duiIndex {
				bkc := cards[indexs[v]]
				cards[indexs[v]], cards[indexs[v]+1] = lzCard, lzCard
				_, subTings := rule.doCheckNTing(chiCards, cards, 0, indexs, 0, all)
				tings = append(tings, subTings...)
				if len(tings) > 0 && !all {
					return 0, tings
				}
				cards[indexs[v]], cards[indexs[v]+1] = bkc, bkc
			}
		}
	}

	if len(tings) > 1 {
		sort.Slice(tings, func(i, j int) bool { return tings[i] < tings[j] })
		tings = util.UniqueSlice(tings, true).([]int32)
	}
	return maxBaseN, tings
}

/** CheckNTing2和CheckNTing的返回结果一致，只是两者的遍历方式不同，CheckNTing采用递归深度优先遍历，CheckNTing2采用循环广度优先遍历
 ** 理论上对于实际向听数小于输入向听数，广度优先遍历效率会远高于深度优先遍历，由于CheckNTing采用预排序和剪枝，所以实际效率两者相差不大
 ** 由于存在拆对子的情况，若忽略此情况，有些极端牌型会导致错误结果，而采用广度优先遍历对于打出一对牌检查听牌难以实现，因此采用CheckNTing
**/
func (rule *RuleMahjong) CheckNTing2(chiCards []*ChiCard, holdCards []int32, N int32, all bool) (int32, []int32) {
	if len(holdCards)%3 != 1 || N < 0 {
		return -1, []int32{}
	}

	tings := rule.CheckTing(chiCards, holdCards, all)
	if len(tings) > 0 {
		return 0, tings
	} else if N == 0 {
		return -1, []int32{}
	}

	cards := make([]int32, len(holdCards))
	copy(cards, holdCards)
	sort.Slice(cards, func(i, j int) bool { return cards[i] < cards[j] })

	indexs := make([]int, 0, len(cards))
	for i := 0; i < len(cards); i++ {
		if rule.IsLaiziCard(cards[i]) || i > 0 && cards[i] == cards[i-1] {
			continue
		}
		indexs = append(indexs, i)
	}

	for i := int32(1); i <= N; i++ {
		if tings := rule.doCheckNTing2(chiCards, cards, i, indexs, all); len(tings) > 0 {
			return i, tings
		}
	}
	return -1, []int32{}
}

func (rule *RuleMahjong) doCheckNTing2(chiCards []*ChiCard, cards []int32, N int32, indexs []int, all bool) []int32 {
	group := make([]int, N)
	for i := 0; i < int(N); i++ {
		group[i] = i
	}
	sumTings := []int32{}
	tmpCards := make([]int32, len(cards))
	for {
		copy(tmpCards, cards)
		for i := 0; i < int(N); i++ {
			tmpCards[indexs[group[i]]] = rule.lzCards[0]
		}
		if tings := rule.CheckTing(chiCards, tmpCards, all); len(tings) > 0 {
			sumTings = append(sumTings, tings...)
		}
		if !rule.nextGroup(group, len(indexs)) {
			break
		}
	}
	if len(sumTings) > 1 {
		sort.Slice(sumTings, func(i, j int) bool { return sumTings[i] < sumTings[j] })
		sumTings = util.UniqueSlice(sumTings, true).([]int32)
	}
	return sumTings
}

func (rule *RuleMahjong) nextGroup(group []int, length int) bool {
	if len(group) == 0 || len(group) > length {
		return false
	}
	if group[0] == length-len(group) {
		return false
	}
	for i := len(group) - 1; i >= 0; i-- {
		if group[i] < length-len(group)+i {
			group[i]++
			for j := i + 1; j < len(group); j++ {
				group[j] = group[j-1] + 1
			}
			break
		}
	}
	return true
}

func (rule *RuleMahjong) removeSolKeSeq(cards []int32) []int32 {
	lzCards := []int32{}
	for i := 0; i < len(cards); {
		if rule.IsLaiziCard(cards[i]) {
			lzCards = append(lzCards, cards[i])
			cards = append(cards[:i], cards[i+1:]...)
		} else {
			i++
		}
	}
	sort.Slice(cards, func(i, j int) bool { return cards[i] < cards[j] })
	for i := 0; i < len(cards)-2; {
		sol := (i == 0 || cards[i-1] < cards[i]-2) && (i == len(cards)-3 || cards[i+3] > cards[i+2]+2)
		if cards[i] == cards[i+1] && cards[i+1] == cards[i+2] && (cards[i]%MAHJONG_MASK >= MAHJONG_DONG || sol) {
			cards = append(cards[:i], cards[i+3:]...) //移除独立的刻子
		} else if c := cards[i] % MAHJONG_MASK; c >= MAHJONG_1 && c <= MAHJONG_7 && cards[i] == cards[i+1]-1 && cards[i+1] == cards[i+2]-1 && sol {
			cards = append(cards[:i], cards[i+3:]...) //移除独立的顺子
		} else {
			i++
		}
	}
	return append(cards, lzCards...)
}
