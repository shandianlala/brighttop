package main

import "fmt"

func main3() {
	a := [][]int{{3, 0, 1, 4, 2}, {5, 6, 3, 2, 1}, {1, 2, 0, 1, 5}, {4, 1, 0, 1, 7}, {1, 0, 3, 0, 5}}
	var numMatrix = Constructor(a)
	region := numMatrix.SumRegion(2, 1, 4, 3) // return 8 (红色矩形框的元素总和)
	fmt.Printf("region1 = %d, must be 8\n", region)

	region = numMatrix.SumRegion(1, 1, 2, 2) // return 11 (绿色矩形框的元素总和)
	fmt.Printf("region = %d, must be 11\n", region)

	region = numMatrix.SumRegion(1, 2, 2, 4) // return 12 (蓝色矩形框的元素总和)
	fmt.Printf("region3 = %d, must be 12\n", region)

	//a := [][]int{{-4, -5}}
	//var numMatrix = Constructor(a)
	//region1 := numMatrix.SumRegion(0, 0, 0, 0) // return -4 (红色矩形框的元素总和)
	//fmt.Printf("[0,0,0,0] region1 = %d, \n", region1)
	//region2 := numMatrix.SumRegion(0, 0, 0, 1) // return -9 (红色矩形框的元素总和)
	//fmt.Printf("[0,0,0,1] region2 = %d, \n", region2)
	//region3 := numMatrix.SumRegion(0, 1, 0, 1) // return -5 (红色矩形框的元素总和)
	//fmt.Printf("[0,1,0,1] region3 = %d, \n", region3)
	// [[[[-4,-5]]],[0,0,0,0],[0,0,0,1],[0,1,0,1]]

}

type NumMatrix struct {
	sum [][]int
}

func Constructor(matrix [][]int) NumMatrix {
	if matrix == nil {
		return NumMatrix{}
	}
	var yLen = len(matrix[0])
	var sum = make([][]int, len(matrix))
	for i, xRows := range matrix {
		var xRow = make([]int, yLen)
		for j, value := range xRows {
			if j == 0 {
				xRow[j] = value
			} else {
				xRow[j] += value + xRow[j-1]
			}
		}
		sum[i] = xRow
	}
	for _, ints := range sum {
		for _, i3 := range ints {
			fmt.Printf("%d \t", i3)
		}
		fmt.Printf("\n")
	}
	var result = NumMatrix{sum: sum}
	return result
}

func (this *NumMatrix) SumRegion(row1 int, col1 int, row2 int, col2 int) int {
	var res int
	for i := row1; i <= row2; i++ {
		if col1 == 0 {
			res += this.sum[i][col2]
		} else {
			var leftCol = col1 - 1
			res += this.sum[i][col2] - this.sum[i][leftCol]
		}
	}
	return res
}

/**
 * Your NumMatrix object will be instantiated and called as such:
 * obj := Constructor(matrix);
 * param_1 := obj.SumRegion(row1,col1,row2,col2);
 */
