package main

import (
	"crypto/sha256"
	"fmt"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"time"
)

type MerkelNode struct {
	left        *MerkelNode
	right       *MerkelNode
	transaction string
}

type MerkelTree struct {
	Root *MerkelNode
}

func addnode(transac string, arr *[10]MerkelTree, merkelindex *int) {

	var ind int
	ind = *merkelindex
	var tree *MerkelNode = arr[ind].Root

	var newitem MerkelNode

	newitem.transaction = transac
	newitem.right = nil
	newitem.left = nil

	if tree == nil {
		fmt.Println("Inserting nil ", transac, "\n")
		arr[ind].Root = &newitem
		ind++

	} else {

		fmt.Println("Inserting ", transac, "\n")
		insert_item(tree, &newitem)
		ind++
	}

	*merkelindex = ind
}

func insert_item(tree *MerkelNode, item *MerkelNode) {

	if len(item.transaction) <= len(tree.transaction) {

		if tree.left == nil {
			tree.left = item
			return
		} else {
			insert_item(tree.left, item)
			return
		}
	} else if len(item.transaction) > len(tree.transaction) {

		if tree.right == nil {
			tree.right = item
			return
		} else {
			insert_item(tree.right, item)
			return
		}

	}
}

func display_tree(tree *MerkelNode) {

	if tree == nil {
		return
	}
	fmt.Println(tree.transaction, "\n")

	if tree.left != nil {
		display_tree(tree.left)
	}

	if tree.right != nil {
		display_tree(tree.right)
	}
}

func update2(tree *MerkelNode, prev string, now string) {

	if tree == nil {
		return
	}
	if prev == tree.transaction {
		tree.transaction = now
	}

	if tree.left != nil {
		update2(tree.left, prev, now)
	}

	if tree.right != nil {
		update2(tree.right, prev, now)
	}

}

func update(arr *[10]MerkelTree, merkelindex *int, prev string, now string) {

	var ind int
	ind = *merkelindex
	var i = 0
	for i = 0; i < ind+1; i++ {
		var tree *MerkelNode = arr[i].Root

		update2(tree, prev, now)

	}

	*merkelindex = ind

}

func traversal2(tree *MerkelNode, tt *string) {

	if tree == nil {
		return
	}
	*tt += tree.transaction

	if tree.left != nil {
		traversal2(tree.left, tt)
	}

	if tree.right != nil {
		traversal2(tree.right, tt)
	}

}

func traversal(arr *[10]MerkelTree, merkelindex *int) string {

	var ind int
	ind = *merkelindex
	var alltransactions string
	var i = 0
	for i = 0; i < ind+1; i++ {
		var tree *MerkelNode = arr[i].Root

		traversal2(tree, &alltransactions)

	}

	*merkelindex = ind

	return alltransactions

}

func displayMerkelTree(arr *[10]MerkelTree, merkelindex *int) {
	var ind int
	ind = *merkelindex
	var i = 0

	fmt.Println("------------------Merkel Tree elements-------------------------------------")
	for i = 0; i < ind+1; i++ {
		var tree *MerkelNode = arr[i].Root

		if tree == nil {
			fmt.Println("arr[", i, "] has no elements\n")

		} else {
			fmt.Println("arr[", i, "] has elements\n")
			display_tree(tree)
		}
	}

	*merkelindex = ind

}

func createMerkelTree(arr *[10]MerkelTree, merkelindex *int) {

	for i := 0; i < 10; i++ {
		arr[i].Root = nil
	}

	addnode("transaction1", arr, merkelindex)
	addnode("transaction2", arr, merkelindex)
	addnode("transaction3", arr, merkelindex)
	addnode("transaction4", arr, merkelindex)
	addnode("transaction5", arr, merkelindex)
	addnode("transaction6", arr, merkelindex)

}

type block struct {
	arr           [10]MerkelTree
	merkelindex   int
	id            int
	nonce         string
	previous_hash string
	current_hash  string
}

type blockchain struct {
	list []*block
}

func newBlock(x int) *block {
	//fmt.Println("------------------------fdsfdsfdsfds-----------------------------------")
	tempblock := new(block)
	tempblock.id = x
	tempblock.nonce = "0"
	tempblock.merkelindex = 0
	createMerkelTree(&tempblock.arr, &tempblock.merkelindex)
	return tempblock
}

func verifyChain(chain *blockchain) bool {
	var temp = ""
	var check = true
	for i := 0; i < len(chain.list); i++ {
		tt := traversal(&chain.list[i].arr, &chain.list[i].merkelindex)

		var attributes string
		attributes += strconv.Itoa(chain.list[i].id)
		attributes += tt + chain.list[i].previous_hash
		total_sum := sha256.Sum256([]byte(attributes))
		temp = fmt.Sprintf("%x", total_sum)

		if temp != chain.list[i].current_hash {
			check = false
			fmt.Printf("Previous block has been tampered, i.e. Block # %d\n", i)
			break

		}
	}

	if check == false {
		fmt.Println("error occured")
	} else {
		fmt.Printf("Blocks verified. No tampering\n")
	}
	return check
}

func Mineblock(blocklist *blockchain) {

	for j := 0; j < len(blocklist.list); j++ {
		print("to match:", blocklist.list[j].current_hash, "\n")
		for i := 0; ; i++ {
			temp := sha256.Sum256([]byte(strconv.Itoa(i)))
			noncex := fmt.Sprintf("%x", temp)
			dum := noncex[:3]
			fmt.Println("dum:", dum)
			fmt.Println(strings.Contains(blocklist.list[j].current_hash, dum))

			if strings.Contains(blocklist.list[j].current_hash, dum) == true {
				blocklist.list[j].nonce = dum
				break

			}

		}

	}

}

func CalculateHash(chain *blockchain) {

	for i := 0; i < len(chain.list); i++ {
		tt := traversal(&chain.list[i].arr, &chain.list[i].merkelindex)
		var attributes string
		attributes += strconv.Itoa(chain.list[i].id)
		attributes += tt + chain.list[i].previous_hash
		total_sum := sha256.Sum256([]byte(attributes))
		chain.list[i].current_hash = fmt.Sprintf("%x", total_sum) // formating to string
		if i < len(chain.list)-1 {
			chain.list[i+1].previous_hash = fmt.Sprintf("%x", total_sum) //storing current block hash to next block in its previous hash var
		}

	}
}

func (blocklist *blockchain) addblock(x int) *block {
	tempblock := newBlock(x)

	if verifyChain(blocklist) {
		blocklist.list = append(blocklist.list, tempblock)
		CalculateHash(blocklist)

		fmt.Printf("block addition in chain successful\n")
	} else {
		fmt.Printf(" error. block addition unsuccessful.\n")
		return nil
	}
	return tempblock
}

func DisplayBlocks(blocklist *blockchain) {
	fmt.Println("")

	for i := 0; i < len(blocklist.list); i++ {
		fmt.Printf("Block id:%d\n\n", blocklist.list[i].id)
		displayMerkelTree(&blocklist.list[i].arr, &blocklist.list[i].merkelindex)
		fmt.Println("nonce value : \n", blocklist.list[i].nonce)
		fmt.Println("current hash: \n", blocklist.list[i].current_hash)
		fmt.Println("previous hash: \n", blocklist.list[i].previous_hash)

	}

	fmt.Println("")

}

func changeBlock(chain *blockchain, x int) { // updating on basis of id value as identifier

	found := false
	for i := 0; i < len(chain.list); i++ {

		if x == chain.list[i].id {

			var now string
			var prev string
			fmt.Println("Enter transaction to change\n")

			fmt.Scanln(&prev)

			fmt.Println("Enter updated value\n")

			fmt.Scanln(&now)
			fmt.Println("updated successfully\n")
			update(&chain.list[i].arr, &chain.list[i].merkelindex, prev, now)
			found = true
		}
	}
	if found == false {
		fmt.Println("error. Couldnt update. block not found")
	}
	return
}

func main() {

	var totalnodes = 5
	var delay = 1
	var nodes []string
	ln, err := net.Listen("tcp", ":8001")
	if err != nil {

	}
	for i := 0; i < totalnodes; i++ {
		conn, err := ln.Accept()
		if err != nil { //connection not successfull
			continue
		}

		Channelnode := make(chan string)

		go handleConnection(conn, Channelnode, nodes) // go routine

		var tempstr string

		for i := range Channelnode {
			tempstr = i
			break
		}
		nodes = append(nodes, tempstr)

		fmt.Println("Node connected to server and registered")
	}
	//var ip string = strings.Split(nodes[0], ":")[0]
	var transactionlist [5]string
	transactionlist[0] = "transaction1"
	transactionlist[1] = "transaction2"
	transactionlist[2] = "transaction3"
	transactionlist[3] = "transaction4"
	transactionlist[4] = "transaction5"

	for j := 0; j < 5; j++ {

		for i := 0; i < 5; i++ {
			var portstr string = strings.Split(nodes[i], ":")[1]
			portstr = "localhost:" + portstr
			fmt.Println("portstr:", portstr)
			conn2, err := net.Dial("tcp", portstr)
			if err != nil {
				// error handling
			}
			fmt.Println("server connected to client 1")
			conn2.Write([]byte(transactionlist[j]))
			time.Sleep(time.Duration(delay) * time.Second)

		}
	}

}

func handleConnection(conn net.Conn, Channelnode chan string, nodes []string) {

	connstr := conn.RemoteAddr().String()
	port := strings.Split(connstr, "]")
	port = strings.Split(port[1], ":")
	var node string
	node = "127.0.0.1:" + port[1] //Add IP and Port in string and send channel
	Channelnode <- node
	var message string = ""
	if len(nodes) == 0 {
		fmt.Print("Port number for first node:", port[1])
		message = "none"
	} else if len(nodes) < 3 { //1,2

		for i := 0; i < len(nodes)-1; i++ { // send addresses
			message += nodes[i]
			message += ","
		}
		message += nodes[len(nodes)-1]

	} else {

		rand.Seed(time.Now().UnixNano()) //give random 2 node address
		var n1, n2 int
		n1 = rand.Intn(100) % len(nodes)
		for {
			n2 = rand.Intn(100) % len(nodes)
			if n2 != n1 {
				break
			}
		}
		message = nodes[n1] + "," + nodes[n2]
		fmt.Println(n1, n2)
	}
	message += "\000"
	conn.Write([]byte(message))
	fmt.Println("Sending: ", message+" to client")

}
