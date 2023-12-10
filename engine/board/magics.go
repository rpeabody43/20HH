package board

import (
	"fmt"
	"math/rand"
)

var RookMagics = [64]uint64{
	0x0600110040a20080, 0x0006043008204028, 0x0100042000084010, 0x1200021040582104,
	0x020002ac08102003, 0x6100008821240012, 0x0500028420c20001, 0xc0800020c0800100,
	0x0012800140001020, 0x0281402010004000, 0x201d004302600010, 0x0813002901201000,
	0x0820808004000800, 0x8800808004008200, 0x0494001004028308, 0x0042000100240382,
	0x0080004000102000, 0x0000c20024860100, 0x0008808010002000, 0x0115010008100020,
	0x1008910008010d00, 0x4102008002800400, 0x0100050100020004, 0x0004008008482100,
	0x0000302e80004000, 0x0420100040400020, 0x0401004100102000, 0x0820100100221900,
	0x0000180080800400, 0x0c02000200291004, 0x2100488400301201, 0x082280810000c422,
	0x0518164000800020, 0x8050004000402008, 0x0112002082004050, 0x3244080280801000,
	0x0000800800804400, 0x0044000600801480, 0x480e000802000104, 0x0000009c02000047,
	0x0800902840008000, 0x04a0005000604008, 0x0012012041820010, 0x0320900100890020,
	0x000208000c008080, 0x0042008850420004, 0x0900860810240003, 0x0811000062408009,
	0x0460422080110200, 0x00810aa040048100, 0x0000220041128200, 0x4000084200211200,
	0x0080300800050100, 0x00002200800c0080, 0x0801000442000300, 0x00003060810c0a00,
	0x0040800247002491, 0x0040400280481121, 0x40001080a008401e, 0x0406281000ac0221,
	0x0002080004229009, 0x0040810804902412, 0x0408080041809004, 0x202a04c420851402,
}

var RookShifts = [64]uint8{
	0x52, 0x52, 0x52, 0x52, 0x52, 0x52, 0x52, 0x52,
	0x52, 0x54, 0x54, 0x54, 0x54, 0x54, 0x54, 0x52,
	0x52, 0x54, 0x54, 0x54, 0x54, 0x54, 0x54, 0x52,
	0x52, 0x54, 0x54, 0x54, 0x54, 0x54, 0x54, 0x52,
	0x52, 0x54, 0x54, 0x54, 0x54, 0x54, 0x54, 0x52,
	0x52, 0x54, 0x54, 0x54, 0x54, 0x54, 0x54, 0x52,
	0x52, 0x54, 0x54, 0x54, 0x54, 0x54, 0x54, 0x52,
	0x52, 0x52, 0x52, 0x52, 0x52, 0x52, 0x52, 0x52,
}

var BishopMagics = [64]uint64{
	0x3440040400802100, 0x2102026212160101, 0x50840c0400411000, 0x00082080e4184224,
	0x0002121020000201, 0x088082a060000000, 0x0010421010090000, 0x812042080c010c00,
	0x0204104488082050, 0x0000084801042220, 0x0128101415802036, 0x004008205040100d,
	0x0021440420000280, 0x8a10820824041004, 0x0001110430048c00, 0x2100004048045031,
	0x0804802020120202, 0x0003002008060884, 0x0010082204001020, 0x0228002082004026,
	0x0304821400a04002, 0x10a2002158020801, 0x4002410212022000, 0x90e0800044008800,
	0x00280c2042040800, 0x0008040408090808, 0x0484020010008011, 0x80e4100808008010,
	0x00ab0600c4008400, 0x08080880091004a2, 0x00822c0226018202, 0x0c21104001220818,
	0x4803604803200821, 0x00820210b0200724, 0x0004014c00084404, 0x0070200901180104,
	0x4040110200090048, 0x018400898028080c, 0x0001480300020100, 0x0101020202208240,
	0x08180208202c8e26, 0x04808a8818002000, 0x4000201410000208, 0x000620420280080c,
	0x2a40280101023010, 0x2024200c01c20408, 0x0020830401000094, 0x0508080043440280,
	0x0000450411401010, 0x4402012082300804, 0x1204020110a84000, 0x8000001060a80800,
	0x030084300a061024, 0x2000060c482e0045, 0x300420040400a000, 0xa00a448400820101,
	0x004d008801013004, 0x0200108404024640, 0x0010084202521817, 0x9843901030218800,
	0xa090040008030400, 0x0180020932180204, 0x80002c10900201a0, 0x8010600200420020,
}

var BishopShifts = [64]uint8{
	0x58, 0x59, 0x59, 0x59, 0x59, 0x59, 0x59, 0x58,
	0x59, 0x59, 0x59, 0x59, 0x59, 0x59, 0x59, 0x59,
	0x59, 0x59, 0x57, 0x57, 0x57, 0x57, 0x59, 0x59,
	0x59, 0x59, 0x57, 0x55, 0x55, 0x57, 0x59, 0x59,
	0x59, 0x59, 0x57, 0x55, 0x55, 0x57, 0x59, 0x59,
	0x59, 0x59, 0x57, 0x57, 0x57, 0x57, 0x59, 0x59,
	0x59, 0x59, 0x59, 0x59, 0x59, 0x59, 0x59, 0x59,
	0x58, 0x59, 0x59, 0x59, 0x59, 0x59, 0x59, 0x58,
}

func sparseRandInt() uint64 {
	return rand.Uint64() & rand.Uint64() & rand.Uint64()
}

func findMagicAtSq(sq Square, rook bool) (uint64, []Bitboard) {
	blockerMask := RookBlockerMasks[sq]
	if !rook {
		blockerMask = BishopBlockerMasks[sq]
	}
	numBlockerBits := countBits(blockerMask)
	permutations := blockerPermutations(blockerMask)
	if rook {
		RookShifts[sq] = 64 - numBlockerBits
	} else {
		BishopShifts[sq] = 64 - numBlockerBits
	}

	foundNum := false
	var magic uint64
	attacksTable := make([]Bitboard, 1<<numBlockerBits)
	for !foundNum {
		magic = sparseRandInt()
		foundNum = true

		for _, perm := range permutations {
			var attacks Bitboard
			var index uint64
			if rook {
				index = (uint64(perm) * magic) >> RookShifts[sq]
				attacks = rookMovesFromBlockers(sq, perm)
			} else {
				index = (uint64(perm) * magic) >> BishopShifts[sq]
				attacks = bishopMovesFromBlockers(sq, perm)
			}

			if attacksTable[index] == 0 {
				attacksTable[index] = attacks
			} else if attacksTable[index] != attacks {
				// muggle number
				foundNum = false
				attacksTable = make([]Bitboard, 1<<numBlockerBits)
				break
			}
		}
	}
	return magic, attacksTable
}

// Copied almost exactly from the method above, just uses
// precomputed magics
func setupRookTable(sq Square) {
	blockerMask := RookBlockerMasks[sq]
	numBlockerBits := countBits(blockerMask)
	permutations := blockerPermutations(blockerMask)
	RookShifts[sq] = 64 - numBlockerBits

	magic := RookMagics[sq]
	RookAttacks[sq] = make([]Bitboard, 1<<numBlockerBits)

	for _, perm := range permutations {
		index := (uint64(perm) * magic) >> RookShifts[sq]
		attacks := rookMovesFromBlockers(sq, perm)

		if RookAttacks[sq][index] == 0 {
			RookAttacks[sq][index] = attacks
		}
	}
}

func setupBishopTable(sq Square) {
	blockerMask := BishopBlockerMasks[sq]
	numBlockerBits := countBits(blockerMask)
	permutations := blockerPermutations(blockerMask)
	BishopShifts[sq] = 64 - numBlockerBits

	magic := BishopMagics[sq]
	BishopAttacks[sq] = make([]Bitboard, 1<<numBlockerBits)

	for _, perm := range permutations {
		index := (uint64(perm) * magic) >> BishopShifts[sq]
		attacks := bishopMovesFromBlockers(sq, perm)

		if BishopAttacks[sq][index] == 0 {
			BishopAttacks[sq][index] = attacks
		}
	}
}

func countBits(bb Bitboard) uint8 {
	ret := uint8(0)
	for bb != 0 {
		bb &= bb - 1
		ret++
	}
	return ret
}

func blockerPermutations(blockerMask Bitboard) []Bitboard {
	bitsInMask := countBits(blockerMask)
	numPermutations := 1 << bitsInMask
	permutations := make([]Bitboard, numPermutations)

	for permIdx := 0; permIdx < numPermutations; permIdx++ {
		permutation := Bitboard(0)
		blockerMaskCopy := blockerMask
		for bitIdx := uint8(0); bitIdx < bitsInMask; bitIdx++ {
			bit := (permIdx >> bitIdx) & 1
			permutation |= Bitboard(bit << blockerMaskCopy.PopLSB())
		}
		permutations[permIdx] = permutation
	}
	return permutations
}

func printMagics(table [64]uint64, varName string) {
	fmt.Printf("var %s = [64]uint64{", varName)
	for sq := 0; sq < 64; sq++ {
		if sq%4 == 0 {
			fmt.Print("\n\t")
		}
		fmt.Printf("0x%016x, ", table[sq])
	}
	fmt.Println("\n}")
	fmt.Println()
}

func printShifts(table [64]uint8, varName string) {
	fmt.Printf("var %s = [64]uint8{", varName)
	for sq := 0; sq < 64; sq++ {
		if sq%8 == 0 {
			fmt.Print("\n\t")
		}
		fmt.Printf("0x%d, ", table[sq])
	}
	fmt.Println("\n}")
	fmt.Println()
}

func printAllMagics() {
	printMagics(RookMagics, "RookMagics")
	printShifts(RookShifts, "RookShifts")
	printMagics(BishopMagics, "BishopMagics")
	printShifts(BishopShifts, "BishopShifts")
}
